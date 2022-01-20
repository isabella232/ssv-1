package exporter

import (
	"encoding/hex"
	"github.com/bloxapp/ssv/eth1"
	"github.com/bloxapp/ssv/eth1/abiparser"
	"github.com/bloxapp/ssv/exporter/api"
	"github.com/bloxapp/ssv/exporter/storage"
	"github.com/bloxapp/ssv/validator"
	validatorstorage "github.com/bloxapp/ssv/validator/storage"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"sync/atomic"
)

// handleEth1Event handles the given eth1 event
func (exp *exporter) handleEth1Event(e *eth1.Event) error {
	var err error = nil
	if validatorAddedEvent, ok := e.Data.(abiparser.ValidatorAddedEvent); ok {
		err = exp.handleValidatorAddedEvent(validatorAddedEvent)
	} else if opertaorAddedEvent, ok := e.Data.(abiparser.OperatorAddedEvent); ok {
		err = exp.handleOperatorAddedEvent(opertaorAddedEvent)
	}
	return err
}

// handleValidatorAddedEvent parses the given event and sync the ibft-data of the validator
func (exp *exporter) handleValidatorAddedEvent(event abiparser.ValidatorAddedEvent) error {
	validatorShare, vi, err := exp.saveNewValidator(event, 0)
	if err != nil {
		return errors.Wrap(err, "failed to add validator share")
	}

	go func() {
		n := exp.ws.BroadcastFeed().Send(api.Message{
			Type:   api.TypeValidator,
			Filter: api.MessageFilter{From: vi.Index, To: vi.Index},
			Data:   []storage.ValidatorInformation{*vi},
		})
		exp.logger.Debug("msg was sent on outbound feed", zap.Int("num of subscribers", n),
			zap.String("eventType", "ValidatorAdded"),
			zap.String("pubKey", hex.EncodeToString(event.PublicKey)))
	}()

	if err = exp.triggerValidator(validatorShare.PublicKey); err != nil {
		return errors.Wrap(err, "failed to trigger ibft sync")
	}

	return nil
}

// registrySyncHandler returns a handler for registry contract events coming from sync
func (exp *exporter) registrySyncHandler(startIndex int64) eth1.EventHandler {
	var counter int64
	return func(e *eth1.Event) error {
		var err error = nil
		var vi *storage.ValidatorInformation
		if validatorAddedEvent, ok := e.Data.(abiparser.ValidatorAddedEvent); ok {
			i := atomic.AddInt64(&counter, 1)
			_, vi, err = exp.saveNewValidator(validatorAddedEvent, i+startIndex)
			exp.logger.Debug("new validator", zap.String("pk", hex.EncodeToString(validatorAddedEvent.PublicKey)),
				zap.Int64("i", i), zap.Int64("index", vi.Index))
		} else if opertaorAddedEvent, ok := e.Data.(abiparser.OperatorAddedEvent); ok {
			err = exp.handleOperatorAddedEvent(opertaorAddedEvent)
		}
		return err
	}
}

// handleValidatorAddedEvent parses the given event and sync the ibft-data of the validator
func (exp *exporter) saveNewValidator(event abiparser.ValidatorAddedEvent, i int64) (*validatorstorage.Share, *storage.ValidatorInformation, error) {
	pubKeyHex := hex.EncodeToString(event.PublicKey)
	// create a share to be used in IBFT
	share, _, err := validator.ShareFromValidatorAddedEvent(event, "")
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not create a share from ValidatorAddedEvent")
	}
	if err := exp.addValidatorShare(share); err != nil {
		return nil, nil, errors.Wrap(err, "failed to add validator share")
	}
	// if share was created, save information for exporting validators
	info, err := toValidatorInformation(event)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not create validator information")
	}
	info.Index = i
	if err := exp.storage.SaveValidatorInformation(info); err != nil {
		return nil, nil, errors.Wrap(err, "failed to save validator information")
	}
	exp.logger.Debug("validator information was saved", zap.Any("index", info.Index), zap.String("pubKey", pubKeyHex))
	return share, info, nil
}

// handleOperatorAddedEvent parses the given event and saves operator information
func (exp *exporter) handleOperatorAddedEvent(event abiparser.OperatorAddedEvent) error {
	logger := exp.logger.With(zap.String("eventType", "OperatorAdded"),
		zap.String("pubKey", string(event.PublicKey)))
	logger.Info("operator added event")
	oi := storage.OperatorInformation{
		PublicKey:    string(event.PublicKey),
		Name:         event.Name,
		OwnerAddress: event.OwnerAddress,
	}
	err := exp.storage.SaveOperatorInformation(&oi)
	if err != nil {
		return err
	}
	logger.Debug("managed to save operator information", zap.Any("value", oi))
	reportOperatorIndex(exp.logger, &oi)

	go func() {
		n := exp.ws.BroadcastFeed().Send(api.Message{
			Type:   api.TypeOperator,
			Filter: api.MessageFilter{From: oi.Index, To: oi.Index},
			Data:   []storage.OperatorInformation{oi},
		})
		logger.Debug("msg was sent on outbound feed", zap.Int("num of subscribers", n))
	}()

	return nil
}

// addValidatorShare is called upon ValidatorAdded to add a new share to storage
func (exp *exporter) addValidatorShare(validatorShare *validatorstorage.Share) error {
	logger := exp.logger.With(zap.String("eventType", "ValidatorAdded"),
		zap.String("pubKey", validatorShare.PublicKey.SerializeToHexStr()))
	// add metadata
	if updated, err := validator.UpdateShareMetadata(validatorShare, exp.beacon); err != nil {
		logger.Warn("could not add validator metadata", zap.Error(err))
	} else if !updated {
		logger.Warn("could not find validator metadata")
	} else {
		logger.Debug("validator metadata was updated", zap.Any("updated", updated))
	}
	if err := exp.validatorStorage.SaveValidatorShare(validatorShare); err != nil {
		return errors.Wrap(err, "failed to save validator share")
	}
	return nil
}

// addValidatorInformation is called upon ValidatorAdded to create and add validator information
func (exp *exporter) addValidatorInformation(event abiparser.ValidatorAddedEvent) (*storage.ValidatorInformation, error) {
	vi, err := toValidatorInformation(event)
	if err != nil {
		return nil, errors.Wrap(err, "could not create ValidatorInformation")
	}
	if err := exp.storage.SaveValidatorInformation(vi); err != nil {
		return nil, errors.Wrap(err, "failed to save validator information")
	}
	return vi, nil
}

// toValidatorInformation converts raw event to ValidatorInformation
func toValidatorInformation(validatorAddedEvent abiparser.ValidatorAddedEvent) (*storage.ValidatorInformation, error) {
	pubKey := &bls.PublicKey{}
	if err := pubKey.Deserialize(validatorAddedEvent.PublicKey); err != nil {
		return nil, errors.Wrap(err, "failed to deserialize validator public key")
	}

	var operators []storage.OperatorNodeLink
	for i, operatorPublicKey := range validatorAddedEvent.OperatorPublicKeys {
		nodeID := uint64(i + 1)
		operators = append(operators, storage.OperatorNodeLink{
			ID: nodeID, PublicKey: string(operatorPublicKey),
		})
	}

	vi := storage.ValidatorInformation{
		PublicKey: pubKey.SerializeToHexStr(),
		Operators: operators,
	}

	return &vi, nil
}

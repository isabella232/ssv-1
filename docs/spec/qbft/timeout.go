package qbft

import "github.com/pkg/errors"

func UponTimout(state State, net Network) error {
	state.SetRound(state.GetRound() + 1)
	roundChange := CreateRoundChange(state)

	if err := net.BroadcastSignedMessage(roundChange); err != nil {
		return errors.Wrap(err, "failed to broadcast round change message")
	}

	return nil
}

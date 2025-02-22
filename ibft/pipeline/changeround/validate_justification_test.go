package changeround

import (
	"encoding/json"
	"github.com/bloxapp/ssv/ibft/proto"
	"github.com/bloxapp/ssv/validator/storage"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

// GenerateNodes generates randomly nodes
func GenerateNodes(cnt int) (map[uint64]*bls.SecretKey, map[uint64]*proto.Node) {
	_ = bls.Init(bls.BLS12_381)
	nodes := make(map[uint64]*proto.Node)
	sks := make(map[uint64]*bls.SecretKey)
	for i := 0; i < cnt; i++ {
		sk := &bls.SecretKey{}
		sk.SetByCSPRNG()

		nodes[uint64(i)] = &proto.Node{
			IbftId: uint64(i),
			Pk:     sk.GetPublicKey().Serialize(),
		}
		sks[uint64(i)] = sk
	}
	return sks, nodes
}

func aggregateSign(t *testing.T, sks map[uint64]*bls.SecretKey, msg *proto.Message) *proto.SignedMessage {
	var aggSignedMsg *proto.SignedMessage
	for index, sk := range sks {
		sigend := SignMsg(t, index, sk, msg)

		if aggSignedMsg == nil {
			aggSignedMsg = sigend
		} else {
			require.NoError(t, aggSignedMsg.Aggregate(sigend))
		}
	}
	return aggSignedMsg
}

// SignMsg signs the given message by the given private key
func SignMsg(t *testing.T, id uint64, sk *bls.SecretKey, msg *proto.Message) *proto.SignedMessage {
	bls.Init(bls.BLS12_381)

	signature, err := msg.Sign(sk)
	require.NoError(t, err)
	return &proto.SignedMessage{
		Message:   msg,
		Signature: signature.Serialize(),
		SignerIds: []uint64{id},
	}
}

func changeRoundDataToBytes(input *proto.ChangeRoundData) []byte {
	ret, _ := json.Marshal(input)
	return ret
}

func TestValidateChangeRound(t *testing.T) {
	sks, nodes := GenerateNodes(4)
	share := &storage.Share{
		Committee: nodes,
	}
	pip := Validate(share)

	validSig := aggregateSign(t, sks, &proto.Message{
		Type:      proto.RoundState_Prepare,
		Round:     2,
		SeqNumber: 12,
		Lambda:    []byte("Lambda"),
		Value:     []byte("value"),
	})

	tests := []struct {
		name string
		err  string
		msg  *proto.SignedMessage
	}{
		{
			"valid nil prepared change round",
			"",
			SignMsg(t, 1, sks[1], &proto.Message{
				Type:      proto.RoundState_ChangeRound,
				Round:     3,
				SeqNumber: 12,
				Lambda:    []byte("Lambda"),
				Value: changeRoundDataToBytes(&proto.ChangeRoundData{
					JustificationMsg: nil,
				}),
			}),
		},
		{
			"valid prepared change round",
			"",
			SignMsg(t, 1, sks[1], &proto.Message{
				Type:      proto.RoundState_ChangeRound,
				Round:     3,
				SeqNumber: 12,
				Lambda:    []byte("Lambda"),
				Value: changeRoundDataToBytes(&proto.ChangeRoundData{
					JustificationMsg: &proto.Message{
						Type:      proto.RoundState_Prepare,
						Round:     2,
						SeqNumber: 12,
						Lambda:    []byte("Lambda"),
						Value:     []byte("value"),
					},
					JustificationSig: validSig.Signature,
					PreparedRound:    2,
					PreparedValue:    []byte("value"),
					SignerIds:        []uint64{0, 1, 2, 3},
				}),
			}),
		},
		{
			"invalid prepared change round",
			"change round justification signature doesn't verify",
			SignMsg(t, 1, sks[1], &proto.Message{
				Type:      proto.RoundState_ChangeRound,
				Round:     3,
				SeqNumber: 12,
				Lambda:    []byte("Lambda"),
				Value: changeRoundDataToBytes(&proto.ChangeRoundData{
					JustificationMsg: &proto.Message{
						Type:      proto.RoundState_Prepare,
						Round:     2,
						SeqNumber: 12,
						Lambda:    []byte("Lambda"),
						Value:     []byte("value"),
					},
					JustificationSig: validSig.Signature,
					PreparedRound:    2,
					PreparedValue:    []byte("value"),
					SignerIds:        []uint64{1, 2, 3},
				}),
			}),
		},
		{
			"wrong sequence prepared change round",
			"change round justification sequence is wrong",
			SignMsg(t, 1, sks[1], &proto.Message{
				Type:      proto.RoundState_ChangeRound,
				Round:     3,
				SeqNumber: 12,
				Lambda:    []byte("Lambda"),
				Value: changeRoundDataToBytes(&proto.ChangeRoundData{
					JustificationMsg: &proto.Message{
						Type:      proto.RoundState_Prepare,
						Round:     2,
						SeqNumber: 11,
						Lambda:    []byte("Lambda"),
						Value:     []byte("value"),
					},
					JustificationSig: validSig.Signature,
					PreparedRound:    2,
					PreparedValue:    []byte("value"),
					SignerIds:        []uint64{0, 1, 2, 3},
				}),
			}),
		},
		{
			"missing change-round-data",
			"change round justification msg is nil",
			SignMsg(t, 1, sks[1], &proto.Message{
				Type:      proto.RoundState_ChangeRound,
				Round:     1,
				SeqNumber: 12,
				Lambda:    []byte("Lambda"),
				Value:     nil,
			}),
		},
		{
			"non-valid change-round-data json",
			"invalid character 'o' in literal null (expecting 'u')",
			SignMsg(t, 1, sks[1], &proto.Message{
				Type:      proto.RoundState_ChangeRound,
				Round:     1,
				SeqNumber: 12,
				Lambda:    []byte("Lambda"),
				Value:     []byte("non-valid"),
			}),
		},
		{
			"justification type is not prepared",
			"change round justification msg type not Prepare",
			SignMsg(t, 1, sks[1], &proto.Message{
				Type:      proto.RoundState_ChangeRound,
				Round:     1,
				SeqNumber: 12,
				Lambda:    []byte("Lambda"),
				Value: changeRoundDataToBytes(&proto.ChangeRoundData{
					JustificationMsg: &proto.Message{
						Type: proto.RoundState_PrePrepare,
					},
					PreparedValue: []byte("value"),
				}),
			}),
		},
		{
			"small justification round",
			"change round justification round lower or equal to message round",
			SignMsg(t, 1, sks[1], &proto.Message{
				Type:      proto.RoundState_ChangeRound,
				Round:     1,
				SeqNumber: 12,
				Lambda:    []byte("Lambda"),
				Value: changeRoundDataToBytes(&proto.ChangeRoundData{
					JustificationMsg: &proto.Message{
						Type:      proto.RoundState_Prepare,
						Round:     1,
						SeqNumber: 12,
					},
					PreparedRound: 0,
					PreparedValue: []byte("value"),
				}),
			}),
		},
		{
			"bad lambda",
			"change round justification msg Lambda not equal to msg Lambda not equal to instance lambda",
			SignMsg(t, 1, sks[1], &proto.Message{
				Type:      proto.RoundState_ChangeRound,
				Round:     3,
				SeqNumber: 12,
				Lambda:    []byte("Lambda"),
				Value: changeRoundDataToBytes(&proto.ChangeRoundData{
					JustificationMsg: &proto.Message{
						Type:      proto.RoundState_Prepare,
						Round:     2,
						SeqNumber: 12,
						Lambda:    []byte("xxx"),
					},
					PreparedRound: 2,
					PreparedValue: []byte("value"),
				}),
			}),
		},
		{
			"bad value",
			"change round prepared value not equal to justification msg value",
			SignMsg(t, 1, sks[1], &proto.Message{
				Type:      proto.RoundState_ChangeRound,
				Round:     3,
				SeqNumber: 12,
				Lambda:    []byte("Lambda"),
				Value: changeRoundDataToBytes(&proto.ChangeRoundData{
					JustificationMsg: &proto.Message{
						Type:      proto.RoundState_Prepare,
						Round:     2,
						SeqNumber: 12,
						Lambda:    []byte("Lambda"),
						Value:     []byte("xxx"),
					},
					PreparedRound: 2,
					PreparedValue: []byte("value"),
				}),
			}),
		},
		{
			"insufficient number of signers",
			"change round justification does not constitute a quorum",
			SignMsg(t, 1, sks[1], &proto.Message{
				Type:      proto.RoundState_ChangeRound,
				Round:     3,
				SeqNumber: 12,
				Lambda:    []byte("Lambda"),
				Value: changeRoundDataToBytes(&proto.ChangeRoundData{
					JustificationMsg: &proto.Message{
						Type:      proto.RoundState_Prepare,
						Round:     2,
						SeqNumber: 12,
						Lambda:    []byte("Lambda"),
						Value:     []byte("value"),
					},
					PreparedRound: 2,
					PreparedValue: []byte("value"),
					SignerIds:     []uint64{1, 2},
				}),
			}),
		},
		{
			"duplicated signers",
			"change round could not verify signature: signers are not unique",
			SignMsg(t, 1, sks[1], &proto.Message{
				Type:      proto.RoundState_ChangeRound,
				Round:     3,
				SeqNumber: 12,
				Lambda:    []byte("Lambda"),
				Value: changeRoundDataToBytes(&proto.ChangeRoundData{
					JustificationMsg: &proto.Message{
						Type:      proto.RoundState_Prepare,
						Round:     2,
						SeqNumber: 12,
						Lambda:    []byte("Lambda"),
						Value:     []byte("value"),
					},
					PreparedRound: 2,
					PreparedValue: []byte("value"),
					SignerIds:     []uint64{1, 2, 2},
				}),
			}),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := pip.Run(test.msg)
			if len(test.err) > 0 {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

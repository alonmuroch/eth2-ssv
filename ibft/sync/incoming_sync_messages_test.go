package sync

import (
	"github.com/bloxapp/ssv/ibft/proto"
	"github.com/bloxapp/ssv/network"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestNewReqHandler(t *testing.T) {
	tests := []struct {
		name            string
		syncChanObj     *network.SyncChanObj
		highestSent     *proto.SignedMessage
		expectedHighest *proto.SignedMessage
	}{
		{
			"valid",
			&network.SyncChanObj{
				Msg: &network.SyncMessage{
					Type: network.Sync_GetHighestType,
				},
				Stream: nil,
			},
			&proto.SignedMessage{
				Message: &proto.Message{
					Type:      proto.RoundState_Decided,
					SeqNumber: 10,
				},
			},
			&proto.SignedMessage{
				Message: &proto.Message{
					Type:      proto.RoundState_Decided,
					SeqNumber: 10,
				},
			},
		},
		{
			"invalid type",
			&network.SyncChanObj{
				Msg: &network.SyncMessage{
					Type: network.Sync_GetInstanceRange,
				},
				Stream: nil,
			},
			&proto.SignedMessage{
				Message: &proto.Message{
					Type:      proto.RoundState_Decided,
					SeqNumber: 10,
				},
			},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			net := NewTestNetwork(t)
			storage := NewTestStorage(test.highestSent)
			handler := NewReqHandler(zap.L(), net, storage)
			//stream := NewTestStream()
			handler.Process(test.syncChanObj)

			if test.expectedHighest == nil {
				require.Nil(t, net.highestDecidedReceived)
			} else {
				require.NotNil(t, net.highestDecidedReceived)
				require.EqualValues(t, test.expectedHighest.Message.Type, net.highestDecidedReceived.Message.Type)
				require.EqualValues(t, test.expectedHighest.Message.SeqNumber, net.highestDecidedReceived.Message.SeqNumber)
			}
		})
	}
}

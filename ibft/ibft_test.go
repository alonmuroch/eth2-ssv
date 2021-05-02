package ibft

import (
	"encoding/hex"
	"github.com/bloxapp/ssv/ibft/proto"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func testIBFTInstance(t *testing.T) *ibftImpl {
	return &ibftImpl{
		instances:           make(map[string]*Instance),
		currentInstanceLock: &sync.Mutex{},
	}
}

func TestCanStartNewInstance(t *testing.T) {
	tests := []struct {
		name          string
		opts          StartOptions
		prevInstances map[string]*Instance
		expectedError string
	}{
		{
			"valid start",
			StartOptions{
				PrevInstance: FirstInstanceIdentifier(),
				Identifier:   []byte{1, 2, 3, 4},
			},
			make(map[string]*Instance),
			"",
		},
		{
			"unknown prev",
			StartOptions{
				PrevInstance: []byte{5, 5, 5, 5},
				Identifier:   []byte{1, 2, 3, 4},
			},
			make(map[string]*Instance),
			"previous instance not found",
		},
		{
			"valid prev",
			StartOptions{
				PrevInstance: []byte{5, 5, 5, 5},
				Identifier:   []byte{1, 2, 3, 4},
			},
			map[string]*Instance{
				hex.EncodeToString([]byte{5, 5, 5, 5}): {
					State: &proto.State{
						Stage: proto.RoundState_Decided,
					},
				},
			},
			"",
		},
		{
			"valid prev but not decided",
			StartOptions{
				PrevInstance: []byte{5, 5, 5, 5},
				Identifier:   []byte{1, 2, 3, 4},
			},
			map[string]*Instance{
				hex.EncodeToString([]byte{5, 5, 5, 5}): {
					State: &proto.State{
						Stage: proto.RoundState_Prepare,
					},
				},
			},
			"previous instance not decided, can't start new instance",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			i := testIBFTInstance(t)
			i.instances = test.prevInstances
			err := i.canStartNewInstance(test.opts)

			if len(test.expectedError) > 0 {
				require.EqualError(t, err, test.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

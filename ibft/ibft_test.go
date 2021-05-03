package ibft

import (
	"github.com/bloxapp/ssv/ibft/proto"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func testIBFTInstance(t *testing.T) *ibftImpl {
	return &ibftImpl{
		instances:           make([]*Instance, 0),
		currentInstanceLock: &sync.Mutex{},
	}
}

func TestCanStartNewInstance(t *testing.T) {
	tests := []struct {
		name          string
		opts          StartOptions
		seqNumber     uint64
		prevInstances []*Instance
		expectedError string
	}{
		{
			"valid start",
			StartOptions{
				PrevInstance: FirstInstanceIdentifier(),
				Identifier:   []byte{1, 2, 3, 4},
			},
			0,
			make([]*Instance, 0),
			"",
		},
		{
			"unknown prev",
			StartOptions{
				PrevInstance: []byte{5, 5, 5, 5},
				Identifier:   []byte{1, 2, 3, 4},
			},
			1,
			make([]*Instance, 0),
			"instance seq invalid",
		},
		{
			"future instance",
			StartOptions{
				PrevInstance: []byte{5, 5, 5, 5},
				Identifier:   []byte{1, 2, 3, 4},
			},
			10,
			make([]*Instance, 0),
			"instance seq invalid",
		},
		{
			"past instance",
			StartOptions{
				PrevInstance: []byte{5, 5, 5, 5},
				Identifier:   []byte{1, 2, 3, 4},
			},
			1,
			[]*Instance{
				{
					State: &proto.State{
						Stage:     proto.RoundState_Decided,
						SeqNumber: 0,
					},
				},
				{
					State: &proto.State{
						Stage:     proto.RoundState_Decided,
						SeqNumber: 1,
					},
				},
				{
					State: &proto.State{
						Stage:     proto.RoundState_Decided,
						SeqNumber: 2,
					},
				},
			},
			"instance seq invalid",
		},
		{
			"valid prev",
			StartOptions{
				PrevInstance: []byte{5, 5, 5, 5},
				Identifier:   []byte{1, 2, 3, 4},
			},
			1,
			[]*Instance{
				{
					State: &proto.State{
						Stage:     proto.RoundState_Decided,
						SeqNumber: 0,
					},
				},
			},
			"",
		},
		{
			"valid prev",
			StartOptions{
				PrevInstance: []byte{5, 5, 5, 5},
				Identifier:   []byte{1, 2, 3, 4},
			},
			1,
			[]*Instance{
				{
					State: &proto.State{
						Stage:     proto.RoundState_Prepare,
						SeqNumber: 0,
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
			instanceOpts := i.instanceOptionsFromStartOptions(test.opts)
			instanceOpts.SeqNumber = test.seqNumber
			err := i.canStartNewInstance(instanceOpts)

			if len(test.expectedError) > 0 {
				require.EqualError(t, err, test.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

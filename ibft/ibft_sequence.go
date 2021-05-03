package ibft

import (
	"github.com/bloxapp/ssv/ibft/proto"
	"github.com/pkg/errors"
)

/**
IBFT Sequence is the equivalent of block number in a blockchain.
An incremental number for a new iBFT instance.
A fully synced iBFT node must have all sequences to be fully synced, no skips or missing sequences.
*/

func (i *ibftImpl) canStartNewInstance(opts InstanceOptions) error {
	if opts.SeqNumber == 0 {
		return nil
	}
	if opts.SeqNumber != uint64(len(i.instances)) {
		return errors.New("instance seq invalid")
	}
	// If previous instance didn't decide, can't start another instance.
	instance := i.instances[opts.SeqNumber-1]
	if instance.Stage() != proto.RoundState_Decided {
		return errors.New("previous instance not decided, can't start new instance")
	}
	return nil
}

// NextSeqNumber returns the previous decided instance seq number + 1
// In case it's the first instance it returns 0
func (i *ibftImpl) NextSeqNumber() uint64 {
	return uint64(len(i.instances))
}

func (i *ibftImpl) instanceOptionsFromStartOptions(opts StartOptions) InstanceOptions {
	return InstanceOptions{
		Logger:         opts.Logger,
		Me:             i.me,
		Network:        i.network,
		Queue:          i.msgQueue,
		ValueCheck:     opts.ValueCheck,
		LeaderSelector: i.leaderSelector,
		Params:         i.params,
		Lambda:         opts.Identifier,
		SeqNumber:      i.NextSeqNumber(),
		PreviousLambda: opts.PrevInstance,
	}
}

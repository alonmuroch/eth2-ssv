package auth

import (
	"github.com/bloxapp/ssv/ibft/pipeline"
	"github.com/bloxapp/ssv/ibft/proto"
	"github.com/pkg/errors"
)

// ValidateSeqNumber validates seq number
func ValidateSeqNumber(state *proto.State) pipeline.Pipeline {
	return pipeline.WrapFunc(func(signedMessage *proto.SignedMessage) error {
		if state.SeqNumber != signedMessage.Message.SeqNumber {
			return errors.Errorf("message seq number (%d) does not equal State seq number (%d)", signedMessage.Message.SeqNumber, state.SeqNumber)
		}

		return nil
	})
}

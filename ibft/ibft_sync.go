package ibft

import (
	"bytes"
	"github.com/bloxapp/ssv/ibft/proto"
)

// SyncIBFT is responsible for syncing and iBFT instance when needed by
// fetching decided messages from the network
type SyncIBFT struct {
}

// Start the sync
func (sync *SyncIBFT) Start() {
	panic("implement SyncIBFT")
}

// FindHighestInstance returns the highest found instance identifier found by asking the P2P network
func (sync *SyncIBFT) FindHighestInstance() []byte {
	return nil
}

// FetchValidateAndSaveInstances fetches, validates and saves decided messages from the P2P network.
func (sync *SyncIBFT) FetchValidateAndSaveInstances(startID []byte, endID []byte) {

}

// ProcessDecidedMessage is responsible for processing an incoming decided message.
// If the decided message is known or belong to the current executing instance, do nothing.
// Else perform a sync operation
/* From https://arxiv.org/pdf/2002.03613.pdf
We can omit this if we assume some mechanism external to the consensus algorithm that ensures
synchronization of decided values.
upon receiving a valid hROUND-CHANGE, λi, −, −, −i message from pj ∧ pi has decided
by calling Decide(λi,− , Qcommit) do
	send Qcommit to process pj
*/
func (i *ibftImpl) ProcessDecidedMessage(msg *proto.SignedMessage) {
	i.currentInstanceLock.Lock()
	defer i.currentInstanceLock.Unlock()

	// TODO - validate msg

	// if we already have this in storage, pass
	if i.storage.GetDecided(msg.Message.Lambda) != nil {
		return
	}

	// If received decided for current instance, let that instance play out.
	// otherwise sync
	// TODO - should we act upon this decided msg and not let it play out?
	if !bytes.Equal(i.currentInstance.State.Lambda, msg.Message.Lambda) {
		// stop current instance
		i.currentInstance.Stop()

		// sync
		s := &SyncIBFT{}
		go s.Start()
	}
}

package ibft

import (
	"bytes"
	"github.com/bloxapp/ssv/ibft/proto"
)

type SyncIBFT struct {
}

func (sync *SyncIBFT) Start() {
	panic("implement SyncIBFT")
}

// FindHighestInstance returns the highest found instance identifier found by asking the P2P network
func (sync *SyncIBFT) FindHighestInstance() []byte {
	return nil
}

// FetchValidateAndSaveInstances fetches, validates and saves decided messages from the P2P network.
func (sync *SyncIBFT) FetchValidateAndSaveInstances(startIdId []byte, endId []byte) {

}

// ProcessDecidedMessage is responsible for processing an incoming decided message.
// If the decided message is known or belong to the current executing instance, do nothing.
// Else perform a sync operation
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

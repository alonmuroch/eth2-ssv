package sync

import "github.com/bloxapp/ssv/network"

// SyncIBFT is responsible for syncing and iBFT instance when needed by
// fetching decided messages from the network
type SyncIBFT struct {
	network network.Network
}

func New(network network.Network) *SyncIBFT {
	return &SyncIBFT{network: network}
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

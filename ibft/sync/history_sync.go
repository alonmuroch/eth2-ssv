package sync

import "github.com/bloxapp/ssv/network"

// Sync is responsible for syncing and iBFT instance when needed by
// fetching decided messages from the network
type Sync struct {
	network network.Network
}

// New returns a new instance of Sync
func New(network network.Network) *Sync {
	return &Sync{network: network}
}

// Start the sync
func (sync *Sync) Start() {
	panic("implement Sync")
}

// FindHighestInstance returns the highest found instance identifier found by asking the P2P network
func (sync *Sync) FindHighestInstance() []byte {

	return nil
}

// FetchValidateAndSaveInstances fetches, validates and saves decided messages from the P2P network.
func (sync *Sync) FetchValidateAndSaveInstances(startID []byte, endID []byte) {

}

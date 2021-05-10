package ibft

import (
	"github.com/bloxapp/ssv/ibft/sync"
	"github.com/bloxapp/ssv/network"
)

func (i *ibftImpl) ProcessSyncMessage(msg *network.SyncChanObj) {
	s := sync.NewReqHandler(i.logger, i.network, i.storage)
	go s.Process(msg)
}

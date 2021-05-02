package p2p

import (
	"encoding/json"
	"github.com/bloxapp/ssv/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// BroadcastSyncMessage broadcasts a sync message to peers.
// If peer list is nil, broadcasts to all.
func (n *p2pNetwork) BroadcastSyncMessage(peers []peer.ID, msg *network.SyncMessage) error {
	s, err := n.host.NewStream(n.ctx, peers[0], syncStreamProtocol)
	if err != nil {
		return err
	}
	// close function for stream
	defer func() {
		if err := s.Close(); err != nil {
			n.logger.Error("could not close peer stream", zap.Error(err))
		}
	}()

	// message to bytes
	msgBytes, err := json.Marshal(network.Message{
		SyncMessage: msg,
		Type:        network.NetworkMsg_SyncType,
	})
	if err != nil {
		return errors.Wrap(err, "failed to marshal message")
	}

	// send message
	if _, err := s.Write(msgBytes); err != nil {
		return err
	}
	return s.Close() // TODO - should we not close it and use the stream struct on the receiving end for the response?
}

// ReceivedSyncMsgChan returns the channel for sync messages
func (n *p2pNetwork) ReceivedSyncMsgChan() <-chan *network.SyncMessage {
	ls := listener{
		syncCh: make(chan *network.SyncMessage, MsgChanSize),
	}

	n.listenersLock.Lock()
	n.listeners = append(n.listeners, ls)
	n.listenersLock.Unlock()

	return ls.syncCh
}

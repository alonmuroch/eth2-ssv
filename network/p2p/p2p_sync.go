package p2p

import (
	"encoding/json"
	"github.com/bloxapp/ssv/network"
	core "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// BroadcastSyncMessage broadcasts a sync message to peers.
// Peer list must not be nil or empty if stream is nil.
// returns a stream closed for writing
func (n *p2pNetwork) sendSyncMessage(stream core.Stream, peers []peer.ID, msg *network.SyncMessage) (core.Stream, error) {
	if stream == nil {
		if peers == nil || len(peers) == 0 {
			return nil, errors.New("peer list is empty or nil")
		}

		s, err := n.host.NewStream(n.ctx, peers[0], syncStreamProtocol)
		if err != nil {
			return nil, err
		}
		stream = s
	}

	// message to bytes
	msgBytes, err := json.Marshal(network.Message{
		SyncMessage: msg,
		Type:        network.NetworkMsg_SyncType,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal message")
	}

	if _, e := stream.Write(msgBytes); e != nil {
		err = e
	}
	if e := stream.CloseWrite(); err != nil {
		err = e
	}
	return stream, nil
}

// BroadcastSyncMessage broadcasts a sync message to peers.
// If peer list is nil, broadcasts to all.
func (n *p2pNetwork) GetHighestDecidedInstance(peers []peer.ID, msg *network.SyncMessage) (*network.Message, error) {
	stream, err := n.sendSyncMessage(nil, peers, msg)
	if err != nil {
		return nil, errors.Wrap(err, "could not send sync msg")
	}

	// close function for stream
	defer func() {
		if err := stream.Close(); err != nil {
			n.logger.Error("could not close peer stream", zap.Error(err))
		}
	}()

	res, err := readMessageData(stream)
	if err != nil {
		err = errors.Wrap(err, "failed to read response for sync message")
	}

	return res, err
}

// RespondToHighestDecidedInstance responds to a GetHighestDecidedInstance
func (n *p2pNetwork) RespondToHighestDecidedInstance(stream core.Stream, msg *network.SyncMessage) error {
	msg.FromPeerID = n.host.ID().String()
	_, err := n.sendSyncMessage(stream, nil, msg)
	return err
}

// ReceivedSyncMsgChan returns the channel for sync messages
func (n *p2pNetwork) ReceivedSyncMsgChan() <-chan *network.SyncChanObj {
	ls := listener{
		syncCh: make(chan *network.SyncChanObj, MsgChanSize),
	}

	n.listenersLock.Lock()
	n.listeners = append(n.listeners, ls)
	n.listenersLock.Unlock()

	return ls.syncCh
}

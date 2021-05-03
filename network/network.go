package network

import (
	"github.com/bloxapp/ssv/ibft/proto"
	core "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/peer"
)

type SyncChanObj struct {
	Msg    *SyncMessage
	Stream core.Stream
}

// Network represents the behavior of the network
type Network interface {
	// Broadcast propagates a signed message to all peers
	Broadcast(msg *proto.SignedMessage) error

	// ReceivedMsgChan is a channel that forwards new propagated messages to a subscriber
	ReceivedMsgChan() <-chan *proto.SignedMessage

	// BroadcastSignature broadcasts the given signature for the given lambda
	BroadcastSignature(msg *proto.SignedMessage) error

	// ReceivedSignatureChan returns the channel with signatures
	ReceivedSignatureChan() <-chan *proto.SignedMessage

	// BroadcastDecided broadcasts a decided instance with collected signatures
	BroadcastDecided(msg *proto.SignedMessage) error

	// ReceivedDecidedChan returns the channel for decided messages
	ReceivedDecidedChan() <-chan *proto.SignedMessage

	// GetHighestDecidedInstance sends a highest decided request to peers and returns answers.
	// If peer list is nil, broadcasts to all.
	GetHighestDecidedInstance(peers []peer.ID, msg *SyncMessage) (*Message, error)

	// RespondToHighestDecidedInstance responds to a GetHighestDecidedInstance
	RespondToHighestDecidedInstance(stream core.Stream, msg *SyncMessage) error

	// ReceivedSyncMsgChan returns the channel for sync messages
	ReceivedSyncMsgChan() <-chan *SyncChanObj
}

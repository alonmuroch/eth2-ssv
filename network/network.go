package network

import (
	"github.com/bloxapp/ssv/ibft/proto"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

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

	GetTopic() *pubsub.Topic
}

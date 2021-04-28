package local

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"sync"

	"github.com/bloxapp/ssv/ibft/proto"
)

// Local implements network.Local interface
type Local struct {
	msgC               []chan *proto.SignedMessage
	sigC               []chan *proto.SignedMessage
	decidedC           []chan *proto.SignedMessage
	createChannelMutex sync.Mutex
}

// GetTopic TBD
func (n *Local) GetTopic() *pubsub.Topic {
	panic("implement me")
}

// NewLocalNetwork creates a new instance of a local network
func NewLocalNetwork() *Local {
	return &Local{
		msgC:     make([]chan *proto.SignedMessage, 0),
		sigC:     make([]chan *proto.SignedMessage, 0),
		decidedC: make([]chan *proto.SignedMessage, 0),
	}
}

// ReceivedMsgChan implements network.Local interface
func (n *Local) ReceivedMsgChan() <-chan *proto.SignedMessage {
	n.createChannelMutex.Lock()
	defer n.createChannelMutex.Unlock()
	c := make(chan *proto.SignedMessage)
	n.msgC = append(n.msgC, c)
	return c
}

// ReceivedSignatureChan returns the channel with signatures
func (n *Local) ReceivedSignatureChan() <-chan *proto.SignedMessage {
	n.createChannelMutex.Lock()
	defer n.createChannelMutex.Unlock()
	c := make(chan *proto.SignedMessage)
	n.sigC = append(n.sigC, c)
	return c
}

// Broadcast implements network.Local interface
func (n *Local) Broadcast(signed *proto.SignedMessage) error {
	go func() {

		// verify node is not prevented from sending msgs
		//for _, id := range signed.SignerIds {
		//	if !n.replay.CanSend(signed.Message.Type, signed.Message.Lambda, signed.Message.Round, id) {
		//		return
		//	}
		//}

		for _, c := range n.msgC {
			//if !n.replay.CanReceive(signed.Message.Type, signed.Message.Lambda, signed.Message.Round, i) {
			//	fmt.Printf("can't receive, node %d, lambda %s\n", i, hex.EncodeToString(signed.Message.Lambda))
			//	continue
			//}

			c <- signed
		}
	}()

	return nil
}

// BroadcastSignature broadcasts the given signature for the given lambda
func (n *Local) BroadcastSignature(msg *proto.SignedMessage) error {
	n.createChannelMutex.Lock()
	go func() {
		for _, c := range n.sigC {
			c <- msg
		}
		n.createChannelMutex.Unlock()
	}()
	return nil
}

// BroadcastDecided broadcasts a decided instance with collected signatures
func (n *Local) BroadcastDecided(msg *proto.SignedMessage) error {
	n.createChannelMutex.Lock()
	go func() {
		for _, c := range n.decidedC {
			c <- msg
		}
		n.createChannelMutex.Unlock()
	}()
	return nil
}

// ReceivedDecidedChan returns the channel for decided messages
func (n *Local) ReceivedDecidedChan() <-chan *proto.SignedMessage {
	n.createChannelMutex.Lock()
	defer n.createChannelMutex.Unlock()
	c := make(chan *proto.SignedMessage)
	n.decidedC = append(n.msgC, c)
	return c
}

package sync

import (
	"github.com/bloxapp/ssv/ibft/proto"
	"github.com/bloxapp/ssv/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"testing"
)

type testNetwork struct {
	t                      *testing.T
	highestDecidedReceived *proto.SignedMessage
}

func NewTestNetwork(t *testing.T) *testNetwork {
	return &testNetwork{t: t}
}

func (n *testNetwork) Broadcast(msg *proto.SignedMessage) error {
	return nil
}

func (n *testNetwork) ReceivedMsgChan() <-chan *proto.SignedMessage {
	return nil
}

func (n *testNetwork) BroadcastSignature(msg *proto.SignedMessage) error {
	return nil
}

func (n *testNetwork) ReceivedSignatureChan() <-chan *proto.SignedMessage {
	return nil
}

func (n *testNetwork) BroadcastDecided(msg *proto.SignedMessage) error {
	return nil
}

func (n *testNetwork) ReceivedDecidedChan() <-chan *proto.SignedMessage {
	return nil
}

func (n *testNetwork) GetHighestDecidedInstance(peers []peer.ID, msg *network.SyncMessage) (*network.Message, error) {
	return nil, nil
}

func (n *testNetwork) RespondToHighestDecidedInstance(stream network.SyncStream, msg *network.SyncMessage) error {
	n.highestDecidedReceived = msg.SignedMessages[0]
	return nil
}

func (n *testNetwork) ReceivedSyncMsgChan() <-chan *network.SyncChanObj {
	return nil
}

type testStorage struct {
	highestDecided *proto.SignedMessage
}

func NewTestStorage(highestDecided *proto.SignedMessage) *testStorage {
	return &testStorage{highestDecided: highestDecided}
}

func (s *testStorage) SavePrepared(signedMsg *proto.SignedMessage) {

}

func (s *testStorage) SaveDecided(signedMsg *proto.SignedMessage) {

}

func (s *testStorage) GetDecided(identifier []byte) *proto.SignedMessage {
	return nil
}

func (s *testStorage) SaveHighestDecidedInstance(signedMsg *proto.SignedMessage) {

}

func (s *testStorage) GetHighestDecidedInstance() *proto.SignedMessage {
	return s.highestDecided
}

type testSyncStream struct {
}

func NewTestStream() network.SyncStream {
	return &testSyncStream{}
}

func (s *testSyncStream) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (s *testSyncStream) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (s *testSyncStream) Close() error {
	return nil
}

func (s *testSyncStream) CloseWrite() error {
	return nil
}

func (s *testSyncStream) RemotePeer() string {
	return ""
}

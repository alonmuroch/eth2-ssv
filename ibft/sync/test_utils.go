package sync

import (
	"encoding/json"
	"github.com/bloxapp/ssv/ibft/proto"
	"github.com/bloxapp/ssv/network"
	core "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"
	"testing"
)

type testNetwork struct {
	t *testing.T
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

func (n *testNetwork) RespondToHighestDecidedInstance(stream core.Stream, msg *network.SyncMessage) error {
	byts, err := json.Marshal(msg)
	require.NoError(n.t, err)

	_, err = stream.Write(byts)
	require.NoError(n.t, err)
	return nil
}

func (n *testNetwork) ReceivedSyncMsgChan() <-chan *network.SyncChanObj {
	return nil
}

type testStorage struct {
}

func NewTestStorage() *testStorage {
	return &testStorage{}
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
	return nil
}

func NewStream() core.Stream {

}

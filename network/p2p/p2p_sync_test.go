package p2p

import (
	"context"
	"github.com/bloxapp/ssv/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"testing"
	"time"
)

func broadcastSyncMsg(t *testing.T, broadcaster network.Network, receiver network.Network) {
	messageToBroadcast := &network.SyncMessage{
		SignedMessages: nil,
		Type:           network.Sync_GetHighestType,
	}

	time.Sleep(time.Millisecond * 100) // important to let nodes reach each other
	err := broadcaster.BroadcastSyncMessage([]peer.ID{
		receiver.(*p2pNetwork).host.ID(),
	}, messageToBroadcast)
	require.NoError(t, err)
	time.Sleep(time.Millisecond * 100)
}

func TestSyncMessageBroadcasting(t *testing.T) {
	logger := zaptest.NewLogger(t)
	topic1 := "test"

	// create 2 peers
	peer1, err := New(context.Background(), logger, &Config{
		DiscoveryType:     "mdns",
		BootstrapNodeAddr: []string{"enr:-LK4QMIAfHA47rJnVBaGeoHwXOrXcCNvUaxFiDEE2VPCxQ40cu_k2hZsGP6sX9xIQgiVnI72uxBBN7pOQCo5d9izhkcBh2F0dG5ldHOIAAAAAAAAAACEZXRoMpD1pf1CAAAAAP__________gmlkgnY0gmlwhH8AAAGJc2VjcDI1NmsxoQJu41tZ3K8fb60in7AarjEP_i2zv35My_XW_D_t6Y1fJ4N0Y3CCE4iDdWRwgg-g"},
		UDPPort:           12000,
		TCPPort:           13000,
		TopicName:         topic1,
	})
	require.NoError(t, err)

	peer2, err := New(context.Background(), logger, &Config{
		DiscoveryType:     "mdns",
		BootstrapNodeAddr: []string{"enr:-LK4QMIAfHA47rJnVBaGeoHwXOrXcCNvUaxFiDEE2VPCxQ40cu_k2hZsGP6sX9xIQgiVnI72uxBBN7pOQCo5d9izhkcBh2F0dG5ldHOIAAAAAAAAAACEZXRoMpD1pf1CAAAAAP__________gmlkgnY0gmlwhH8AAAGJc2VjcDI1NmsxoQJu41tZ3K8fb60in7AarjEP_i2zv35My_XW_D_t6Y1fJ4N0Y3CCE4iDdWRwgg-g"},
		UDPPort:           12001,
		TCPPort:           13001,
		TopicName:         topic1,
	})
	require.NoError(t, err)

	// set receivers
	peer1Chan := peer1.ReceivedSyncMsgChan()
	peer2Chan := peer2.ReceivedSyncMsgChan()

	peer1Verified := false
	go func() {
		msgFromPeer1 := <-peer1Chan
		require.IsType(t, network.SyncMessage{}, *msgFromPeer1)
		require.EqualValues(t, peer2.(*p2pNetwork).host.ID().String(), msgFromPeer1.FromPeerID)
		require.EqualValues(t, network.Sync_GetHighestType, msgFromPeer1.Type)
		peer1Verified = true
	}()

	peer2Verified := false
	go func() {
		msgFromPeer2 := <-peer2Chan
		require.IsType(t, network.SyncMessage{}, *msgFromPeer2)
		require.EqualValues(t, peer1.(*p2pNetwork).host.ID().String(), msgFromPeer2.FromPeerID)
		require.EqualValues(t, network.Sync_GetHighestType, msgFromPeer2.Type)
		peer2Verified = true

		broadcastSyncMsg(t, peer2, peer1)
	}()

	broadcastSyncMsg(t, peer1, peer2)
	time.Sleep(time.Millisecond * 300) // important to let msgs propagate

	// verify
	require.True(t, peer1Verified, "did not verify peer 1 streamed msg")
	require.True(t, peer2Verified, "did not verify peer 2 streamed msg")
}

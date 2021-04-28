package p2p

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	p2pHost "github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p/p2p/protocol/identify"
	"github.com/prysmaticlabs/prysm/beacon-chain/p2p/peers"
	"github.com/prysmaticlabs/prysm/beacon-chain/p2p/peers/scorers"
	"github.com/prysmaticlabs/prysm/shared/runutil"
	"net"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/bloxapp/ssv/ibft/proto"
	"github.com/bloxapp/ssv/network"
)

const (
	// DiscoveryInterval is how often we re-publish our mDNS records.
	DiscoveryInterval = time.Second

	// DiscoveryServiceTag is used in our mDNS advertisements to discover other chat peers.
	DiscoveryServiceTag = "bloxstaking.ssv"

	// MsgChanSize is the buffer size of the message channel
	MsgChanSize = 128

	topicFmt = "bloxstaking.ssv.%s"
)

type listener struct {
	msgCh     chan *proto.SignedMessage
	sigCh     chan *proto.SignedMessage
	decidedCh chan *proto.SignedMessage
}

// p2pNetwork implements network.Network interface using P2P
type p2pNetwork struct {
	ctx           context.Context
	cfg           *Config
	listenersLock sync.Locker
	dv5Listener   iListener
	listeners     []listener
	logger        *zap.Logger
	privKey       *ecdsa.PrivateKey
	peers         *peers.Status
	host          p2pHost.Host
	pubsub        *pubsub.PubSub
}

// New is the constructor of p2pNetworker
func New(ctx context.Context, logger *zap.Logger, cfg *Config) (network.Network, error) {
	n := &p2pNetwork{
		ctx:           ctx,
		cfg:           cfg,
		listenersLock: &sync.Mutex{},
		logger:        logger,
	}

	var _ipAddr net.IP

	if cfg.DiscoveryType == "mdns" { // use mdns discovery
		// Create a new libp2p Host that listens on a random TCP port
		host, err := libp2p.New(ctx, libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
		if err != nil {
			return nil, errors.Wrap(err, "failed to create a new P2P host")
		}
		n.host = host
		n.cfg.HostID = host.ID()
	} else if cfg.DiscoveryType == "discv5" {
		dv5Nodes := n.parseBootStrapAddrs(n.cfg.BootstrapNodeAddr)
		n.cfg.Discv5BootStrapAddr = dv5Nodes

		_ipAddr = n.ipAddr()
		//_ipAddr = net.ParseIP("127.0.0.1")
		logger.Info("Ip Address", zap.Any("ip", _ipAddr))

		privKey, err := privKey()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to generate p2p private key")
		}
		n.privKey = privKey
		opts := n.buildOptions(_ipAddr, privKey)
		host, err := libp2p.New(ctx, opts...)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create p2p host")
		}
		host.RemoveStreamHandler(identify.IDDelta)
		n.host = host
	} else {
		logger.Error("Unsupported discovery flag")
		return nil, errors.New("Unsupported discovery flag")
	}

	n.logger = logger.With(zap.String("id", n.host.ID().String()), zap.String("Topic", n.cfg.TopicName))
	n.logger.Info("New peer created")

	// Gossipsub registration is done before we add in any new peers
	// due to libp2p's gossipsub implementation not taking into
	// account previously added peers when creating the gossipsub
	// object.
	psOpts := []pubsub.Option{
		//pubsub.WithMessageSignaturePolicy(pubsub.StrictNoSign),
		//pubsub.WithNoAuthor(),
		//pubsub.WithMessageIdFn(msgIDFunction),
		//pubsub.WithSubscriptionFilter(s),
		pubsub.WithPeerOutboundQueueSize(256),
		pubsub.WithValidateQueueSize(256),
	}

	setPubSubParameters()

	// Create a new PubSub service using the GossipSub router
	gs, err := pubsub.NewGossipSub(ctx, n.host, psOpts...)
	if err != nil {
		n.logger.Error("Failed to start pubsub")
		return nil, err
	}
	n.pubsub = gs

	if cfg.DiscoveryType == "mdns" { // use mdns discovery {
		// Setup Local mDNS discovery
		if err := setupDiscovery(ctx, logger, n.host); err != nil {
			return nil, errors.Wrap(err, "failed to setup discovery")
		}
	} else if cfg.DiscoveryType == "discv5" {
		n.peers = peers.NewStatus(ctx, &peers.StatusConfig{
			PeerLimit: 45,
			ScorerParams: &scorers.Config{
				BadResponsesScorerConfig: &scorers.BadResponsesScorerConfig{
					Threshold:     5,
					DecayInterval: time.Hour,
				},
			},
		})

		listener, err := n.startDiscoveryV5(_ipAddr, n.privKey)
		if err != nil {
			n.logger.Error("Failed to start discovery", zap.Error(err))
			return nil, err
		}
		n.dv5Listener = listener

		err = n.connectToBootnodes()
		if err != nil {
			n.logger.Error("Could not add bootnode to the exclusion list", zap.Error(err))
			return nil, err
		}

		go n.listenForNewNodes()

		if n.cfg.HostAddress != "" {
			//logExternalIPAddr(s.host.ID(), p2pHostAddress, p2pTCPPort)
			a := net.JoinHostPort(n.cfg.HostAddress, fmt.Sprintf("%d", n.cfg.TCPPort))
			conn, err := net.DialTimeout("tcp", a, time.Second*10)
			if err != nil {
				n.logger.Error("IP address is not accessible", zap.Error(err))
			}
			if err := conn.Close(); err != nil {
				n.logger.Error("Could not close connection", zap.Error(err))
			}
		}
	}

	// Join the pubsub Topic
	topic, err := n.pubsub.Join(getTopic(n.cfg.TopicName))
	if err != nil {
		return nil, errors.Wrap(err, "failed to join to Topic")
	}

	n.cfg.Topic = topic

	// And subscribe to it
	sub, err := topic.Subscribe()
	if err != nil {
		return nil, errors.Wrap(err, "failed to subscribe on Topic")
	}

	n.cfg.Sub = sub

	go n.listen()

	runutil.RunEvery(n.ctx, 1*time.Minute, func() {
		n.logger.Info("Current peers status", zap.Any("peers", n.GetTopic().ListPeers()))
	})

	return n, nil
}

func (n *p2pNetwork) GetTopic() *pubsub.Topic {
	return n.cfg.Topic
}

// Broadcast propagates a signed message to all peers
func (n *p2pNetwork) Broadcast(msg *proto.SignedMessage) error {
	msgBytes, err := json.Marshal(network.Message{
		Lambda: msg.Message.Lambda,
		Msg:    msg,
		Type:   network.IBFTBroadcastingType,
	})
	if err != nil {
		return errors.Wrap(err, "failed to marshal message")
	}

	n.logger.Debug("Broadcasting to topic", zap.Any("topic", n.cfg.Sub.Topic()), zap.Any("peers", n.cfg.Topic.ListPeers()))
	return n.cfg.Topic.Publish(n.ctx, msgBytes)
}

// ReceivedMsgChan return a channel with messages
func (n *p2pNetwork) ReceivedMsgChan() <-chan *proto.SignedMessage {
	ls := listener{
		msgCh: make(chan *proto.SignedMessage, MsgChanSize),
	}

	n.listenersLock.Lock()
	n.listeners = append(n.listeners, ls)
	n.listenersLock.Unlock()

	return ls.msgCh
}

// BroadcastSignature broadcasts the given signature for the given lambda
func (n *p2pNetwork) BroadcastSignature(msg *proto.SignedMessage) error {
	msgBytes, err := json.Marshal(network.Message{
		Lambda: msg.Message.Lambda,
		Msg:    msg,
		Type:   network.SignatureBroadcastingType,
	})
	if err != nil {
		return errors.Wrap(err, "failed to marshal message")
	}

	return n.cfg.Topic.Publish(n.ctx, msgBytes)
}

// ReceivedSignatureChan returns the channel with signatures
func (n *p2pNetwork) ReceivedSignatureChan() <-chan *proto.SignedMessage {
	ls := listener{
		sigCh: make(chan *proto.SignedMessage, MsgChanSize),
	}

	n.listenersLock.Lock()
	n.listeners = append(n.listeners, ls)
	n.listenersLock.Unlock()

	return ls.sigCh
}

// BroadcastDecided broadcasts a decided instance with collected signatures
func (n *p2pNetwork) BroadcastDecided(msg *proto.SignedMessage) error {
	msgBytes, err := json.Marshal(network.Message{
		Lambda: msg.Message.Lambda,
		Msg:    msg,
		Type:   network.DecidedBroadcastingType,
	})
	if err != nil {
		return errors.Wrap(err, "failed to marshal message")
	}

	return n.cfg.Topic.Publish(n.ctx, msgBytes)
}

// ReceivedDecidedChan returns the channel for decided messages
func (n *p2pNetwork) ReceivedDecidedChan() <-chan *proto.SignedMessage {
	ls := listener{
		decidedCh: make(chan *proto.SignedMessage, MsgChanSize),
	}

	n.listenersLock.Lock()
	n.listeners = append(n.listeners, ls)
	n.listenersLock.Unlock()

	return ls.sigCh
}

// ReceivedMsgChan return a channel with messages
func (n *p2pNetwork) listen() {
	for {
		select {
		case <-n.ctx.Done():
			if err := n.cfg.Topic.Close(); err != nil {
				n.logger.Error("failed to close Topic", zap.Error(err))
			}

			n.cfg.Sub.Cancel()
		default:
			msg, err := n.cfg.Sub.Next(n.ctx)
			if err != nil {
				n.logger.Error("failed to get message from subscription Topic", zap.Error(err))
				return
			}

			var cm network.Message
			if err := json.Unmarshal(msg.Data, &cm); err != nil {
				n.logger.Error("failed to unmarshal message", zap.Error(err))
				continue
			}

			n.logger.Debug("Got message from peer", zap.String("sender peerId", msg.ReceivedFrom.String()), zap.Any("msg", cm))

			for _, ls := range n.listeners {
				go func(ls listener) {

					switch cm.Type {
					case network.IBFTBroadcastingType:
						ls.msgCh <- cm.Msg
					case network.SignatureBroadcastingType:
						ls.sigCh <- cm.Msg
					case network.DecidedBroadcastingType:
						ls.decidedCh <- cm.Msg
					}
				}(ls)
			}
		}
	}
}

func getTopic(topicName string) string {
	return fmt.Sprintf(topicFmt, topicName)
}

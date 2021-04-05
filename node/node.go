package node

import (
	"context"
	"github.com/bloxapp/ssv/network/msgqueue"
	"time"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/herumi/bls-eth-go-binary/bls"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"go.uber.org/zap"

	"github.com/bloxapp/ssv/beacon"
	"github.com/bloxapp/ssv/ibft"
	"github.com/bloxapp/ssv/network"
	"github.com/bloxapp/ssv/slotqueue"
)

// Options contains options to create the node
type Options struct {
	NodeID                     uint64
	ValidatorPubKey            *bls.PublicKey
	PrivateKey                 *bls.SecretKey
	ETHNetwork                 core.Network
	Network                    network.Network
	Queue                      *msgqueue.MessageQueue
	Consensus                  string
	Beacon                     beacon.Beacon
	IBFT                       ibft.IBFT
	Logger                     *zap.Logger
	SignatureCollectionTimeout time.Duration
}

// Node represents the behavior of SSV node
type Node interface {
	// Start starts the SSV node
	Start(ctx context.Context) error
}

// ssvNode implements Node interface
type ssvNode struct {
	nodeID          uint64
	validatorPubKey *bls.PublicKey
	privateKey      *bls.SecretKey
	ethNetwork      core.Network
	network         network.Network
	queue           *msgqueue.MessageQueue
	consensus       string
	slotQueue       slotqueue.Queue
	beacon          beacon.Beacon
	iBFT            ibft.IBFT
	logger          *zap.Logger

	// timeouts
	signatureCollectionTimeout time.Duration
}

// New is the constructor of ssvNode
func New(opts Options) Node {
	return &ssvNode{
		nodeID:                     opts.NodeID,
		validatorPubKey:            opts.ValidatorPubKey,
		privateKey:                 opts.PrivateKey,
		ethNetwork:                 opts.ETHNetwork,
		network:                    opts.Network,
		queue:                      opts.Queue,
		consensus:                  opts.Consensus,
		slotQueue:                  slotqueue.New(opts.ETHNetwork),
		beacon:                     opts.Beacon,
		iBFT:                       opts.IBFT,
		logger:                     opts.Logger,
		signatureCollectionTimeout: opts.SignatureCollectionTimeout,
	}
}

// Start implements Node interface
func (n *ssvNode) Start(ctx context.Context) error {
	go n.startSlotQueueListener(ctx)
	go n.listenToNetworkMessages()

	streamDuties, err := n.beacon.StreamDuties(ctx, n.validatorPubKey.Serialize())
	if err != nil {
		n.logger.Fatal("failed to open duties stream", zap.Error(err))
	}

	n.logger.Info("start streaming duties")
	for duty := range streamDuties {
		go func(duty *ethpb.DutiesResponse_Duty) {
			slots := collectSlots(duty)
			if len(slots) == 0 {
				n.logger.Debug("no slots found for the given duty")
				return
			}

			for _, slot := range slots {
				go func(slot uint64) {
					n.logger.Info("scheduling duty processing start for slot",
						zap.Time("start_time", n.getSlotStartTime(slot)),
						zap.Uint64("committee_index", duty.GetCommitteeIndex()),
						zap.Uint64("slot", slot))

					if err := n.slotQueue.Schedule(n.validatorPubKey.Serialize(), slot, duty); err != nil {
						n.logger.Error("failed to schedule slot")
					}
				}(slot)
			}
		}(duty)
	}

	return nil
}

func (n *ssvNode) listenToNetworkMessages() {
	sigChan := n.network.ReceivedSignatureChan()
	for {
		select {
		case sigMsg := <-sigChan:
			n.queue.AddMessage(&network.Message{
				Lambda: sigMsg.Message.Lambda,
				Msg:    sigMsg,
				Type:   network.SignatureBroadcastingType,
			})
		}
	}
}

// startSlotQueueListener starts slot queue listener
func (n *ssvNode) startSlotQueueListener(ctx context.Context) {
	n.logger.Info("start listening slot queue")

	prevIdentifier := ibft.FirstInstanceIdentifier()
	for {
		slot, duty, ok, err := n.slotQueue.Next(n.validatorPubKey.Serialize())
		if err != nil {
			n.logger.Error("failed to get next slot data", zap.Error(err))
			continue
		}

		if !ok {
			n.logger.Debug("no duties for slot scheduled")
			continue
		}
		go n.executeDuty(ctx, prevIdentifier, slot, duty)
	}
}

// getSlotStartTime returns the start time for the given slot
func (n *ssvNode) getSlotStartTime(slot uint64) time.Time {
	timeSinceGenesisStart := slot * uint64(n.ethNetwork.SlotDurationSec().Seconds())
	start := time.Unix(int64(n.ethNetwork.MinGenesisTime()+timeSinceGenesisStart), 0)
	return start
}

// collectSlots collects slots from the given duty
func collectSlots(duty *ethpb.DutiesResponse_Duty) []uint64 {
	var slots []uint64
	slots = append(slots, duty.GetAttesterSlot())
	slots = append(slots, duty.GetProposerSlots()...)
	return slots
}

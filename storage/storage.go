package storage

import "github.com/bloxapp/ssv/ibft/proto"

// Storage is an interface for persisting chain data
type Storage interface {
	// SavePrepared saves a signed message for an ibft instance with prepared justification
	SavePrepared(signedMsg *proto.SignedMessage)
	// SaveDecided saves a signed message for an ibft instance with decided justification
	SaveDecided(signedMsg *proto.SignedMessage)
	// GetDecided returns a signed message for an ibft instance which decided by identifier
	GetDecided(identifier []byte) *proto.SignedMessage
	// SaveHighestDecidedInstance saves a signed message for an ibft instance which is currently highest
	SaveHighestDecidedInstance(signedMsg *proto.SignedMessage)
	// GetHighestDecidedInstance gets a signed message for an ibft instance which is the highest
	GetHighestDecidedInstance() *proto.SignedMessage
}

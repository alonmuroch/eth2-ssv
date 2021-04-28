package inmem

import (
	"github.com/bloxapp/ssv/ibft/proto"
	"github.com/bloxapp/ssv/storage"
)

// inMemStorage implements storage.Storage interface
type inMemStorage struct {
}

// New is the constructor of inMemStorage
func New() storage.Storage {
	return &inMemStorage{}
}

// SavePrepared saves a signed message for an ibft instance with prepared justification
func (s *inMemStorage) SavePrepared(signedMsg *proto.SignedMessage) {
	// TODO: Implement
}

// SaveDecided saves a signed message for an ibft instance with decided justification
func (s *inMemStorage) SaveDecided(signedMsg *proto.SignedMessage) {
	// TODO: Implement
}

// GetDecided returns a signed message for an ibft instance which decided by identifier
func (s *inMemStorage) GetDecided(identifier []byte) *proto.SignedMessage {
	return nil
}

// SaveHighestDecidedInstance saves a signed message for an ibft instance which is currently highest
func (s *inMemStorage) SaveHighestDecidedInstance(signedMsg *proto.SignedMessage) {
	// TODO: Implement
}

// GetHighestDecidedInstance gets a signed message for an ibft instance which is the highest
func (s *inMemStorage) GetHighestDecidedInstance() *proto.SignedMessage {
	// TODO: Implement
	return nil
}

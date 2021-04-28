package proto

// PreviouslyPrepared checks if state prepare round and value are set
// TODO - remove, not used
func (s *State) PreviouslyPrepared() bool {
	return s.PreparedRound != 0 && s.PreparedValue != nil
}

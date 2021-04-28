package ibft

func (i *ibftImpl) SyncIBFT() {
	i.syncIbftMutex.Lock()
	defer i.syncIbftMutex.Unlock()

	panic("implement SyncIBFT")
}

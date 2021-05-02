package p2p

import (
	"encoding/json"
	"github.com/bloxapp/ssv/network"
	core "github.com/libp2p/go-libp2p-core"
	"go.uber.org/zap"
	"io/ioutil"
)

func readMessageData(stream core.Stream) (*network.Message, error) {
	data := &network.Message{}
	buf, err := ioutil.ReadAll(stream)
	if err != nil {
		return nil, err
	}

	// unmarshal
	if err := json.Unmarshal(buf, data); err != nil {
		return nil, err
	}
	return data, nil
}

// handleStream sets a stream handler for the host to process streamed messages
func (n *p2pNetwork) handleStream() {
	n.host.SetStreamHandler(syncStreamProtocol, func(stream core.Stream) {
		cm, err := readMessageData(stream)
		if err != nil {
			n.logger.Error("could not read and parse stream", zap.Error(err))
			return
		}

		// send to listeners
		for _, ls := range n.listeners {
			go func(ls listener) {
				switch cm.Type {
				case network.NetworkMsg_SyncType:
					cm.SyncMessage.FromPeerID = stream.Conn().RemotePeer().String()
					ls.syncCh <- cm.SyncMessage
				}
			}(ls)
		}
	})
}

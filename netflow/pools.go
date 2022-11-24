package netflow

import "sync"

var packetPool = sync.Pool{
	New: func() interface{} {
		return new(Packet)
	},
}

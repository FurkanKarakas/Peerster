package gossiper

import "net"

//Global variables
var (
	RumorID      map[string]uint32
	PeerNames    map[string]string
	GossiperAddr *net.UDPAddr
	C            chan GossipPacket
)

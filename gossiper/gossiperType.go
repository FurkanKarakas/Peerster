package gossiper

import "net"

//Gossiper contains all relevant inputs to the gossiper program.
type Gossiper struct {
	UIPort          string
	GossipAddr      string
	Name            string
	Peers           string
	Simple          bool
	KnownAddresses  []string
	GossiperUDPConn *net.UDPConn
}

package gossiper

import "net"

//NewGossiper creates a new gossiper with the given parameters.
func NewGossiper(UIPort, gossipAddr, name, peers string,
	simple bool, knownAddresses []string, gossiperUDPConn *net.UDPConn) *Gossiper {
	return &Gossiper{
		UIPort:          UIPort,
		GossipAddr:      gossipAddr,
		Name:            name,
		Peers:           peers,
		Simple:          simple,
		KnownAddresses:  knownAddresses,
		GossiperUDPConn: gossiperUDPConn,
	}
}

package gossiper

import "fmt"

//addNewAddrAndPrintPeers adds a new address to the known addresses list if it was not previously known.
func (g *Gossiper) addNewAddrAndPrintPeers(clientAddress string) {
	newelement := true

	if g.KnownAddresses != nil {
		for _, element := range g.KnownAddresses {
			if element == clientAddress {
				newelement = false
				break
			}
		}
	}

	if newelement {
		if len(g.Peers) == 0 {
			g.Peers = clientAddress
		} else {
			g.Peers = g.Peers + "," + clientAddress
		}
		g.KnownAddresses = append(g.KnownAddresses, clientAddress)
		fmt.Println("PEERS", g.Peers)
	} else {
		fmt.Println("PEERS", g.Peers)
	}
}

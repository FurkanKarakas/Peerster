package gossiper

import (
	"fmt"

	"github.com/dedis/protobuf"
)

//ListenPeers listens to all of the known peers.
func (g *Gossiper) ListenPeers() {
	for {
		buffer := make([]byte, 4096)
		n, temp, err := g.GossiperUDPConn.ReadFromUDP(buffer)
		if n == 0 {
			continue
		}
		clientAddress := temp.String()
		//println(clientAddress)
		if err != nil {
			fmt.Println("ERROR:", err)
			fmt.Println("Bytes read:", n)
			panic(err)
		}
		//var gossippacket GossipPacket
		gossippacket := GossipPacket{}
		err = protobuf.Decode(buffer[0:n], &gossippacket)
		if err != nil {
			fmt.Println("ERROR11:", err)
			panic(err)
		}
		if gossippacket.Simple != nil {
			fmt.Println("SIMPLE MESSAGE origin", gossippacket.Simple.OriginalName,
				"from", gossippacket.Simple.RelayPeerAddr, "contents", gossippacket.Simple.Contents)
			g.addNewAddrAndPrintPeers(clientAddress)
			//Store the original name in a string map
			PeerNames[clientAddress] = gossippacket.Simple.OriginalName
			if g.Simple {
				gossippacket.Simple.RelayPeerAddr = g.GossipAddr
				g.sendGossip(gossippacket, clientAddress)
			}

		} else if gossippacket.Rumor != nil {
			//Store the original name in a string map
			PeerNames[clientAddress] = gossippacket.Rumor.Origin
			//Save the rumor ID
			//rumorID[gossippacket.Rumor.Origin]
			fmt.Println("RUMOR origin", gossippacket.Rumor.Origin,
				"from", clientAddress, "ID", gossippacket.Rumor.ID,
				"contents", gossippacket.Rumor.Text)

			g.addNewAddrAndPrintPeers(clientAddress)

			if !g.Simple {
				var psSlice []PeerStatus
				if g.KnownAddresses != nil {
					for _, address := range g.KnownAddresses {
						if PeerNames[address] != "" {
							psSlice = append(psSlice, PeerStatus{
								Identifier: PeerNames[address],
								NextID:     RumorID[PeerNames[address]] + 1,
							})
						}
					}
				}
				sp := StatusPacket{Want: psSlice}
				gp := GossipPacket{Status: &sp}
				g.sendGossip(gp, clientAddress)
			}

		} else if gossippacket.Status != nil { //Need to implement status packet logic.
			fmt.Print("STATUS from ", clientAddress)
			for _, element := range gossippacket.Status.Want {
				fmt.Print(" peer ", element.Identifier, " nextID ", element.NextID)
			}
			fmt.Println()
			fmt.Println("PEERS", g.Peers)
		} else {
			fmt.Println("you must be trolling dude")
		}

	}
}

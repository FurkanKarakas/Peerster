package gossiper

import (
	"fmt"
	"net"
	"strconv"
)

//ListenClient listens to the client in an infinite loop.
func (g *Gossiper) ListenClient() {
	portNr, err := strconv.Atoi(g.UIPort)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		panic(err)
	}
	addr := net.UDPAddr{
		Port: portNr,
		IP:   net.ParseIP("127.0.0.1"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		panic(err)
	}
	defer conn.Close()

	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("ERROR WHILE READING BUFFER:", err)
			panic(err)
		}
		buffer = buffer[0:n]
		fmt.Println("CLIENT MESSAGE", string(buffer))
		switch g.Simple {
		case true:
			var simplemessage SimpleMessage = SimpleMessage{
				OriginalName:  g.Name,
				RelayPeerAddr: g.GossipAddr,
				Contents:      string(buffer),
			}
			var gp = GossipPacket{Simple: &simplemessage}
			g.sendGossip(gp, g.GossipAddr)

		case false:
			RumorID[g.Name]++
			var rumormessage RumorMessage = RumorMessage{
				Origin: g.Name,
				ID:     RumorID[g.Name],
				Text:   string(buffer),
			}
			var gp = GossipPacket{Rumor: &rumormessage}
			g.sendGossip(gp, g.GossipAddr)

		default: //Some extra logic here if necessary
			println("Sorry, no known meaning of that flag!")
			return
		}

	}
}

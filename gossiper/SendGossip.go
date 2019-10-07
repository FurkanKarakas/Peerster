package gossiper

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"

	"github.com/dedis/protobuf"
)

//sendGossip spreads the message to the neighbor peers.
func (g *Gossiper) sendGossip(gp GossipPacket, spreader string) {
	//fmt.Println(g.knownAddresses, spreader)
	if gp.Simple != nil {
		if len(g.KnownAddresses) == 0 {
			return
		}
		for _, addr := range g.KnownAddresses {
			if addr == spreader {
				continue
			}
			//Prepare the destination address
			tempDestAddr := strings.Split(addr, ":")
			i, err := strconv.Atoi(tempDestAddr[1])
			if err != nil {
				fmt.Println("ERROR:", err)
				panic(err)
			}
			tempDestIPAddr := net.UDPAddr{
				IP:   net.ParseIP(tempDestAddr[0]),
				Port: i,
			}
			packetBytes, err := protobuf.Encode(&gp)
			if err != nil {
				fmt.Println("ERROR: ", err)
				panic(err)
			}
			n, err := g.GossiperUDPConn.WriteToUDP(packetBytes, &tempDestIPAddr)
			if err != nil {
				fmt.Println("ERROR: ", err)
				fmt.Println("Number of bytes written:", n)
				panic(err)
			}
		}
	} else if gp.Rumor != nil {
		if len(g.KnownAddresses) == 0 {
			return
		}
		randomaddress := g.KnownAddresses[rand.Intn(len(g.KnownAddresses))]
		//Prepare the destination address
		tempDestAddr := strings.Split(randomaddress, ":")
		i, err := strconv.Atoi(tempDestAddr[1])
		if err != nil {
			fmt.Println("ERROR:", err)
			panic(err)
		}
		tempDestIPAddr := net.UDPAddr{
			IP:   net.ParseIP(tempDestAddr[0]),
			Port: i,
		}
		fmt.Println("MONGERING with", randomaddress)

		packetBytes, err := protobuf.Encode(&gp)
		if err != nil {
			fmt.Println("ERROR: ", err)
			panic(err)
		}
		n, err := g.GossiperUDPConn.WriteToUDP(packetBytes, &tempDestIPAddr)
		if err != nil {
			fmt.Println("ERROR: ", err)
			fmt.Println("Number of bytes written:", n)
			panic(err)
		}
		/*var spreaderIndex int = -1
		for i, addr := range g.knownAddresses {
			if addr == spreader {
				spreaderIndex = i
				break
			}
		}
		if spreaderIndex == -1 {
			randomaddress := g.knownAddresses[rand.Intn(len(g.knownAddresses))]
			//Prepare the destination address
			tempDestAddr := strings.Split(randomaddress, ":")
			i, err := strconv.Atoi(tempDestAddr[1])
			if err != nil {
				fmt.Println("ERROR:", err)
				panic(err)
			}
			tempDestIPAddr := net.UDPAddr{
				IP:   net.ParseIP(tempDestAddr[0]),
				Port: i,
			}
			fmt.Println("MONGERING with", randomaddress)

			packetBytes, err := protobuf.Encode(&gp)
			if err != nil {
				fmt.Println("ERROR: ", err)
				panic(err)
			}
			n, err := g.gossiperUDPConn.WriteToUDP(packetBytes, &tempDestIPAddr)
			if err != nil {
				fmt.Println("ERROR: ", err)
				fmt.Println("Number of bytes written:", n)
				panic(err)
			}

		} else if len(g.knownAddresses) == 1 {
			return
		} else {
			var randomaddress string
			randomint := rand.Intn(len(g.knownAddresses) - 1)
			if randomint < spreaderIndex {
				randomaddress = g.knownAddresses[randomint]
			} else {
				randomaddress = g.knownAddresses[randomint+1]
			}
			//Prepare the destination address
			tempDestAddr := strings.Split(randomaddress, ":")
			i, err := strconv.Atoi(tempDestAddr[1])
			if err != nil {
				fmt.Println("ERROR:", err)
				panic(err)
			}
			tempDestIPAddr := net.UDPAddr{
				IP:   net.ParseIP(tempDestAddr[0]),
				Port: i,
			}
			fmt.Println("MONGERING with", randomaddress)

			packetBytes, err := protobuf.Encode(&gp)
			if err != nil {
				fmt.Println("ERROR: ", err)
				panic(err)
			}
			n, err := g.gossiperUDPConn.WriteToUDP(packetBytes, &tempDestIPAddr)
			if err != nil {
				fmt.Println("ERROR: ", err)
				fmt.Println("Number of bytes written:", n)
				panic(err)
			}

		}*/

	} else if gp.Status != nil {
		//Prepare the destination address
		tempDestAddr := strings.Split(spreader, ":")
		i, err := strconv.Atoi(tempDestAddr[1])
		if err != nil {
			fmt.Println("ERROR:", err)
			panic(err)
		}
		tempDestIPAddr := net.UDPAddr{
			IP:   net.ParseIP(tempDestAddr[0]),
			Port: i,
		}

		packetBytes, err := protobuf.Encode(&gp)
		if err != nil {
			fmt.Println("ERROR: ", err)
			panic(err)
		}
		n, err := g.GossiperUDPConn.WriteToUDP(packetBytes, &tempDestIPAddr)
		if err != nil {
			fmt.Println("ERROR: ", err)
			fmt.Println("Number of bytes written:", n)
			panic(err)
		}

	} else {
		fmt.Println("What kind of message is this, dude? U wot m8")
	}

}

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"

	"github.com/dedis/protobuf"
)

//Global variables
var (
	rumorID      map[string]uint32
	peerNames    map[string]string
	GossiperAddr *net.UDPAddr
	c            chan GossipPacket
)

//Gossiper contains all relevant inputs to the gossiper program.
type Gossiper struct {
	UIPort          string
	gossipAddr      string
	name            string
	peers           string
	simple          bool
	knownAddresses  []string
	gossiperUDPConn *net.UDPConn
}

//SimpleMessage contains all relevant information about the message.
type SimpleMessage struct {
	OriginalName  string
	RelayPeerAddr string
	Contents      string
}

//RumorMessage contains all relevant information about the rumor message.
type RumorMessage struct {
	Origin string
	ID     uint32
	Text   string
}

//PeerStatus contains information about the status of the peer.
type PeerStatus struct {
	Identifier string
	NextID     uint32
}

//StatusPacket is the status.
type StatusPacket struct {
	Want []PeerStatus
}

//GossipPacket To provide compatibility with future versions
type GossipPacket struct {
	Simple *SimpleMessage
	Rumor  *RumorMessage
	Status *StatusPacket
}

//addNewAddrAndPrintPeers adds a new address to the known addresses list if it was not previously known.
func (g *Gossiper) addNewAddrAndPrintPeers(clientAddress string) {
	newelement := true

	if g.knownAddresses != nil {
		for _, element := range g.knownAddresses {
			if element == clientAddress {
				newelement = false
				break
			}
		}
	}

	if newelement {
		if len(g.peers) == 0 {
			g.peers = clientAddress
		} else {
			g.peers = g.peers + "," + clientAddress
		}
		g.knownAddresses = append(g.knownAddresses, clientAddress)
		fmt.Println("PEERS", g.peers)
	} else {
		fmt.Println("PEERS", g.peers)
	}
}

//listenClient listens to the client in an infinite loop.
func (g *Gossiper) listenClient() {
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
		switch g.simple {
		case true:
			var simplemessage SimpleMessage = SimpleMessage{
				OriginalName:  g.name,
				RelayPeerAddr: g.gossipAddr,
				Contents:      string(buffer),
			}
			var gp = GossipPacket{Simple: &simplemessage}
			g.sendGossip(gp, g.gossipAddr)

		case false:
			rumorID[g.name]++
			var rumormessage RumorMessage = RumorMessage{
				Origin: g.name,
				ID:     rumorID[g.name],
				Text:   string(buffer),
			}
			var gp = GossipPacket{Rumor: &rumormessage}
			g.sendGossip(gp, g.gossipAddr)

		default: //Some extra logic here if necessary
			println("Sorry, no known meaning of that flag!")
			return
		}

	}
}

//listenPeers listens to all of the known peers.
func (g *Gossiper) listenPeers() {
	for {
		buffer := make([]byte, 4096)
		n, temp, err := g.gossiperUDPConn.ReadFromUDP(buffer)
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
			peerNames[clientAddress] = gossippacket.Simple.OriginalName
			if g.simple {
				gossippacket.Simple.RelayPeerAddr = g.gossipAddr
				g.sendGossip(gossippacket, clientAddress)
			}

		} else if gossippacket.Rumor != nil {
			//Store the original name in a string map
			peerNames[clientAddress] = gossippacket.Rumor.Origin
			//Save the rumor ID
			//rumorID[gossippacket.Rumor.Origin]
			fmt.Println("RUMOR origin", gossippacket.Rumor.Origin,
				"from", clientAddress, "ID", gossippacket.Rumor.ID,
				"contents", gossippacket.Rumor.Text)

			g.addNewAddrAndPrintPeers(clientAddress)

			if !g.simple {
				var psSlice []PeerStatus
				if g.knownAddresses != nil {
					for _, address := range g.knownAddresses {
						if peerNames[address] != "" {
							psSlice = append(psSlice, PeerStatus{
								Identifier: peerNames[address],
								NextID:     rumorID[peerNames[address]] + 1,
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
			fmt.Println("PEERS", g.peers)
		} else {
			fmt.Println("you must be trolling dude")
		}

	}
}

//sendGossip spreads the message to the neighbor peers.
func (g *Gossiper) sendGossip(gp GossipPacket, spreader string) {
	//fmt.Println(g.knownAddresses, spreader)
	if gp.Simple != nil {
		if len(g.knownAddresses) == 0 {
			return
		}
		for _, addr := range g.knownAddresses {
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
			n, err := g.gossiperUDPConn.WriteToUDP(packetBytes, &tempDestIPAddr)
			if err != nil {
				fmt.Println("ERROR: ", err)
				fmt.Println("Number of bytes written:", n)
				panic(err)
			}
		}
	} else if gp.Rumor != nil {
		if len(g.knownAddresses) == 0 {
			return
		}
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
		n, err := g.gossiperUDPConn.WriteToUDP(packetBytes, &tempDestIPAddr)
		if err != nil {
			fmt.Println("ERROR: ", err)
			fmt.Println("Number of bytes written:", n)
			panic(err)
		}

	} else {
		fmt.Println("What kind of message is this, dude? U wot m8")
	}

}

func main() {
	var g Gossiper
	flag.StringVar(&g.UIPort, "UIPort", "8080", "port for the UI client (default 8080)")
	flag.StringVar(&g.gossipAddr, "gossipAddr", "127.0.0.1:5000",
		"port for the gossiper (default \"127.0.0.1:5000\"")
	flag.StringVar(&g.name, "name", "", "name of the gossiper")
	flag.StringVar(&g.peers, "peers", "", "comma separated list of peers of the form ip:port")
	flag.BoolVar(&g.simple, "simple", false, "run gossiper in simple broadcast mode")
	flag.Parse()

	// A bunch of useful variables
	tempGossiperAddr := strings.Split(g.gossipAddr, ":")
	portNr, err := strconv.Atoi(tempGossiperAddr[1])
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		panic(err)
	}
	GossiperAddr = &net.UDPAddr{
		Port: portNr,
		IP:   net.ParseIP(tempGossiperAddr[0]),
	}
	rumorID = make(map[string]uint32)
	peerNames = make(map[string]string)

	if g.peers != "" {
		g.knownAddresses = strings.Split(g.peers, ",")
	} else {
		g.knownAddresses = nil
	}
	//Open the UDP socket
	g.gossiperUDPConn, err = net.ListenUDP("udp4", GossiperAddr)
	if err != nil {
		fmt.Println("ERROR:", err)
		panic(err)
	}
	defer g.gossiperUDPConn.Close()

	go g.listenClient()
	go g.listenPeers()
	for {
	}
}

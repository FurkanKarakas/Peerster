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
var rumorID uint32 = 1

//Gossiper contains all relevant inputs to the gossiper program.
type Gossiper struct {
	UIPort         string
	gossipAddr     string
	name           string
	peers          string
	simple         bool
	knownAddresses []string
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

//listenClient listens to the client in an infinite loop.
func (g *Gossiper) listenClient() {
	portNr, err := strconv.Atoi(g.UIPort)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return
	}
	addr := net.UDPAddr{
		Port: portNr,
		IP:   net.ParseIP("127.0.0.1"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return
	}
	defer conn.Close()

	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		fmt.Println(conn.LocalAddr())
		if err != nil {
			fmt.Println("ERROR WHILE READING BUFFER:", err)
			return
		}
		buffer = buffer[0:n]
		var messageType byte = buffer[len(buffer)-1]
		switch messageType {
		case 48: //48 corresponds to 0
			fmt.Println("CLIENT MESSAGE " + string(buffer[0:len(buffer)-1]))
			fmt.Println("PEERS " + g.peers)
			var simplemessage SimpleMessage = SimpleMessage{
				OriginalName:  g.name,
				RelayPeerAddr: g.gossipAddr,
				Contents:      string(buffer[0 : len(buffer)-1]),
			}
			var gp = GossipPacket{Simple: &simplemessage}
			g.sendGossip(gp, g.gossipAddr)

		case 49: //49 corresponds to 1
			fmt.Println("CLIENT MESSAGE " + string(buffer[0:len(buffer)-1]))
			fmt.Println("PEERS " + g.peers)
			var rumormessage RumorMessage = RumorMessage{
				Origin: g.name,
				ID:     rumorID,
				Text:   string(buffer[0 : len(buffer)-1]),
			}
			rumorID++
			var gp = GossipPacket{Rumor: &rumormessage}
			g.sendGossip(gp, g.gossipAddr)

		default: //Some extra logic here if necessary
			println("Sorry, no known meaning of that last character!")
			return
		}

	}
}

//listenPeers listens to all of the known peers.
func (g *Gossiper) listenPeers() {
	gossipaddress := strings.Split(g.gossipAddr, ":")
	i2, err := strconv.Atoi(gossipaddress[1])
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	addr := net.UDPAddr{
		Port: i2,
		IP:   net.ParseIP(gossipaddress[0]),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	defer conn.Close()

	for { //TODO: ListenPeers modification for Rumor mongering
		buffer := make([]byte, 4096)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("ERROR:", err)
			fmt.Println("Bytes read:", n)
			return
		}
		var gossippacket GossipPacket
		err = protobuf.Decode(buffer[0:n], &gossippacket)
		if err != nil {
			fmt.Println("ERROR11:", err)
			return
		}
		fmt.Println("SIMPLE MESSAGE origin", gossippacket.Simple.OriginalName,
			"from", gossippacket.Simple.RelayPeerAddr, "contents", gossippacket.Simple.Contents)
		newelement := true
		for _, element := range g.knownAddresses {
			if element == gossippacket.Simple.RelayPeerAddr {
				newelement = false
				break
			}
		}
		temp := gossippacket.Simple.RelayPeerAddr
		gossippacket.Simple.RelayPeerAddr = g.gossipAddr
		if newelement {
			if len(g.peers) == 0 {
				g.peers = temp
			} else {
				g.peers = g.peers + "," + temp
			}
			g.knownAddresses = append(g.knownAddresses, temp)
			fmt.Println("PEERS", g.peers)
		} else {
			fmt.Println("PEERS", g.peers)
		}
		g.sendGossip(gossippacket, temp)

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
			conn, err := net.Dial("udp", addr)
			if err != nil {
				fmt.Println("ERROR:", err)
				return
			}
			defer conn.Close()

			packetBytes, err := protobuf.Encode(&gp)
			if err != nil {
				fmt.Println("ERROR: ", err)
				return
			}
			n, err := conn.Write(packetBytes)
			if err != nil {
				fmt.Println("ERROR: ", err)
				fmt.Println("Number of bytes written:", n)
				return
			}
		}
	} else if gp.Rumor != nil {
		if len(g.knownAddresses) == 0 {
			return
		}
		var spreaderIndex int = -1
		for i, addr := range g.knownAddresses {
			if addr == spreader {
				spreaderIndex = i
				break
			}
		}
		if spreaderIndex == -1 {
			randomaddress := g.knownAddresses[rand.Intn(len(g.knownAddresses))]
			conn, err := net.Dial("udp", randomaddress)
			if err != nil {
				fmt.Println("ERROR:", err)
				return
			}
			defer conn.Close()

			packetBytes, err := protobuf.Encode(&gp)
			if err != nil {
				fmt.Println("ERROR: ", err)
				return
			}
			n, err := conn.Write(packetBytes)
			if err != nil {
				fmt.Println("ERROR: ", err)
				fmt.Println("Number of bytes written:", n)
				return
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
			conn, err := net.Dial("udp", randomaddress)
			if err != nil {
				fmt.Println("ERROR:", err)
				return
			}
			defer conn.Close()

			packetBytes, err := protobuf.Encode(&gp)
			if err != nil {
				fmt.Println("ERROR: ", err)
				return
			}
			n, err := conn.Write(packetBytes)
			if err != nil {
				fmt.Println("ERROR: ", err)
				fmt.Println("Number of bytes written:", n)
				return
			}
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
	if g.peers != "" {
		g.knownAddresses = strings.Split(g.peers, ",")
	} else {
		g.knownAddresses = nil
	}

	go g.listenClient()
	go g.listenPeers()
	for {
	}
}

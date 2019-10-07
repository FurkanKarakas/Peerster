package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/FurkanKarakas/Peerster/gossiper"
)

func main() {
	var (
		UIPort, gossipAddr, name, peers string
		simple                          bool
		knownAddresses                  []string
		gossiperUDPConn                 *net.UDPConn
	)
	flag.StringVar(&UIPort, "UIPort", "8080", "port for the UI client (default 8080)")
	flag.StringVar(&gossipAddr, "gossipAddr", "127.0.0.1:5000",
		"port for the gossiper (default \"127.0.0.1:5000\"")
	flag.StringVar(&name, "name", "", "name of the gossiper")
	flag.StringVar(&peers, "peers", "", "comma separated list of peers of the form ip:port")
	flag.BoolVar(&simple, "simple", false, "run gossiper in simple broadcast mode")
	flag.Parse()

	// A bunch of useful variables
	tempGossiperAddr := strings.Split(gossipAddr, ":")
	portNr, err := strconv.Atoi(tempGossiperAddr[1])
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		panic(err)
	}
	GossiperAddr := &net.UDPAddr{
		Port: portNr,
		IP:   net.ParseIP(tempGossiperAddr[0]),
	}
	gossiper.RumorID = make(map[string]uint32)
	gossiper.PeerNames = make(map[string]string)

	if peers != "" {
		knownAddresses = strings.Split(peers, ",")
	} else {
		knownAddresses = nil
	}
	//Open the UDP socket
	gossiperUDPConn, err = net.ListenUDP("udp4", GossiperAddr)
	if err != nil {
		fmt.Println("ERROR:", err)
		panic(err)
	}
	var g *gossiper.Gossiper = gossiper.NewGossiper(UIPort, gossipAddr,
		name, peers, simple, knownAddresses, gossiperUDPConn)
	defer g.GossiperUDPConn.Close()

	go g.ListenClient()
	go g.ListenPeers()
	//Be stuck forever
	for {
	}
}

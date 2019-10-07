package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
)

// ClientArgs contains the relevant information, which the client needs.
type ClientArgs struct {
	UIPort  string
	msg     string
	msgType string
}

func main() {
	var c ClientArgs
	flag.StringVar(&c.UIPort, "UIPort", "8080", "port for the UI client (default 8080)")
	flag.StringVar(&c.msg, "msg", "", "message to be sent")
	flag.Parse()

	i, err := strconv.Atoi(c.UIPort)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	udpaddr := net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: i,
	}
	udpconn, err := net.DialUDP("udp4", nil, &udpaddr)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return
	}
	defer udpconn.Close()

	n, err := fmt.Fprintf(udpconn, c.msg)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		fmt.Println("Number of bytes written:", n)
		return
	}

}

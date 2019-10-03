package main

import (
	"flag"
	"fmt"
	"net"
)

// ClientArgs contains the relevant information, which the client needs.
type ClientArgs struct {
	UIPort string
	msg    string
}

func main() {
	var c ClientArgs
	flag.StringVar(&c.UIPort, "UIPort", "8080", "port for the UI client (default 8080)")
	flag.StringVar(&c.msg, "msg", "", "message to be sent")
	flag.Parse()

	conn, err := net.Dial("udp", "127.0.0.1:"+c.UIPort)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return
	}
	defer conn.Close()

	n, err := fmt.Fprintf(conn, c.msg)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		fmt.Println("Number of bytes written: %i", n)
		return
	}

}

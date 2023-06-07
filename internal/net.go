package internal

import (
	"log"
	"net"
)

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	// this does *not* establish a connection.
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// watch "ps -eo pid,comm | grep exe/udp_read | awk '{ print \$1 }' | xargs ps M | wc -l"
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	maxport  = 65535
	numConns = 100
)

func main() {
	for conns, port := 0, 1024; conns < numConns && port < maxport; port++ {
		laddr := &net.UDPAddr{
			Port: port,
			IP:   net.ParseIP("127.0.0.1"),
		}
		conn, err := net.ListenUDP("udp", laddr)
		if err != nil {
			fmt.Printf("Error listening on port %d: %s\n", port, err)
			continue
		}
		defer conn.Close()

		go readConnection(conn)
		conns++

		if conns%100 == 0 {
			fmt.Printf("Opened %d connections...\n", conns)
		}
	}

	fmt.Println("Press enter to continue...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func readConnection(conn *net.UDPConn) {
	bytes := make([]byte, 32)
	for {
		conn.ReadFromUDP(bytes)
	}
}

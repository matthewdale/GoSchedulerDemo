// watch "ps -eo pid,comm | grep exe/udp_read | awk '{ print \$1 }' | xargs ps M | wc -l"
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
)

const maxport = 65535

func main() {
	numConns, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	for conns, port := 0, 1024; conns < numConns && port < maxport; port++ {
		laddr := &net.UDPAddr{
			Port: port,
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
	for {
		bytes := make([]byte, 32)
		conn.ReadFromUDP(bytes)
	}
}

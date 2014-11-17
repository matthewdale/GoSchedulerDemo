// watch "ps -eo pid,comm | grep exe/net_read | awk '{ print \$1 }' | xargs ps M | wc -l"
package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	laddr := &net.TCPAddr{
		Port: port,
	}
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		panic(fmt.Sprintf("Error listening on port %d: %s\n", port, err))
	}
	defer listener.Close()

	fmt.Printf("Listening for connections on port %d...\n", port)
	conns := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting new connection on %s: %s", listener.Addr().String(), err)
			continue
		}
		go readConnection(conn)

		conns++

		if (conns+1)%100 == 0 {
			fmt.Printf("Accepted %d connections...\n", conns)
		}
	}
}

func readConnection(conn net.Conn) {
	for {
		bytes := make([]byte, 32)
		conn.Read(bytes)
	}
}

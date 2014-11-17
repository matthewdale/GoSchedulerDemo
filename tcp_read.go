// watch "ps -eo pid,comm | grep exe/net_read | awk '{ print \$1 }' | xargs ps M | wc -l"
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
	numConns, err := strconv.Atoi(os.Args[0])
	if err != nil {
		panic(err)
	}

	for conns, port := 0, 1024; conns < numConns && port < maxport; port++ {
		laddr := &net.TCPAddr{
			Port: port,
		}
		listener, err := net.ListenTCP("tcp", laddr)
		if err != nil {
			fmt.Printf("Error listening on port %d: %s\n", port, err)
			continue
		}
		defer listener.Close()

		go acceptConnections(listener)
		conns++
	}

	fmt.Println("Press enter to continue...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func acceptConnections(listener *net.TCPListener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting new connection on %s: %s", listener.Addr().String(), err)
			continue
		}
		go readConnection(conn)
	}
}

func readConnection(conn *net.Conn) {
	for {
		bytes := make([]byte, 32)
		conn.Read(bytes)
	}
}

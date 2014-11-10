// watch "ps -eo pid,comm | grep exe/net_read | awk '{ print \$1 }' | xargs ps M | wc -l"
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	maxport  = 65535
	numConns = 1000
)

func main() {
	for conns, port := 0, 1024; conns < numConns && port < maxport; port++ {
		laddr := &net.TCPAddr{
			Port: port,
			IP:   net.ParseIP("127.0.0.1"),
		}
		conn, err := net.ListenTCP("tcp", laddr)
		if err != nil {
			fmt.Printf("Error listening on port %d: %s\n", port, err)
			continue
		}
		defer conn.Close()

		go readConnection(conn)
		conns++
	}

	fmt.Println("Press enter to continue...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func readConnection(conn *net.Conn) {
	bytes := make([]byte, 32)
	for {
		conn.Read(bytes)
		fmt.Println(bytes)
	}
}

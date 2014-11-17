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
	host := os.Args[0]
	numConns, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	for conns, port := 0, 1024; conns < numConns && port < maxport; port++ {
		raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%p", host, port))
		if err != nil {
			panic(err)
		}

		conn, err := net.DialTCP("tcp", nil, raddr)
		if err != nil {
			fmt.Printf("Error dialing port %d: %s\n", port, err)
			continue
		}
		defer conn.Close()

		go writeConnection(conn)
		conns++

		if conns%100 == 0 {
			fmt.Printf("Opened %d connections...\n", conns)
		}
	}

	fmt.Println("Press enter to continue...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func writeConnection(conn *net.TCPConn) {
	bytes := []byte("derping!")
	for {
		conn.Write(bytes)
	}
}

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
	host := os.Args[1]
	numConns, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	for conns, port := 0, 1024; conns < numConns && port < maxport; port++ {
		raddr := net.ResolveTCPAddr("udp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			panic(err)
		}

		conn, err := net.DialUDP("udp", nil, raddr)
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

func writeConnection(conn *net.UDPConn) {
	bytes := []byte("derping!")
	for {
		conn.Write(bytes)
	}
}

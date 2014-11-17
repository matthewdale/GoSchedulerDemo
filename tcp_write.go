package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {
	connString := os.Args[1]
	maxConns, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	raddr, err := net.ResolveTCPAddr("tcp", connString)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Dialing %s...\n", raddr.String())
	for conns := 0; conns < maxConns; conns++ {
		conn, err := net.DialTCP("tcp", nil, raddr)
		if err != nil {
			panic(fmt.Sprintf("Error dialing %s: %s\n", connString, err))
		}
		defer conn.Close()

		go writeConnection(conn)

		if (conns+1)%10 == 0 {
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

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {
	directory := os.Args[0]
	numConns, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	for i := 0; i < numConns; i++ {
		socketPath := directory + "/socket" + strconv.Itoa(i)
		raddr := &net.TCPAddr{
			Name: socketPath,
			Net:  "unix",
		}

		conn, err := net.DialUnix("unix", nil, raddr)
		if err != nil {
			fmt.Printf("Error dialing socket %s: %s\n", socketPath, err)
			continue
		}
		defer conn.Close()

		go writeConnection(conn)
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

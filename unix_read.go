// watch "ps -eo pid,comm | grep exe/unix_read | awk '{ print \$1 }' | xargs ps M | wc -l"
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
		laddr := &net.UnixAddr{
			Name: socketPath,
			Net:  "unix",
		}
		listener, err := net.ListenUnix("unix", laddr)
		if err != nil {
			fmt.Printf("Error opening socket %s: %s\n", socketPath, err)
			continue
		}
		defer listener.Close()
		defer os.Remove(socketPath)

		go acceptConnections(listener)
	}

	fmt.Println("Press enter to continue...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func acceptConnections(listener *net.UnixListener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting new connection on %s: %s", listener.Addr().String(), err)
			continue
		}
		go readConnection(conn)
	}
}

func readConnection(conn *net.UnixConn) {
	for {
		bytes := make([]byte, 32)
		conn.ReadFromUnix(bytes)
	}
}

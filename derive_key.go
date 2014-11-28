package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/gob"
	"net"
	"os"

	"code.google.com/p/go.crypto/pbkdf2"
	"github.com/coreos/go-log/log"
)

type DeriveKeyRequest struct {
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

var bufferedLog *log.Logger

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("First argument must be socket path.")
	}
	socket := os.Args[1]

	// create a buffered writer for stdout
	buffer := bufio.NewWriter(os.Stdout)
	defer buffer.Flush()
	bufferedLog = log.NewSimple(log.CombinedSink(buffer, log.BasicFormat, log.BasicFields))

	laddr := &net.UnixAddr{
		Name: socket,
		Net:  "unix",
	}
	listener, err := net.ListenUnix("unix", laddr)
	if err != nil {
		log.Fatalf("Error listening on %s: %s", socket, err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Errorf("Error accepting connection on socket %s: %s", socket, err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	for {
		dec := gob.NewDecoder(conn)
		var request DeriveKeyRequest
		dec.Decode(&request)

		deriveKey(request.Password, request.Salt)
	}
}

func deriveKey(password, salt string) {
	dk := pbkdf2.Key([]byte(password), []byte(salt), 8192, 256, sha1.New)

	bufferedLog.Info("Derived key ", base64.StdEncoding.EncodeToString(dk))
}

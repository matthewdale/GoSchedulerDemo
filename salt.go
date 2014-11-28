package main

import (
	"bufio"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	mrand "math/rand"
	"net"
	"net/http"
	"os"

	"github.com/coreos/go-log/log"
)

const SaltLength = 32

type Request struct {
	Password string `json:"password"`
}

type Response struct {
	Salt string `json:"salt"`
}

type DeriveKeyRequest struct {
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

func init() {
	// seed math/rand with some bytes from crypto/rand
	bytes := make([]byte, 8)
	crand.Read(bytes)
	seed := int64(binary.LittleEndian.Uint64(bytes))
	mrand.Seed(seed)
}

var bufferedLog *log.Logger
var deriveKeyConnection net.Conn

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("First argument must be socket path.")
	}
	socket := os.Args[1]
	var err error
	deriveKeyConnection, err = net.Dial("unix", socket)
	if err != nil {
		log.Fatalf("Error dialing socket %s: %s", socket, err)
	}

	// create a buffered writer for stdout
	buffer := bufio.NewWriter(os.Stdout)
	defer buffer.Flush()
	bufferedLog = log.NewSimple(log.CombinedSink(buffer, log.BasicFormat, log.BasicFields))

	http.HandleFunc("/", generateSalt)
	http.ListenAndServe(":8080", nil)
}

func generateSalt(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var requestBody Request
	decoder.Decode(&requestBody)

	bytes := make([]byte, SaltLength)
	randomBytes(bytes)
	salt := base64.StdEncoding.EncodeToString(bytes)
	bufferedLog.Info("Generated salt ", salt)

	dkRequest := DeriveKeyRequest{
		Password: requestBody.Password,
		Salt:     salt,
	}
	enc := gob.NewEncoder(deriveKeyConnection)
	err := enc.Encode(dkRequest)
	if err != nil {
		log.Error("Error encoding gob: ", err)
	}

	response := Response{
		Salt: salt,
	}
	js, err := json.Marshal(response)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

func randomBytes(p []byte) int {
	for i := range p {
		p[i] = byte(mrand.Int63() & 0xff)
	}
	return len(p)
}

// watch "ps -eo pid,comm | grep exe/starvation | awk '{ print \$1 }' | xargs ps M"
package main

import (
	"bufio"
	crand "crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	mrand "math/rand"
	"net/http"
	"os"

	"code.google.com/p/go.crypto/pbkdf2"
	"github.com/coreos/go-log/log"
)

const SaltLength = 32

type Request struct {
	Password string `json:"password"`
}

type Response struct {
	Salt string `json:"salt"`
}

func init() {
	// seed math/rand with some bytes from crypto/rand
	bytes := make([]byte, 8)
	crand.Read(bytes)
	seed := int64(binary.LittleEndian.Uint64(bytes))
	mrand.Seed(seed)
}

var logger *log.Logger

func main() {
	// create a buffered writer for stdout
	buffer := bufio.NewWriter(os.Stdout)
	defer buffer.Flush()
	logger = log.NewSimple(log.CombinedSink(buffer, log.BasicFormat, log.BasicFields))

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
	logger.Info("Generated salt ", salt)

	go deriveKey(requestBody.Password, salt)

	response := Response{
		Salt: string(salt),
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

func deriveKey(password, salt string) {
	dk := pbkdf2.Key([]byte(password), []byte(salt), 8192, 256, sha1.New)

	logger.Info("Derived key ", base64.StdEncoding.EncodeToString(dk))
}

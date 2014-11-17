// watch "ps -eo pid,comm | grep exe/starvation | awk '{ print \$1 }' | xargs ps M"
package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"code.google.com/p/go.crypto/pbkdf2"
)

type Message struct {
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

func main() {
	http.HandleFunc("/", handleMessage)
	http.ListenAndServe(":8080", nil)
}

func handleMessage(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var message Message
	decoder.Decode(&message)
	go deriveKey([]byte(message.Password), []byte(message.Salt))
}

func deriveKey(password, salt []byte) {
	dk := pbkdf2.Key(password, salt, 4096, 256, sha1.New)
	encoder := base64.NewEncoder(base64.StdEncoding, os.Stdout)
	encoder.Write(dk)
	encoder.Close()

	fmt.Println("")
}

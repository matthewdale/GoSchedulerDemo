// watch "ps -eo pid,comm | grep exe/file | awk '{ print \$1 }' | xargs ps M | wc -l"
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

const numFiles = 1000

func main() {
	fmt.Printf("Creating %d files in %s\n", numFiles, os.TempDir())

	var files [numFiles]*os.File
	for i := 0; i < numFiles; i++ {
		file, err := ioutil.TempFile(os.TempDir(), "test")
		if err != nil {
			fmt.Printf("Error creating tempfile: %s\n", err)
			continue
		}
		files[i] = file

		go readFile(file)

		defer file.Close()
		defer os.Remove(file.Name())
	}

	fmt.Println("Press enter to continue...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func readFile(file *os.File) {
	bytes := []byte("derping!")
	file.Write(bytes)
	for {
		read := make([]byte, 32)
		file.Read(read)
	}
}

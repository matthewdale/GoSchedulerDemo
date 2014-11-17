// watch "ps -eo pid,comm | grep exe/file | awk '{ print \$1 }' | xargs ps M | wc -l"
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	numFiles, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	fmt.Printf("Creating %d files in %s\n", numFiles, os.TempDir())

	files := make([]*os.File, numFiles)
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
	bytes := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	// file.Write(bytes)
	for {
		file.Write(bytes)
		read := make([]byte, 64)
		file.Read(read)
	}
}

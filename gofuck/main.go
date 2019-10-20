package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/chrissexton/gofuck"
)

func main() {
	m := gofuck.NewStdin()
	instructions := make([]byte, 0)

	// read in the instructions
	if len(os.Args) < 2 {
		fmt.Printf("No input given")
		return
	}

	file, err := os.Open(os.Args[1])
	reader := bufio.NewReader(file)
	instructions, err = ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	if err := m.Run(instructions); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

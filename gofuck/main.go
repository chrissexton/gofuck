package main

import (
	"bufio"
	"io/ioutil"
	"os"

	"github.com/chrissexton/gofuck"
)

func main() {
	m := gofuck.NewStdin()
	instructions := make([]byte, 0)

	// read in the instructions
	if len(os.Args) < 2 {
		panic("No input given")
	}

	file, err := os.Open(os.Args[1])
	reader := bufio.NewReader(file)
	instructions, err = ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	m.Run(instructions)
}

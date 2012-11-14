package gofuck

import (
	"bufio"
	"fmt"
	"os"
)

const (
	MEM_MAX = 3000000
	MEM_STD = 30000
)

type Machine struct {
	array  []byte
	ptr    int
	reader *bufio.Reader
}

// Returns a new machine with standard memory size
func New() *Machine {
	bytes := make([]byte, MEM_STD)
	return &Machine{
		array:  bytes,
		ptr:    0,
		reader: bufio.NewReader(os.Stdin),
	}
}

// Implements '>'
func (m *Machine) PtrIncr() {
	if m.ptr >= len(m.array)-1 {
		m.array = append(m.array, 0)
	}
	m.ptr += 1
	if m.ptr > MEM_MAX {
		panic("Memory overflow")
	}
}

// Implements '<'
func (m *Machine) PtrDecr() {
	if m.ptr == 0 {
		panic("Memory underflow")
	}
	m.ptr -= 1
}

// Implements '+'
func (m *Machine) ByteIncr() {
	if m.ptr > len(m.array)-1 {
		fmt.Printf("Memory overflow, ptr=%d\n", m.ptr)
	}
	m.array[m.ptr] += 1
}

// Implements '-'
func (m *Machine) ByteDecr() {
	m.array[m.ptr] -= 1
}

// Implements the '.' command
func (m *Machine) Output() {
	fmt.Printf("%c", m.array[m.ptr])
}

// Implements the ',' command
func (m *Machine) Input() {
	input, err := m.reader.ReadByte()
	if err != nil {
		m.array[m.ptr] = 0
	} else {
		m.array[m.ptr] = input
	}
}

func (m *Machine) value() byte {
	return m.array[m.ptr]
}

// Run the whole program specified by input
func (m *Machine) Run(input []byte) {
	// execute the program
	for ip := 0; ip < len(input); ip++ {
		instr := input[ip]
		// if *ptr == 0, jump to ]
		if instr == '[' && m.value() == 0 {
			for lc := 1; lc > 0; {
				ip += 1
				if input[ip] == ']' {
					lc -= 1
				} else if input[ip] == '[' {
					lc += 1
				}
			}
		} else if instr == ']' && m.value() != 0 {
			// if *ptr != 0, go back to [
			for lc := 1; lc > 0; {
				ip -= 1
				if input[ip] == ']' {
					lc += 1
				} else if input[ip] == '[' {
					lc -= 1
				}
			}
		} else if instr == '>' {
			m.PtrIncr()
		} else if instr == '<' {
			m.PtrDecr()
		} else if instr == '+' {
			m.ByteIncr()
		} else if instr == '-' {
			m.ByteDecr()
		} else if instr == '.' {
			m.Output()
		} else if instr == ',' {
			m.Input()
		}

	}
}

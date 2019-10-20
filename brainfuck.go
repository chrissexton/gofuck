// gofuck is a simple brainfuck interpreter written in Go
package gofuck

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	MEM_STD = 30000
)

type ErrMemoryOverflow struct {}
func (e ErrMemoryOverflow) Error() string { return "Memory Overflow" }
type ErrMemoryUnderflow struct {}
func (e ErrMemoryUnderflow) Error() string { return "Memory Underflow" }
type ErrInstructionLimit struct {}
func (e ErrInstructionLimit) Error() string { return "Instruction Limit Reached" }

type Machine struct {
	array  []byte
	ptr    int
	reader *bufio.Reader
	writer io.Writer

	// InstructionLimit prevents a runaway execution if running under a controlled environment
	InstructionLimit int

	// MemMax is the limit of how large memory can expand
	MemMax int
}

// Returns a new machine with standard memory size
func NewStdin() *Machine {
	return New(os.Stdin, os.Stdout)
}

func New(in io.Reader, out io.Writer) *Machine {
	bytes := make([]byte, MEM_STD)
	return &Machine{
		array:  bytes,
		ptr:    0,
		reader: bufio.NewReader(in),
		writer: out,
		InstructionLimit: 0,
		MemMax: 3000000,
	}
}

// Implements '>'
func (m *Machine) PtrIncr() error {
	if m.ptr >= len(m.array)-1 {
		m.array = append(m.array, 0)
	}
	m.ptr += 1
	if m.ptr > m.MemMax {
		return ErrMemoryOverflow{}
	}
	return nil
}

// Implements '<'
func (m *Machine) PtrDecr() error {
	if m.ptr == 0 {
		return ErrMemoryUnderflow{}
	}
	m.ptr -= 1
	return nil
}

// Implements '+'
func (m *Machine) ByteIncr() {
	m.array[m.ptr] += 1
}

// Implements '-'
func (m *Machine) ByteDecr() {
	m.array[m.ptr] -= 1
}

// Implements the '.' command
func (m *Machine) Output() {
	fmt.Fprintf(m.writer, "%c", m.array[m.ptr])
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
func (m *Machine) Run(input []byte) error {
	instrCount := 0
	for ip := 0; ip < len(input); ip++ {
		instrCount++
		if m.InstructionLimit > 0 && instrCount > m.InstructionLimit {
			return ErrInstructionLimit{}
		}
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
			if err := m.PtrIncr(); err != nil {
				return err
			}
		} else if instr == '<' {
			if err := m.PtrDecr(); err != nil {
				return err
			}
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

	return nil
}

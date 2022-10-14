package unittest_test

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	read2()
}

func read2() {
	s := strings.NewReader("ABCDEFG")
	br := bufio.NewReader(s)
	b, _ := br.Peek(2)
	fmt.Printf("%s\n", b) // AB
	br.Discard(len(b))

	c, _ := br.Peek(2)
	fmt.Printf("%s\n", c) // CD

	b2 := make([]byte, 2)
	br.Read(b2)
	fmt.Printf("%s\n", b2) // CD
}

func read1() {
	s := strings.NewReader("ABCDEFG")
	br := bufio.NewReader(s)

	b, _ := br.Peek(5)
	fmt.Printf("%s\n", b)
	// ABCDE

	b[0] = 'a'
	b, _ = br.Peek(5)
	fmt.Printf("%s\n", b)
}

package bencode

import (
	"io"
	"strings"
	"testing"
)

func TestBencode(t *testing.T) {
	s := strings.NewReader("Hello World!\n")
	i := io.Reader(s)
	bReader := NewBencodeReader(&i)

	bReader.DecodeStream()
}

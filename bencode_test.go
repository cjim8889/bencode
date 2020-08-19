package bencode

import (
	"io"
	"strings"
	"testing"
)

func TestBencode(t *testing.T) {
	s := strings.NewReader("i-51e")
	i := io.Reader(s)
	bReader := NewBencodeReader(&i)

	result, err := bReader.DecodeStream()
	if err != nil {
		t.Error(err.Error())
	}

	if result != -51 {
		t.Error("Error Parsing Int")
	}

}

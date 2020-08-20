package bencode

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestParseInt(t *testing.T) {
	s := strings.NewReader("i-500000000000000e")
	i := io.Reader(s)
	bReader := NewBencodeReader(&i)

	result, err := bReader.DecodeStream()
	if err != nil {
		t.Error(err.Error())
	}

	if result != -500000000000000 {
		t.Error("Error Parsing Int")
	}
}

func TestParseIntNegative(t *testing.T) {
	s := strings.NewReader("i--0e")
	i := io.Reader(s)
	bReader := NewBencodeReader(&i)

	_, err := bReader.DecodeStream()
	if err == nil {
		t.Error("Double negation not being detected")
	}
}

func TestParseBytes(t *testing.T) {
	s := strings.NewReader("2:ab")
	i := io.Reader(s)
	bReader := NewBencodeReader(&i)

	result, err := bReader.DecodeStream()
	if err != nil {
		t.Error(err.Error())
	}

	if string(result.([]byte)) != "ab" {
		t.Error("Parse bytes error")
	}
}

func TestParseList(t *testing.T) {
	s := strings.NewReader("l2:abi5ee")
	i := io.Reader(s)
	bReader := NewBencodeReader(&i)

	result, err := bReader.DecodeStream()
	if err != nil {
		t.Error(err.Error())
	}

	fmt.Println(result)
}
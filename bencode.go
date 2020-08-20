package bencode

import (
	"bufio"
	"io"
)

type BencodeReader struct {
	reader *bufio.Reader
}

func NewBencodeReader(r *io.Reader) BencodeReader {
	return BencodeReader{bufio.NewReader(*r)}
}

func (r *BencodeReader) DecodeStream() (interface{}, error)  {
	result, err := Parse(r)
	if err != nil {
		return nil, err
	}


	return result, nil
}


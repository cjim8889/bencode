package bencode

import (
	"bufio"
	"io"
)

type Reader struct {
	reader *bufio.Reader
}

func NewBencodeReader(r io.Reader) Reader {
	return Reader{bufio.NewReader(r)}
}

func (r *Reader) DecodeStream() (interface{}, error)  {
	result, err := Parse(r)
	if err != nil {
		return nil, err
	}


	return result, nil
}


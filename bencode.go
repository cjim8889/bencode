package bencode

import (
	"bufio"
	"fmt"
	"io"
)

type BencodeReader struct {
	reader *bufio.Reader
}

func NewBencodeReader(r *io.Reader) BencodeReader {
	return BencodeReader{bufio.NewReader(*r)}
}

func (r *BencodeReader) DecodeStream() (interface{}, error)  {
	bytes := make([]byte, 64)
	currentHead := 0
	readMore := func () (int, error) {
		return r.reader.Read(bytes)
	}

	for {
		numOfBytes, err := readMore()
		if err != nil || numOfBytes == 0 {
			break
		}

		switch b := bytes[currentHead]; b {
		case '':
			
		
		}

	}

	return nil, nil
}

func parse(r *BencodeReader) (interface{}, error) {
	bytes := make([]byte, 64)

	//currentHead := 0
	readMore := func () (int, error) {
		return r.reader.Read(bytes)
	}

	for {
		numOfBytes, err := readMore()
		if err != nil || numOfBytes == 0 {
			break
		}

		fmt.Printf("%q\n", bytes)
	}

	return nil, nil
}


func init() {
}

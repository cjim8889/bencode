package bencode

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type BencodeReader struct {
	reader *bufio.Reader
}

type ParseIntError struct {
	error string
}

func (e ParseIntError) Error() string {
	return fmt.Sprintf("Parse Int Error: %v\n", e.error)
}

func NewBencodeReader(r *io.Reader) BencodeReader {
	return BencodeReader{bufio.NewReader(*r)}
}

func (r *BencodeReader) DecodeStream() (interface{}, error)  {
	for {
		b, err := r.reader.Peek(1)
		if err != nil {
			break
		}

		switch b[0] {
		case 'i': {
			result, err := ParseInt(r)
			if err != nil {
				return nil, err
			}

			return result, nil
		}
		}
	}

	return nil, nil
}

func ParseInt(r *BencodeReader) (int, error) {
	numberBuffer := make([]byte, 12)
	bufferHead := 0
	confirmation, err := r.reader.ReadByte()
	negativeFlag := -1
	if err != nil || confirmation != 'i' {
		return 0, ParseIntError{}
	}

	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return 0, ParseIntError{err.Error()}
		}

		if b == 'e' {
			break
		}

		if b == '-' {
			if negativeFlag != -1 {
				return 0, ParseIntError{"Duplicate hyphen"}
			}

			if bufferHead != 0 {
				return 0, ParseIntError{"Misplaced hyphen"}
			}

			negativeFlag = 1
			continue
		}

		if b < 48 || b > 57 {
			return 0, ParseIntError{}
		}

		if bufferHead < cap(numberBuffer) {
			numberBuffer[bufferHead] = b
		} else {
			numberBuffer = append(numberBuffer, b)
		}

		bufferHead += 1
	}

	if bufferHead == 0 {
		return 0, ParseIntError{}
	}

	strRep := string(numberBuffer[:bufferHead])
	if strings.HasPrefix(strRep, "0") && (len(strRep) > 1 || negativeFlag != -1) {
		return 0, ParseIntError{}
	}

	result, err := strconv.Atoi(strRep[:bufferHead])
	if err != nil {
		return 0, ParseIntError{}
	}

	if negativeFlag == 1 {
		return -result, nil
	} else {
		return result, nil
	}
}


func init() {
}

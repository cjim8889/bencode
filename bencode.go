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

type ParseError struct {
	error string
}

type BencodeListCell struct {
	value interface{}
}

func (e ParseError) Error() string {
	return fmt.Sprintf("Parse Int Error: %v\n", e.error)
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

func ParseList(r *BencodeReader) ([]BencodeListCell, error) {
	confirmation, err := r.reader.ReadByte()
	if err != nil {
		return nil, ParseError{err.Error()}
	}

	if confirmation != 'l' {
		panic("FATAL ERROR: PARSE LIST")
	}

	var result []BencodeListCell

	for {
		c, err := Parse(r)
		if err != nil {
			return nil, ParseError{fmt.Sprintf("ParseList Error: %v\n", err.Error())}
		}

		result = append(result, BencodeListCell{c})

		p, err := r.reader.Peek(1)
		if err != nil {
			if err == io.EOF {
				return nil, ParseError{"ParseList Error: EOF"}
			}

			return nil, ParseError{"ParseList Error: Forward peek failed"}
		}

		if p[0] == 'e' {
			break
		}
	}

	return result, nil
}

func Parse(r *BencodeReader) (interface{}, error) {
	b, err := r.reader.Peek(1)
	if err != nil {
		return nil, err
	}

	switch c := b[0]; {
	case c == 'i': {
		result, err := ParseInt(r)
		if err != nil {
			return nil, err
		}

		return result, nil
	}
	case c >= 48 && c <= 57: {
		result, err := ParseByteString(r)
		if err != nil {
			return nil, err
		}

		return result, nil
	}
	case c == 'l': {
		result, err := ParseList(r)
		if err != nil {
			return nil, err
		}

		return result, nil
	}
	}

	return nil, ParseError{}
}

func ParseByteString(r *BencodeReader) ([]byte, error) {
	byteCountA, err := r.reader.ReadBytes(':')
	if err != nil {
		return nil, ParseError{}
	}

	strRep := string(byteCountA[:len(byteCountA) - 1])
	byteCount, err := strconv.Atoi(strRep)
	if err != nil || byteCount <= 0 {
		return nil, ParseError{}
	}

	result := make([]byte, byteCount)
	n, err := r.reader.Read(result)
	if err != nil {
		return nil, ParseError{}
	}

	if byteCount != n {
		return nil, ParseError{}
	}

	return result, nil
}

func ParseInt(r *BencodeReader) (int, error) {
	var numberBuffer []byte
	confirmation, err := r.reader.ReadByte()

	if err != nil || confirmation != 'i' {
		return 0, ParseError{}
	}

	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return 0, ParseError{err.Error()}
		}

		if b == 'e' {
			break
		}

		if b < 48 || b > 57 {
			if b != '-' {
				return 0, ParseError{"Illegal ascii code"}
			}
		}

		numberBuffer = append(numberBuffer, b)
	}

	strRep := string(numberBuffer)

	if strings.HasPrefix(strRep, "0") && len(strRep) > 1 {
		return 0, ParseError{"Number can not start with 0"}
	}

	if strings.HasPrefix(strRep, "-0") {
		return 0, ParseError{"-0 is not permitted"}
	}

	result, err := strconv.Atoi(strRep)

	if err != nil {
		return 0, ParseError{err.Error()}
	}

	return result, nil
}


func init() {
}

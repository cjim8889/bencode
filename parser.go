package bencode

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type ParserError struct {
	error string
}

type BencodeCell struct {
	Value interface{}
}

func (e ParserError) Error() string {
	return fmt.Sprintf("Parser Error: %v\n", e.error)
}

func ParseList(r *Reader) ([]BencodeCell, error) {
	confirmation, err := r.reader.ReadByte()
	if err != nil {
		return nil, ParserError{err.Error()}
	}

	if confirmation != 'l' {
		return nil, ParserError{"FATAL ERROR: PARSE LIST"}
	}

	var result []BencodeCell

	for {
		c, err := Parse(r)
		if err != nil {
			return nil, ParserError{fmt.Sprintf("ParseList Error: %v\n", err.Error())}
		}

		result = append(result, BencodeCell{c})

		p, err := r.reader.Peek(1)
		if err != nil {
			if err == io.EOF {
				return nil, ParserError{"ParseList Error: EOF"}
			}

			return nil, ParserError{"ParseList Error: Forward peek failed"}
		}

		if p[0] == 'e' {
			r.reader.ReadByte()
			break
		}
	}

	return result, nil
}

func Parse(r *Reader) (interface{}, error) {
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
	case c == 'd': {
		result, err := ParseDictionary(r)
		if err != nil {
			return nil, err
		}

		return result, nil
	}
	}

	return nil, ParserError{"Fatal error"}
}

func ParseByteString(r *Reader) (string, error) {
	byteCountA, err := r.reader.ReadBytes(':')
	if err != nil {
		return "", ParserError{}
	}

	strRep := string(byteCountA[:len(byteCountA) - 1])
	byteCount, err := strconv.Atoi(strRep)
	if err != nil || byteCount <= 0 {
		return "", ParserError{}
	}

	result, err := readFullByteString(r, byteCount)
	if err != nil {
		return "", err
	}

	if byteCount != len(*result) {
		return "", ParserError{"ByteCount does not match with the length."}
	}

	return string(*result), nil
}

func readFullByteString(r *Reader, length int) (*[]byte, error) {
	result := make([]byte, 0, length)
	var tempLength int
	if 512 <= length {
		tempLength = 512
	} else {
		tempLength = length
	}


	temp := make([]byte, tempLength)
	totalByteRead := 0

	for {
		n, err := r.reader.Read(temp)
		if err != nil {
			return nil, ParserError{"ByteString reading failed"}
		}

		result = append(result, temp[:n]...)
		totalByteRead += n

		if length - totalByteRead < 512 {
			temp = make([]byte, length - totalByteRead)
		}

		if totalByteRead == length {
			break
		}
	}

	if totalByteRead != length {
		return nil, ParserError{""}
	}

	return &result, nil
}

func ParseDictionary(r *Reader) (map[string]BencodeCell, error) {
	confirmation, err := r.reader.ReadByte()
	if err != nil {
		return nil, ParserError{err.Error()}
	}

	if confirmation != 'd' {
		return nil, ParserError{"FATAL ERROR: Parse Dictionary"}
	}

	var lastIndex string
	result := make(map[string]BencodeCell)

	for {
		i, err := ParseByteString(r)
		if err != nil {
			return nil, ParserError{fmt.Sprintf("ParseDictionary Error: Error parsing index %v\n", err.Error())}
		}

		strComp := strings.Compare(lastIndex, i)
		if lastIndex != "" && strComp == 1 {
			return nil, ParserError{"ParseDictionary Error: All keys must be byte strings and must appear in lexicographical order"}
		}

		v, err := Parse(r)
		if err != nil {
			return nil, ParserError{"ParseDictionary Error: Error parsing Value"}
		}

		result[i] = BencodeCell{v}

		p, err := r.reader.Peek(1)
		if err != nil {
			return nil, ParserError{"ParseDictionary Error: Error peeking"}
		}

		if p[0] == 'e' {
			r.reader.ReadByte()
			break
		}
	}

	return result, nil
}


func ParseInt(r *Reader) (int, error) {
	var numberBuffer []byte
	confirmation, err := r.reader.ReadByte()

	if err != nil || confirmation != 'i' {
		return 0, ParserError{}
	}

	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return 0, ParserError{err.Error()}
		}

		if b == 'e' {
			break
		}

		if b < 48 || b > 57 {
			if b != '-' {
				return 0, ParserError{"Illegal ascii code"}
			}
		}

		numberBuffer = append(numberBuffer, b)
	}

	strRep := string(numberBuffer)

	if strings.HasPrefix(strRep, "0") && len(strRep) > 1 {
		return 0, ParserError{"Number can not start with 0"}
	}

	if strings.HasPrefix(strRep, "-0") {
		return 0, ParserError{"-0 is not permitted"}
	}

	result, err := strconv.Atoi(strRep)

	if err != nil {
		return 0, ParserError{err.Error()}
	}

	return result, nil
}
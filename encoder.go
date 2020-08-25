package bencode

import (
	"fmt"
	"sort"
	"strconv"
)

func EncodeInt(i int) []byte {
	s := fmt.Sprintf("i%ve", strconv.Itoa(i))
	return []byte(s)
}

func EncodeString(str string) []byte {
	bytes := []byte(str)
	length := []byte(fmt.Sprintf("%v:", len(bytes)))

	return append(length, bytes...)
}

type EncoderError struct {
	error string
}

func (e EncoderError) Error() string {
	return fmt.Sprintf("Parser Error: %v\n", e.error)
}

func EncodeDictionary(d map[string]BencodeCell) ([]byte, error) {
	result := make([]byte, 0, len(d)*2)
	result = append(result, 'd')
	k := getAllKeysSorted(d)
	for _, v := range k {
		mappedVal, ok := d[v]
		if ok != true {
			return nil, EncoderError{"Key not found"}
		}

		encodedVal, err := Encode(mappedVal.Value)
		if err != nil {
			return nil, err
		}

		result = append(result, EncodeString(v)...)
		result = append(result, encodedVal...)
	}

	result = append(result, 'e')

	return result, nil
}

func getAllKeysSorted(m map[string]BencodeCell) []string {
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

func Encode(data interface{}) ([]byte, error) {
	switch data.(type) {
	case int:
		return EncodeInt(data.(int)), nil
	case string:
		return EncodeString(data.(string)), nil
	case []BencodeCell: {
		r, err := EncodeList(data.([]BencodeCell))
		if err != nil {
			return nil, err
		}

		return r, nil
	}
	case map[string]BencodeCell: {
		r, err := EncodeDictionary(data.(map[string]BencodeCell))
		if err != nil {
			return nil, err
		}

		return r, nil
	}
	default:
		return nil, EncoderError{"Unknown type of List"}
	}
}

func EncodeList(l []BencodeCell) ([]byte, error) {
	result := []byte("l")

	for _, v := range l {
		r, err := Encode(v.Value)
		if err != nil {
			return nil, err
		}

		result = append(result, r...)
	}

	result = append(result, 'e')
	return result, nil
}
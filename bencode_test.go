package bencode

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestParseInt(t *testing.T) {
	s := strings.NewReader("i-500000000000000e")
	i := io.Reader(s)
	bReader := NewBencodeReader(i)

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
	bReader := NewBencodeReader(i)

	_, err := bReader.DecodeStream()
	if err == nil {
		t.Error("Double negation not being detected")
	}
}

func TestParseBytes(t *testing.T) {
	s := strings.NewReader("2:ab")
	i := io.Reader(s)
	bReader := NewBencodeReader(i)

	result, err := bReader.DecodeStream()
	if err != nil {
		t.Error(err.Error())
	}

	if result.(string) != "ab" {
		t.Error("Parse bytes error")
	}
}

func TestParseList(t *testing.T) {
	s := strings.NewReader("l2:abi5ee")
	i := io.Reader(s)
	bReader := NewBencodeReader(i)

	result, err := bReader.DecodeStream()
	if err != nil {
		t.Error(err.Error())
	}

	r := result.([]BencodeCell)
	if r[0].value.(string) != "ab" && r[1].value.(int) != 5 {
		t.Error("ParseList failed")
	}
}

func TestParseDictionary(t *testing.T) {
	s := strings.NewReader("d3:bar4:spam3:fooi42ee")
	i := io.Reader(s)
	bReader := NewBencodeReader(i)

	result, err := bReader.DecodeStream()
	if err != nil {
		t.Error(err.Error())
	}

	r := result.(map[string]BencodeCell)
	if r["bar"].value.(string) != "spam" {
		t.Error("Test parse Dictionary failed")
	}
}

func TestParseListOfList(t *testing.T) {
	s := strings.NewReader("d8:announce35:http://tracker.tfile.co:80/announce13:announce-listll35:http://tracker.tfile.co:80/announceel33:udp://open.stealth.si:80/announceeee")
	i := io.Reader(s)
	bReader := NewBencodeReader(i)

	result, err := bReader.DecodeStream()
	if err != nil {
		t.Error(err.Error())
	}

	r := result.(map[string]BencodeCell)
	if len(r) != 2 {
		t.Error("Error parsing list of list")
	}
}

func TestEncoder(t *testing.T) {
	s := []BencodeCell{{10}, {"nima"}}

	r, err := EncodeList(s)
	if err != nil {
		t.Error(err.Error())
	}

	if string(r) != "li10e4:nimae" {
		t.Error("Encoder test failed")
	}
}

func TestEncodeDictionary(t *testing.T) {
	s := make(map[string]BencodeCell)
	s["hello"] = BencodeCell{1}
	s["world"] = BencodeCell{"2"}

	r, err := Encode(s)
	if err != nil {
		t.Error(err.Error())
	}

	if string(r) != "d5:helloi1e5:world1:2e" {
		t.Error("Encoder Dictionary test failed")
	}

}

func TestBencodeParse(t *testing.T) {
	file, err := os.Open("/Users/wuhaochen/go/src/github.com/cjim8889/bencode/test.torrent")
	if err != nil {
		t.Error(err.Error())
	}

	bReader := NewBencodeReader(file)
	r, err := bReader.DecodeStream()
	if err != nil {
		t.Error(err.Error())
	}

	_, ok := r.(map[string]BencodeCell)
	if !ok {
		t.Error("Error reading torrent")
	}
}
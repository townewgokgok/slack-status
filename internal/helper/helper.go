package helper

import (
	"bytes"
	"io/ioutil"
	"strconv"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func ParseFloat64(str string) float64 {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return .0
	}
	return val
}

func ConvertFromShiftJIS(sjis []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(sjis), japanese.ShiftJIS.NewDecoder())
	utf8, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(utf8), nil
}

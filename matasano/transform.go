package matasano

import (
	"encoding/hex"
	"encoding/base64"
	"bytes"
	"errors"
)

func HexToBase64(input string) (output string, err error) {
	bytes, decodeError := hex.DecodeString(input)
	if decodeError != nil {
		return "", decodeError
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func FixedXor(input string, mask string) (output string, err error) {
	var decodeErr error
	var inputBytes, maskBytes []byte
	inputBytes, decodeErr = hex.DecodeString(input)
	if decodeErr != nil {
		return "", decodeErr
	}
	maskBytes, decodeErr = hex.DecodeString(mask)
	if decodeErr != nil {
		return "", decodeErr
	}

	resultBytes := bytes.Map(makeRepeatingXorClojure(maskBytes), inputBytes)
	return hex.EncodeToString(resultBytes), nil
}

func FixedXorWithSingleByteMask(input []byte, mask byte) []byte {
	return bytes.Map(makeRepeatingXorClojure([]byte{mask}), input)
}

func RepeatingXor(input []byte, mask []byte) []byte {
	return bytes.Map(makeRepeatingXorClojure(mask), input)
}

func makeRepeatingXorClojure(mask [] byte) func(r rune) rune {
	i := 0
	maskLength := len(mask)
	return func(r rune) rune {
		result := rune(byte(r) ^ mask[i])
		i = (i+1)%maskLength
		return result
	}
}

func HemingDistance(a []byte, b []byte) (uint, error) {
	if len(a) != len(b) {
		return 0, errors.New("Both arrays should be the same length.")
	}
	diffBits := uint(0)
	for i := 0 ; i < len(a) ; i++ {
		diff := a[i] ^ b[i]
		for j := uint(0) ; j < 8 ; j++ {
			diffBits += uint((diff >> j) & 1)
		}
	}
	return diffBits, nil
}

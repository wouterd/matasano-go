package matasano

import (
	"encoding/hex"
	"encoding/base64"
	"bytes"
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

	resultBytes := bytes.Map(makeFixedXorClojure(maskBytes), inputBytes)
	return hex.EncodeToString(resultBytes), nil
}

func makeFixedXorClojure(mask []byte) func(r rune) rune {
	i := 0
	return func(r rune) rune {
		result := rune(byte(r) ^ mask[i])
		i++
		return result
	}
}

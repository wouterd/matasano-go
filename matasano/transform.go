package matasano

import (
	"encoding/hex"
	"encoding/base64"
	"bytes"
	"errors"
	"crypto/aes"
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

func DecryptAES128ECB(encrypted []byte, key []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	decrypted := new(bytes.Buffer)
	input := bytes.NewBuffer(encrypted)
	bufIn := make([]byte, cipher.BlockSize())
	bufOut := make([]byte, cipher.BlockSize())
	for {
		_, err := input.Read(bufIn)
		if err != nil {
			// EOF
			break
		}
		cipher.Decrypt(bufOut, bufIn)
		decrypted.Write(bufOut)
	}
	return decrypted.Bytes(), nil
}

/*
	Counts the amount of matching blocks in a byte array given a certain block size (in bytes).
	Every match will be single-counted, so if block a matches block b, then block b matching block a will not be
	counted.
 */
func CountMatchingBlocks(data []byte, blockSize int) int {
	blocks := len(data) / blockSize
	matchingBlocks := 0
	for i := 0 ; i < blocks ; i++ {
		offset := i * blockSize
		this := data[offset:offset+blockSize]
		for j := i + 1 ; j < blocks ; j++ {
			offset := j * blockSize
			that := data[offset:offset+blockSize]
			if bytes.Equal(this, that) {
				matchingBlocks++
			}
		}
	}
	return matchingBlocks
}

func PadBufferUsingPKCS7(input []byte, blocksize int) []byte {
	if len(input) >= blocksize {
		return input
	}
	padding := byte(blocksize - len(input))
	var output bytes.Buffer
	output.Write(input)
	for i := byte(0) ; i < padding ; i++ {
		output.WriteByte(padding)
	}
	return output.Bytes()
}

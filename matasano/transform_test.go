package matasano

import (
	"testing"
	"github.com/stretchr/testify/require"
	"encoding/hex"
	"encoding/base64"
)

func TestHexToBase64(t *testing.T) {
	input := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	expected := "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
	result, convertError := HexToBase64(input)
	if convertError != nil {
		t.Error("Error occured when converting")
	}
	if result != expected {
		t.Error("Expected ", expected, ", got ", result)
	}
}

func TestFixedXor(t *testing.T) {
	input := "1c0111001f010100061a024b53535009181c"
	mask := "686974207468652062756c6c277320657965"
	expected := "746865206b696420646f6e277420706c6179"
	result, err := FixedXor(input, mask)
	if err != nil {
		t.Error("Error occured when converting from hex")
	}
	if result != expected {
		t.Error("Expected ", expected, ", but got ", result)
	}
}

func TestFixedXorWithSingleByteMaskForZero(t *testing.T) {
	input := "Hi there"
	result := FixedXorWithSingleByteMask([]byte(input), byte(0))
	expected := []byte(input)
	require.Equal(t, expected, result, "Did not get the same result back after XOR with 0")
}

func TestRepeatingXor(t *testing.T) {
	input := "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal"
	expected := "0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272"+
			"a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f"
	result := RepeatingXor([]byte(input), []byte("ICE"))
	actual := hex.EncodeToString(result)
	if actual != expected {
		t.Error("When encrypting", input, "result was", actual, "but expected", expected)
	}
}

func TestRepeatingXorTwiceNetsInput(t *testing.T) {
	input := "Hi there, \nI'm a sailor!\nA pirate!"
	cypher := []byte("ICE ICE BABY")
	pass1 := RepeatingXor([]byte(input), cypher)
	result := string(RepeatingXor(pass1, cypher))
	if result != input {
		t.Error("XORing twice should return the input, expected", input, "but got", result)
	}
}

func TestHemingDistance(t *testing.T) {
	a := "this is a test"
	b := "wokka wokka!!!"
	distance, err := HemingDistance(([]byte)(a), ([]byte)(b))
	require.Nil(t, err, "There should be no error, but there was: {}", err)
	require.Equal(t, 37, distance, "Distance should be 37")
}

func TestDecryptAES128ECB(t *testing.T) {
	key   := "YELLOW SUBMARINE"
	input := "CRIwqt4+szDbqkNY+I0qbDe3LQz0wiw0SuxBQtAM5TDdMbjCMD/venUDW9BL"
	bytes, err := base64.StdEncoding.DecodeString(input)
	require.Nil(t, err)
	result, err := DecryptAES128ECB(bytes, []byte(key))
	require.Nil(t, err)
	actual := string(result)
	require.Contains(t, actual, "I'm back and I'm ringin' the bel")
}

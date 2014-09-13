package main

import (
	"flag"
	"github.com/wouterd/matasano-go/matasano"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
)

type Candidate struct {
	englishness float64
	phrase      string
	cypher      byte
}

func main() {
	flag.Parse()
	encrypted := flag.Arg(0)
	bytes, err := hex.DecodeString(encrypted)
	if err != nil {
		fmt.Println("Input message is not hex encoded")
		return
	}
	var current Candidate
	for i := 0 ; i < 256 ; i++ {
		cypher := byte(i)
		decoded := matasano.FixedXorWithSingleByteMask(bytes, cypher)
		phrase := string(decoded)
		englishness := stdDevFromCharFrequencies(phrase)
		if current.englishness < englishness {
			current = Candidate{englishness, phrase, cypher}
		}
	}
	hexOfCypherByte := "0x" + hex.EncodeToString([]byte{current.cypher})
	fmt.Println("Most english phrase: "+current.phrase+"\n, with cypher", hexOfCypherByte)
}

func stdDevFromCharFrequencies(phrase string) float64 {
	englishCharFreqs := map[string]float64{
		"a": 0.08167, "b": 0.01492, "c": 0.02782, "d": 0.04253, "e": 0.12702, "f": 0.02228, "g": 0.02015, "h": 0.06094,
		"i": 0.06966, "j": 0.00153, "k": 0.00772, "l": 0.04025, "m": 0.02406, "n": 0.06749, "o": 0.07507, "p": 0.01929,
		"q": 0.00095, "r": 0.05987, "s": 0.06327, "t": 0.09056, "u": 0.02758, "v": 0.00978, "w": 0.02360, "x": 0.00150,
		"y": 0.01974, "z": 0.00074, " ": 0.15000,
	}
	lowerCased := strings.ToLower(phrase)
	s2 := float64(0)
	total := float64(0)
	for key, _ := range englishCharFreqs {
		total += float64(strings.Count(lowerCased, key))
	}
	for key, value := range englishCharFreqs {
		freq := float64(strings.Count(lowerCased, key)) / float64(total)
		s2 += math.Pow(freq-value, 2)
	}
	return 1 / math.Sqrt(s2)
}

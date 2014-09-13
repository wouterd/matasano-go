package main

import (
	"flag"
	"fmt"
	"math"
	"strings"
	"io/ioutil"
	"github.com/wouterd/matasano-go/matasano"
	"encoding/hex"
)

type Candidate struct {
	englishness float64
	phrase      string
	cypher      byte
	line        int
}

func main() {
	filenamep := flag.String("file", "", "The file to read the encrypted lines from (hex encoded byte buffers)")
	flag.Parse()
	encrypted, err := getEncryptedText(*filenamep)
	if err != nil {
		return
	}
	lines := strings.Split(encrypted, "\n")

	candidates := make(chan Candidate, 20)
	best := make(chan Candidate)

	go resolveBestCandidate(candidates, len(lines), best)
	fmt.Println("MAIN: Started go routine to resolve best candidate..")

	for lineIdx, line := range lines {
		bytes, err := hex.DecodeString(line)
		lineNr := lineIdx + 1
		if err != nil {
			fmt.Println("Line", lineNr, "was not hex encoded.")
			return
		}
		go findCandidate(candidates, bytes, lineNr)
		fmt.Println("MAIN: Started go routine 'findCandidate' to resolve line", lineNr, "..")
	}

	result := <-best
	hexOfCypherByte := "0x" + hex.EncodeToString([]byte{result.cypher})
	fmt.Println("MAIN: Most english phrase: "+result.phrase+"\n, with cypher", hexOfCypherByte, "at line", result.line)
}

func resolveBestCandidate(in<- chan Candidate, amCandidates int, out chan <- Candidate) {
	candidates := 0
	var current Candidate
	for candidates < amCandidates {
		candidate := <-in
		fmt.Println("resolveBestCandidate: Received a candidate..")
		candidates++
		if candidate.englishness > current.englishness {
			current = candidate
		}
	}
	out <- current
}

func findCandidate(c chan <- Candidate, bytes []byte, line int) {
	var best Candidate
	for i := 0 ; i < 256 ; i++ {
		cypher := byte(i)
		decoded := matasano.FixedXorWithSingleByteMask(bytes, cypher)
		phrase := string(decoded)
		englishness := stdDevFromCharFrequencies(phrase)
		if best.englishness < englishness {
			best = Candidate{englishness, phrase, cypher, line}
		}
	}
	c <- best
}

func getEncryptedText(filename string) (string, error) {
	if filename == "" {
		return flag.Arg(0), nil
	} else {
		contents, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println("Error reading contents from input file", err)
			return "", err
		}
		return string(contents), nil
	}
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

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"encoding/base64"
	"github.com/wouterd/matasano-go/matasano"
	"sort"
	"strings"
)

func main() {
	flag.Parse()
	filename := flag.Arg(0)
	if filename == "" {
		fmt.Println("Please specify the file to read on the command line.")
		return
	}
	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Problem reading file", filename, ":", err)
		return
	}
	decodeBuffer := make([]byte, base64.StdEncoding.DecodedLen(len(fileContents)))
	nWritten, err := base64.StdEncoding.Decode(decodeBuffer, fileContents)
	if err != nil {
		fmt.Println("Unable to decode base64 data:", err)
	}
	encryptedBuffer := decodeBuffer[:nWritten]

	keySize := determineBestKeysize(encryptedBuffer)

	cypher := resolveBestRepeatingXorKey(encryptedBuffer, keySize)
	fmt.Println("For keysize", keySize, "key is '"+string(cypher)+"'")

	decrypted := matasano.RepeatingXor(encryptedBuffer, cypher)
	fmt.Println("Message:\n", string(decrypted))
}

func resolveBestRepeatingXorKey(data []byte, keysize uint) []byte {
	blockLen := uint(len(data)) / keysize
	blocks := make([][]byte, keysize)
	for block := uint(0) ; block < keysize ; block++ {
		blocks[block] = make([]byte, blockLen)
		for i := uint(0) ; i < blockLen ; i++ {
			blocks[block][i] = data[i*keysize+block]
		}
	}
	key := make([]byte, keysize)
	for i := uint(0) ; i < keysize ; i++ {
		key[i] = findBestCypher(blocks[i])
	}
	return key
}

//--- Struct that holds key size statistics and can be sorted by normalized distance
type KeysizeStatistics struct {
	keysize      uint
	normDistance float64
}

type ByNormDistance []KeysizeStatistics

func (this ByNormDistance) Len() int {
	return len(this)
}

func (this ByNormDistance) Less(i, j int) bool {
	return this[i].normDistance < this[j].normDistance
}

func (this ByNormDistance) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

//---------------

func calcNormalizedHemingDistance(a, b []byte) float64 {
	distance, err := matasano.HemingDistance(a, b)
	if err != nil {
		panic(err)
	}
	return float64(distance) / float64(len(a))
}

func determineBestKeysize(buffer []byte) uint {
	minKeySize := uint(2)
	maxKeySize := uint(40)
	keySizeStats := make([]KeysizeStatistics, maxKeySize-minKeySize+1)

	for keysize := minKeySize ; keysize <= maxKeySize ; keysize++ {
		normDistance := float64(0)
		testChunks := uint(len(buffer)) / 2 / keysize
		for i := uint(0) ; i < testChunks ; i++ {
			offset := i * keysize * 2
			slice1 := buffer[offset:offset+keysize]
			slice2 := buffer[offset+keysize:offset+keysize*2]
			normDistance += calcNormalizedHemingDistance(slice1, slice2)/float64(testChunks)
			keySizeStats[keysize-minKeySize] = KeysizeStatistics{keysize, normDistance}
		}
	}

	sort.Sort(ByNormDistance(keySizeStats))
	return keySizeStats[0].keysize
}

type Candidate struct {
	englishness float64
	phrase      string
	cypher      byte
	line        int
}

func findBestCypher(bytes []byte) byte {
	var best Candidate
	for i := 0 ; i < 256 ; i++ {
		cypher := byte(i)
		decoded := matasano.FixedXorWithSingleByteMask(bytes, cypher)
		phrase := string(decoded)
		englishness := getLetterDensity(phrase)
		if best.englishness < englishness {
			best = Candidate{englishness, phrase, cypher, 0}
		}
	}
	return best.cypher
}

func getLetterDensity(phrase string) float64 {
	letters := "0123456789abcdefghijklmnopqrstuvwxyz !,.?"
	lowerCase := strings.ToLower(phrase)
	totalChars := 0
	for _, char := range lowerCase {
		if strings.ContainsRune(letters, char) {
			totalChars++
		}
	}
	return float64(totalChars) / float64(len(phrase))
}

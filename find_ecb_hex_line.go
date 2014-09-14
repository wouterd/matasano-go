package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"github.com/wouterd/matasano-go/matasano"
	"encoding/hex"
)

type Candidate struct {
	matches  int
	line     int
	contents string
}

func main() {
	blockSize := flag.Int("blocksize", 0, "The block size to use for scanning")
	flag.Parse()
	filename := flag.Arg(0)
	if filename == "" {
		fmt.Println("Please specify a file to search on the command line")
		return
	}
	if *blockSize == 0 {
		fmt.Println("Please specify a block size on the command line")
		return
	}
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	lines := strings.Split((string)(contents), "\n")

	var best Candidate
	for i, line := range lines {
		bytes, err := hex.DecodeString(line)
		if err != nil {
			panic(err)
		}
		matches := matasano.CountMatchingBlocks(bytes, *blockSize)
		if best.matches < matches {
			best = Candidate{matches, i, line}
		}
	}

	fmt.Println("Most likely candidate for ECB encoding:")
	fmt.Println("line:", best.line, "matching blocks:", best.matches)
	fmt.Println("Contents:", best.contents)
}


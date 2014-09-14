package main

import (
	"fmt"
	"encoding/base64"
	"io/ioutil"
	"flag"
	"github.com/wouterd/matasano-go/matasano"
)

func main() {
	filenamep := flag.String("file", "", "Specify the file to decrypt")
	passwordp := flag.String("pwd", "", "Specify the password")
	flag.Parse()
	if *filenamep == "" {
		fmt.Println("Please specify the file to read on the command line.")
		return
	}
	if *passwordp == "" {
		fmt.Println("Please specify the password to use for decryption")
		return
	}
	fileContents, err := ioutil.ReadFile(*filenamep)
	if err != nil {
		fmt.Println("Problem reading file", *filenamep, ":", err)
		return
	}
	decodeBuffer := make([]byte, base64.StdEncoding.DecodedLen(len(fileContents)))
	nWritten, err := base64.StdEncoding.Decode(decodeBuffer, fileContents)
	if err != nil {
		fmt.Println("Unable to decode base64 data:", err)
	}
	encryptedBuffer := decodeBuffer[:nWritten]

	decrypted, err := matasano.DecryptAES128ECB(encryptedBuffer, []byte(*passwordp))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(decrypted))
}


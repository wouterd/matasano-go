package matasano

import "encoding/hex"
import "encoding/base64"

func HexToBase64(input string) (output string, err error) {
  bytes, decodeError := hex.DecodeString(input)
  if decodeError != nil {
    return "", decodeError
  }
  return base64.StdEncoding.EncodeToString(bytes), nil
}

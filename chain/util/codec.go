package util

import (
	"encoding/hex"
	"strings"
)

func HexDecodeString(s string) ([]byte, error) {
	s = strings.TrimPrefix(s, "0x")

	if len(s)%2 != 0 {
		s = "0" + s
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func HexEncodeToString(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

func TextToHex(s string) string {
	return "0x" + hex.EncodeToString([]byte(s))
}

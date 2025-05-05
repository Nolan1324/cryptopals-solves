package util

import (
	"bufio"
	"cryptopals/internal/enc"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
)

func IsAlpha(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

func IsNumeric(c byte) bool {
	return '0' <= c && c <= '9'
}

func IsAlphanumeric(c byte) bool {
	return IsAlpha(c) || IsNumeric(c)
}

func IsPunctuation(c byte) bool {
	return strings.Contains(" !@#$%^&*()-=_+[]{}\\|;':\",./<>?", string(rune(c)))
}

func IsSentenceAscii(c byte) bool {
	return IsAlphanumeric(c) || strings.Contains(" \"',.\n\r", string(rune(c)))
}

func IsVisibleAscii(c byte) bool {
	return (32 <= c && c <= 126) || c == '\r' || c == '\n'
}

func ReadBase64File(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("error opening file: %v\n", err)
	}
	defer file.Close()

	var buf []byte

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		decoded, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			log.Fatalf("error \n")
		}
		buf = append(buf, decoded...)
	}

	return buf
}

func ReadBase64ListFile(filename string) ([][]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("when opening file: %w", err)
	}
	defer file.Close()

	var strings [][]byte

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		decoded, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			return nil, fmt.Errorf("when decoding file: %w", err)
		}
		strings = append(strings, decoded)
	}

	return strings, nil
}

func ReadHexListFile(filename string) [][]byte {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	var strings [][]byte

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		bytes := enc.HexDecode([]byte(line))
		strings = append(strings, bytes)
	}

	return strings
}

package pkg

import (
	"bytes"
	"encoding/base64"
	"filippo.io/age"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

type Decoder struct {
	identity *age.X25519Identity
}

func CreateDecoder(privateKeyFilePath string) (Decoder, error) {
	privateKey, err := ioutil.ReadFile(privateKeyFilePath)
	if err != nil {
		return Decoder{}, fmt.Errorf("failed to read file '%s'. Caused by : '%v'", privateKeyFilePath, err)
	}

	trimmed := strings.TrimSuffix(string(privateKey), "\n")
	decoded, err := decode(trimmed)
	if err != nil {
		return Decoder{}, err
	}
	identity, err := age.ParseX25519Identity(decoded)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	return Decoder{identity: identity}, err
}

func (e Decoder) Decrypt(data string) ([]byte, error) {
	r, err := age.Decrypt(strings.NewReader(data), e.identity)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data")
	}
	out := &bytes.Buffer{}
	if _, err := io.Copy(out, r); err != nil {
		return nil, fmt.Errorf("failed to decrypt")
	}
	return out.Bytes(), nil
}

func decode(input string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", fmt.Errorf("failed to decode identity file content")
	}
	return string(decoded), nil
}

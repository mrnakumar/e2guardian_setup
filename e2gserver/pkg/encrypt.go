package pkg

import (
	"bytes"
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
	identity, err := age.ParseX25519Identity(strings.TrimSuffix(string(privateKey), "\n"))
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

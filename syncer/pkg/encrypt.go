package pkg

import (
	"bytes"
	"filippo.io/age"
	"fmt"
	"github.com/mrnakumar/e2g_utils"
	"io/ioutil"
	"strings"
)

type Encryptor struct {
	recipient *age.X25519Recipient
}

func CreateEncryptor(recipientKeyPath string) (Encryptor, error) {
	recipientKey, err := ioutil.ReadFile(recipientKeyPath)
	if err != nil {
		return Encryptor{}, fmt.Errorf("failed to read file '%s'. Caused by : '%v'", recipientKeyPath, err)
	}
	trimmed := strings.TrimSuffix(string(recipientKey), "\n")
	decoded, err := e2g_utils.Base64Decode(trimmed)
	if err != nil {
		return Encryptor{}, fmt.Errorf("failed to decode recepient key path. Caused by: '%v'", err)
	}
	publicKey, err := age.ParseX25519Recipient(decoded)
	return Encryptor{recipient: publicKey}, err
}

func (e Encryptor) Encrypt(data []byte) ([]byte, error) {
	out := &bytes.Buffer{}
	w, err := age.Encrypt(out, e.recipient)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	err = w.Close()
	return out.Bytes(), err
}

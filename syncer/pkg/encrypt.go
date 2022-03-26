package pkg

import (
	"bytes"
	"filippo.io/age"
)

type Encryptor struct {
	recipient *age.X25519Recipient
}

func CreateEncryptor(recipient string) (Encryptor, error) {
	publicKey, err := age.ParseX25519Recipient(recipient)
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

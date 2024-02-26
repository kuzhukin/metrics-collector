package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type Decryptor struct {
	key *rsa.PrivateKey
}

func NewDecryptor(keyPath string) (*Decryptor, error) {
	key, err := newPrivateKeyFromFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("new private key from file, err=%w", err)
	}

	return &Decryptor{key: key}, nil
}

func (d *Decryptor) Decrypt(data []byte) ([]byte, error) {
	msgLen := len(data)
	step := d.key.PublicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(sha512.New(), rand.Reader, d.key, data[start:finish], nil)
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}

type Encryptor struct {
	key *rsa.PublicKey
}

func NewEncryptor(keyPath string) (*Encryptor, error) {
	key, err := newPublicKeyFromFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("new public key from file, err=%w", err)
	}

	return &Encryptor{key: key}, nil
}

func (e *Encryptor) Encrypt(data []byte) ([]byte, error) {
	msgLen := len(data)
	step := e.key.Size() - 2*sha512.Size - 2
	var encryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlockBytes, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, e.key, data[start:finish], nil)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}

func newPrivateKeyFromFile(filename string) (*rsa.PrivateKey, error) {
	keyBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	blockBytes := block.Bytes

	privateKey, err := x509.ParsePKCS1PrivateKey(blockBytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func newPublicKeyFromFile(filename string) (*rsa.PublicKey, error) {
	keyBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	blockBytes := block.Bytes

	publicKey, err := x509.ParsePKCS1PublicKey(blockBytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

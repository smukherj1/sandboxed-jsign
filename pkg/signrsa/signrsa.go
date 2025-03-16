package signrsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type Signer struct {
	key *rsa.PrivateKey
}

func (s *Signer) Sign(payload []byte) ([]byte, error) {
	sig, err := rsa.SignPKCS1v15(rand.Reader, s.key, crypto.SHA256, payload)
	if err != nil {
		return nil, fmt.Errorf("error signing: %w", err)
	}
	return sig, nil
}

func NewSigner(keyFile string) (*Signer, error) {
	keyBytes, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read RSA key file %v: %w", keyFile, err)
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode key file %v, got PEM Block type %v, want PRIVATE KEY", keyFile, block.Type)
	}

	k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCSI private key from PEM block in key file %v: %w", keyFile, err)
	}
	rk, ok := k.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("did not find a RSA private key in key file %v: %w", keyFile, err)
	}
	return &Signer{key: rk}, nil
}

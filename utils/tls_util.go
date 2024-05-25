package utils

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

func LoadCerts(certFile string) ([][]byte, error) {
	if certFile == "" {
		return nil, nil
	}

	certPEM, err := os.ReadFile(certFile)
	if err != nil {
		return nil, err
	}

	certs, err := parseCertsPEM(certPEM)
	if err != nil {
		return nil, err
	}

	return certs, nil
}

func LoadPrivateKey(keyFile string) (any, error) {
	if keyFile == "" {
		return nil, nil
	}

	keyPEM, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	key, err := parsePrivateKeyPEM(keyPEM)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func parseCertsPEM(certsPEM []byte) ([][]byte, error) {
	var certs [][]byte
	var cert []byte
	for len(certsPEM) > 0 {
		var block *pem.Block
		block, certsPEM = pem.Decode(certsPEM)
		if block == nil {
			break
		}
		cert = append(cert, block.Bytes...)
		if len(certsPEM) == 0 {
			certs = append(certs, cert)
		}
	}
	return certs, nil
}

func parsePrivateKeyPEM(keyPEM []byte) (any, error) {
	block, _ := pem.Decode(keyPEM)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(block.Bytes)
	default:
		return nil, errors.New("unsupported private key type")
	}
}

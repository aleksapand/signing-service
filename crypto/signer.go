package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
	VerifySignature(data []byte, base64Signature []byte) bool
	GetPublicKey() crypto.PublicKey
	GetAlgorithm() string
}

// RSASigner stores RSA Keys and handles data signing and verifying.
type RSASigner struct {
	keyPair *RSAKeyPair
}

// ECCSigner stores ECDS Keys and handles data signing and verifying.
type ECCSigner struct {
	keyPair *ECCKeyPair
}

// NewRSASigner is a factory to instantiate a new RSASigner.
func NewRSASigner() *RSASigner {
	keyGenerator := RSAGenerator{}
	keyPair, _ := keyGenerator.Generate()
	return &RSASigner{keyPair}
}

// Sign of RSASigner signs the data with RSA algorithm.
func (signer *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hashed := sha256.Sum256(dataToBeSigned)
	signature, err := rsa.SignPKCS1v15(rand.Reader, signer.keyPair.Private, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign data: %v", err)
	}
	base64Bytes := make([]byte, base64.StdEncoding.EncodedLen(len(signature)))
	base64.StdEncoding.Encode(base64Bytes, signature)
	return base64Bytes, nil
}

// VerifySignature of RSASigner verifies the signature of given data with RSA algorithm.
func (signer *RSASigner) VerifySignature(data []byte, base64Signature []byte) bool {
	signatureBytes, err := base64.StdEncoding.DecodeString(string(base64Signature))
	if err != nil {
		fmt.Println("Error decoding signature:", err)
		return false
	}

	hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(signer.keyPair.Public, crypto.SHA256, hashed[:], signatureBytes)
	return err == nil
}

// GetPublicKey of RSASigner returns RSA Public Key.
func (signer *RSASigner) GetPublicKey() crypto.PublicKey {
	return signer.keyPair.Public
}

func (signer RSASigner) GetAlgorithm() string {
	return "RSA"
}

// NewECCSigner is a factory to instantiate a new ECCSigner.
func NewECCSigner() *ECCSigner {
	keyGenerator := ECCGenerator{}
	keyPair, _ := keyGenerator.Generate()
	return &ECCSigner{keyPair}
}

// Sign of ECCSigner signs the data with ECDS algorithm.
func (signer *ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hashed := sha256.Sum256(dataToBeSigned)
	r, s, err := ecdsa.Sign(rand.Reader, signer.keyPair.Private, hashed[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign data: %v", err)
	}

	signature := append(r.Bytes(), s.Bytes()...)
	base64Bytes := make([]byte, base64.StdEncoding.EncodedLen(len(signature)))
	base64.StdEncoding.Encode(base64Bytes, signature)

	return base64Bytes, nil
}

// VerifySignature of ECCSigner verifies the signature of given data with ECDS algorithm.
func (signer *ECCSigner) VerifySignature(data []byte, base64Signature []byte) bool {
	decodedSignature := make([]byte, base64.StdEncoding.DecodedLen(len(base64Signature)))
	n, err := base64.StdEncoding.Decode(decodedSignature, base64Signature)
	if err != nil {
		fmt.Println("Error decoding signature:", err)
		return false
	}

	decodedSignature = decodedSignature[:n]
	rBytes := decodedSignature[:len(decodedSignature)/2]
	sBytes := decodedSignature[len(decodedSignature)/2:]
	var r, s big.Int
	r.SetBytes(rBytes)
	s.SetBytes(sBytes)
	hashed := sha256.Sum256(data)

	return ecdsa.Verify(signer.keyPair.Public, hashed[:], &r, &s)
}

// GetPublicKey of ECCSigner returns ECDS Public Key
func (signer *ECCSigner) GetPublicKey() crypto.PublicKey {
	return signer.keyPair.Public
}

func (signer ECCSigner) GetAlgorithm() string {
	return "ECC"
}

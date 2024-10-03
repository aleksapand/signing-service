package domain

import (
	"bytes"
	"encoding/base64"
	"strings"
	"sync"
	"testing"

	crypto2 "github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

func TestDevice_CreateDevice(t *testing.T) {
	algorithm := "RSA"
	signer, _ := crypto2.SignerFactory(algorithm)
	signatureDevice := NewSignatureDevice("", signer)

	if signatureDevice.SignatureCounter != 0 {
		t.Error("Device not properly initialized")
	}

	base64EncodedId := make([]byte, base64.StdEncoding.EncodedLen(len(signatureDevice.Id[:])))
	base64.StdEncoding.Encode(base64EncodedId, signatureDevice.Id[:])

	if !bytes.Equal(signatureDevice.LastSig, base64EncodedId) {
		t.Error("Device not properly initialized")
	}
}

func TestDevice_CreateAndSign(t *testing.T) {
	algorithm := "RSA"
	signer, _ := crypto2.SignerFactory(algorithm)
	signatureDevice := NewSignatureDevice("", signer)

	data, signature, _ := signatureDevice.SignData([]byte("Hello World!"))

	if !signatureDevice.Signer.VerifySignature(data, signature) {
		t.Error("Signature verification failed")
	}

	if signatureDevice.SignatureCounter != 1 {
		t.Error("Device state not properly updated")
	}

	if !bytes.Equal(signatureDevice.LastSig, signature) {
		t.Error("Device state not properly updated")
	}
}

// go test -race
func TestDevice_ConcurrentSignatures(t *testing.T) {
	algorithm := "RSA"
	signer, _ := crypto2.SignerFactory(algorithm)
	signatureDevice := NewSignatureDevice("", signer)

	var wg sync.WaitGroup
	goroutinesNum := 10

	wg.Add(goroutinesNum)
	for i := 0; i < goroutinesNum; i++ {
		go func() {
			defer wg.Done()
			signatureDevice.SignData([]byte("Hello World!"))
		}()
	}
}

func TestDevice_PrepareData(t *testing.T) {
	data := prepareData(123, []byte("Hello World!"), []byte("42"))
	referenceData := []byte(strings.Join([]string{"123_Hello World!", base64.StdEncoding.EncodeToString([]byte("42"))}, "_"))
	if !bytes.Equal(data, referenceData) {
		t.Error("Data preparation failed")
	}
}

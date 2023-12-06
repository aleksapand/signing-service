package api

import (
	"bytes"
	"encoding/base64"
	crypto2 "github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"strings"
	"testing"
)

func TestDevice_CreateDevice(t *testing.T) {
	algorithm := "RSA"
	signer, _ := crypto2.SignerFactory(algorithm)
	signatureDevice := domain.NewSignatureDevice("", algorithm, signer)

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
	signatureDevice := domain.NewSignatureDevice("", algorithm, signer)

	data := PrepareData(signatureDevice.SignatureCounter, []byte("Hello World!"), signatureDevice.LastSig)
	signature, _ := signatureDevice.Signer.Sign(data)
	signatureDevice.SetLastSignature(signature)

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

func TestDevice_PrepareData(t *testing.T) {
	data := PrepareData(123, []byte("Hello World!"), []byte("42"))
	referenceData := []byte(strings.Join([]string{"123_Hello World!", base64.StdEncoding.EncodeToString([]byte("42"))}, "_"))
	if !bytes.Equal(data, referenceData) {
		t.Error("Data preparation failed")
	}
}

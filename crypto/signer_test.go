package crypto

import (
	"testing"
)

func TestRSASigner_Sign(t *testing.T) {
	signer := NewRSASigner()
	dataToBeSigned := []byte("Hello, World!")
	signature, _ := signer.Sign(dataToBeSigned)

	if !signer.VerifySignature(dataToBeSigned, signature) {
		t.Error("Signature verification failed")
	}
}

func TestRSASigner_SignEmpty(t *testing.T) {
	signer := NewRSASigner()
	dataToBeSigned := []byte("")
	signature, _ := signer.Sign(dataToBeSigned)

	if !signer.VerifySignature(dataToBeSigned, signature) {
		t.Error("Signature verification failed")
	}
}

func TestRSASigner_SignVerifyWrongData(t *testing.T) {
	signer := NewRSASigner()
	dataToBeSigned := []byte("Hello, World!")
	signature, _ := signer.Sign(dataToBeSigned)

	dataToBeSigned = []byte("Hello, World!!!")
	if signer.VerifySignature(dataToBeSigned, signature) {
		t.Error("Signature verification should not succeed")
	}
}

func TestECCSigner_Sign(t *testing.T) {
	signer := NewECCSigner()
	dataToBeSigned := []byte("Hello, World!")
	signature, _ := signer.Sign(dataToBeSigned)
	if !signer.VerifySignature(dataToBeSigned, signature) {
		t.Error("Signature verification failed")
	}
}

func TestECCSigner_SignEmpty(t *testing.T) {
	signer := NewECCSigner()
	dataToBeSigned := []byte("")
	signature, _ := signer.Sign(dataToBeSigned)
	if !signer.VerifySignature(dataToBeSigned, signature) {
		t.Error("Signature verification failed")
	}
}

func TestSignerFactory_RSASigner(t *testing.T) {
	signer, _ := SignerFactory("rsa")

	if _, isRSASigner := signer.(*RSASigner); !isRSASigner {
		t.Error("Signer Factory failed.")
	}
}

func TestSignerFactory_ECCSigner(t *testing.T) {
	signer, _ := SignerFactory("ECC")

	if _, isECCSigner := signer.(*ECCSigner); !isECCSigner {
		t.Error("Signer Factory failed.")
	}
}

func TestSignerFactory_XYZSigner(t *testing.T) {
	signer, _ := SignerFactory("XYZ")

	if signer != nil {
		t.Error("Signer Factory failed.")
	}
}

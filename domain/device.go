package domain

import (
	"encoding/base64"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/google/uuid"
)

// SignatureDevice is struct holding all signature device data.
type SignatureDevice struct {
	Id               uuid.UUID
	Label            string
	Algorithm        string
	SignatureCounter uint64
	Signer           crypto.Signer
	LastSig          []byte
}

// NewSignatureDevice is factory that initializes signature device.
func NewSignatureDevice(label string, algorithm string, signer crypto.Signer) *SignatureDevice {
	id := uuid.New()
	base64EncodedId := make([]byte, base64.StdEncoding.EncodedLen(len(id[:])))
	base64.StdEncoding.Encode(base64EncodedId, id[:])
	if label == "" {
		label = id.String()
	}
	signatureDevice := SignatureDevice{id, label, algorithm, 0, signer, base64EncodedId}
	return &signatureDevice
}

// SetLastSignature updates last signature and increments signature counter
func (device *SignatureDevice) SetLastSignature(lastSig []byte) {
	device.SignatureCounter += 1
	device.LastSig = lastSig
}

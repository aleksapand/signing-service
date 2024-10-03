package domain

import (
	"encoding/base64"
	"strconv"
	"strings"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/google/uuid"
)

// SignatureDevice is struct holding all signature device data.
type SignatureDevice struct {
	Id               uuid.UUID
	Label            string
	SignatureCounter uint64
	Signer           crypto.Signer
	LastSig          []byte
	sigMutex         sync.RWMutex
}

// NewSignatureDevice is factory that initializes signature device.
func NewSignatureDevice(label string, signer crypto.Signer) *SignatureDevice {
	id := uuid.New()
	base64EncodedId := make([]byte, base64.StdEncoding.EncodedLen(len(id[:])))
	base64.StdEncoding.Encode(base64EncodedId, id[:])
	if label == "" {
		label = id.String()
	}
	signatureDevice := SignatureDevice{id, label, 0, signer, base64EncodedId, sync.RWMutex{}}
	return &signatureDevice
}

func (device *SignatureDevice) SignData(rawData []byte) ([]byte, []byte, error) {
	device.sigMutex.Lock()
	defer device.sigMutex.Unlock()
	data := prepareData(device.SignatureCounter, rawData, device.LastSig)
	signature, err := device.Signer.Sign(data)
	device.setLastSignature(signature)

	return data, signature, err
}

// SetLastSignature updates last signature and increments signature counter
func (device *SignatureDevice) setLastSignature(lastSig []byte) {
	device.SignatureCounter += 1
	device.LastSig = lastSig
}

// PrepareData appends and prepends id and last signature to the data from sign request.
func prepareData(counter uint64, data []byte, lastSig []byte) []byte {

	formattedData := strings.Join([]string{strconv.FormatUint(counter, 10), string(data),
		base64.StdEncoding.EncodeToString(lastSig)}, "_")
	return []byte(formattedData)

}

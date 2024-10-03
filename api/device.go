package api

import (
	"crypto"
	"encoding/json"
	"net/http"
	"strings"

	crypto2 "github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
)

// SignatureDeviceResponse is response for newly created signature device.
type SignatureDeviceResponse struct {
	Id        uuid.UUID        `json:"id"`
	Label     string           `json:"label"`
	Algorithm string           `json:"algorithm"`
	PublicKey crypto.PublicKey `json:"publicKey"`
}

// SignatureDeviceRequest is a request with data needed for signature device creation. Label is optional.
type SignatureDeviceRequest struct {
	Algorithm string `json:"algorithm"`
	Label     string `json:"label"`
}

// SignDataRequest is a request for data signing.
type SignDataRequest struct {
	Id   uuid.UUID `json:"id"`
	Data string    `json:"data"`
}

// SignDataResponse holds a signed data.
type SignDataResponse struct {
	Signature  []byte `json:"signature"`
	SignedData string `json:"signed_data"`
}

// CreateSignatureDevice handles a request for new signature device creation.
// It parses a request which holds information about algorithm and optional label for the device.
func (s *Server) CreateSignatureDevice(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	var requestData SignatureDeviceRequest
	err := json.NewDecoder(request.Body).Decode(&requestData)
	if err != nil {
		http.Error(response, "Failed to decode Request", http.StatusBadRequest)
		return
	}
	signer, err := crypto2.SignerFactory(requestData.Algorithm)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	newSignatureDevice := domain.NewSignatureDevice(requestData.Label, signer)
	s.db.Set(newSignatureDevice.Id, newSignatureDevice)

	newDeviceResponse := SignatureDeviceResponse{
		Id:        newSignatureDevice.Id,
		Label:     newSignatureDevice.Label,
		Algorithm: newSignatureDevice.Signer.GetAlgorithm(),
		PublicKey: newSignatureDevice.Signer.GetPublicKey(),
	}
	WriteAPIResponse(response, http.StatusOK, newDeviceResponse)
}

// SignData handles request for signing the data. It parses request for the data and id of a signature device.
func (s *Server) SignData(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	var requestData SignDataRequest
	err := json.NewDecoder(request.Body).Decode(&requestData)
	if err != nil {
		http.Error(response, "Failed to decode Request", http.StatusBadRequest)
		return
	}

	signatureDevice, exists := s.db.Get(requestData.Id)
	if !exists {
		WriteErrorResponse(response, http.StatusNotFound, []string{
			"No device found under provided id.",
		})
	}

	data, signature, err := signatureDevice.SignData([]byte(requestData.Data))

	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{
			"Signing data failed",
		})
		return
	}

	signedDataResponse := SignDataResponse{
		signature,
		string(data),
	}

	WriteAPIResponse(response, http.StatusOK, signedDataResponse)
}

// GetDevices handles get request for all created devices.
func (s *Server) GetDevices(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	allDevices := s.db.GetAll()

	var allDevicesResponse []SignatureDeviceResponse
	for _, device := range allDevices {
		allDevicesResponse = append(allDevicesResponse, SignatureDeviceResponse{
			Id:        device.Id,
			Label:     device.Label,
			Algorithm: device.Signer.GetAlgorithm(),
			PublicKey: device.Signer.GetPublicKey(),
		})
	}

	WriteAPIResponse(response, http.StatusOK, allDevicesResponse)
}

// GetDevice handles a request for one specific device.
func (s *Server) GetDevice(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	parts := strings.Split(request.URL.Path, "/")
	deviceIDStr := parts[len(parts)-1]
	id, err := uuid.Parse(deviceIDStr)
	if err != nil {
		http.Error(response, "Invalid device ID", http.StatusBadRequest)
		return
	}

	device, exists := s.db.Get(id)
	if !exists {
		WriteErrorResponse(response, http.StatusNotFound, []string{
			"No device found under provided id.",
		})
	}

	deviceResponse := SignatureDeviceResponse{
		Id:        device.Id,
		Label:     device.Label,
		Algorithm: device.Signer.GetAlgorithm(),
		PublicKey: device.Signer.GetPublicKey(),
	}

	WriteAPIResponse(response, http.StatusOK, deviceResponse)
}

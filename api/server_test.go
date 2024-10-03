package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

func initServer(db *persistence.InMemoryDB) *httptest.Server {
	server := NewServer(":8080", db)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/new", server.CreateSignatureDevice)
	mux.HandleFunc("/api/v0/devices", server.GetDevices)
	mux.HandleFunc("/api/v0/devices/", server.GetDevice)
	mux.HandleFunc("/api/v0/sign", server.SignData)
	return httptest.NewServer(mux)
}

func sendPostRequest(path string, body []byte) (*http.Response, error) {
	req, _ := http.NewRequest("POST", path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	return res, err
}

func sendGetRequest(path string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", path, nil)

	client := &http.Client{}
	res, err := client.Do(req)
	return res, err
}

func TestServer_CreateDevice(t *testing.T) {
	db := persistence.GetInMemoryDB()
	ts := initServer(db)
	defer ts.Close()

	requestData := SignatureDeviceRequest{
		"RSA",
		"Device1",
	}
	body, _ := json.Marshal(requestData)
	res, err := sendPostRequest(ts.URL+"/api/v0/new", body)
	if err != nil {
		t.Fatal(err)
	}

	var response Response
	_ = json.NewDecoder(res.Body).Decode(&response)

	var deviceResponse SignatureDeviceResponse
	dataBytes, _ := json.Marshal(response.Data)
	_ = json.Unmarshal(dataBytes, &deviceResponse)

	if deviceResponse.Label != "Device1" {
		t.Error("Device not properly created")
	}

	if deviceResponse.Algorithm != "RSA" {
		t.Error("Device not properly created")
	}

}

func TestServer_CreateAndGetDevices(t *testing.T) {
	db := persistence.GetInMemoryDB()
	ts := initServer(db)
	defer ts.Close()

	requestData1 := SignatureDeviceRequest{
		"RSA",
		"Device1",
	}

	requestData2 := SignatureDeviceRequest{
		"ECC",
		"Device2",
	}
	body, _ := json.Marshal(requestData1)
	res, err := sendPostRequest(ts.URL+"/api/v0/new", body)
	if err != nil {
		t.Fatal(err)
	}

	body, _ = json.Marshal(requestData2)
	res, err = sendPostRequest(ts.URL+"/api/v0/new", body)
	if err != nil {
		t.Fatal(err)
	}

	res, _ = sendGetRequest(ts.URL + "/api/v0/devices")

	var response Response
	_ = json.NewDecoder(res.Body).Decode(&response)

	var deviceResponse []SignatureDeviceResponse
	dataBytes, _ := json.Marshal(response.Data)
	_ = json.Unmarshal(dataBytes, &deviceResponse)

	if len(deviceResponse) != 2 {
		t.Error("There must be two devices.")
	}
}

func TestServer_CreateAndGetSpecificDevice(t *testing.T) {
	db := persistence.GetInMemoryDB()
	ts := initServer(db)
	defer ts.Close()

	requestData1 := SignatureDeviceRequest{
		"RSA",
		"Device1",
	}

	requestData2 := SignatureDeviceRequest{
		"ECC",
		"Device2",
	}
	body, _ := json.Marshal(requestData1)
	res, err := sendPostRequest(ts.URL+"/api/v0/new", body)
	if err != nil {
		t.Fatal(err)
	}

	body, _ = json.Marshal(requestData2)
	res, err = sendPostRequest(ts.URL+"/api/v0/new", body)
	if err != nil {
		t.Fatal(err)
	}
	var response Response
	_ = json.NewDecoder(res.Body).Decode(&response)

	var deviceResponse SignatureDeviceResponse
	dataBytes, _ := json.Marshal(response.Data)
	_ = json.Unmarshal(dataBytes, &deviceResponse)

	id := deviceResponse.Id
	res, _ = sendGetRequest(ts.URL + "/api/v0/device/" + id.String())

	var response2 Response
	_ = json.NewDecoder(res.Body).Decode(&response2.Data)
	var deviceResponse2 SignatureDeviceResponse
	_ = json.Unmarshal(dataBytes, &deviceResponse2)

	if id != deviceResponse.Id {
		t.Error("Incorrect device retrieved.")
	}
}

func TestServer_CreateDevicesAndSignatures(t *testing.T) {
	db := persistence.GetInMemoryDB()
	ts := initServer(db)
	defer ts.Close()

	requestData1 := SignatureDeviceRequest{
		"RSA",
		"Device1",
	}

	requestData2 := SignatureDeviceRequest{
		"ECC",
		"Device2",
	}
	body, _ := json.Marshal(requestData1)
	res, err := sendPostRequest(ts.URL+"/api/v0/new", body)
	if err != nil {
		t.Fatal(err)
	}

	var response Response
	_ = json.NewDecoder(res.Body).Decode(&response)
	var deviceResponse SignatureDeviceResponse
	dataBytes, _ := json.Marshal(response.Data)
	_ = json.Unmarshal(dataBytes, &deviceResponse)
	idDevice1 := deviceResponse.Id

	body, _ = json.Marshal(requestData2)
	res, err = sendPostRequest(ts.URL+"/api/v0/new", body)
	if err != nil {
		t.Fatal(err)
	}

	_ = json.NewDecoder(res.Body).Decode(&response)
	dataBytes, _ = json.Marshal(response.Data)
	_ = json.Unmarshal(dataBytes, &deviceResponse)
	idDevice2 := deviceResponse.Id

	requestSignData1 := SignDataRequest{
		idDevice1,
		"Hello World to Device1",
	}
	body, _ = json.Marshal(requestSignData1)
	res, err = sendPostRequest(ts.URL+"/api/v0/sign", body)
	if err != nil {
		t.Fatal(err)
	}

	_ = json.NewDecoder(res.Body).Decode(&response)
	var signDataResponse SignDataResponse
	dataBytes, _ = json.Marshal(response.Data)
	_ = json.Unmarshal(dataBytes, &signDataResponse)

	device1, _ := db.Get(idDevice1)
	if !device1.Signer.VerifySignature([]byte(signDataResponse.SignedData), signDataResponse.Signature) {
		t.Error("Failed to verify signature from device1.")
	}

	if !bytes.Equal(device1.LastSig, signDataResponse.Signature) {
		t.Error("Device1 state not updated properly.")
	}

	if device1.SignatureCounter != 1 {
		t.Error("Device1 state not updated properly.")
	}

	requestSignData1 = SignDataRequest{
		idDevice1,
		"Hello World to Device1 again!",
	}
	body, _ = json.Marshal(requestSignData1)
	res, err = sendPostRequest(ts.URL+"/api/v0/sign", body)
	if err != nil {
		t.Fatal(err)
	}

	_ = json.NewDecoder(res.Body).Decode(&response)
	dataBytes, _ = json.Marshal(response.Data)
	_ = json.Unmarshal(dataBytes, &signDataResponse)

	if !device1.Signer.VerifySignature([]byte(signDataResponse.SignedData), signDataResponse.Signature) {
		t.Error("Failed to verify signature from device1.")
	}

	if !bytes.Equal(device1.LastSig, signDataResponse.Signature) {
		t.Error("Device1 state not updated properly.")
	}

	if device1.SignatureCounter != 2 {
		t.Error("Device1 state not updated properly.")
	}

	requestSignData2 := SignDataRequest{
		idDevice2,
		"Hello World to Device2",
	}
	body, _ = json.Marshal(requestSignData2)
	res, err = sendPostRequest(ts.URL+"/api/v0/sign", body)
	if err != nil {
		t.Fatal(err)
	}

	_ = json.NewDecoder(res.Body).Decode(&response)
	dataBytes, _ = json.Marshal(response.Data)
	_ = json.Unmarshal(dataBytes, &signDataResponse)

	device2, _ := db.Get(idDevice2)
	if !device2.Signer.VerifySignature([]byte(signDataResponse.SignedData), signDataResponse.Signature) {
		t.Error("Failed to verify signature from device2.")
	}

	if !bytes.Equal(device2.LastSig, signDataResponse.Signature) {
		t.Error("Device2 state not updated properly.")
	}

	if device2.SignatureCounter != 1 {
		t.Error("Device2 state not updated properly.")
	}

	requestSignData2 = SignDataRequest{
		idDevice2,
		"Hello World to Device2 again!",
	}
	body, _ = json.Marshal(requestSignData2)
	res, err = sendPostRequest(ts.URL+"/api/v0/sign", body)
	if err != nil {
		t.Fatal(err)
	}

	_ = json.NewDecoder(res.Body).Decode(&response)
	dataBytes, _ = json.Marshal(response.Data)
	_ = json.Unmarshal(dataBytes, &signDataResponse)

	if !device2.Signer.VerifySignature([]byte(signDataResponse.SignedData), signDataResponse.Signature) {
		t.Error("Failed to verify signature from device2.")
	}

	if !bytes.Equal(device2.LastSig, signDataResponse.Signature) {
		t.Error("Device1 state not updated properly.")
	}

	if device2.SignatureCounter != 2 {
		t.Error("Device1 state not updated properly.")
	}

	// Test concurrent calls to sign endpoint targeting same device
	var wg sync.WaitGroup
	numOfRequests := 100

	wg.Add(numOfRequests)

	for i := 0; i < numOfRequests; i++ {
		go func() {
			defer wg.Done()
			go sendPostRequest(ts.URL+"/api/v0/sign", body)
		}()
	}

}

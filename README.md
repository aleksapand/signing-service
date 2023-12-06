# Signing-service
This simple Signature Service supports the creation of signing device and signing the data with either RSA or ECDS algorithms.

# REST API

The REST API of the Signature Service is described below.

## Create a new signature device

### Request
Supported algorithms are RSA and ECC. Label is optional and can be left out. In that case generated id is used as a label.

`POST api/v0/new`

    curl --location 'localhost:8080/api/v0/new' \
    --header 'Content-Type: application/json' \
    --data '{
    "algorithm": "ECC",
    "label": "device1"
    }'

### Response

    {
        "data": {
            "id": "ab7717f8-7d47-4b79-b2de-619b9fdcbff0",
            "label": "Device1",
            "algorithm": "ECC",
            "publicKey": {
                "Curve": {},
                "X": 14150058114304752430485154291362069177671805864696455352452630952178867615047346244714779898186202487256956604212952,
                "Y": 19750439714822972420105680895248542285472999723123076309020756177625985895936909923909381762510581648614629924037463
            }
        }
    }

## Sign data

### Request

`POST api/v0/sign`

    curl --location 'localhost:8080/api/v0/sign' \
    --header 'Content-Type: application/json' \
    --data '{
    "id": "ab7717f8-7d47-4b79-b2de-619b9fdcbff0",
    "data": "Hello World"
    }'

### Response

    {
        "data": {
            "signature": "UzU1YmVwdzhWb1NCeXNFNVV5RkRHck96cDV6d3VCOVR4RTduSzh5WG14VjZGRStRb0kvL1VodmZKSFc3TFdKUFZPOTBKWDl2Q3VRRkE1c0p6Rm9GRmZaRGdYSUtmU0hPMmdpcDRNMEZuMWxsb1lub3JPK3dKZ1NTbGZTbmswbz0=",
            "signed_data": "0_Hello World_cTNjWCtIMUhTM215M21HYm45eS84QT09"
        }
    }

## Get all created devices

### Request

`GET api/v0/devices`

    curl --location 'localhost:8080/api/v0/devices' --data 

### Response

    {
        "data": [
            {
                "id": "ab7717f8-7d47-4b79-b2de-619b9fdcbff0",
                "label": "Device1",
                "algorithm": "ECC",
                "publicKey": {
                "Curve": {},
                "X": 14150058114304752430485154291362069177671805864696455352452630952178867615047346244714779898186202487256956604212952,
                "Y": 19750439714822972420105680895248542285472999723123076309020756177625985895936909923909381762510581648614629924037463
                }
            }
        ]
    }

## Get specific device

### Request

`GET api/v0/devices/{id}`

    curl --location 'localhost:8080/api/v0/devices/ab7717f8-7d47-4b79-b2de-619b9fdcbff0' \
    --data ''

### Response
    {
        "data": {
        "id": "ab7717f8-7d47-4b79-b2de-619b9fdcbff0",
        "label": "Device1",
        "algorithm": "ECC",
        "publicKey": {
            "Curve": {},
            "X": 14150058114304752430485154291362069177671805864696455352452630952178867615047346244714779898186202487256956604212952,
            "Y": 19750439714822972420105680895248542285472999723123076309020756177625985895936909923909381762510581648614629924037463
            }
        }
    }

# Tests

Unit Tests are located in respective packages under `signer_test.go`, `device_test.go`.
Integration Tests of server can be found under `server_test.go`.

package crypto

import (
	"errors"
	"strings"
)

// SignerFactory instantiates the correct signer given the name of algorithm.
func SignerFactory(algorithm string) (Signer, error) {

	algorithm = strings.ToUpper(algorithm)
	switch algorithm {
	case "RSA":
		return NewRSASigner(), nil
	case "ECC":
		return NewECCSigner(), nil
	default:
		return nil, errors.New("unsupported algorithm")
	}
}

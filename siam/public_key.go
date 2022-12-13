package siam

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"
)

/*
In here we implement the management of the SIAM public keys

They will be reloaded from the SIAM endpoint in fixed intervals, to get new generated ones.

// location of the SIAM key endpoint
// https://itdoc.schwarz/display/IAM/SIAM+IDP+Endpoints
*/
const jwksURL string = "https://federation.auth.schwarz/nidp/oauth/nam/keys"

var (
	myPublicKey PublicKey
)

// keep modulo and exponent
type PublicKey struct {
	Modulo     big.Int
	Exponent   int
	LastUpdate int64
	JwksURL    string
}

// JSON element for a single public key
type Key struct {
	Alg string `json:"alg"`
	Kty string `json:"kyt"`
	Mod string `json:"n"`
	Exp string `json:"e"`
}

// outer array of the json public keys
type Keys struct {
	Keys []Key `json:"keys"`
}

// constructor to keep a instance with an always valid public key
func NewPublicKey() (*PublicKey, error) {
	return &myPublicKey, nil
}

// must be refactored to auto reload the key every hour in the background
func init() {
	response, err := http.Get(jwksURL)
	manageError(err)

	data, err := io.ReadAll(response.Body)
	manageError(err)

	var values Keys
	err = json.Unmarshal(data, &values)
	manageError(err)

	n, ok := base64.RawURLEncoding.DecodeString(values.Keys[0].Mod)
	if ok != nil {
		manageError(fmt.Errorf("jwt verification error: It was not possible convert modulo string into byte array. Error was: '%s'", ok))
	}
	myPublicKey.Modulo.SetBytes(n)

	e, ok := base64.RawURLEncoding.DecodeString(values.Keys[0].Exp)
	if ok != nil {
		manageError(fmt.Errorf("jwt verification error: It was not possible convert exponent string into byte array. Error was: '%s'", ok))
	}
	var bufferExp bytes.Buffer
	bufferExp.WriteByte(0)
	bufferExp.Write(e)
	myPublicKey.Exponent = int(binary.BigEndian.Uint32(bufferExp.Bytes()))

	myPublicKey.LastUpdate = time.Now().Unix()
	fmt.Printf("Updated public SIAM key at %v\n", time.Now())
}

// Must be refactored into proper error handling
func manageError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

}

package oauth

import (
	"bytes"
	"crypto"
	"crypto/rsa"
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
In here you'll find the offline check if a jwt access_token was signed from SIAM

1. we use the public key in form of modulo and exponent provided via the SIAM key endpoint
2. we check if the decoded signature is exactly the same as the header and payload part

Hint:
The public key must be downloaded continuously as it will be generated in fixed intervals.
*/

// location of the SIAM key endpoint
// https://itdoc.schwarz/display/IAM/SIAM+IDP+Endpoints
const jwksURL string = "https://federation.auth.schwarz/nidp/oauth/nam/keys"

var (
	myPublicKey Verifier
)

// keep modulo and exponent
type Verifier struct {
	mod        big.Int
	exp        int
	lastUpdate int64
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
func NewVerifier() (*Verifier, error) {
	return &myPublicKey, nil
}

func (v Verifier) Verify(jwt Jwt) error {

	// Calculate public key with help of modulo and exponent
	publicKey := &rsa.PublicKey{N: &v.mod, E: v.exp}
	// Sign the hash of the message
	sha256Hasher := crypto.SHA256.New()
	sha256Hasher.Write(jwt.Message)

	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, sha256Hasher.Sum(nil), jwt.Signature)

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
	myPublicKey.mod.SetBytes(n)

	e, ok := base64.RawURLEncoding.DecodeString(values.Keys[0].Exp)
	if ok != nil {
		manageError(fmt.Errorf("jwt verification error: It was not possible convert exponent string into byte array. Error was: '%s'", ok))
	}
	var bufferExp bytes.Buffer
	bufferExp.WriteByte(0)
	bufferExp.Write(e)
	myPublicKey.exp = int(binary.BigEndian.Uint32(bufferExp.Bytes()))

	myPublicKey.lastUpdate = time.Now().Unix()
	fmt.Printf("Updated public SIAM key at %v\n", time.Now())
}

// Must be refactored into proper error handling
func manageError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

}

package siam

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/robfig/cron"
)

/*
In here we implement the management of the SIAM public keys

They will be reloaded from the SIAM endpoint in fixed intervals, to get new generated ones.

// location of the SIAM key endpoint
// https://itdoc.schwarz/display/IAM/SIAM+IDP+Endpoints
*/
const JwksURL string = "https://federation.auth.schwarz/nidp/oauth/nam/keys"

// interval in seconds to reload the public key (Mind to set it to a realistic value)
const reloadInterval time.Duration = 1

var (
	myPublicKey PublicKey
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	client      = http.Client{
		Timeout: 2 * time.Second,
	}
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
func GetPublicKey() (*PublicKey, error) {
	return &myPublicKey, nil
}

// must be refactored to auto reload the key every hour in the background
func init() {
	loadKey(&myPublicKey)
	c := cron.New()
	c.AddFunc("@every 2s", func() { loadKey(&myPublicKey) })
	c.Start()
}

func loadKey(key *PublicKey) {
	response, err := client.Get(JwksURL)
	if err != nil {
		ErrorLogger.Printf("Wasn't able to get the public key from '%s'. Error was: %v\n", myPublicKey.JwksURL, err)
		return
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		ErrorLogger.Printf("Wasn't able to get the result out the response from '%s'. Error was: %v\n", myPublicKey.JwksURL, err)
		return
	}

	var values Keys
	err = json.Unmarshal(data, &values)
	if err != nil {
		ErrorLogger.Printf("Wasn't able to unmarshal the public key from '%s'. Error was: %v\n", myPublicKey.JwksURL, err)
		return
	}

	n, ok := base64.RawURLEncoding.DecodeString(values.Keys[0].Mod)
	if ok != nil {
		ErrorLogger.Printf("jwt verification error: It was not possible convert modulo string into byte array. Error was: '%s'", ok)
		return
	}
	key.Modulo.SetBytes(n)

	e, ok := base64.RawURLEncoding.DecodeString(values.Keys[0].Exp)
	if ok != nil {
		ErrorLogger.Printf("jwt verification error: It was not possible convert exponent string into byte array. Error was: '%s'", ok)
		return
	}
	var bufferExp bytes.Buffer
	bufferExp.WriteByte(0)
	bufferExp.Write(e)
	key.Exponent = int(binary.BigEndian.Uint32(bufferExp.Bytes()))

	key.JwksURL = JwksURL

	key.LastUpdate = time.Now().Unix()
	log.Default().Printf("Updated public SIAM key at %v\n", time.Unix(key.LastUpdate, 0))
}

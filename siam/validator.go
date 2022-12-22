package siam

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

/*
In here we will do the jwt access token validation work

1. Verifying the signature
2. Validation the attributes

*/

// the SIAM client (same as client id)
// This will prove that the access token was issued by "your" SIAM client
const aud = "9de0efb9-767e-453b-84fa-a3b0d76c97be"
const secret = "MOUO38WXiSPeNHNwZmpo234SwNJZvxWopKZmX9XuE1anK1Qpd_RyLDOlJQSP_qxahf-ewCjxRWYhMcola2egpg"
const introspectURL = "https://federation.auth.schwarz/nidp/oauth/v1/nam/introspect"

/*
Will check that the decrypted signature is the same as header + payload
*/
func (jwt Jwt) VerifySignature(key *PublicKey) error {
	same := strings.HasPrefix(key.JwksURL, jwt.Payload.Iss)
	if !same {
		return fmt.Errorf("jwt verification error: The access token was issued by an unknown idp. We found '%s' but need '%s'", jwt.Payload.Iss, key.JwksURL)
	}

	// Calculate public key with help of modulo and exponent
	publicKey := &rsa.PublicKey{N: &key.Modulo, E: key.Exponent}
	// Sign the hash of the message
	sha256Hasher := crypto.SHA256.New()
	sha256Hasher.Write(jwt.Message)

	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, sha256Hasher.Sum(nil), jwt.Signature)
}

/*
Will validate the creation date and the expiry date. For SIAM this is enough as we do SIAM role managed access
*/
func (jwt Jwt) ValidatePayload() error {
	now := time.Now().Unix()

	// access token must be created in the past
	if jwt.Payload.Iat >= now {
		return fmt.Errorf("jwt verification error: The access token was created in the future. was created '%v' but now we have '%v'", time.Unix(jwt.Payload.Iat, 0), time.Unix(now, 0))
	}

	//access token must still be valid (we add 10 minutes)
	if jwt.Payload.Exp+(60*10) < now {
		return fmt.Errorf("jwt verification error: The access token already outdated. It was valid since '%v' but now we already have '%v'", time.Unix(jwt.Payload.Exp, 0), time.Unix(now, 0))
	}

	//access token must be created just for our application
	if jwt.Payload.Aud != aud {
		return fmt.Errorf("jwt verification error: The access token has the wrong audience. The audience of this token is '%s'", jwt.Payload.Aud)
	}

	return nil
}

func (jwt Jwt) Introspect() error {
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	req, err := http.NewRequest("POST", introspectURL, strings.NewReader(fmt.Sprintf("token=%s", jwt.Token)))
	if err != nil {
		return fmt.Errorf("jwt verification error: Not able to build the introspection POST call '%s'", err)
	}

	req.Header = http.Header{
		"Content-Type":  {"application/x-www-form-urlencoded"},
		"Authorization": {"Basic: " + base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", aud, secret)))},
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("jwt verification error: Not able to call the introspection URL '%s'", err)
	}

	array, _ := io.ReadAll(res.Body)
	type intro struct {
		Active bool `json:"active"`
	}
	var result intro
	err = json.Unmarshal(array, &result)
	if err != nil {
		return fmt.Errorf("jwt verification error: Not able to call the introspection URL '%s'", err)
	}

	if !result.Active {
		return fmt.Errorf("jwt verification error: The token is not valid any more (Introspection result) '%s'", err)
	}
	fmt.Println(result.Active)
	return nil
}

package siam

import (
	"crypto"
	"crypto/rsa"
	"fmt"
	"strings"
	"time"
)

/*
In here we will do the jwt access token validation work

1. Verifying the signature
2. Validation the attributes

*/

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

	if jwt.Payload.Iat >= now {
		return fmt.Errorf("jwt verification error: The access token was created in the future. was created '%v' but now we have '%v'", time.Unix(jwt.Payload.Iat, 0), time.Unix(now, 0))
	}

	//we add 10 minutes
	if jwt.Payload.Exp < now+(60*10) {
		return fmt.Errorf("jwt verification error: The access token already outdated. It was valid since but '%v' now we already have '%v'", time.Unix(jwt.Payload.Exp, 0), time.Unix(now, 0))
	}

	return nil
}

package siam

import (
	"crypto"
	"crypto/rsa"
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
	// Calculate public key with help of modulo and exponent
	publicKey := &rsa.PublicKey{N: &key.Modulo, E: key.Exponent}
	// Sign the hash of the message
	sha256Hasher := crypto.SHA256.New()
	sha256Hasher.Write(jwt.Message)

	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, sha256Hasher.Sum(nil), jwt.Signature)
}

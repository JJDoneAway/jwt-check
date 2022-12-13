package siam

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

/*
This object will represent a SIAM JWT access token.
*/
type Jwt struct {
	Header  Header
	Payload Payload
	//Header and payload as decoded byte array
	Message []byte
	//Signature as decoded string
	Signature []byte
}

// Just the important header fields
type Header struct {
	Alg  string `json:"alg"`
	Type string `json:"type"`
}

// the payload contains SIAM / SIT specific claims
type Payload struct {
	// who issued (created) this token
	Iss string `json:"iss"`
	// resource owner id
	Uid string `json:"uid"`
	// SIAM roles in RAW format
	GroupMembership []string `json:"groupMembership"`
	// until when is the token valid
	Exp int64 `json:"exp"`
	// when was the token created
	Iat int64 `json:"iat"`
}

func DecodeJWT(jwt string) (Jwt, error) {
	var ret Jwt
	tokenParts := strings.Split(jwt, ".")
	if len(tokenParts) != 3 {
		return ret, fmt.Errorf("jwt verification error: JWT access_token string must contain out of three parts, separated with dots. We found %d parts", len(tokenParts))
	}

	ret.Message = []byte(strings.Join(tokenParts[0:2], "."))

	signature, ok := base64.RawURLEncoding.DecodeString(tokenParts[2])
	if ok != nil {
		return ret, fmt.Errorf("jwt verification error: It was not possible to extract the signature. Error was: '%s'", ok)
	}
	ret.Signature = signature

	if dec, ok := base64.RawURLEncoding.DecodeString(tokenParts[0]); ok != nil {
		return ret, fmt.Errorf("jwt verification error: It was not possible to decode the header. Error was: '%s'", ok)
	} else {
		var header Header
		if ok := json.Unmarshal([]byte(dec), &header); ok != nil {
			return ret, fmt.Errorf("jwt verification error: It was not possible to convert the header into JSON. Error was: '%s'", ok)
		}
		ret.Header = header
	}

	if dec, ok := base64.RawURLEncoding.DecodeString(tokenParts[1]); ok != nil {
		return ret, fmt.Errorf("jwt verification error: It was not possible to decode the payload. Error was: '%s'", ok)
	} else {
		var payload Payload
		if ok := json.Unmarshal([]byte(dec), &payload); ok != nil {
			return ret, fmt.Errorf("jwt verification error: It was not possible to convert the payload into JSON. Error was: '%s'", ok)
		}
		ret.Payload = payload
	}

	return ret, nil
}

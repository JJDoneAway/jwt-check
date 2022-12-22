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
	//the original token
	Token string
}

// Just the important header fields
type Header struct {
	Alg  string `json:"alg"`
	Type string `json:"type"`
}

/*
	 the payload contains SIAM / SIT specific claims
		iss => who issued (created) this token
		uid => resource owner id
		aud => the audience of this token (the API)
		groupMembership => SIAM roles in RAW Active Directory format
		exp => until when is the token valid
		iat => when was the token created
*/
type Payload struct {
	Iss             string   `json:"iss"`
	Uid             string   `json:"uid"`
	FullName        string   `json:"fullName"`
	Email           string   `json:"mail"`
	Aud             string   `json:"aud"`
	GroupMembership []string `json:"groupMembership"`
	Exp             int64    `json:"exp"`
	Iat             int64    `json:"iat"`
}

func DecodeJWT(jwt string) (Jwt, error) {
	var ret Jwt
	ret.Token = jwt
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

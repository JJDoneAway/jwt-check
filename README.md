# JWT access token check in SIAM context
--------------------------------------

This code should show how a JWT access token can be validated offline in front of an API

It has the following parts:

## Public Key
------------
It download the public key from the jwks endpoint. The key is necessary for the verification of the signature

[public_key.go](./siam/public_key.go) 

##Convert JWT string into go slice
----------------------------------
It will extract all necessary attributes out of the JWT raw string

[access_token.go](./siam/access_token.go)

## Validation and versification of the JWT
-----------------------------------------
* The signature will be decoded with the public key and compared with the header and payload string
* The attributes out of the payload will be validated

[validator.go](./siam/validator.go)

## Get the user context
----------------------
* All attributes to authN and authZ the user will be extracted and stored in a slice
* Role management is implemented in here

[user.go](./siam/user.go)

## Tests
-------
To verify the implementation, all functions are tested in here

[jwt_test.go](./siam/jwt_test.go)

## Example
---------
find in the main a example on how to use this code

[main.go](main.go)

```
go run jwt-check
```
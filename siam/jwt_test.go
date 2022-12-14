package siam

import (
	"log"
	"math/big"
	"strings"
	"testing"
	"time"
)

const mod string = "20466271965265116130583795161097203145303306928701644085814824022456203387753085999066651433237924598404905480315483953783694420153052459265864491099133746066901565585850053961086116903370201389520483061300470438385181474233866798864623730594882212749344038269851645961975139482722223220011857493000632753075599487805668830984497051624604467790565811184854803668274055877525097055279257874736042335331693026265806740719499552110805089548573791581810134106296583472674555845220401251517475586374584992120421312677336731731082174474644305471083669011559949118852777633366451767100723452697621007609610525634699195720077"
const exp int = 65537
const lastUpdated int64 = 1671019458

const jwtTestRaw string = "eyJraWQiOiI5MDk3MDM3OTI2MzI2NjAxMzUxOTgwOTkzNzQwODA3MjEzODkwOCIsInR5cCI6IkpXVCIsImFsZyI6IlJTMjU2In0.eyJpc3MiOiJodHRwczovL2ZlZGVyYXRpb24uYXV0aC5zY2h3YXJ6L25pZHAvb2F1dGgvbmFtIiwianRpIjoiNmUxZTg4MTMtNmUxMi00ZDNjLWI2NDktZDVmY2Q2NjE5ZWZlIiwiYXVkIjoiOWRlMGVmYjktNzY3ZS00NTNiLTg0ZmEtYTNiMGQ3NmM5N2JlIiwiZXhwIjoxNjcwNzg1NTMwLCJpYXQiOjE2NzA3ODE5MzAsIm5iZiI6MTY3MDc4MTkwMCwic3ViIjoiYzMxZTBlYmMxODhhYWY0NWJiYmFjMzFlMGViYzE4OGEiLCJfcHZ0IjoiTGRJSEdueDRQcVUwNFBCcFE1TkpzNnp5dmR0MGRKUnYyNzRMRjBaczhSTUlCMXNXMlNMNzBKRzduUnJaVWY2OWZ4OFovNG56THpMRGV1eVhGSDZDNlNCRTJXK3BQcUNJbTF2b2QxcHVZQUYzWllEV2trTVFMUzFvSzYwSU5zM1M4TkdXdUxZSWFub3hwQ0c2L20zY3UxU2kwM0dSbW5PT01BRUR5TDFaZzd1WkZDRjBHRjM3YlFOdzE5RW5Fc2w1VzZYVkJCYlBrU2pjUW1TcDVYNDN3TnBVcVB1QkNCSnZwRmwxVzlvOVVZcFViQlBjTFEyZHV0NlQ2bEJjd1lzQmoyWTNhK1FkU3krUG1NNmtvZ0I4OW9aMkxpNlFzTTFQekpYQXFLeU1LdndBOFRuczE5K2VZMTl5ZXRLVGluQXBGRUhrVUI5WkJ5OTVyZ2RuODNLNnRJUTdMNitYYWNsaFhPQ2tTNXdmNGZNaHFIM1BFZVpMbE9wc2lHTk1ETUpFb3JBM0dyb1FFNFdSUmxWdExhUCt3VTNIMkxEM2tSbXdqOUZ6eGRpZTl5UWtsUkY4ZUd1TjA4L0xrNHVZTThBbWF4UENMYTJVNGRGSk5JQVM2STNjdkk2U3M0RDNvbW5aQmdpc0N5ODlUb3AvK3VBWWNKNkh1L3dZS3NRWTl0ekhabDR2eTZST0VmN3VtUGhIRXFYNmo3cEc0YlpCc0pWKzBuRUJQcUFlOGp5Q3pHNU0xeTdYTk5BK0xucHBVbkJXLjkiLCJzY29wZSI6WyJzaWFtIiwicHJvZmlsZSJdLCJ1aWQiOiJob2VobmVqbyIsIm1haWwiOiJKb2hhbm5lcy5Ib2VobmVAbWFpbC5zY2h3YXJ6IiwiZ2l2ZW5OYW1lIjoiSm9oYW5uZXMiLCJuaWNrbmFtZSI6IjEwMDI0Njk3NjIiLCJmdWxsTmFtZSI6IkpvaGFubmVzIEjDtmhuZSIsInNuIjoiSMO2aG5lIiwiQ2xvdWRMb2dpbk5hbWUiOiJKb2hhbm5lcy5Ib2VobmVAbWFpbC5zY2h3YXJ6Iiwia2xDbG91ZExvZ2luTmFtZSI6IkpvaGFubmVzLkhvZWhuZUBtYWlsLnNjaHdhcnoiLCJncm91cE1lbWJlcnNoaXAiOlsiY249c2l0LW5ld3MtbWVtYmVyLG91PWFwcDR5b3Usb3U9YXBwcyxvPWdsb2JhbCIsImNuPWFsbC1zZGwsb3U9YXBwNHlvdSxvdT1hcHBzLG89Z2xvYmFsIiwiY249YWxsLXN6ZCxvdT1hcHA0eW91LG91PWFwcHMsbz1nbG9iYWwiLCJjbj1zaXQtbWVtYmVyLG91PWFwcDR5b3Usb3U9YXBwcyxvPWdsb2JhbCIsImNuPWdldHNtYXJ0LW5nLWFwcC1zYm94LG91PWdldHNtYXJ0LG91PWFwcHMsbz1nbG9iYWwiLCJjbj1hcmliYS1pbnQtaW50ZXJuLG91PWFyaWJhLG91PWFwcHMsbz1nbG9iYWwiLCJjbj1hbGwtc2l0LG91PWFwcDR5b3Usb3U9YXBwcyxvPWdsb2JhbCIsImNuPWdldHNtYXJ0LW5nLWxpZGwtaW50LWNvcGUsb3U9Z2V0c21hcnQsb3U9YXBwcyxvPWdsb2JhbCIsImNuPWZpbGU0eW91LWludC1BY2Nlc3Msb3U9ZmlsZTR5b3Usb3U9YXBwcyxvPWdsb2JhbCIsImNuPW9kai1pbnQtc2l0LWJjZC1vZGotdXNlcixvdT1vZGosb3U9YXBwcyxvPWdsb2JhbCIsImNuPXN0YWNraXQtcG9ydGFsLXh4LWFjY291bnQtbWVtYmVyLG91PXN0YWNraXQtcG9ydGFsLG91PWFwcHMsbz1nbG9iYWwiLCJjbj1zbnlrLXh4LXNpdC1tYWdpYy1tb25pdG9yaW5nLWFkbSxvdT1zbnlrLG91PWFwcHMsbz1nbG9iYWwiLCJjbj1hcnRpZmFjdG9yeS14eC1zaXQtbWFnaWMtbW9uaXRvcmluZy1yZHIsb3U9YXJ0aWZhY3Rvcnksb3U9YXBwcyxvPWdsb2JhbCIsImNuPW9kai1pbnQtc2l0LWJjZC1zdWItc3ViLWVhbS1kZXYsb3U9b2RqLG91PWFwcHMsbz1nbG9iYWwiLCJjbj1zbnlrLXh4LXNpdC1vZGotYXBpZW5hYmxlbWVudC1hZG0sb3U9c255ayxvdT1hcHBzLG89Z2xvYmFsIiwiY249c29uYXJxdWJlLXh4LXNpdC1vZGotYXBpZW5hYmxlbWVudC1jb24sb3U9c29uYXJxdWJlLG91PWFwcHMsbz1nbG9iYWwiLCJjbj1hcnRpZmFjdG9yeS14eC1zaXQtb2RqLWFwaWVuYWJsZW1lbnQtYWRtLG91PWFydGlmYWN0b3J5LG91PWFwcHMsbz1nbG9iYWwiLCJjbj1naXRlYS14eC12Y3MtYWNjZXNzLG91PWdpdGVhLG91PWFwcHMsbz1nbG9iYWwiLCJjbj1wdWxzZS12cG4tbC1pbnQsb3U9cHVsc2Usb3U9YXBwcyxvPWdsb2JhbCIsImNuPWFydGlmYWN0b3J5LXh4LXNpdC1vZGotam9oYW5uZXMtbWVldHMtb2RqLWFkbSxvdT1hcnRpZmFjdG9yeSxvdT1hcHBzLG89Z2xvYmFsIiwiY249c255ay14eC1zaXQtb2RqLWpvaGFubmVzLW1lZXRzLW9kai1hZG0sb3U9c255ayxvdT1hcHBzLG89Z2xvYmFsIiwiY249c29uYXJxdWJlLXh4LXNpdC1vZGotam9oYW5uZXMtbWVldHMtb2RqLWNvbixvdT1zb25hcnF1YmUsb3U9YXBwcyxvPWdsb2JhbCIsImNuPW9kai1pbnQtc2l0LWJjZC1zdWItc3ViLW1vbml0b3Jpbmctc2Nod2Fyei1kZXYsb3U9b2RqLG91PWFwcHMsbz1nbG9iYWwiLCJjbj1zbnlrLXh4LXNpdC1vZGotYWxlcnRpbmctc2Nod2Fyei1hZG0sb3U9c255ayxvdT1hcHBzLG89Z2xvYmFsIiwiY249c29uYXJxdWJlLXh4LXNpdC1vZGotYWxlcnRpbmctc2Nod2Fyei1jb24sb3U9c29uYXJxdWJlLG91PWFwcHMsbz1nbG9iYWwiLCJjbj1hcnRpZmFjdG9yeS14eC1zaXQtb2RqLWFsZXJ0aW5nLXNjaHdhcnotYWRtLG91PWFydGlmYWN0b3J5LG91PWFwcHMsbz1nbG9iYWwiLCJjbj1hcnRpZmFjdG9yeS14eC1zaXQtbW9uaXRvcmluZy1zY2h3YXJ6LWFkbSxvdT1hcnRpZmFjdG9yeSxvdT1hcHBzLG89Z2xvYmFsIiwiY249c255ay14eC1zaXQtb2RqLW1vbml0b3Jpbmctc2Nod2Fyei1hZG0sb3U9c255ayxvdT1hcHBzLG89Z2xvYmFsIiwiY249c255ay14eC1zaXQtb2RqLW1vbml0b3Jpbmctc2Nod2Fyei1jb24sb3U9c255ayxvdT1hcHBzLG89Z2xvYmFsIiwiY249c255ay14eC1zaXQtb2RqLWVhbS1iYXNpYy1hZG0sb3U9c255ayxvdT1hcHBzLG89Z2xvYmFsIiwiY249c29uYXJxdWJlLXh4LXNpdC1vZGotdGVzdHRsMTIzNC1jb24sb3U9c29uYXJxdWJlLG91PWFwcHMsbz1nbG9iYWwiLCJjbj1zbnlrLXh4LXNpdC1vZGotY29udGFjdHVpMDgxNS1hZG0sb3U9c255ayxvdT1hcHBzLG89Z2xvYmFsIiwiY249c255ay14eC1zaXQtb2RqLWVhbS1hei1hLWFkbSxvdT1zbnlrLG91PWFwcHMsbz1nbG9iYWwiLCJjbj1zb25hcnF1YmUteHgtc2l0LW9kai1jb250YWN0dWkwODE1LWNvbixvdT1zb25hcnF1YmUsb3U9YXBwcyxvPWdsb2JhbCIsImNuPXNvbmFycXViZS14eC1zaXQtb2RqLWVhbS1iYXNpYy1jb24sb3U9c29uYXJxdWJlLG91PWFwcHMsbz1nbG9iYWwiLCJjbj1zbnlrLXh4LXNpdC1vZGotdGVzdHRsMTIzNC1hZG0sb3U9c255ayxvdT1hcHBzLG89Z2xvYmFsIiwiY249c29uYXJxdWJlLXh4LXNpdC1vZGotbW9uaXRvcmluZy1zY2h3YXJ6LWNvbixvdT1zb25hcnF1YmUsb3U9YXBwcyxvPWdsb2JhbCIsImNuPXNueWsteHgtc2l0LW9kai1hZGRyZXNzYm9va3Bhc2NhbC1hZG0sb3U9c255ayxvdT1hcHBzLG89Z2xvYmFsIiwiY249YXJ0aWZhY3RvcnkteHgtc2l0LWFkZHJlc3Nib29rcGFzY2FsLWFkbSxvdT1hcnRpZmFjdG9yeSxvdT1hcHBzLG89Z2xvYmFsIiwiY249YXJ0aWZhY3RvcnkteHgtc2l0LWNsb3Vkc3VydmV5cy1hZG0sb3U9YXJ0aWZhY3Rvcnksb3U9YXBwcyxvPWdsb2JhbCIsImNuPXNueWsteHgtc2l0LW9kai1jbG91ZHN1cnZleXMtYWRtLG91PXNueWssb3U9YXBwcyxvPWdsb2JhbCIsImNuPXNvbmFycXViZS14eC1zaXQtb2RqLWNsb3Vkc3VydmV5cy1jb24sb3U9c29uYXJxdWJlLG91PWFwcHMsbz1nbG9iYWwiLCJjbj1jb21tY2Fycy14eC1kd2sta29uZmlndXJpZXJlci1mZTQsb3U9Y29tbWNhcnMsb3U9YXBwcyxvPWdsb2JhbCIsImNuPWVmcy14eC1uZXRkcml2ZS1hZHMtc2Nod2Fyei1vLG91PWVmcyxvdT1hcHBzLG89Z2xvYmFsIl0sIndvcmtmb3JjZUlEIjoiMTAwMjQ2OTc2MiIsIl90YXJnZXQiOiJJZGVudGl0eVByb3ZpZGVyUlNVRSJ9.Sz3MCQOpjKeZF6Tgg8zsuHlItYEGIBf1J4aMrSaU_ImZxoDNw0A8Bw2ihCcjTpCrZui-Z5GbRRnUowFQPcPk9OtRTG0INbP2GKrzoq7AN9g2bGFfPwlXmcyiedJLwGzhxgd2k8rcCsAu79LGPtTFXrtQ_2oIfPhBzoj5gyEG4QqJaQ8UXOpZTBcNYtaTB4I-FUumkFGmE72ljhj4lWQcOfp8rAZJRH35KI_veFvFJlGjVb3H1fy2AhQ2gTXS5Y5jIMIt0zFj7yv1XuauXgnmGgm3OIvosyj5cgqsXKTbFDLS8IquqKzGLusCaeq7fx4wjcKy2tDZvBGJbGqkGRRzlg"

var (
	publicKey PublicKey
)

func setupSuite(tb testing.TB) {
	log.Println("Setup public key")
	publicKey.JwksURL = JwksURL
	publicKey.LastUpdate = lastUpdated
	publicKey.Exponent = exp

	val, _ := new(big.Int).SetString(mod, 10)
	publicKey.Modulo = *val
}

func TestAccessToken(t *testing.T) {
	setupSuite(t)
	t.Run("Test convertion of JWT string into struct", func(t *testing.T) {

		jwt, ok := DecodeJWT(jwtTestRaw)

		if ok != nil {
			t.Errorf("Convertion into Jwt failed: %v", ok)
		}

		if jwt.Header.Alg != "RS256" {
			t.Error("Alg is not mapped")
		}

		if jwt.Payload.Aud != "9de0efb9-767e-453b-84fa-a3b0d76c97be" {
			t.Error("aud is missing")
		}

		if jwt.Payload.FullName != "Johannes Höhne" {
			t.Error("Name is not mapped")
		}

		if jwt.Payload.Exp != 1670785530 {
			t.Errorf("Exp is not mapped %d", jwt.Payload.Exp)
		}

		if jwt.Payload.Iat != 1670781930 {
			t.Errorf("Iat is not mapped %d", jwt.Payload.Iat)
		}

	})

}

func TestVerifySignature(t *testing.T) {
	setupSuite(t)
	jwt, _ := DecodeJWT(jwtTestRaw)

	t.Run("Test verify signature", func(t *testing.T) {
		ok := jwt.VerifySignature(&publicKey)

		//positive one
		if ok != nil {
			t.Errorf("Verification of the signature failed: %v", ok)
		}

		//change signature
		jwt.Signature[0] = 10
		if ok = jwt.VerifySignature(&publicKey); ok.Error() != "crypto/rsa: verification error" {
			t.Errorf("Verification of the signature mus fail failed: %v", ok)
		}

		//change message
		jwt, _ = DecodeJWT(jwtTestRaw)
		jwt.Message[10] = 12
		if ok = jwt.VerifySignature(&publicKey); ok.Error() != "crypto/rsa: verification error" {
			t.Errorf("Verification of the signature mus fail failed: %v", ok)
		}

	})

}

func TestValidateJwt(t *testing.T) {
	jwt, _ := DecodeJWT(jwtTestRaw)

	t.Run("Test validate JWT", func(t *testing.T) {
		ok := jwt.ValidatePayload()

		//it is already to old one
		if ok == nil {
			t.Errorf("Validation of the payload must faile: %v", ok)
		}

		//correct time
		jwt.Payload.Iat = time.Now().Unix() - 10
		jwt.Payload.Exp = time.Now().Unix() + 10
		if ok = jwt.ValidatePayload(); ok != nil {
			t.Errorf("Validation of the payload must not faile: %v", ok)
		}

		//future valid
		jwt.Payload.Iat = time.Now().Unix() + 10
		jwt.Payload.Exp = time.Now().Unix() + 100
		if ok = jwt.ValidatePayload(); ok == nil {
			t.Errorf("Validation of the payload must faile (future valid): %v", ok)
		}

		//outdated
		jwt.Payload.Iat = time.Now().Unix() - 10
		jwt.Payload.Exp = time.Now().Unix() - 5*60*10
		if ok = jwt.ValidatePayload(); ok == nil {
			t.Errorf("Validation of the payload must faile (out dated): %v", ok)
		}

		//outdated
		jwt.Payload.Iat = time.Now().Unix() - 10
		jwt.Payload.Exp = time.Now().Unix() + 10
		jwt.Payload.Aud = "Hallo"
		if ok = jwt.ValidatePayload(); ok == nil {
			t.Errorf("Validation of the payload must failed (wrong aud): %v", ok)
		}

	})

}

func TestJwtUser(t *testing.T) {
	jwt, _ := DecodeJWT(jwtTestRaw)

	t.Run("Test user", func(t *testing.T) {
		user := jwt.GetUser()

		if user.Name != "Johannes Höhne" {
			t.Errorf("User name is not mapped: %v", user)
		}

		if user.Email != "Johannes.Hoehne@mail.schwarz" {
			t.Errorf("User mail is not mapped: %v", user)
		}

		if user.Uid != "hoehnejo" {
			t.Errorf("User uid is not mapped: %v", user)
		}

		if user.Roles[0] != "sit-news-member" {
			t.Errorf("Wrong role. Found %s but want %s", user.Roles[0], "sit-news-member")
		}

		if user.Roles[len(user.Roles)-1] != "efs-xx-netdrive-ads-schwarz-o" {
			t.Errorf("Wrong role. Found %s but want %s", user.Roles[len(user.Roles)-1], "efs-xx-netdrive-ads-schwarz-o")
		}

		if len(user.Roles) != 44 {
			t.Errorf("Missing roles. Found %d but want %d", len(user.Roles), 44)
		}

	})

}

func TestFindUser(t *testing.T) {
	jwt, _ := DecodeJWT(jwtTestRaw)

	t.Run("Test user", func(t *testing.T) {
		user := jwt.GetUser()

		if user.HasRole("LaberRababer") {
			t.Errorf("User must not be in the list: %v", user)
		}

		if !user.HasRole("sit-news-member") {
			t.Errorf("Role %s must be in the list", "sit-news-member")
		}

		if !user.HasRole("efs-xx-netdrive-ads-schwarz-o") {
			t.Errorf("Role %s must be in the list", "efs-xx-netdrive-ads-schwarz-o")
		}

		if !user.HasRole(strings.ToUpper("efs-xx-netdrive-ads-schwarz-o")) {
			t.Errorf("Role %s must be in the list", "efs-xx-netdrive-ads-schwarz-o")
		}

	})

}

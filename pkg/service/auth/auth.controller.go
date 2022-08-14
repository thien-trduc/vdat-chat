package auth

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/dgrijalva/jwt-go"
	"github.com/getsentry/sentry-go"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/cors"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/utils"
	"log"
	"net/http"
	"strings"
)

var Jwtkey = (`-----BEGIN CERTIFICATE-----
MIICpTCCAY0CBgFrPLdvYjANBgkqhkiG9w0BAQsFADAWMRQwEgYDVQQDDAt2ZGF0bGFiLmNvbTAeFw0xOTA2MDkxNDQ4MDNaFw0yOTA2MDkxNDQ5NDNaMBYxFDASBgNVBAMMC3ZkYXRsYWIuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAleIkHoO6Q0GRQ4POIAKmN5Ev3zfAm8raTJQ1e/CbTXW4FQ0kDS9YPhLXPwcdnbxiL3rSGgz7+iWcq/Ix7yExuNbSyqDUjLUJSU6I9JvB1YP8GSaO8d996+TVCDC8E/VSID6wmfWbMNb5Ns6Y7YY/HAhj9zc73ObErvi0NV0BjeYAVOBqJKKgl9cHfyBshr+kpC/7nrbTRnAP7JQhKrQF6wBTKQiuJlEyYqvi1ugCRBYg2BZLPtTry+Kineb1DT8ynmxJjKMtr9hU0dsLPJpqW/4DWwNOarLOBP/K9WkfR2LUxbrm41goSTjJbz6s7f/Mvn/gDLjGjIsdlFP3Y7I2lwIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQBXU5Awwhv/cYJKCdSUzmtXpXty8KrdrHaDNa8potDXlEc2JrK3wHyFRwwfpBhkaicP0LllxRHUGUNsWFnggae1fudc75fysZ16NPH7VJlUuyV96K06K4v1aM5VCSWl5djky7rtyfi2W9iH2ddWZvCeSyFsSgCD4P5GjgYpsLy27g/cvdJJAdp/b7bweVDI1grlBtnInxLUPhJ4cnoNw3crh7twqKgG6F3GmZc2Hjl45LdlxBFfftDUYH66D1X0mdoipQCbg4JWlIxUZHVjJDIrSIlwnRMwjzCm7MUYv0ySmvsxgoNVI2NuFU6A/F7zlyVkDkmO4ilp4BueRtBKb7yR
-----END CERTIFICATE-----`)

type UserClaims struct {
	jwt.StandardClaims
	FullName   string `json:"name"`
	UserName   string `json:"preferred_username"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}

const UserKey = "Subject"

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cors.SetupResponse(&w, r)

		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			utils.ResponseErr(w, http.StatusUnauthorized)
		}
		if ValiToken(tokenHeader) {
			user := JWTparseOwner(tokenHeader)
			if len(user) < 0 {
				utils.ResponseErr(w, http.StatusUnauthorized)
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, UserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			utils.ResponseErr(w, http.StatusUnauthorized)
		}
	}
}

func AuthenMiddleJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cors.SetupResponse(&w, r)

		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			utils.ResponseErr(w, http.StatusUnauthorized)
		}
		if ValiToken(tokenHeader) {
			next.ServeHTTP(w, r)
		} else {
			utils.ResponseErr(w, http.StatusUnauthorized)
		}
	}

}
func ValiToken(tokenHeader string) bool {
	splitted := strings.Split(tokenHeader, " ") // Bearer jwt_token
	if len(splitted) != 2 {
		return false
	}
	tokenPart := splitted[1]

	block, _ := pem.Decode([]byte(Jwtkey))
	var cert *x509.Certificate
	cert, _ = x509.ParseCertificate(block.Bytes)
	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)
	tk := &jwt.StandardClaims{}
	token, err1 := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return rsaPublicKey, nil
	})
	if err1 != nil {
		return false
	}
	if token.Valid {
		return true
	} else {
		return false
	}
}
func JWTparseOwner(tokenHeader string) string {
	splitted := strings.Split(tokenHeader, " ") // Bearer jwt_token
	if len(splitted) != 2 {
		return ""
	}

	block, _ := pem.Decode([]byte(Jwtkey))
	var cert *x509.Certificate
	cert, _ = x509.ParseCertificate(block.Bytes)

	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)
	tokenPart := splitted[1]
	tk := &UserClaims{}
	_, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return rsaPublicKey, nil
	})
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s \n", err)
		return ""
	}
	return tk.Subject
}

func JWTparseOwner2(tokenHeader string) string {

	block, _ := pem.Decode([]byte(Jwtkey))
	var cert *x509.Certificate
	cert, _ = x509.ParseCertificate(block.Bytes)

	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)

	tk := &UserClaims{}
	_, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
		return rsaPublicKey, nil
	})
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s \n", err)
		return ""
	}
	return tk.Subject
}

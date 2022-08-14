package utils

import (
	"encoding/json"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func CheckTokenExp(tokenStr string) bool {
	//tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})

	if err != nil {
		log.Fatal(err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Fatalf("Can't convert token's claims to standard claims")
	}

	var tm time.Time
	switch iat := claims["iat"].(type) {
	case float64:
		tm = time.Unix(int64(iat), 0)
	case json.Number:
		v, _ := iat.Int64()
		tm = time.Unix(v, 0)
	}
	now := time.Now()
	diff := now.Sub(tm)
	if diff >= time.Duration(3500)*time.Second {
		return false
	}
	return true
}

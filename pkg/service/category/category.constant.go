package category

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("Internal Server Error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("Your requested Item is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("Your Item already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("Given Param is not valid")
)

type Token struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  string `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

func Connect() string {
	const (
		clientSecret string = "7161982e-cabe-44d3-ade1-324698d2f5d8"
		clientId     string = "chat.services.vdatlab.com"
		urlHost      string = "https://accounts.vdatlab.com/auth/realms/vdatlab.com/protocol/openid-connect/token"
	)

	client := &http.Client{}
	data := url.Values{}
	data.Set("client_id", clientId)
	data.Add("client_secret", clientSecret)
	//data.Add("scope","profile")
	data.Add("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", urlHost, bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	if err != nil {
		log.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	f, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var token Token
	json.Unmarshal(f, &token)
	//fmt.Print(token.AccessToken)
	//fmt.Println(string(f))

	return token.AccessToken
}

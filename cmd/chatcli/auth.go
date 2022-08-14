package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"

	//"github.com/oklog/oklog/pkg/groups"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
	"github.com/skratchdot/open-golang/open"
)

type WellKnownConfig struct {
	Issuer                                     string   `json:"issuer,omitempty"`
	AuthorizationEndpoint                      string   `json:"authorization_endpoint,omitempty"`
	TokenEndpoint                              string   `json:"token_endpoint,omitempty"`
	TokenIntrospectionEndpoint                 string   `json:"token_introspection_endpoint,omitempty"`
	EndSessionEndpoint                         string   `json:"end_session_endpoint,omitempty"`
	JwksUri                                    string   `json:"jwks_uri,omitempty"`
	GrantTypesSupported                        []string `json:"grant_types_supported,omitempty"`
	ResponseTypesSupported                     []string `json:"response_types_supported,omitempty"`
	ResponseModesSupported                     []string `json:"response_modes_supported,omitempty"`
	RegistrationEndpoint                       string   `json:"registration_endpoint,omitempty"`
	TokenEndpointAuthMethodsSupported          []string `json:"token_endpoint_auth_methods_supported,omitempty"`
	TokenEndpointAuthSigningAlgValuesSupported []string `json:"token_endpoint_auth_signing_alg_values_supported,omitempty"`
	ScopesSupported                            []string `json:"scopes_supported,omitempty"`
	IntrospectionEndpoint                      string   `json:"introspection_endpoint,omitempty"`
}

func getWellKnownConfig(serverUrl string) (*WellKnownConfig, error) {
	resp, err := http.Get(fmt.Sprintf("%s/.well-known/uma2-configuration", serverUrl))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var c WellKnownConfig
	err = json.Unmarshal(body, &c)
	if err != nil {
		panic(err)
	}
	return &c, nil
}

// AuthorizeUser implements the PKCE OAuth2 flow.
func login() (*oauth2.Token, error) {
	var (
		authServerUrl = "https://accounts.vdatlab.com/auth/realms/vdatlab.com"
		clientID      = "chat.apps.vdatlab.com"
		redirectURL   = "http://127.0.0.1:12345/auth/callback"
		token         *oauth2.Token
	)

	authServerWellKnownCfg, err := getWellKnownConfig(authServerUrl)
	if err != nil {
		panic(err)
	}

	// initialize the code verifier
	var CodeVerifier, _ = cv.CreateCodeVerifier()

	// Create code_challenge with S256 method
	codeChallenge := CodeVerifier.CodeChallengeS256()

	// construct the authorization URL (with Auth0 as the authorization provider)
	authorizationURL := fmt.Sprintf(
		"%s?audience="+
			"&scope=openid"+
			"&response_type=code"+
			"&client_id=%s"+
			"&code_challenge=%s"+
			"&code_challenge_method=S256"+
			"&redirect_uri=%s",
		authServerWellKnownCfg.AuthorizationEndpoint, clientID, codeChallenge, redirectURL)

	// start a web server to listen on a callback URL
	server := &http.Server{Addr: redirectURL}

	// define a handler that will get the authorization code, call the token endpoint, and close the HTTP server
	http.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		// get the authorization code
		code := r.URL.Query().Get("code")
		if code == "" {
			fmt.Println("snap: Url Param 'code' is missing")
			io.WriteString(w, "Error: could not find 'code' URL parameter\n")

			// close the HTTP server and return
			cleanup(server)
			return
		}

		// trade the authorization code and the code verifier for an access token
		codeVerifier := CodeVerifier.String()
		token, err = getAccessToken(authServerWellKnownCfg.TokenEndpoint, clientID, codeVerifier, code, redirectURL)
		if err != nil {
			fmt.Println("snap: could not get access token")
			io.WriteString(w, "Error: could not retrieve access token\n")

			// close the HTTP server and return
			cleanup(server)
			return
		}

		// return an indication of success to the caller
		io.WriteString(w, `
		<html>
			<body>
				<h1>Login successful!</h1>
			</body>
			<script type='text/javascript'>
				 self.close();
			</script>
		</html>`)

		fmt.Println("Login Successfully")

		// close the HTTP server
		cleanup(server)
	})

	// parse the redirect URL for the port number
	u, err := url.Parse(redirectURL)
	if err != nil {
		fmt.Printf("snap: bad redirect URL: %s\n", err)
		os.Exit(1)
	}

	// set up a listener on the redirect port
	port := fmt.Sprintf(":%s", u.Port())
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("snap: can't listen to port %s: %s\n", port, err)
		os.Exit(1)
	}

	// open a browser window to the authorizationURL
	err = open.Start(authorizationURL)
	if err != nil {
		fmt.Printf("snap: can't open browser to URL %s: %s\n", authorizationURL, err)
		os.Exit(1)
	}

	// start the blocking web server loop
	// this will exit when the handler gets fired and calls server.Close()
	server.Serve(l)

	return token, nil
}

// getAccessToken trades the authorization code retrieved from the first OAuth2 leg for an access token
func getAccessToken(tokenEndpoint, clientID string, codeVerifier string, authorizationCode string, callbackURL string) (*oauth2.Token, error) {
	// set the url and form-encoded data for the POST to the access token endpoint
	data := fmt.Sprintf(
		"grant_type=authorization_code&client_id=%s"+
			"&code_verifier=%s"+
			"&code=%s"+
			"&redirect_uri=%s",
		clientID, codeVerifier, authorizationCode, callbackURL)
	payload := strings.NewReader(data)

	// create the request and execute it
	req, _ := http.NewRequest("POST", tokenEndpoint, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("snap: HTTP error: %s", err)
		return nil, err
	}

	// process the response
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	// unmarshal the json into a string map
	var token oauth2.Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		fmt.Printf("snap: JSON error: %s", err)
		return nil, err
	}

	return &token, nil
}

// cleanup closes the HTTP server
func cleanup(server *http.Server) {
	// we run this as a goroutine so that this function falls through and
	// the socket to the browser gets flushed/closed before the server goes away
	go server.Close()
}

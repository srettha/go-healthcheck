package oauth

import (
	"challenge/go-healthcheck/client"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"

	social "github.com/kkdai/line-login-sdk-go"
	"github.com/skratchdot/open-golang/open"
)

const (
	LOGIN_SCOPE  string = "profile openid"
	LOGIN_PROMPT string = "consent"
)

var accessToken string

// Cleanup closes the HTTP server
func Cleanup(server *http.Server) {
	// we run this as a goroutine so that this function falls through and
	// the socket to the browser gets flushed/closed before the server goes away
	go server.Close()
}

// Get access token after getting code from Login URL
func GetAccessToken(socialClient *client.SocialClient, redirectURL string, code string) (string, error) {
	token, err := socialClient.GetAccessToken(redirectURL, code)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Authorize user with Line Login API and redirect URL
func AuthorizeUser(socialClient *client.SocialClient, baseRedirectURL string) string {
	redirectURL := fmt.Sprintf("%s/oauth/callback", baseRedirectURL)

	// start a web server to listen on a callback URL
	server := &http.Server{Addr: redirectURL}

	// define a handler that will get the authorization code, call the token endpoint, and close the HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// get the authorization code
		code := req.URL.Query().Get("code")
		if code == "" {
			io.WriteString(w, "Error: could not find 'code' URL parameter\n")

			// close the HTTP server and return
			Cleanup(server)
			return
		}

		// trade the authorization code and the code verifier for an access token
		token, err := GetAccessToken(socialClient, redirectURL, code)
		if err != nil {
			io.WriteString(w, "Error: could not retrieve access token\n")

			// close the HTTP server and return
			Cleanup(server)
			return
		}

		accessToken = token

		// return an indication of success to the caller
		io.WriteString(w, `
		<html>
			<body>
				<h1>Login successful!</h1>
				<h2>You can close this window and return CLI.</h2>
			</body>
		</html>`)

		// close the HTTP server
		Cleanup(server)
	})

	// parse the redirect URL for the port number
	u, err := url.Parse(redirectURL)
	if err != nil {
		fmt.Printf("snap: bad redirect URL: %s\n", err)
		os.Exit(1)
	}

	// set up a listener on the redirect port
	port := fmt.Sprintf(":%s", u.Port())
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("snap: can't listen to port %s: %s\n", port, err)
		os.Exit(1)
	}

	// start the blocking web server loop
	// this will exit when the handler gets fired and calls server.Close()
	server.Serve(listener)

	return accessToken
}

// Get line login URL
func GetLineLoginURL(socialClient *client.SocialClient, baseRedirectURL string) string {
	redirectURL := fmt.Sprintf("%s/oauth/callback", baseRedirectURL)

	return socialClient.GetLineLoginURL(redirectURL, social.GenerateNonce(), LOGIN_SCOPE, client.SocialClientAuthOptions{Nonce: social.GenerateNonce(), Prompt: LOGIN_PROMPT})
}

// Login user by opening browser to login url
func LoginUser(socialClient *client.SocialClient, baseRedirectURL string) error {
	loginURL := GetLineLoginURL(socialClient, baseRedirectURL)

	return open.Start(loginURL)
}

package client

import (
	"net/http"
	"time"

	social "github.com/kkdai/line-login-sdk-go"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (*http.Response, error)
}

type SocialClientAuthOptions struct {
	Nonce  string
	Prompt string
}

type SocialClientGetAccessToken func(redirectURL string, code string) (string, error)
type SocialClientGetWebLoginURL func(redirectURL string, state string, scope string, options SocialClientAuthOptions) string

type SocialClient struct {
	GetAccessToken  SocialClientGetAccessToken
	GetLineLoginURL SocialClientGetWebLoginURL
}

// Get http client
func GetHttpClient() HttpClient {
	httpClient := &http.Client{
		Timeout: 15 * time.Second,
	}

	return httpClient
}

// Get social client
func GetSocialClient(channelID string, channelSecret string) *SocialClient {
	socialClient, _ := social.New(channelID, channelSecret)

	return &SocialClient{
		GetAccessToken: func(redirectURL string, code string) (string, error) {
			token, err := socialClient.GetAccessToken(redirectURL, code).Do()
			if err != nil {
				return "", err
			}

			return token.AccessToken, nil
		},
		GetLineLoginURL: func(redirectURL, state, scope string, options SocialClientAuthOptions) string {
			return socialClient.GetWebLoinURL(redirectURL, state, scope, social.AuthRequestOptions{Nonce: options.Nonce, Prompt: options.Prompt})
		},
	}
}

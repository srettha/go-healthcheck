package oauth

import (
	"challenge/go-healthcheck/client"
	"testing"
)

func TestGetAccessToken(t *testing.T) {
	redirectURL := "__REDIRECT_URL__"
	code := "__CODE__"

	t.Run("it should return access token", func(t *testing.T) {
		want := "__ACCESS_TOKEN__"
		mockSocialClient := &client.SocialClient{
			GetAccessToken: func(redirectURL string, code string) (string, error) {
				return want, nil
			},
		}

		_, err := GetAccessToken(mockSocialClient, redirectURL, code)
		if err != nil {
			t.Error("Expected no err, but there is err")
		}
	})
}

func TestGetLineLoginURL(t *testing.T) {
	t.Run("it should return LINE login URL", func(t *testing.T) {
		want := "__REDIRECT_URL__"
		mockSocialClient := &client.SocialClient{
			GetLineLoginURL: func(redirectURL string, state string, scope string, options client.SocialClientAuthOptions) string {
				return want
			},
		}

		got := GetLineLoginURL(mockSocialClient, "__BASE_REDIRECT_URL__")
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

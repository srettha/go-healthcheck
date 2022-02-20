package client

import (
	"testing"
)

func TestGetHttpClient(t *testing.T) {
	t.Run("it should return http client", func(t *testing.T) {
		got := GetHttpClient()
		if got == nil {
			t.Error("Expected http client, but nothing nothing")
		}
	})
}

func TestGetSocialClient(t *testing.T) {
	t.Run("it should return social client", func(t *testing.T) {
		channelID := "__CHANNEL_ID__"
		channelSecret := "__CHANNEL_SECRET__"

		got := GetSocialClient(channelID, channelSecret)
		if got == nil {
			t.Error("Expected social client, but got nothing")
		}
	})
}

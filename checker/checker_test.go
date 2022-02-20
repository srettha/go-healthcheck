package checker

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type MockDoType func(req *http.Request) (*http.Response, error)
type MockGetType func(url string) (*http.Response, error)

// MockClient is the mock client
type MockClient struct {
	MockDo  MockDoType
	MockGet MockGetType
}

// Overriding what the Do function should "do" in our MockClient
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}

func (m *MockClient) Get(url string) (*http.Response, error) {
	return m.MockGet(url)
}

func TestPingUrl(t *testing.T) {
	url := "__URL__"

	t.Run("it should return error if system fails to ping given url", func(t *testing.T) {
		json := `{"message": "success"}`
		r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
		mockClient := &MockClient{
			MockGet: func(url string) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       r,
				}, nil
			},
		}

		err := PingUrl(mockClient, url)
		if err != nil {
			t.Error("Expected err, but got nothing")
		}
	})

	t.Run("it should return true", func(t *testing.T) {
		mockClient := &MockClient{
			MockGet: func(url string) (*http.Response, error) {
				return &http.Response{
					StatusCode: 500,
					Body:       nil,
				}, errors.New("Something went wrong")
			},
		}

		err := PingUrl(mockClient, url)
		if err == nil {
			t.Error("Expected no err, but got err")
		}
	})
}

func TestFormatPingResult(t *testing.T) {
	totalWebsites := 2
	success := 1
	failure := 1
	totalTime := time.Duration(123450000)

	t.Run("it should return formatted ping result", func(t *testing.T) {
		want := Result{
			TotalWebsites: 2,
			Success:       1,
			Failure:       1,
			TotalTime:     123,
		}

		got := FormatPingResult(totalWebsites, success, failure, totalTime)
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

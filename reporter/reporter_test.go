package reporter

import (
	"bytes"
	"challenge/go-healthcheck/checker"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
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

func TestCreateRequest(t *testing.T) {
	accessToken := "__ACCESS_TOKEN__"
	pingResult := checker.Result{
		TotalWebsites: 2,
		Success:       1,
		Failure:       1,
		TotalTime:     1000,
	}

	t.Run("it should return request", func(t *testing.T) {
		got, err := CreateRequest(accessToken, pingResult)
		if err != nil {
			t.Error("Expected no err, but there is err")
		}

		if got.Header.Get("Authorization") != "Bearer __ACCESS_TOKEN__" {
			t.Errorf("got %q, want \"Bearer __ACCESS_TOKEN__\"", got.Header.Get("Authorization"))
		}
	})
}

func TestReport(t *testing.T) {
	accessToken := "__ACCESS_TOKEN__"
	pingResult := checker.Result{
		TotalWebsites: 2,
		Success:       1,
		Failure:       1,
		TotalTime:     1000,
	}

	t.Run("it should fail to send report to Healthcheck report API", func(t *testing.T) {
		mockClient := &MockClient{
			MockDo: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 404,
					Body:       nil,
				}, errors.New("Something went wrong")
			},
		}

		err := Report(mockClient, accessToken, pingResult)
		if err == nil {
			t.Error("Expected err, but there is no err")
		}
	})

	t.Run("it should send report to Healthcheck report API", func(t *testing.T) {
		json := `{"message": "success"}`
		r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
		mockClient := &MockClient{
			MockDo: func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       r,
				}, nil
			},
		}

		err := Report(mockClient, accessToken, pingResult)
		if err != nil {
			t.Error("Expected no err, but there is err")
		}
	})
}

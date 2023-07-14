package rdw_opendata_go

import (
	"context"
	"errors"
	"github.com/jarcoal/httpmock"
	"net/http"
	"testing"
)

func TestClient_Do_InvalidToken(t *testing.T) {
	t.Run("should return an error when the API token is invalid", func(t *testing.T) {
		setup()
		defer teardown()

		httpmock.RegisterResponder("GET", "https://opendata.rdw.nl/resource/m9d7-ebf2.json", func(request *http.Request) (*http.Response, error) {
			if request.Header.Get("X-App-Token") == "somevalidtoken" {
				return nil, errors.New("should not have a valid token")
			}

			resp, _ := httpmock.NewBytesResponder(403, loadFixture("403_invalid_token.json"))(request)
			return resp, nil
		})

		ctx := context.Background()

		client := NewClient(nil, "someinvalidtoken")
		req, _ := client.NewRequest(ctx, "GET", "/resource/m9d7-ebf2.json", nil, RegisteredVehiclesListOptions{})
		err := client.Do(req, nil)
		if err == nil {
			t.Error("should return an error")
		}

		expectedErr := &AppTokenInvalidError{}

		if err.Error() != expectedErr.Error() {
			t.Errorf("should return an error of type AppTokenInvalidError, returned: %v", err)
		}
	})
}

func TestClient_Do_EmptyToken(t *testing.T) {
	t.Run("should return an error when the API token is empty", func(t *testing.T) {
		setup()
		defer teardown()

		ctx := context.Background()

		client := NewClient(nil, "")
		registeredVehicles := NewRegisteredVehicles(client)

		_, err := registeredVehicles.List(ctx, RegisteredVehiclesListOptions{})
		if err == nil {
			t.Error("should return an error")
		}

		expectedErr := &AppTokenMissingError{}

		if err.Error() != expectedErr.Error() {
			t.Errorf("should return an error of type %v, returned: %v", expectedErr, err)
		}
	})
}

func TestClient_Do_ApiError(t *testing.T) {
	t.Run("should return an error when the API token is invalid", func(t *testing.T) {
		setup()
		defer teardown()

		httpmock.RegisterResponder("GET", "https://opendata.rdw.nl/resource/m9d7-ebf2.json", httpmock.NewStringResponder(500, "Internal server error"))

		ctx := context.Background()

		client := NewClient(nil, "somevalidtoken")

		req, _ := client.NewRequest(ctx, "GET", "/resource/m9d7-ebf2.json", nil, RegisteredVehiclesListOptions{})
		err := client.Do(req, nil)
		if err == nil {
			t.Error("should return an error")
		}

		expectedErr := &ApiError{
			StatusCode: 500,
			Message:    "Internal server error",
		}

		if err.Error() != expectedErr.Error() {
			t.Errorf("should return an error: %v, returned: %v", expectedErr, err)
		}
	})
}

func TestClient_Do_ApiLimitExceeded(t *testing.T) {
	t.Run("should return an error when the API limit is exceeded", func(t *testing.T) {
		setup()
		defer teardown()

		httpmock.RegisterResponder("GET", "https://opendata.rdw.nl/resource/m9d7-ebf2.json", func(request *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewBytesResponder(429, loadFixture("429_rate_limit_exceeded.json"))(request)
			return resp, nil
		})

		ctx := context.Background()

		client := NewClient(nil, "somevalidtoken")

		req, _ := client.NewRequest(ctx, "GET", "/resource/m9d7-ebf2.json", nil, RegisteredVehiclesListOptions{})
		err := client.Do(req, nil)
		if err == nil {
			t.Error("should return an error")
		}

		expectedErr := &ApiLimitExceededError{}

		if err.Error() != expectedErr.Error() {
			t.Errorf("should return an error, returned: %v", err)
		}
	})
}

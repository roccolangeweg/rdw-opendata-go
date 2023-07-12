package rdw_opendata_go

import (
	"context"
	"errors"
	"fmt"
	"github.com/jarcoal/httpmock"
	"net/http"
	"os"
	"testing"
)

func loadFixture(filename string) []byte {
	f, err := os.ReadFile("fixtures/" + filename)
	if err != nil {
		panic(fmt.Sprintf("Cannot load fixture %v", filename))
	}
	return f
}

func setup() {
	httpmock.Activate()
}

func teardown() {
	httpmock.DeactivateAndReset()
}

func TestRegisteredVehicles_ListAll(t *testing.T) {

	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", "https://opendata.rdw.nl/resource/m9d7-ebf2.json", func(request *http.Request) (*http.Response, error) {
		if request.URL.RawQuery != "" {
			return nil, errors.New("should not have query parameters")
		}
		resp, _ := httpmock.NewBytesResponder(200, loadFixture("registered_vehicles.json"))(request)
		return resp, nil
	})

	t.Run("should return a list of registered vehicles", func(t *testing.T) {
		ctx := context.Background()

		client := NewClient(nil, "somevalidtoken")
		registeredVehicles := NewRegisteredVehicles(client)

		vehicles, err := registeredVehicles.List(ctx, RegisteredVehiclesListOptions{})

		if err != nil {
			t.Error("should not return an error")
		}

		if len(vehicles) == 0 {
			t.Error("should return a list of vehicles")
		}

		if vehicles[0].LicensePlate != "AAAAAA" {
			t.Errorf("should return a vehicle with license plate AAAAAA, returned: %v", vehicles[0].LicensePlate)
		}
	})
}

func TestRegisteredVehicles_ListByLicensePlate(t *testing.T) {
	setup()
	defer teardown()

	httpmock.RegisterResponder("GET", "https://opendata.rdw.nl/resource/m9d7-ebf2.json", func(request *http.Request) (*http.Response, error) {
		if request.URL.Query().Get("kenteken") != "BBBBBB" {
			return nil, errors.New("should have license plate filter")
		}
		resp, _ := httpmock.NewBytesResponder(200, loadFixture("registered_vehicles_license_filter.json"))(request)
		return resp, nil
	})

	t.Run("should return a list of registered vehicles with a license plate filter", func(t *testing.T) {
		ctx := context.Background()

		client := NewClient(nil, "somevalidtoken")
		registeredVehicles := NewRegisteredVehicles(client)

		vehicles, err := registeredVehicles.List(ctx, RegisteredVehiclesListOptions{
			ListOptions: ListOptions{
				Limit: 1,
			},
			LicensePlate: "BBBBBB",
		})

		if err != nil {
			t.Errorf("should not return an error, returned: %v", err)
		}

		if len(vehicles) != 1 {
			t.Error("should return a list of vehicles")
		}

		if vehicles[0].LicensePlate != "BBBBBB" {
			t.Errorf("should return a vehicle with license plate R814ZT, returned: %v", vehicles[0].LicensePlate)
		}
	})
}

func TestRegisteredVehicles_InvalidToken(t *testing.T) {
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
		registeredVehicles := NewRegisteredVehicles(client)

		_, err := registeredVehicles.List(ctx, RegisteredVehiclesListOptions{})
		if err == nil {
			t.Error("should return an error")
		}

		if !errors.Is(err, &AppTokenInvalidError{}) {
			t.Errorf("should return an error of type AppTokenInvalidError, returned: %v", err)
		}
	})
}

func TestRegisteredVehicles_EmptyToken(t *testing.T) {
	t.Run("should return an error when the API token is invalid", func(t *testing.T) {
		setup()
		defer teardown()

		ctx := context.Background()

		client := NewClient(nil, "")
		registeredVehicles := NewRegisteredVehicles(client)

		_, err := registeredVehicles.List(ctx, RegisteredVehiclesListOptions{})
		if err == nil {
			t.Error("should return an error")
		}

		if !errors.Is(err, &AppTokenMissingError{}) {
			t.Errorf("should return an error of type AppTokenMissingError, returned: %v", err)
		}
	})
}

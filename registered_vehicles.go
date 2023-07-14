package rdw_opendata_go

import (
	"context"
	"fmt"
)

type RegisteredVehicles struct {
	client *Client
}

func NewRegisteredVehicles(client *Client) *RegisteredVehicles {
	return &RegisteredVehicles{
		client: client,
	}
}

type RegisteredVehiclesListOptions struct {
	ListOptions
	LicensePlate string `url:"kenteken,omitempty"`
	Brand        string `url:"merk,omitempty"`
	Model        string `url:"handelsbenaming,omitempty"`
	VehicleType  string `url:"voertuigsoort,omitempty"`
}

func (r *RegisteredVehicles) List(ctx context.Context, options RegisteredVehiclesListOptions) ([]RegisteredVehicle, error) {
	req, err := r.client.NewRequest(ctx, "GET", "/resource/m9d7-ebf2.json", nil, options)
	if err != nil {
		return nil, err
	}

	var registeredVehicles []RegisteredVehicle
	err = r.client.Do(req, &registeredVehicles)

	if err != nil {
		return nil, fmt.Errorf("error while requesting registered vehicles: %w", err)
	}

	return registeredVehicles, nil
}

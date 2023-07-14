package rdw_opendata_go

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

type ListOptions struct {
	Offset int `url:"$offset,omitempty"`
	Limit  int `url:"$limit,omitempty"`
}

type Client struct {
	client   *http.Client
	baseURL  *url.URL
	appToken string

	RegisteredVehicles *RegisteredVehicles
}

func NewClient(client *http.Client, appToken string) *Client {

	if client == nil {
		client = &http.Client{}
	}

	c := &Client{
		client:   client,
		appToken: appToken,
		baseURL: &url.URL{
			Scheme: "https",
			Host:   "opendata.rdw.nl",
		},
	}

	c.RegisteredVehicles = NewRegisteredVehicles(c)

	return c
}

func (c *Client) NewRequest(ctx context.Context, method, relPath string, body, options interface{}) (*http.Request, error) {
	rel, err := url.Parse(relPath)
	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(rel)

	// Add custom options
	if options != nil {
		optionsQuery, err := query.Values(options)
		if err != nil {
			return nil, err
		}

		u.RawQuery = optionsQuery.Encode()
	}

	// A bit of JSON ceremony
	var js []byte = nil

	if body != nil {
		js, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewBuffer(js))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "rdw-opendata-go")

	if c.appToken == "" {
		return nil, &AppTokenMissingError{}
	}
	req.Header.Add("X-App-Token", c.appToken)

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		if resp.StatusCode == http.StatusTooManyRequests {
			return &ApiLimitExceededError{}
		} else if resp.StatusCode == http.StatusForbidden {
			return &AppTokenInvalidError{}
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %w", err)
		}

		return &ApiError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	defer resp.Body.Close()

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			return err
		}
	}

	return nil
}

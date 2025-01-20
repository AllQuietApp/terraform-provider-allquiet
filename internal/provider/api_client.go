package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type AuthTransport struct {
	APIKey    string
	Transport http.RoundTripper
	BasicAuth *BasicAuth
}

type BasicAuth struct {
	Username string
	Password string
}

func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	if t.BasicAuth != nil {
		req.SetBasicAuth(t.BasicAuth.Username, t.BasicAuth.Password)
	}

	req.Header.Add("X-Authorization", t.APIKey)
	return t.Transport.RoundTrip(req)
}

type AllQuietAPIClient struct {
	APIKey      string
	EndpointURL string
	HTTPClient  *http.Client
}

func NewAllQuietAPIClient(apiKey, endpointURL string, basicAuth *BasicAuth) *AllQuietAPIClient {
	return &AllQuietAPIClient{
		APIKey:      apiKey,
		EndpointURL: endpointURL,
		HTTPClient: &http.Client{
			Transport: &AuthTransport{
				APIKey:    apiKey,
				Transport: http.DefaultTransport,
				BasicAuth: basicAuth,
			},
		},
	}
}

// newRequest creates a new HTTP request with the base URL and provided path.
func (c *AllQuietAPIClient) newRequest(method, path string, data interface{}) (*http.Request, error) {
	var buf bytes.Buffer
	if data != nil {
		err := json.NewEncoder(&buf).Encode(data)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, c.EndpointURL+path, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// post sends a POST request with the given data as JSON.
func (c *AllQuietAPIClient) post(ctx context.Context, path string, data interface{}) (*http.Response, error) {

	tflog.Trace(ctx, "%POST "+path)
	req, err := c.newRequest("POST", path, data)
	if err != nil {
		return nil, err
	}

	return c.HTTPClient.Do(req)
}

// post sends a PUT request with the given data as JSON.
func (c *AllQuietAPIClient) put(ctx context.Context, path string, data interface{}) (*http.Response, error) {
	tflog.Trace(ctx, "PUT "+path)
	req, err := c.newRequest("PUT", path, data)
	if err != nil {
		return nil, err
	}

	return c.HTTPClient.Do(req)
}

// post sends a POST request with the given data as JSON.
func (c *AllQuietAPIClient) get(ctx context.Context, path string) (*http.Response, error) {

	tflog.Trace(ctx, "GET "+path)

	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	return c.HTTPClient.Do(req)
}

// post sends a DELETE request with the given data as JSON.
func (c *AllQuietAPIClient) delete(ctx context.Context, path string) (*http.Response, error) {
	tflog.Trace(ctx, "DELETE "+path)
	req, err := c.newRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	return c.HTTPClient.Do(req)
}

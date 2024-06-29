package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type outboundIntegrationResponse struct {
	Id          string
	DisplayName string
	TeamId      string
	Type        string
}

type outboundIntegrationCreateRequest struct {
	DisplayName string `json:"displayName"`
	TeamId      string `json:"teamId"`
	Type        string `json:"type"`
}

func mapOutboundIntegrationCreateRequest(plan *OutboundIntegrationModel) *outboundIntegrationCreateRequest {
	return &outboundIntegrationCreateRequest{
		DisplayName: plan.DisplayName.ValueString(),
		TeamId:      plan.TeamId.ValueString(),
		Type:        plan.Type.ValueString(),
	}
}

func (c *AllQuietAPIClient) CreateOutboundIntegrationResource(ctx context.Context, data *OutboundIntegrationModel) (*outboundIntegrationResponse, error) {
	reqBody := mapOutboundIntegrationCreateRequest(data)

	url := "/outbound-integration"
	httpResp, err := c.post(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		logErrorResponse(httpResp)
		return nil, fmt.Errorf("non-200 response from API for POST %s: %d", url, httpResp.StatusCode)
	}

	var result outboundIntegrationResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) DeleteOutboundIntegrationResource(ctx context.Context, id string) error {
	url := fmt.Sprintf("/outbound-integration/%s", url.PathEscape(id))
	httpResp, err := c.delete(ctx, url)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		logErrorResponse(httpResp)
		return fmt.Errorf("non-200 response from API for DELETE %s: %d", url, httpResp.StatusCode)
	}

	return nil
}

func (c *AllQuietAPIClient) UpdateOutboundIntegrationResource(ctx context.Context, id string, data *OutboundIntegrationModel) (*outboundIntegrationResponse, error) {
	reqBody := mapOutboundIntegrationCreateRequest(data)

	url := fmt.Sprintf("/outbound-integration/%s", url.PathEscape(id))
	httpResp, err := c.put(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		logErrorResponse(httpResp)
		return nil, fmt.Errorf("non-200 response from API for PUT %s: %d", url, httpResp.StatusCode)
	}

	var result outboundIntegrationResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) GetOutboundIntegrationResource(ctx context.Context, id string) (*outboundIntegrationResponse, error) {
	url := fmt.Sprintf("/outbound-integration/%s", url.PathEscape(id))
	httpResp, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		logErrorResponse(httpResp)
		return nil, fmt.Errorf("non-200 response from API for GET %s: %d", url, httpResp.StatusCode)
	}

	var result outboundIntegrationResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

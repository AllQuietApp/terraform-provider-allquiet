package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type integrationResponse struct {
	Id              string
	DisplayName     string
	TeamId          string
	IsMuted         bool
	IsInMaintenance bool
	Type            string
}

type integrationCreateRequest struct {
	DisplayName     string `json:"displayName"`
	TeamId          string `json:"teamId"`
	IsMuted         bool   `json:"isMuted"`
	IsInMaintenance bool   `json:"isInMaintenance"`
	Type            string `json:"type"`
}

func mapIntegrationCreateRequest(plan *IntegrationModel) *integrationCreateRequest {
	return &integrationCreateRequest{
		DisplayName:     plan.DisplayName.ValueString(),
		TeamId:          plan.TeamId.ValueString(),
		IsMuted:         plan.IsMuted.ValueBool(),
		IsInMaintenance: plan.IsInMaintenance.ValueBool(),
		Type:            plan.Type.ValueString(),
	}
}

func (c *AllQuietAPIClient) CreateIntegrationResource(ctx context.Context, data *IntegrationModel) (*integrationResponse, error) {
	reqBody := mapIntegrationCreateRequest(data)

	url := "/inbound-integration"
	httpResp, err := c.post(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		logErrorResponse(httpResp)
		return nil, fmt.Errorf("non-200 response from API for POST %s: %d", url, httpResp.StatusCode)
	}

	var result integrationResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) DeleteIntegrationResource(ctx context.Context, id string) error {
	url := fmt.Sprintf("/inbound-integration/%s", url.PathEscape(id))
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

func (c *AllQuietAPIClient) UpdateIntegrationResource(ctx context.Context, id string, data *IntegrationModel) (*integrationResponse, error) {
	reqBody := mapIntegrationCreateRequest(data)

	url := fmt.Sprintf("/inbound-integration/%s", url.PathEscape(id))
	httpResp, err := c.put(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		logErrorResponse(httpResp)
		return nil, fmt.Errorf("non-200 response from API for PUT %s: %d", url, httpResp.StatusCode)
	}

	var result integrationResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) GetIntegrationResource(ctx context.Context, id string) (*integrationResponse, error) {
	url := fmt.Sprintf("/inbound-integration/%s", url.PathEscape(id))
	httpResp, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		logErrorResponse(httpResp)
		return nil, fmt.Errorf("non-200 response from API for GET %s: %d", url, httpResp.StatusCode)
	}

	var result integrationResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

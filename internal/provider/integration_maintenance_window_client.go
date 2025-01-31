package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type integrationMaintenanceWindowResponse struct {
	Id            string  `json:"id"`
	IntegrationId string  `json:"integrationId"`
	Start         *string `json:"start"`
	End           *string `json:"end"`
	Description   *string `json:"description"`
	Type          string  `json:"type"`
}

type integrationMaintenanceWindowCreateRequest struct {
	IntegrationId string  `json:"integrationId"`
	Start         *string `json:"start"`
	End           *string `json:"end"`
	Description   *string `json:"description"`
	Type          string  `json:"type"`
}

func (c *AllQuietAPIClient) CreateIntegrationMaintenanceWindowResource(ctx context.Context, data *IntegrationMaintenanceWindowModel) (*integrationMaintenanceWindowResponse, error) {
	reqBody := mapIntegrationMaintenanceWindowCreateRequest(data)

	url := "/inbound-integration-maintenance-windows"
	httpResp, err := c.post(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result integrationMaintenanceWindowResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func mapIntegrationMaintenanceWindowCreateRequest(plan *IntegrationMaintenanceWindowModel) *integrationMaintenanceWindowCreateRequest {
	var req integrationMaintenanceWindowCreateRequest
	req.IntegrationId = plan.IntegrationId.ValueString()
	req.Start = plan.Start.ValueStringPointer()
	req.End = plan.End.ValueStringPointer()
	req.Description = plan.Description.ValueStringPointer()
	req.Type = plan.Type.ValueString()

	return &req
}

func (c *AllQuietAPIClient) DeleteIntegrationMaintenanceWindowResource(ctx context.Context, id string) error {
	url := fmt.Sprintf("/inbound-integration-maintenance-windows/%s", url.PathEscape(id))
	httpResp, err := c.delete(ctx, url)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return logErrorResponse(httpResp, nil)
	}

	return nil
}

func (c *AllQuietAPIClient) UpdateIntegrationMaintenanceWindowResource(ctx context.Context, id string, data *IntegrationMaintenanceWindowModel) (*integrationMaintenanceWindowResponse, error) {
	reqBody := mapIntegrationMaintenanceWindowCreateRequest(data)

	url := fmt.Sprintf("/inbound-integration-maintenance-windows/%s", url.PathEscape(id))
	httpResp, err := c.put(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result integrationMaintenanceWindowResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) GetIntegrationMaintenanceWindowResource(ctx context.Context, id string) (*integrationMaintenanceWindowResponse, error) {
	url := fmt.Sprintf("/inbound-integration-maintenance-windows/%s", url.PathEscape(id))
	httpResp, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result integrationMaintenanceWindowResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

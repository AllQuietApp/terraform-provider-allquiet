package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type integrationResponse struct {
	Id              string                  `json:"id"`
	DisplayName     string                  `json:"displayName"`
	TeamId          string                  `json:"teamId"`
	IsMuted         bool                    `json:"isMuted"`
	IsInMaintenance bool                    `json:"isInMaintenance"`
	Type            string                  `json:"type"`
	WebhookUrl      *string                 `json:"webhookUrl"`
	SnoozeSettings  *snoozeSettingsResponse `json:"snoozeSettings"`
}

type integrationCreateRequest struct {
	DisplayName     string                  `json:"displayName"`
	TeamId          string                  `json:"teamId"`
	IsMuted         bool                    `json:"isMuted"`
	IsInMaintenance bool                    `json:"isInMaintenance"`
	Type            string                  `json:"type"`
	SnoozeSettings  *snoozeSettingsResponse `json:"snoozeSettings"`
}

type snoozeSettingsResponse struct {
	SnoozeWindowInMinutes *int64 `json:"snoozeWindowInMinutes"`
}

func mapIntegrationCreateRequest(plan *IntegrationModel) *integrationCreateRequest {
	return &integrationCreateRequest{
		DisplayName:     plan.DisplayName.ValueString(),
		TeamId:          plan.TeamId.ValueString(),
		IsMuted:         plan.IsMuted.ValueBool(),
		IsInMaintenance: plan.IsInMaintenance.ValueBool(),
		Type:            plan.Type.ValueString(),
		SnoozeSettings:  mapSnoozeSettingsCreateRequest(plan.SnoozeSettings),
	}
}

func mapSnoozeSettingsCreateRequest(plan *SnoozeSettingsModel) *snoozeSettingsResponse {
	if plan == nil {
		return nil
	}

	return &snoozeSettingsResponse{
		SnoozeWindowInMinutes: plan.SnoozeWindowInMinutes.ValueInt64Pointer(),
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
		return nil, logErrorResponse(httpResp, reqBody)
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
		return logErrorResponse(httpResp, nil)
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
		return nil, logErrorResponse(httpResp, reqBody)
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
		return nil, logErrorResponse(httpResp, nil)
	}

	var result integrationResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

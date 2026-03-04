package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type serviceResponse struct {
	Id                     string                  `json:"id"`
	DisplayName            string                  `json:"displayName"`
	PublicTitle            string                  `json:"publicTitle"`
	PublicDescription      *string                 `json:"publicDescription"`
	Templates              *[]serviceTemplate      `json:"templates"`
	Integrations           *[]serviceIntegration   `json:"integrations"`
	TeamConnectionSettings *teamConnectionSettings `json:"teamConnectionSettings"`
}

type serviceTemplate struct {
	Id          *string `json:"id"`
	DisplayName string  `json:"displayName"`
	Message     string  `json:"message"`
}

type serviceIntegration struct {
	Id            *string   `json:"id,omitempty"`
	IntegrationId string    `json:"integrationId"`
	Severities    *[]string `json:"severities"`
}

type serviceCreateRequest struct {
	DisplayName            string                  `json:"displayName"`
	PublicTitle            string                  `json:"publicTitle"`
	PublicDescription      *string                 `json:"publicDescription"`
	Templates              *[]serviceTemplate      `json:"templates"`
	Integrations           *[]serviceIntegration   `json:"integrations"`
	TeamConnectionSettings *teamConnectionSettings `json:"teamConnectionSettings"`
}

func mapServiceCreateRequest(plan *ServiceModel) *serviceCreateRequest {
	return &serviceCreateRequest{
		DisplayName:            plan.DisplayName.ValueString(),
		PublicTitle:            plan.PublicTitle.ValueString(),
		PublicDescription:      plan.PublicDescription.ValueStringPointer(),
		Templates:              mapTemplates(plan.Templates),
		Integrations:           mapServiceIntegrationsToRequest(plan.Integrations),
		TeamConnectionSettings: MapTeamConnectionSettingsToRequest(plan.TeamConnectionSettings),
	}
}

func mapServiceIntegrationsToRequest(integrations *[]ServiceIntegrationModel) *[]serviceIntegration {
	if integrations == nil {
		return nil
	}
	result := make([]serviceIntegration, len(*integrations))
	for i, si := range *integrations {
		idStr := si.Id.ValueString()
		var id *string
		if idStr != "" {
			id = &idStr
		}
		result[i] = serviceIntegration{
			Id:            id,
			IntegrationId: si.IntegrationId.ValueString(),
			Severities:    ListToStringArray(si.Severities),
		}
	}
	return &result
}

func mapTemplates(templates *[]ServiceTemplateModel) *[]serviceTemplate {
	if templates == nil {
		return nil
	}

	var result []serviceTemplate
	for _, template := range *templates {
		result = append(result, serviceTemplate{
			DisplayName: template.DisplayName.ValueString(),
			Message:     template.Message.ValueString(),
		})
	}

	return &result
}

func (c *AllQuietAPIClient) CreateServiceResource(ctx context.Context, data *ServiceModel) (*serviceResponse, error) {
	reqBody := mapServiceCreateRequest(data)

	url := "/service"
	httpResp, err := c.post(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result serviceResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) DeleteServiceResource(ctx context.Context, id string) error {
	url := fmt.Sprintf("/service/%s", url.PathEscape(id))
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

func (c *AllQuietAPIClient) UpdateServiceResource(ctx context.Context, id string, data *ServiceModel) (*serviceResponse, error) {
	reqBody := mapServiceCreateRequest(data)

	url := fmt.Sprintf("/service/%s", url.PathEscape(id))
	httpResp, err := c.put(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result serviceResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) GetServiceResource(ctx context.Context, id string) (*serviceResponse, error) {
	url := fmt.Sprintf("/service/%s", url.PathEscape(id))
	httpResp, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result serviceResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type serviceResponse struct {
	Id                string
	DisplayName       string
	PublicTitle       string
	PublicDescription *string
	Templates         *[]serviceTemplate
}

type serviceTemplate struct {
	Id          *string
	DisplayName string
	Message     string
}

type serviceCreateRequest struct {
	DisplayName       string             `json:"displayName"`
	PublicTitle       string             `json:"publicTitle"`
	PublicDescription *string            `json:"publicDescription"`
	Templates         *[]serviceTemplate `json:"templates"`
}

func mapServiceCreateRequest(plan *ServiceModel) *serviceCreateRequest {
	return &serviceCreateRequest{
		DisplayName:       plan.DisplayName.ValueString(),
		PublicTitle:       plan.PublicTitle.ValueString(),
		PublicDescription: plan.PublicDescription.ValueStringPointer(),
		Templates:         mapTemplates(plan.Templates),
	}
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

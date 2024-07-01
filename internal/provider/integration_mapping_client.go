package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type integrationMappingResponse struct {
	Id                string                      `json:"id"`
	IntegrationId     string                      `json:"integrationId"`
	AttributesMapping integrationAttributeMapping `json:"attributesMapping"`
}

type integrationMappingCreateRequest struct {
	AttributesMapping integrationAttributeMapping `json:"attributesMapping"`
}

type integrationAttributeMapping struct {
	Attributes []attribute `json:"attributes"`
}

type attribute struct {
	Name     string    `json:"name"`
	Mappings []mapping `json:"mappings"`
}

type mapping struct {
	XPath    *string `json:"xPath"`
	JSONPath *string `json:"jsonPath"`
	Regex    *string `json:"regex"`
	Replace  *string `json:"replace"`
	Map      *string `json:"map"`
	Static   *string `json:"static"`
}

func (c *AllQuietAPIClient) CreateIntegrationMappingResource(ctx context.Context, data *IntegrationMappingModel) (*integrationMappingResponse, error) {
	reqBody := mapIntegrationMappingCreateRequest(data)

	url := fmt.Sprintf("/inbound-integration/%s/mapping", url.PathEscape(data.IntegrationId.ValueString()))
	httpResp, err := c.post(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp)
	}

	var result integrationMappingResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func mapIntegrationMappingCreateRequest(plan *IntegrationMappingModel) *integrationMappingCreateRequest {
	var req integrationMappingCreateRequest
	req.AttributesMapping.Attributes = make([]attribute, len(plan.AttributesMapping.Attributes))

	for i, attribute := range plan.AttributesMapping.Attributes {
		req.AttributesMapping.Attributes[i].Name = attribute.Name.ValueString()
		req.AttributesMapping.Attributes[i].Mappings = make([]mapping, len(attribute.Mappings))

		for j, m := range attribute.Mappings {
			req.AttributesMapping.Attributes[i].Mappings[j] = mapping{
				XPath:    m.XPath.ValueStringPointer(),
				JSONPath: m.JSONPath.ValueStringPointer(),
				Regex:    m.Regex.ValueStringPointer(),
				Replace:  m.Replace.ValueStringPointer(),
				Map:      m.Map.ValueStringPointer(),
				Static:   m.Static.ValueStringPointer(),
			}
		}
	}

	return &req
}

func (c *AllQuietAPIClient) DeleteIntegrationMappingResource(ctx context.Context, id string) error {
	url := fmt.Sprintf("/inbound-integration/%s/mapping", url.PathEscape(id))
	httpResp, err := c.delete(ctx, url)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return logErrorResponse(httpResp)
	}

	return nil
}

func (c *AllQuietAPIClient) UpdateIntegrationMappingResource(ctx context.Context, id string, data *IntegrationMappingModel) (*integrationMappingResponse, error) {
	reqBody := mapIntegrationMappingCreateRequest(data)

	url := fmt.Sprintf("/inbound-integration/%s/mapping", url.PathEscape(data.IntegrationId.ValueString()))
	httpResp, err := c.put(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp)
	}

	var result integrationMappingResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) GetIntegrationMappingResource(ctx context.Context, id string) (*integrationMappingResponse, error) {
	url := fmt.Sprintf("/inbound-integration/%s/mapping", url.PathEscape(id))
	httpResp, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp)
	}

	var result integrationMappingResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

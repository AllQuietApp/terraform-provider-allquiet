package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type integrationMappingResponse struct {
	AttributesMapping integrationAttributeMappingResponse
}

type integrationAttributeMappingResponse struct {
	Attributes []attributeResponse
}

type attributeResponse struct {
	Name     string
	Mappings []mappingResponse
}

type mappingResponse struct {
	XPath    *string
	JSONPath *string
	Regex    *string
	Replace  *string
	Map      *string
	Static   *string
}

type integrationMappingCreateRequest struct {
	AttributesMapping attributesMappingCreateRequest `json:"attributesMapping"`
}

type attributesMappingCreateRequest struct {
	Attributes []attributeCreateRequest `json:"attributes"`
}

type attributeCreateRequest struct {
	Name     string                 `json:"name"`
	Mappings []mappingCreateRequest `json:"mappings"`
}

type mappingCreateRequest struct {
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
		logErrorResponse(httpResp)
		return nil, fmt.Errorf("non-200 response from API for POST %s: %d", url, httpResp.StatusCode)
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
	req.AttributesMapping.Attributes = make([]attributeCreateRequest, len(plan.AttributesMapping.Attributes))

	for i, attribute := range plan.AttributesMapping.Attributes {
		req.AttributesMapping.Attributes[i].Name = attribute.Name.ValueString()
		req.AttributesMapping.Attributes[i].Mappings = make([]mappingCreateRequest, len(attribute.Mappings))

		for j, mapping := range attribute.Mappings {
			req.AttributesMapping.Attributes[i].Mappings[j] = mappingCreateRequest{
				XPath:    mapping.XPath.ValueStringPointer(),
				JSONPath: mapping.JSONPath.ValueStringPointer(),
				Regex:    mapping.Regex.ValueStringPointer(),
				Replace:  mapping.Replace.ValueStringPointer(),
				Map:      mapping.Map.ValueStringPointer(),
				Static:   mapping.Static.ValueStringPointer(),
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
		logErrorResponse(httpResp)
		return fmt.Errorf("non-200 response from API for DELETE %s: %d", url, httpResp.StatusCode)
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
		logErrorResponse(httpResp)
		return nil, fmt.Errorf("non-200 response from API for PUT %s: %d", url, httpResp.StatusCode)
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
		logErrorResponse(httpResp)
		return nil, fmt.Errorf("non-200 response from API for GET %s: %d", url, httpResp.StatusCode)
	}

	var result integrationMappingResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

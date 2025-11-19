package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type onCallOverrideResponse struct {
	Id                 string    `json:"id"`
	UserId             string    `json:"userId"`
	TeamId             *string   `json:"teamId"`
	Type               string    `json:"type"`
	Start              string    `json:"start"`
	End                string    `json:"end"`
	ReplacementUserIds *[]string `json:"replacementUserIds"`
}

type onCallOverrideCreateRequest struct {
	UserId             string    `json:"userId"`
	TeamId             *string   `json:"teamId,omitempty"`
	Type               string    `json:"type"`
	Start              string    `json:"start"`
	End                string    `json:"end"`
	ReplacementUserIds *[]string `json:"replacementUserIds"`
}

func mapOnCallOverrideCreateRequest(plan *OnCallOverrideModel) *onCallOverrideCreateRequest {
	var req onCallOverrideCreateRequest

	req.UserId = plan.UserId.ValueString()
	req.TeamId = plan.TeamId.ValueStringPointer()
	req.Type = plan.Type.ValueString()
	req.Start = plan.Start.ValueString()
	req.End = plan.End.ValueString()
	req.ReplacementUserIds = ListToStringArray(plan.ReplacementUserIds)

	return &req
}

func (c *AllQuietAPIClient) CreateOnCallOverrideResource(ctx context.Context, data *OnCallOverrideModel) (*onCallOverrideResponse, error) {
	reqBody := mapOnCallOverrideCreateRequest(data)

	url := "/on-call-override"
	httpResp, err := c.post(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result onCallOverrideResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) DeleteOnCallOverrideResource(ctx context.Context, id string) error {
	url := fmt.Sprintf("/on-call-override/%s", url.PathEscape(id))
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

func (c *AllQuietAPIClient) UpdateOnCallOverrideResource(ctx context.Context, id string, data *OnCallOverrideModel) (*onCallOverrideResponse, error) {
	reqBody := mapOnCallOverrideCreateRequest(data)

	url := fmt.Sprintf("/on-call-override/%s", url.PathEscape(id))
	httpResp, err := c.put(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result onCallOverrideResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) GetOnCallOverrideResource(ctx context.Context, id string) (*onCallOverrideResponse, error) {
	url := fmt.Sprintf("/on-call-override/%s", url.PathEscape(id))
	httpResp, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result onCallOverrideResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

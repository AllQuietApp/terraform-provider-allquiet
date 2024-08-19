package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type teamMembershipResponse struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
	Role   string `json:"role"`
	TeamId string `json:"teamId"`
}

type teamMembershipCreateRequest struct {
	UserId string `json:"userId"`
	Role   string `json:"role"`
	TeamId string `json:"teamId"`
}

func mapTeamMembershipCreateRequest(plan *TeamMembershipModel) *teamMembershipCreateRequest {
	var req teamMembershipCreateRequest

	req.Role = plan.Role.ValueString()
	req.TeamId = plan.TeamId.ValueString()
	req.UserId = plan.UserId.ValueString()

	return &req
}

func (c *AllQuietAPIClient) CreateTeamMembershipResource(ctx context.Context, data *TeamMembershipModel) (*teamMembershipResponse, error) {
	reqBody := mapTeamMembershipCreateRequest(data)

	url := "/team-membership"
	httpResp, err := c.post(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp)
	}

	var result teamMembershipResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) DeleteTeamMembershipResource(ctx context.Context, id string) error {
	url := fmt.Sprintf("/team-membership/%s", url.PathEscape(id))
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

func (c *AllQuietAPIClient) UpdateTeamMembershipResource(ctx context.Context, id string, data *TeamMembershipModel) (*teamMembershipResponse, error) {
	reqBody := mapTeamMembershipCreateRequest(data)

	url := fmt.Sprintf("/team-membership/%s", url.PathEscape(id))
	httpResp, err := c.put(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp)
	}

	var result teamMembershipResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) GetTeamMembershipResource(ctx context.Context, id string) (*teamMembershipResponse, error) {
	url := fmt.Sprintf("/team-membership/%s", url.PathEscape(id))
	httpResp, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp)
	}

	var result teamMembershipResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

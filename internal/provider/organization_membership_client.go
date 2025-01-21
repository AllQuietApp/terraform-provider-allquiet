package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type organizationMembershipResponse struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
	Role   string `json:"role"`
}

type organizationMembershipCreateRequest struct {
	UserId string `json:"userId"`
	Role   string `json:"role"`
}

func mapOrganizationMembershipCreateRequest(plan *OrganizationMembershipModel) *organizationMembershipCreateRequest {
	var req organizationMembershipCreateRequest

	req.Role = plan.Role.ValueString()
	req.UserId = plan.UserId.ValueString()

	return &req
}

func (c *AllQuietAPIClient) CreateOrganizationMembershipResource(ctx context.Context, data *OrganizationMembershipModel) (*organizationMembershipResponse, error) {
	reqBody := mapOrganizationMembershipCreateRequest(data)

	url := "/organization-membership"
	httpResp, err := c.post(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result organizationMembershipResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) DeleteOrganizationMembershipResource(ctx context.Context, id string) error {
	url := fmt.Sprintf("/organization-membership/%s", url.PathEscape(id))
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

func (c *AllQuietAPIClient) UpdateOrganizationMembershipResource(ctx context.Context, id string, data *OrganizationMembershipModel) (*organizationMembershipResponse, error) {
	reqBody := mapOrganizationMembershipCreateRequest(data)

	url := fmt.Sprintf("/organization-membership/%s", url.PathEscape(id))
	httpResp, err := c.put(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result organizationMembershipResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) GetOrganizationMembershipResource(ctx context.Context, id string) (*organizationMembershipResponse, error) {
	url := fmt.Sprintf("/organization-membership/%s", url.PathEscape(id))
	httpResp, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result organizationMembershipResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

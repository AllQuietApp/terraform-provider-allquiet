package provider

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type teamMembershipsDataSourceResponse struct {
	TeamMemberships []teamMembershipDataSourceResponse `json:"teamMemberships"`
}

func (c *AllQuietAPIClient) GetTeamMembershipsDataSource(ctx context.Context, teamMembershipsDataSource *TeamMembershipsDataSourceModel, diagnostics *diag.Diagnostics) (*teamMembershipsDataSourceResponse, error) {

	url := getTeamMembershipsUrl(teamMembershipsDataSource, diagnostics)
	if url == nil {
		return nil, nil
	}

	httpResp, err := c.get(ctx, *url)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result teamMembershipsDataSourceResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func getTeamMembershipsUrl(teamMembershipsDataSource *TeamMembershipsDataSourceModel, diagnostics *diag.Diagnostics) *string {

	url := "/team-membership/search/list"
	if teamMembershipsDataSource.UserId.ValueStringPointer() != nil {
		url = AddQueryParam(url, "userId", teamMembershipsDataSource.UserId.ValueString())
	}

	if teamMembershipsDataSource.TeamId.ValueStringPointer() != nil {
		url = AddQueryParam(url, "teamId", teamMembershipsDataSource.TeamId.ValueString())
	}

	if teamMembershipsDataSource.Role.ValueStringPointer() != nil {
		url = AddQueryParam(url, "role", teamMembershipsDataSource.Role.ValueString())
	}

	return &url
}

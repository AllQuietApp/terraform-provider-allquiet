package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type teamMembershipDataSourceResponse struct {
	Id          string `json:"id"`
	UserId      string `json:"userId"`
	TeamId      string `json:"teamId"`
	Role        string `json:"role"`
	ActivatedAt string `json:"activatedAt"`
}

func (c *AllQuietAPIClient) GetTeamMembershipDataSource(ctx context.Context, teamMembershipDataSource *TeamMembershipDataSourceModel, diagnostics *diag.Diagnostics) (*teamMembershipDataSourceResponse, error) {

	url := getTeamMembershipUrl(teamMembershipDataSource)
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

	var result teamMembershipDataSourceResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func getTeamMembershipUrl(teamMembershipDataSource *TeamMembershipDataSourceModel) *string {

	if teamMembershipDataSource.Id.ValueStringPointer() != nil {
		url := fmt.Sprintf("/team-membership/%s", url.PathEscape(teamMembershipDataSource.Id.ValueString()))
		return &url
	}

	url := "/team-membership/search"
	if teamMembershipDataSource.UserId.ValueStringPointer() != nil {
		url = AddQueryParam(url, "userId", teamMembershipDataSource.UserId.ValueString())
	}

	if teamMembershipDataSource.TeamId.ValueStringPointer() != nil {
		url = AddQueryParam(url, "teamId", teamMembershipDataSource.TeamId.ValueString())
	}

	if teamMembershipDataSource.Role.ValueStringPointer() != nil {
		url = AddQueryParam(url, "role", teamMembershipDataSource.Role.ValueString())
	}

	return &url
}

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type teamDataSourceResponse struct {
	Id          string    `json:"id"`
	DisplayName string    `json:"displayName"`
	TimeZoneId  string    `json:"timeZoneId"`
	Labels      *[]string `json:"labels"`
}

func (c *AllQuietAPIClient) GetTeamDataSource(ctx context.Context, teamDataSource *TeamDataSourceModel, diagnostics *diag.Diagnostics) (*teamDataSourceResponse, error) {

	url := getTeamUrl(teamDataSource, diagnostics)
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

	var result teamDataSourceResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func getTeamUrl(teamDataSource *TeamDataSourceModel, diagnostics *diag.Diagnostics) *string {

	if teamDataSource.Id.ValueStringPointer() != nil {
		url := fmt.Sprintf("/team/%s", url.PathEscape(teamDataSource.Id.ValueString()))
		return &url
	}

	if teamDataSource.DisplayName.ValueStringPointer() != nil {
		url := fmt.Sprintf("/team/search?displayName=%s", url.QueryEscape(teamDataSource.DisplayName.ValueString()))
		return &url
	}

	diagnostics.AddError("Client Error", "You need to provide either an id or display name to look up a team")
	return nil
}

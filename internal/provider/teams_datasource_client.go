package provider

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type teamsDataSourceResponse struct {
	Teams []teamDataSourceResponse `json:"teams"`
}

func (c *AllQuietAPIClient) GetTeamsDataSource(ctx context.Context, teamsDataSource *TeamsDataSourceModel, diagnostics *diag.Diagnostics) (*teamsDataSourceResponse, error) {

	url := getTeamsUrl(teamsDataSource, diagnostics)
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

	var result teamsDataSourceResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func getTeamsUrl(teamsDataSource *TeamsDataSourceModel, diagnostics *diag.Diagnostics) *string {

	url := "/team/search/list"
	if teamsDataSource.DisplayName.ValueStringPointer() != nil {
		url = AddQueryParam(url, "displayName", teamsDataSource.DisplayName.ValueString())
	}

	return &url
}

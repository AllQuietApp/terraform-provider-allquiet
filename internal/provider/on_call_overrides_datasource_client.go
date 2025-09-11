package provider

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type onCallOverridesDataSourceResponse struct {
	OnCallOverrides []onCallOverrideDataSourceResponse `json:"onCallOverrides"`
}

type onCallOverrideDataSourceResponse struct {
	Id                 string    `json:"id"`
	UserId             string    `json:"userId"`
	Type               string    `json:"type"`
	Start              string    `json:"start"`
	End                string    `json:"end"`
	ReplacementUserIds *[]string `json:"replacementUserIds"`
}

func (c *AllQuietAPIClient) GetOnCallOverridesDataSource(ctx context.Context, onCallOverridesDataSource *OnCallOverridesDataSourceModel, diagnostics *diag.Diagnostics) (*onCallOverridesDataSourceResponse, error) {

	url := getOnCallOverridesUrl(onCallOverridesDataSource)
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

	var result onCallOverridesDataSourceResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func getOnCallOverridesUrl(onCallOverridesDataSource *OnCallOverridesDataSourceModel) *string {

	url := "/on-call-override/search/list"
	if onCallOverridesDataSource.UserId.ValueStringPointer() != nil {
		url = AddQueryParam(url, "userId", onCallOverridesDataSource.UserId.ValueString())
	}

	return &url
}

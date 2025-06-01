package provider

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type usersDataSourceResponse struct {
	Users []userDataSourceResponse `json:"users"`
}

func (c *AllQuietAPIClient) GetUsersDataSource(ctx context.Context, data *UsersDataSourceModel, diagnostics *diag.Diagnostics) (*usersDataSourceResponse, error) {

	url := "/user/search/list"
	if data.Email.ValueString() != "" {
		url = AddQueryParam(url, "email", data.Email.ValueString())
	}
	if data.DisplayName.ValueString() != "" {
		url = AddQueryParam(url, "displayName", data.DisplayName.ValueString())
	}

	httpResp, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result usersDataSourceResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type userDataSourceResponse struct {
	Id             string `json:"id"`
	DisplayName    string `json:"displayName"`
	Email          string `json:"email"`
	ScimExternalId string `json:"scimExternalId"`
}

func (c *AllQuietAPIClient) GetUserDataSource(ctx context.Context, userDataSource *UserDataSourceModel, diagnostics *diag.Diagnostics) (*userDataSourceResponse, error) {

	url := getUserUrl(userDataSource, diagnostics)
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

	var result userDataSourceResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func getUserUrl(userDataSource *UserDataSourceModel, diagnostics *diag.Diagnostics) *string {

	if userDataSource.Id.ValueStringPointer() != nil {
		url := fmt.Sprintf("/user/%s", url.PathEscape(userDataSource.Id.ValueString()))
		return &url
	}

	if userDataSource.Email.ValueStringPointer() != nil {
		url := fmt.Sprintf("/user/search?email=%s", url.QueryEscape(userDataSource.Email.ValueString()))
		return &url
	}

	if userDataSource.DisplayName.ValueStringPointer() != nil {
		url := fmt.Sprintf("/user/search?displayName=%s", url.QueryEscape(userDataSource.DisplayName.ValueString()))
		return &url
	}

	if userDataSource.ScimExternalId.ValueStringPointer() != nil {
		url := fmt.Sprintf("/user/search?scimExternalId=%s", url.QueryEscape(userDataSource.ScimExternalId.ValueString()))
		return &url
	}

	diagnostics.AddError("Client Error", "You need to provide either an id, email, or display name to look up a user")
	return nil
}

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type statusPageResponse struct {
	Id                            string
	DisplayName                   string
	PublicTitle                   string
	PublicDescription             *string
	Slug                          *string
	ServiceIds                    *[]string
	PublicCompanyUrl              *string
	PublicCompanyName             *string
	PublicSupportUrl              *string
	PublicSupportEmail            *string
	HistoryInDays                 int64
	TimeZoneId                    *string
	DisablePublicSubscription     bool
	PublicSeverityMappingMinor    *string
	PublicSeverityMappingWarning  *string
	PublicSeverityMappingCritical *string
	BannerBackgroundColor         *string
	BannerBackgroundColorDarkMode *string
	BannerTextColor               *string
	BannerTextColorDarkMode       *string
	CustomHostSettings            *customHostSettingsResponse
}

type customHostSettingsResponse struct {
	Host string `json:"host"`
}

type statusPageCreateRequest struct {
	DisplayName                   string                     `json:"displayName"`
	PublicTitle                   string                     `json:"publicTitle"`
	PublicDescription             *string                    `json:"publicDescription"`
	Slug                          *string                    `json:"slug"`
	ServiceIds                    *[]string                  `json:"serviceIds"`
	PublicCompanyUrl              *string                    `json:"publicCompanyUrl"`
	PublicCompanyName             *string                    `json:"publicCompanyName"`
	PublicSupportUrl              *string                    `json:"publicSupportUrl"`
	PublicSupportEmail            *string                    `json:"publicSupportEmail"`
	HistoryInDays                 int64                      `json:"historyInDays"`
	TimeZoneId                    *string                    `json:"timeZoneId"`
	DisablePublicSubscription     bool                       `json:"disablePublicSubscription"`
	PublicSeverityMappingMinor    *string                    `json:"publicSeverityMappingMinor"`
	PublicSeverityMappingWarning  *string                    `json:"publicSeverityMappingWarning"`
	PublicSeverityMappingCritical *string                    `json:"publicSeverityMappingCritical"`
	BannerBackgroundColor         *string                    `json:"bannerBackgroundColor"`
	BannerBackgroundColorDarkMode *string                    `json:"bannerBackgroundColorDarkMode"`
	BannerTextColor               *string                    `json:"bannerTextColor"`
	BannerTextColorDarkMode       *string                    `json:"bannerTextColorDarkMode"`
	CustomHostSettings            *customHostSettingsRequest `json:"customHostSettings"`
}

type customHostSettingsRequest struct {
	Host string `json:"host"`
}

func mapStatusPageCreateRequest(plan *StatusPageModel) *statusPageCreateRequest {
	return &statusPageCreateRequest{
		DisplayName:                   plan.DisplayName.ValueString(),
		PublicTitle:                   plan.PublicTitle.ValueString(),
		PublicDescription:             plan.PublicDescription.ValueStringPointer(),
		Slug:                          plan.Slug.ValueStringPointer(),
		ServiceIds:                    ListToStringArray(plan.Services),
		PublicCompanyUrl:              plan.PublicCompanyUrl.ValueStringPointer(),
		PublicCompanyName:             plan.PublicCompanyName.ValueStringPointer(),
		PublicSupportUrl:              plan.PublicSupportUrl.ValueStringPointer(),
		PublicSupportEmail:            plan.PublicSupportEmail.ValueStringPointer(),
		HistoryInDays:                 plan.HistoryInDays.ValueInt64(),
		TimeZoneId:                    plan.TimeZoneId.ValueStringPointer(),
		DisablePublicSubscription:     plan.DisablePublicSubscription.ValueBool(),
		PublicSeverityMappingMinor:    plan.PublicSeverityMappingMinor.ValueStringPointer(),
		PublicSeverityMappingWarning:  plan.PublicSeverityMappingWarning.ValueStringPointer(),
		PublicSeverityMappingCritical: plan.PublicSeverityMappingCritical.ValueStringPointer(),
		BannerBackgroundColor:         plan.BannerBackgroundColor.ValueStringPointer(),
		BannerBackgroundColorDarkMode: plan.BannerBackgroundColorDarkMode.ValueStringPointer(),
		BannerTextColor:               plan.BannerTextColor.ValueStringPointer(),
		BannerTextColorDarkMode:       plan.BannerTextColorDarkMode.ValueStringPointer(),
		CustomHostSettings:            mapCustomHostSettingsRequestToModel(plan.CustomHostSettings),
	}
}

func mapCustomHostSettingsRequestToModel(request *CustomHostSettings) *customHostSettingsRequest {
	if request == nil {
		return nil
	}

	return &customHostSettingsRequest{
		Host: request.Host.ValueString(),
	}
}

func (c *AllQuietAPIClient) CreateStatusPageResource(ctx context.Context, data *StatusPageModel) (*statusPageResponse, error) {
	reqBody := mapStatusPageCreateRequest(data)

	url := "/status-page"
	httpResp, err := c.post(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result statusPageResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) DeleteStatusPageResource(ctx context.Context, id string) error {
	url := fmt.Sprintf("/status-page/%s", url.PathEscape(id))
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

func (c *AllQuietAPIClient) UpdateStatusPageResource(ctx context.Context, id string, data *StatusPageModel) (*statusPageResponse, error) {
	reqBody := mapStatusPageCreateRequest(data)

	url := fmt.Sprintf("/status-page/%s", url.PathEscape(id))
	httpResp, err := c.put(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result statusPageResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) GetStatusPageResource(ctx context.Context, id string) (*statusPageResponse, error) {
	url := fmt.Sprintf("/status-page/%s", url.PathEscape(id))
	httpResp, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result statusPageResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

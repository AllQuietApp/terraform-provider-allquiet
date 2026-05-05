package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type userResponse struct {
	Id                           string                                `json:"id"`
	DisplayName                  string                                `json:"displayName"`
	Email                        string                                `json:"email"`
	PhoneNumber                  *string                               `json:"phoneNumber"`
	TimeZoneId                   string                                `json:"timeZoneId"`
	IncidentNotificationSettings *incidentNotificationSettingsResponse `json:"incidentNotificationSettings"`
}

type incidentNotificationSettingsResponse struct {
	ShouldSendSMS bool      `json:"shouldSendSms"`
	DelayInMinSMS int64     `json:"delayInMinSms"`
	SeveritiesSMS *[]string `json:"severitiesSms"`

	ShouldCallVoice bool      `json:"shouldCallVoice"`
	DelayInMinVoice int64     `json:"delayInMinVoice"`
	SeveritiesVoice *[]string `json:"severitiesVoice"`

	ShouldSendPush bool      `json:"shouldSendPush"`
	DelayInMinPush int64     `json:"delayInMinPush"`
	SeveritiesPush *[]string `json:"severitiesPush"`

	ShouldSendEmail bool      `json:"shouldSendEmail"`
	DelayInMinEmail int64     `json:"delayInMinEmail"`
	SeveritiesEmail *[]string `json:"severitiesEmail"`

	DisabledIntentsEmail *[]string `json:"disabledIntentsEmail"`
	DisabledIntentsVoice *[]string `json:"disabledIntentsVoice"`
	DisabledIntentsPush  *[]string `json:"disabledIntentsPush"`
	DisabledIntentsSMS   *[]string `json:"disabledIntentsSMS"`
}

type userCreateRequest struct {
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	TimeZoneId  string `json:"timeZoneId"`
}

func mapUserCreateRequest(plan *UserModel) *userCreateRequest {
	return &userCreateRequest{
		DisplayName: plan.DisplayName.ValueString(),
		Email:       plan.Email.ValueString(),
		TimeZoneId:  plan.TimeZoneId.ValueString(),
	}
}

func (c *AllQuietAPIClient) CreateUserResource(ctx context.Context, data *UserModel) (*userResponse, error) {
	reqBody := mapUserCreateRequest(data)

	url := "/user"
	httpResp, err := c.post(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result userResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) DeleteUserResource(ctx context.Context, id string) error {
	url := fmt.Sprintf("/user/%s", url.PathEscape(id))
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

func (c *AllQuietAPIClient) UpdateUserResource(ctx context.Context, id string, data *UserModel) (*userResponse, error) {
	reqBody := mapUserCreateRequest(data)

	url := fmt.Sprintf("/user/%s", url.PathEscape(id))
	httpResp, err := c.put(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result userResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) GetUserResource(ctx context.Context, id string) (*userResponse, error) {
	url := fmt.Sprintf("/user/%s", url.PathEscape(id))
	httpResp, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result userResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

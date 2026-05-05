package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type userIncidentNotificationSettingsResponse struct {
	UserId string `json:"userId"`

	PhoneNumber *string `json:"phoneNumber"`

	ShouldSendSMS *bool     `json:"shouldSendSMS"`
	DelayInMinSMS *int64    `json:"delayInMinSMS"`
	SeveritiesSMS *[]string `json:"severitiesSMS"`

	ShouldCallVoice *bool     `json:"shouldCallVoice"`
	DelayInMinVoice *int64    `json:"delayInMinVoice"`
	SeveritiesVoice *[]string `json:"severitiesVoice"`

	ShouldSendPush *bool     `json:"shouldSendPush"`
	DelayInMinPush *int64    `json:"delayInMinPush"`
	SeveritiesPush *[]string `json:"severitiesPush"`

	ShouldSendEmail *bool     `json:"shouldSendEmail"`
	DelayInMinEmail *int64    `json:"delayInMinEmail"`
	SeveritiesEmail *[]string `json:"severitiesEmail"`

	DisabledIntentsEmail *[]string `json:"disabledIntentsEmail"`
	DisabledIntentsVoice *[]string `json:"disabledIntentsVoice"`
	DisabledIntentsPush  *[]string `json:"disabledIntentsPush"`
	DisabledIntentsSMS   *[]string `json:"disabledIntentsSMS"`
}

type userIncidentNotificationSettingsRequest struct {
	PhoneNumber *string `json:"phoneNumber"`

	ShouldSendSMS bool      `json:"shouldSendSMS"`
	DelayInMinSMS int64     `json:"delayInMinSMS"`
	SeveritiesSMS *[]string `json:"severitiesSMS"`

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

func mapUserIncidentNotificationSettingsRequest(plan *UserIncidentNotificationSettingsModel) *userIncidentNotificationSettingsRequest {
	return &userIncidentNotificationSettingsRequest{
		PhoneNumber: plan.PhoneNumber.ValueStringPointer(),

		ShouldSendSMS: plan.ShouldSendSMS.ValueBool(),
		DelayInMinSMS: plan.DelayInMinSMS.ValueInt64(),
		SeveritiesSMS: ListToStringArray(plan.SeveritiesSMS),

		ShouldCallVoice: plan.ShouldCallVoice.ValueBool(),
		DelayInMinVoice: plan.DelayInMinVoice.ValueInt64(),
		SeveritiesVoice: ListToStringArray(plan.SeveritiesVoice),

		ShouldSendPush: plan.ShouldSendPush.ValueBool(),
		DelayInMinPush: plan.DelayInMinPush.ValueInt64(),
		SeveritiesPush: ListToStringArray(plan.SeveritiesPush),

		ShouldSendEmail: plan.ShouldSendEmail.ValueBool(),
		DelayInMinEmail: plan.DelayInMinEmail.ValueInt64(),
		SeveritiesEmail: ListToStringArray(plan.SeveritiesEmail),

		DisabledIntentsEmail: ListToNonNullableStringArray(plan.DisabledIntentsEmail),
		DisabledIntentsVoice: ListToNonNullableStringArray(plan.DisabledIntentsVoice),
		DisabledIntentsPush:  ListToNonNullableStringArray(plan.DisabledIntentsPush),
		DisabledIntentsSMS:   ListToNonNullableStringArray(plan.DisabledIntentsSMS),
	}
}

func userIncidentNotificationSettingsPath(userId string) string {
	return fmt.Sprintf("/user/%s/incident-notification-settings", url.PathEscape(userId))
}

func (c *AllQuietAPIClient) CreateUserIncidentNotificationSettingsResource(ctx context.Context, userId string, data *UserIncidentNotificationSettingsModel) (*userIncidentNotificationSettingsResponse, error) {
	reqBody := mapUserIncidentNotificationSettingsRequest(data)

	httpResp, err := c.post(ctx, userIncidentNotificationSettingsPath(userId), reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, reqBody)
	}

	var result userIncidentNotificationSettingsResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *AllQuietAPIClient) UpdateUserIncidentNotificationSettingsResource(ctx context.Context, userId string, data *UserIncidentNotificationSettingsModel) (*userIncidentNotificationSettingsResponse, error) {
	reqBody := mapUserIncidentNotificationSettingsRequest(data)

	httpResp, err := c.put(ctx, userIncidentNotificationSettingsPath(userId), reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, reqBody)
	}

	var result userIncidentNotificationSettingsResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *AllQuietAPIClient) GetUserIncidentNotificationSettingsResource(ctx context.Context, userId string) (*userIncidentNotificationSettingsResponse, error) {
	httpResp, err := c.get(ctx, userIncidentNotificationSettingsPath(userId))
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

	var result userIncidentNotificationSettingsResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *AllQuietAPIClient) DeleteUserIncidentNotificationSettingsResource(ctx context.Context, userId string) error {
	httpResp, err := c.delete(ctx, userIncidentNotificationSettingsPath(userId))
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return logErrorResponse(httpResp, nil)
	}
	return nil
}

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type integrationResponse struct {
	Id                    string                         `json:"id"`
	DisplayName           string                         `json:"displayName"`
	TeamId                string                         `json:"teamId"`
	IsMuted               bool                           `json:"isMuted"`
	IsInMaintenance       bool                           `json:"isInMaintenance"`
	Type                  string                         `json:"type"`
	WebhookUrl            *string                        `json:"webhookUrl"`
	SnoozeSettings        *snoozeSettingsResponse        `json:"snoozeSettings"`
	WebhookAuthentication *webhookAuthenticationResponse `json:"webhookAuthentication"`
	IntegrationSettings   *integrationSettingsResponse   `json:"integrationSettings"`
}

type integrationCreateRequest struct {
	DisplayName           string                         `json:"displayName"`
	TeamId                string                         `json:"teamId"`
	IsMuted               bool                           `json:"isMuted"`
	IsInMaintenance       bool                           `json:"isInMaintenance"`
	Type                  string                         `json:"type"`
	SnoozeSettings        *snoozeSettingsResponse        `json:"snoozeSettings"`
	WebhookAuthentication *webhookAuthenticationResponse `json:"webhookAuthentication"`
	IntegrationSettings   *integrationSettingsResponse   `json:"integrationSettings"`
}

type webhookAuthenticationResponse struct {
	Type   string                               `json:"type"`
	Bearer *webhookAuthenticationBearerResponse `json:"bearer"`
}

type webhookAuthenticationBearerResponse struct {
	Token string `json:"token"`
}

type integrationSettingsResponse struct {
	HttpMonitoring *httpMonitoringResponse `json:"httpMonitoring"`
}

type httpMonitoringResponse struct {
	Url                                string             `json:"url"`
	Method                             string             `json:"method"`
	TimeoutInMilliseconds              int64              `json:"timeoutInMilliseconds"`
	IntervalInSeconds                  int64              `json:"intervalInSeconds"`
	AuthenticationType                 *string            `json:"authenticationType"`
	BasicAuthenticationUsername        *string            `json:"basicAuthenticationUsername"`
	BasicAuthenticationPassword        *string            `json:"basicAuthenticationPassword"`
	BearerAuthenticationToken          *string            `json:"bearerAuthenticationToken"`
	Headers                            *map[string]string `json:"headers"`
	Body                               *string            `json:"body"`
	IsPaused                           bool               `json:"isPaused"`
	ContentTest                        *string            `json:"contentTest"`
	SSLCertificateMaxAgeInDaysDegraded *int64             `json:"sslCertificateMaxAgeInDaysDegraded"`
	SSLCertificateMaxAgeInDaysDown     *int64             `json:"sslCertificateMaxAgeInDaysDown"`
	SeverityDegraded                   *string            `json:"severityDegraded"`
	SeverityDown                       *string            `json:"severityDown"`
}

type snoozeSettingsResponse struct {
	SnoozeWindowInMinutes *int64                  `json:"snoozeWindowInMinutes"`
	Filters               *[]snoozeFilterResponse `json:"filters"`
}

type snoozeFilterResponse struct {
	SelectedDays          *[]string `json:"selectedDays"`
	From                  *string   `json:"from"`
	Until                 *string   `json:"until"`
	SnoozeWindowInMinutes *int64    `json:"snoozeWindowInMinutes"`
	SnoozeUntilAbsolute   *string   `json:"snoozeUntilAbsolute"`
}

func mapIntegrationCreateRequest(ctx context.Context, plan *IntegrationModel) *integrationCreateRequest {
	return &integrationCreateRequest{
		DisplayName:           plan.DisplayName.ValueString(),
		TeamId:                plan.TeamId.ValueString(),
		IsMuted:               plan.IsMuted.ValueBool(),
		IsInMaintenance:       plan.IsInMaintenance.ValueBool(),
		Type:                  plan.Type.ValueString(),
		SnoozeSettings:        mapSnoozeSettingsCreateRequest(plan.SnoozeSettings),
		WebhookAuthentication: mapWebhookAuthenticationCreateRequest(plan.WebhookAuthentication),
		IntegrationSettings:   mapIntegrationSettingsCreateRequest(plan.IntegrationSettings),
	}
}

func mapIntegrationSettingsCreateRequest(plan *IntegrationSettingsModel) *integrationSettingsResponse {
	if plan == nil {
		return nil
	}

	return &integrationSettingsResponse{
		HttpMonitoring: mapHttpMonitoringCreateRequest(plan.HttpMonitoring),
	}
}

func mapHttpMonitoringCreateRequest(plan *HttpMonitoringModel) *httpMonitoringResponse {
	if plan == nil {
		return nil
	}

	return &httpMonitoringResponse{
		Url:                                plan.Url.ValueString(),
		Method:                             plan.Method.ValueString(),
		TimeoutInMilliseconds:              plan.TimeoutInMilliseconds.ValueInt64(),
		IntervalInSeconds:                  plan.IntervalInSeconds.ValueInt64(),
		AuthenticationType:                 plan.AuthenticationType.ValueStringPointer(),
		BasicAuthenticationUsername:        plan.BasicAuthenticationUsername.ValueStringPointer(),
		BasicAuthenticationPassword:        plan.BasicAuthenticationPassword.ValueStringPointer(),
		BearerAuthenticationToken:          plan.BearerAuthenticationToken.ValueStringPointer(),
		Headers:                            mapHeadersCreateRequest(plan.Headers),
		Body:                               plan.Body.ValueStringPointer(),
		IsPaused:                           plan.IsPaused.ValueBool(),
		ContentTest:                        plan.ContentTest.ValueStringPointer(),
		SSLCertificateMaxAgeInDaysDegraded: plan.SSLCertificateMaxAgeInDaysDegraded.ValueInt64Pointer(),
		SSLCertificateMaxAgeInDaysDown:     plan.SSLCertificateMaxAgeInDaysDown.ValueInt64Pointer(),
		SeverityDegraded:                   plan.SeverityDegraded.ValueStringPointer(),
		SeverityDown:                       plan.SeverityDown.ValueStringPointer(),
	}
}

func mapHeadersCreateRequest(plan types.Map) *map[string]string {
	if plan.IsNull() || plan.IsUnknown() {
		return nil
	}

	headers := make(map[string]string)
	for k, v := range plan.Elements() {
		if v.IsNull() || v.IsUnknown() {
			continue
		}
		str, ok := v.(types.String)
		if !ok {
			continue
		}
		headers[k] = str.ValueString()
	}

	return &headers
}

func mapWebhookAuthenticationCreateRequest(plan *WebhookAuthenticationModel) *webhookAuthenticationResponse {
	if plan == nil {
		return nil
	}

	return &webhookAuthenticationResponse{
		Type:   plan.Type.ValueString(),
		Bearer: mapWebhookAuthenticationBearerCreateRequest(plan.Bearer),
	}
}

func mapWebhookAuthenticationBearerCreateRequest(plan *BearerModel) *webhookAuthenticationBearerResponse {
	if plan == nil {
		return nil
	}

	return &webhookAuthenticationBearerResponse{
		Token: plan.Token.ValueString(),
	}
}

func mapSnoozeSettingsCreateRequest(plan *SnoozeSettingsModel) *snoozeSettingsResponse {
	if plan == nil {
		return nil
	}

	return &snoozeSettingsResponse{
		SnoozeWindowInMinutes: plan.SnoozeWindowInMinutes.ValueInt64Pointer(),
		Filters:               mapSnoozeFiltersCreateRequest(plan.Filters),
	}
}

func mapSnoozeFiltersCreateRequest(plan *[]SnoozeFilterModel) *[]snoozeFilterResponse {
	if plan == nil {
		return nil
	}

	filters := make([]snoozeFilterResponse, len(*plan))
	for i, filter := range *plan {
		filters[i] = *mapSnoozeFilterCreateRequest(&filter)
	}
	return &filters
}

func mapSnoozeFilterCreateRequest(plan *SnoozeFilterModel) *snoozeFilterResponse {
	if plan == nil {
		return nil
	}

	return &snoozeFilterResponse{
		SelectedDays:          ListToStringArray(plan.SelectedDays),
		From:                  plan.From.ValueStringPointer(),
		Until:                 plan.Until.ValueStringPointer(),
		SnoozeWindowInMinutes: plan.SnoozeWindowInMinutes.ValueInt64Pointer(),
		SnoozeUntilAbsolute:   plan.SnoozeUntilAbsolute.ValueStringPointer(),
	}
}

func (c *AllQuietAPIClient) CreateIntegrationResource(ctx context.Context, plan *IntegrationModel) (*integrationResponse, error) {
	request := &integrationCreateRequest{
		DisplayName:           plan.DisplayName.ValueString(),
		TeamId:                plan.TeamId.ValueString(),
		IsMuted:               plan.IsMuted.ValueBool(),
		IsInMaintenance:       plan.IsInMaintenance.ValueBool(),
		Type:                  plan.Type.ValueString(),
		SnoozeSettings:        mapSnoozeSettingsCreateRequest(plan.SnoozeSettings),
		WebhookAuthentication: mapWebhookAuthenticationCreateRequest(plan.WebhookAuthentication),
		IntegrationSettings:   mapIntegrationSettingsCreateRequest(plan.IntegrationSettings),
	}

	url := "/inbound-integration"
	httpResp, err := c.post(ctx, url, request)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, request)
	}

	var result integrationResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) DeleteIntegrationResource(ctx context.Context, id string) error {
	url := fmt.Sprintf("/inbound-integration/%s", url.PathEscape(id))
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

func (c *AllQuietAPIClient) UpdateIntegrationResource(ctx context.Context, id string, data *IntegrationModel) (*integrationResponse, error) {
	reqBody := mapIntegrationCreateRequest(ctx, data)

	url := fmt.Sprintf("/inbound-integration/%s", url.PathEscape(id))
	httpResp, err := c.put(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, reqBody)
	}

	var result integrationResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) GetIntegrationResource(ctx context.Context, id string) (*integrationResponse, error) {
	url := fmt.Sprintf("/inbound-integration/%s", url.PathEscape(id))
	httpResp, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result integrationResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

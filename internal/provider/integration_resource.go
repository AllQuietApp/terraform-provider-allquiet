// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &Integration{}
var _ resource.ResourceWithImportState = &Integration{}

func NewIntegration() resource.Resource {
	return &Integration{}
}

// Integration defines the resource implementation.
type Integration struct {
	client *AllQuietAPIClient
}

// IntegrationModel describes the resource data model.
type IntegrationModel struct {
	Id                    types.String                `tfsdk:"id"`
	DisplayName           types.String                `tfsdk:"display_name"`
	TeamId                types.String                `tfsdk:"team_id"`
	IsMuted               types.Bool                  `tfsdk:"is_muted"`
	IsInMaintenance       types.Bool                  `tfsdk:"is_in_maintenance"`
	Type                  types.String                `tfsdk:"type"`
	WebhookUrl            types.String                `tfsdk:"webhook_url"`
	SnoozeSettings        *SnoozeSettingsModel        `tfsdk:"snooze_settings"`
	WebhookAuthentication *WebhookAuthenticationModel `tfsdk:"webhook_authentication"`
	IntegrationSettings   *IntegrationSettingsModel   `tfsdk:"integration_settings"`
}

type IntegrationSettingsModel struct {
	HttpMonitoring   *HttpMonitoringModel   `tfsdk:"http_monitoring"`
	HeartbeatMonitor *HeartbeatMonitorModel `tfsdk:"heartbeat_monitor"`
	CronjobMonitor   *CronjobMonitorModel   `tfsdk:"cronjob_monitor"`
}

type HttpMonitoringModel struct {
	Url                                types.String `tfsdk:"url"`
	Method                             types.String `tfsdk:"method"`
	TimeoutInMilliseconds              types.Int64  `tfsdk:"timeout_in_milliseconds"`
	IntervalInSeconds                  types.Int64  `tfsdk:"interval_in_seconds"`
	AuthenticationType                 types.String `tfsdk:"authentication_type"`
	BasicAuthenticationUsername        types.String `tfsdk:"basic_authentication_username"`
	BasicAuthenticationPassword        types.String `tfsdk:"basic_authentication_password"`
	BearerAuthenticationToken          types.String `tfsdk:"bearer_authentication_token"`
	Headers                            types.Map    `tfsdk:"headers"`
	Body                               types.String `tfsdk:"body"`
	IsPaused                           types.Bool   `tfsdk:"is_paused"`
	ContentTest                        types.String `tfsdk:"content_test"`
	SSLCertificateMaxAgeInDaysDegraded types.Int64  `tfsdk:"ssl_certificate_max_age_in_days_degraded"`
	SSLCertificateMaxAgeInDaysDown     types.Int64  `tfsdk:"ssl_certificate_max_age_in_days_down"`
	SeverityDegraded                   types.String `tfsdk:"severity_degraded"`
	SeverityDown                       types.String `tfsdk:"severity_down"`
}

type HeartbeatMonitorModel struct {
	IntervalInSec    types.Int64  `tfsdk:"interval_in_sec"`
	GracePeriodInSec types.Int64  `tfsdk:"grace_period_in_sec"`
	Severity         types.String `tfsdk:"severity"`
}

type CronjobMonitorModel struct {
	CronExpression   types.String `tfsdk:"cron_expression"`
	GracePeriodInSec types.Int64  `tfsdk:"grace_period_in_sec"`
	Severity         types.String `tfsdk:"severity"`
	TimeZoneId       types.String `tfsdk:"time_zone_id"`
}

type WebhookAuthenticationModel struct {
	Type   types.String `tfsdk:"type"`
	Bearer *BearerModel `tfsdk:"bearer"`
}

type BearerModel struct {
	Token types.String `tfsdk:"token"`
}

type SnoozeSettingsModel struct {
	SnoozeWindowInMinutes types.Int64          `tfsdk:"snooze_window_in_minutes"`
	Filters               *[]SnoozeFilterModel `tfsdk:"filters"`
}

type SnoozeFilterModel struct {
	SelectedDays          types.List   `tfsdk:"selected_days"`
	From                  types.String `tfsdk:"from"`
	Until                 types.String `tfsdk:"until"`
	SnoozeWindowInMinutes types.Int64  `tfsdk:"snooze_window_in_minutes"`
	SnoozeUntilAbsolute   types.String `tfsdk:"snooze_until_absolute"`
}

func (r *Integration) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration"
}

func (r *Integration) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The `integration` resource represents an integration in All Quiet. Integrations are used to receive alerts from external systems like Datadog or Prometheus.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the integration",
				Required:            true,
			},
			"team_id": schema.StringAttribute{
				MarkdownDescription: "The team id of the integration",
				Required:            true,
			},
			"is_muted": schema.BoolAttribute{
				MarkdownDescription: "If the integration is muted. Deprecated: Use resource `allquiet_integration_maintenance_window` instead.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_in_maintenance": schema.BoolAttribute{
				MarkdownDescription: "If the integration is in maintenance mode. Deprecated: Use resource `allquiet_integration_maintenance_window` instead.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"snooze_settings": schema.SingleNestedAttribute{
				MarkdownDescription: "The snooze settings of the integration",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"snooze_window_in_minutes": schema.Int64Attribute{
						MarkdownDescription: "The snooze window in minutes. If your integration is flaky and you'd like to reduce noise, you can set a snooze window. This will keep the incident snoozed for the specified time period and only alert you once the snooze window is over and the incident has not been resolved yet. Max 1440 minutes (24 hours).",
						Optional:            true,
						Validators: []validator.Int64{
							int64validator.Between(0, 1440),
						},
					},
					"filters": schema.ListNestedAttribute{
						MarkdownDescription: "The snooze filters of the integration. Only the first matching filter will be applied. Filters are applied in the order they are defined.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"selected_days": schema.ListAttribute{
									Optional:            true,
									MarkdownDescription: "Days of the week. Possible values are: " + strings.Join(ValidDaysOfWeek, ", "),
									ElementType:         types.StringType,
									Validators: []validator.List{
										listvalidator.ValueStringsAre(DaysOfWeekValidator("Not a valid day of week")),
									},
								},
								"from": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "From time of the time filter. Format: HH:mm",
									Validators:          []validator.String{TimeValidator("Not a valid time")},
								},
								"until": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Until time of the time filter. Format: HH:mm",
									Validators:          []validator.String{TimeValidator("Not a valid time")},
								},
								"snooze_window_in_minutes": schema.Int64Attribute{
									MarkdownDescription: "The snooze window in minutes. If your integration is flaky and you'd like to reduce noise, you can set a snooze window. This will keep the incident snoozed for the specified time period and only alert you once the snooze window is over and the incident has not been resolved yet. Max 1440 minutes (24 hours).",
									Optional:            true,
									Validators: []validator.Int64{
										int64validator.Between(0, 1440),
									},
								},
								"snooze_until_absolute": schema.StringAttribute{
									MarkdownDescription: "The absolute time to snooze the integration until. Format:HH:mm. Examples: When the incident happens at 01 am in the night, and the snooze until absolute is set to 07:00, the incident will be snoozed until 07:00 the same night. If the incident happens at 14:00, it will be snoozed until 07:00 the next day.",
									Optional:            true,
									Validators:          []validator.String{TimeValidator("Not a valid time")},
								},
							},
						},
					},
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the integration. See all types here: https://allquiet.app/api/public/v1/inbound-integration/types",
				Required:            true,
			},
			"webhook_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The webhook url of the integration if it is a webhook-like integration e.g. Amazon CloudWatch",
			},
			"webhook_authentication": schema.SingleNestedAttribute{
				MarkdownDescription: "The webhook authentication of the integration",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the webhook authentication. Possible values are: " + strings.Join(ValidWebhookAuthenticationTypes, ", "),
						Required:            true,
						Validators:          []validator.String{WebhookAuthenticationTypeValidator("Not a valid webhook authentication type")},
					},
					"bearer": schema.SingleNestedAttribute{
						MarkdownDescription: "The bearer token of the webhook authentication",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"token": schema.StringAttribute{
								MarkdownDescription: "The token of the webhook authentication",
								Required:            true,
								Sensitive:           true,
							},
						},
					},
				},
			},
			"integration_settings": schema.SingleNestedAttribute{
				MarkdownDescription: "The integration settings of the integration",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"http_monitoring": schema.SingleNestedAttribute{
						MarkdownDescription: "The http monitoring of the integration",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"url": schema.StringAttribute{
								MarkdownDescription: "The url of the http monitoring",
								Required:            true,
							},
							"method": schema.StringAttribute{
								MarkdownDescription: "The method of the http monitoring. Possible values are: " + strings.Join(ValidHttpMonitoringMethods, ", "),
								Required:            true,
								Validators: []validator.String{
									HttpMonitoringMethodValidator("Not a valid method"),
								},
							},
							"timeout_in_milliseconds": schema.Int64Attribute{
								MarkdownDescription: "The timeout in milliseconds of the http monitoring. Min 50, max 60000.",
								Required:            true,
								Validators: []validator.Int64{
									int64validator.Between(50, 60000),
								},
							},
							"interval_in_seconds": schema.Int64Attribute{
								MarkdownDescription: "The interval in seconds of the http monitoring. Valid values are: " + strings.Join(ValidIntervalsInSecondsAsString, ", "),
								Required:            true,
								Validators: []validator.Int64{
									IntervalInSecondsValidator("Not a valid interval in seconds"),
								},
							},
							"authentication_type": schema.StringAttribute{
								MarkdownDescription: "The authentication type of the http monitoring. Possible values are: " + strings.Join(ValidHttpMonitoringAuthenticationTypes, ", "),
								Optional:            true,
								Validators:          []validator.String{HttpMonitoringAuthenticationTypeValidator("Not a valid authentication type")},
							},
							"basic_authentication_username": schema.StringAttribute{
								MarkdownDescription: "The basic authentication username of the http monitoring",
								Optional:            true,
								Sensitive:           true,
							},
							"basic_authentication_password": schema.StringAttribute{
								MarkdownDescription: "The basic authentication password of the http monitoring",
								Optional:            true,
								Sensitive:           true,
							},
							"bearer_authentication_token": schema.StringAttribute{
								MarkdownDescription: "The bearer authentication token of the http monitoring",
								Optional:            true,
								Sensitive:           true,
							},
							"headers": schema.MapAttribute{
								MarkdownDescription: "The headers of the http monitoring",
								Optional:            true,
								Sensitive:           true,
								ElementType:         types.StringType,
							},
							"body": schema.StringAttribute{
								MarkdownDescription: "The body to send in the http request",
								Optional:            true,
								Sensitive:           true,
							},
							"is_paused": schema.BoolAttribute{
								MarkdownDescription: "If the http monitoring is paused",
								Optional:            true,
								Computed:            true,
								Default:             booldefault.StaticBool(false),
							},
							"content_test": schema.StringAttribute{
								MarkdownDescription: "The content test of the http monitoring",
								Optional:            true,
							},
							"ssl_certificate_max_age_in_days_degraded": schema.Int64Attribute{
								MarkdownDescription: "The ssl certificate max age in days degraded of the http monitoring",
								Optional:            true,
							},
							"ssl_certificate_max_age_in_days_down": schema.Int64Attribute{
								MarkdownDescription: "The ssl certificate max age in days down of the http monitoring",
								Optional:            true,
							},
							"severity_degraded": schema.StringAttribute{
								MarkdownDescription: "The severity degraded of the http monitoring. Possible values are: " + strings.Join(ValidSeverities, ", "),
								Optional:            true,
								Validators: []validator.String{
									SeverityValidator("Not a valid severity"),
								},
							},
							"severity_down": schema.StringAttribute{
								MarkdownDescription: "The severity down of the http monitoring. Possible values are: " + strings.Join(ValidSeverities, ", "),
								Optional:            true,
								Validators: []validator.String{
									SeverityValidator("Not a valid severity"),
								},
							},
						},
					},
					"heartbeat_monitor": schema.SingleNestedAttribute{
						MarkdownDescription: "The heartbeat monitor of the integration",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"interval_in_sec": schema.Int64Attribute{
								MarkdownDescription: "The interval in seconds of the heartbeat monitor",
								Required:            true,
							},
							"grace_period_in_sec": schema.Int64Attribute{
								MarkdownDescription: "The grace period in seconds of the heartbeat monitor",
								Required:            true,
							},
							"severity": schema.StringAttribute{
								MarkdownDescription: "The severity of the heartbeat monitor. Possible values are: " + strings.Join(ValidSeverities, ", "),
								Required:            true,
								Validators: []validator.String{
									SeverityValidator("Not a valid severity"),
								},
							},
						},
					},
					"cronjob_monitor": schema.SingleNestedAttribute{
						MarkdownDescription: "The cronjob monitor of the integration",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"cron_expression": schema.StringAttribute{
								MarkdownDescription: "The cron expression of the cronjob monitor",
								Required:            true,
							},
							"grace_period_in_sec": schema.Int64Attribute{
								MarkdownDescription: "The grace period in seconds of the cronjob monitor",
								Required:            true,
							},
							"severity": schema.StringAttribute{
								MarkdownDescription: "The severity of the cronjob monitor. Possible values are: " + strings.Join(ValidSeverities, ", "),
								Required:            true,
								Validators: []validator.String{
									SeverityValidator("Not a valid severity"),
								},
							},
							"time_zone_id": schema.StringAttribute{
								MarkdownDescription: "The time zone id of the cronjob monitor. Find all timezone ids [here](https://allquiet.app/api/public/v1/timezone)",
								Optional:            true,
							},
						},
					},
				},
			},
		},
	}
}

func (r *Integration) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*AllQuietAPIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *Integration) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IntegrationModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integrationResponse, err := r.client.CreateIntegrationResource(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create integration resource, got error: %s", err))
		return
	}

	mapIntegrationResponseToModel(ctx, integrationResponse, &data)

	tflog.Trace(ctx, "created integration resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Integration) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IntegrationModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	integrationResponse, err := r.client.GetIntegrationResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get integration resource, got error: %s", err))
		return
	}

	if integrationResponse == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to get integration resource, got nil response")
		return
	}

	mapIntegrationResponseToModel(ctx, integrationResponse, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Integration) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IntegrationModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	integrationResponse, err := r.client.UpdateIntegrationResource(ctx, data.Id.ValueString(), &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update integration resource, got error: %s", err))
		return
	}

	mapIntegrationResponseToModel(ctx, integrationResponse, &data)

	tflog.Trace(ctx, "updated integration resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Integration) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IntegrationModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteIntegrationResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update integration resource, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted integration resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Integration) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapIntegrationResponseToModel(ctx context.Context, response *integrationResponse, data *IntegrationModel) {

	data.Id = types.StringValue(response.Id)
	data.DisplayName = types.StringValue(response.DisplayName)
	data.TeamId = types.StringValue(response.TeamId)
	data.IsMuted = types.BoolValue(response.IsMuted)
	data.IsInMaintenance = types.BoolValue(response.IsInMaintenance)
	data.Type = types.StringValue(response.Type)
	data.WebhookUrl = types.StringPointerValue(response.WebhookUrl)
	data.SnoozeSettings = mapSnoozeSettingsResponseToModel(ctx, response.SnoozeSettings)
	data.WebhookAuthentication = mapWebhookAuthenticationResponseToModel(response.WebhookAuthentication)
	data.IntegrationSettings = mapIntegrationSettingsResponseToModel(response.IntegrationSettings)
}

func mapIntegrationSettingsResponseToModel(response *integrationSettingsResponse) *IntegrationSettingsModel {
	if response == nil {
		return nil
	}

	return &IntegrationSettingsModel{
		HttpMonitoring:   mapHttpMonitoringResponseToModel(response.HttpMonitoring),
		HeartbeatMonitor: mapHeartbeatMonitorResponseToModel(response.HeartbeatMonitor),
		CronjobMonitor:   mapCronjobMonitorResponseToModel(response.CronjobMonitor),
	}
}

func mapHeartbeatMonitorResponseToModel(response *heartbeatMonitorResponse) *HeartbeatMonitorModel {
	if response == nil {
		return nil
	}

	return &HeartbeatMonitorModel{
		IntervalInSec:    types.Int64Value(response.IntervalInSec),
		GracePeriodInSec: types.Int64Value(response.GracePeriodInSec),
		Severity:         types.StringValue(response.Severity),
	}
}

func mapCronjobMonitorResponseToModel(response *cronjobMonitorResponse) *CronjobMonitorModel {

	if response == nil {
		return nil
	}

	return &CronjobMonitorModel{
		CronExpression:   types.StringValue(response.CronExpression),
		GracePeriodInSec: types.Int64Value(response.GracePeriodInSec),
		Severity:         types.StringValue(response.Severity),
		TimeZoneId:       types.StringPointerValue(response.TimeZoneId),
	}
}

func mapHttpMonitoringResponseToModel(response *httpMonitoringResponse) *HttpMonitoringModel {
	if response == nil {
		return nil
	}

	return &HttpMonitoringModel{
		Url:                                types.StringValue(response.Url),
		Method:                             types.StringValue(response.Method),
		TimeoutInMilliseconds:              types.Int64Value(response.TimeoutInMilliseconds),
		IntervalInSeconds:                  types.Int64Value(response.IntervalInSeconds),
		AuthenticationType:                 types.StringPointerValue(response.AuthenticationType),
		BasicAuthenticationUsername:        types.StringPointerValue(response.BasicAuthenticationUsername),
		BasicAuthenticationPassword:        types.StringPointerValue(response.BasicAuthenticationPassword),
		BearerAuthenticationToken:          types.StringPointerValue(response.BearerAuthenticationToken),
		Headers:                            mapHeadersResponseToModel(response.Headers),
		Body:                               types.StringPointerValue(response.Body),
		IsPaused:                           types.BoolValue(response.IsPaused),
		ContentTest:                        types.StringPointerValue(response.ContentTest),
		SSLCertificateMaxAgeInDaysDegraded: types.Int64PointerValue(response.SSLCertificateMaxAgeInDaysDegraded),
		SSLCertificateMaxAgeInDaysDown:     types.Int64PointerValue(response.SSLCertificateMaxAgeInDaysDown),
		SeverityDegraded:                   types.StringPointerValue(response.SeverityDegraded),
		SeverityDown:                       types.StringPointerValue(response.SeverityDown),
	}
}
func mapHeadersResponseToModel(response *map[string]string) types.Map {
	if response == nil {
		return types.MapNull(types.StringType)
	}

	// Convert map[string]string to map[string]attr.Value
	elements := make(map[string]attr.Value, len(*response))
	for k, v := range *response {
		elements[k] = types.StringValue(v)
	}

	val, diags := types.MapValue(types.StringType, elements)
	if diags.HasError() {
		// In production code, you'd probably want to return an error here
		// but for now just fallback to null if conversion fails
		return types.MapNull(types.StringType)
	}

	return val
}

func mapWebhookAuthenticationResponseToModel(response *webhookAuthenticationResponse) *WebhookAuthenticationModel {
	if response == nil {
		return nil
	}

	return &WebhookAuthenticationModel{
		Type:   types.StringValue(response.Type),
		Bearer: mapBearerResponseToModel(response.Bearer),
	}
}

func mapBearerResponseToModel(response *webhookAuthenticationBearerResponse) *BearerModel {
	if response == nil {
		return nil
	}

	return &BearerModel{
		Token: types.StringValue(response.Token),
	}
}

func mapSnoozeSettingsResponseToModel(ctx context.Context, response *snoozeSettingsResponse) *SnoozeSettingsModel {
	if response == nil {
		return nil
	}

	return &SnoozeSettingsModel{
		SnoozeWindowInMinutes: types.Int64PointerValue(response.SnoozeWindowInMinutes),
		Filters:               mapSnoozeFiltersResponseToModel(ctx, response.Filters),
	}
}

func mapSnoozeFiltersResponseToModel(ctx context.Context, response *[]snoozeFilterResponse) *[]SnoozeFilterModel {
	if response == nil {
		return nil
	}

	filters := make([]SnoozeFilterModel, len(*response))
	for i, filter := range *response {
		filters[i] = *mapSnoozeFilterResponseToModel(ctx, &filter)
	}
	return &filters
}

func mapSnoozeFilterResponseToModel(ctx context.Context, response *snoozeFilterResponse) *SnoozeFilterModel {
	return &SnoozeFilterModel{
		SelectedDays:          MapNullableList(ctx, response.SelectedDays),
		From:                  types.StringPointerValue(response.From),
		Until:                 types.StringPointerValue(response.Until),
		SnoozeWindowInMinutes: types.Int64PointerValue(response.SnoozeWindowInMinutes),
		SnoozeUntilAbsolute:   types.StringPointerValue(response.SnoozeUntilAbsolute),
	}
}

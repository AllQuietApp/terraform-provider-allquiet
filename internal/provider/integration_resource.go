// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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

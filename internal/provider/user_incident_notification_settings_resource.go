// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &UserIncidentNotificationSettings{}
var _ resource.ResourceWithImportState = &UserIncidentNotificationSettings{}

func NewUserIncidentNotificationSettings() resource.Resource {
	return &UserIncidentNotificationSettings{}
}

type UserIncidentNotificationSettings struct {
	client *AllQuietAPIClient
}

type UserIncidentNotificationSettingsModel struct {
	Id     types.String `tfsdk:"id"`
	UserId types.String `tfsdk:"user_id"`

	PhoneNumber types.String `tfsdk:"phone_number"`

	ShouldSendSMS types.Bool  `tfsdk:"should_send_sms"`
	DelayInMinSMS types.Int64 `tfsdk:"delay_in_min_sms"`
	SeveritiesSMS types.List  `tfsdk:"severities_sms"`

	ShouldCallVoice types.Bool  `tfsdk:"should_call_voice"`
	DelayInMinVoice types.Int64 `tfsdk:"delay_in_min_voice"`
	SeveritiesVoice types.List  `tfsdk:"severities_voice"`

	ShouldSendPush types.Bool  `tfsdk:"should_send_push"`
	DelayInMinPush types.Int64 `tfsdk:"delay_in_min_push"`
	SeveritiesPush types.List  `tfsdk:"severities_push"`

	ShouldSendEmail types.Bool  `tfsdk:"should_send_email"`
	DelayInMinEmail types.Int64 `tfsdk:"delay_in_min_email"`
	SeveritiesEmail types.List  `tfsdk:"severities_email"`

	DisabledIntentsEmail types.List `tfsdk:"disabled_intents_email"`
	DisabledIntentsVoice types.List `tfsdk:"disabled_intents_voice"`
	DisabledIntentsPush  types.List `tfsdk:"disabled_intents_push"`
	DisabledIntentsSMS   types.List `tfsdk:"disabled_intents_sms"`
}

func (r *UserIncidentNotificationSettings) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_incident_notification_settings"
}

func (r *UserIncidentNotificationSettings) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a user's incident notification settings as a standalone resource. " +
			"While this resource exists, the corresponding user cannot edit their own notification settings " +
			"in the All Quiet UI \u2014 Terraform owns them. Removing this resource releases the lock again " +
			"and preserves the last applied settings.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Id of the resource (matches the user id).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Id of the user whose incident notification settings are managed by this resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{GuidValidator("user_id must be a valid UUID")},
			},

			"phone_number": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Phone number used for SMS and voice notifications, in international format (for example `+12035479055`). " +
					"Strict ownership: omitting this attribute clears the user's phone number on the backend. " +
					"While this resource exists the user cannot change their phone number in the All Quiet UI.",
				Validators: []validator.String{stringvalidator.RegexMatches(
					regexp.MustCompile(`^\+\d+$`),
					"must contain phone number in international format matching the pattern '^\\+\\d+$'",
				)},
				Sensitive: true,
			},

			"should_send_sms": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Should send SMS notifications",
			},
			"delay_in_min_sms": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Delay in minutes for SMS notifications",
				Validators:          []validator.Int64{int64validator.Between(0, 60)},
			},
			"severities_sms": schema.ListAttribute{
				Required:            true,
				MarkdownDescription: "Severities for SMS notifications. Possible values are: " + strings.Join(ValidSeverities, ", "),
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(SeverityValidator("Not a valid severity")),
				},
			},

			"should_call_voice": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Should send Voice Call notifications",
			},
			"delay_in_min_voice": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Delay in minutes for Voice Call notifications",
				Validators:          []validator.Int64{int64validator.Between(0, 60)},
			},
			"severities_voice": schema.ListAttribute{
				Required:            true,
				MarkdownDescription: "Severities for Voice Call notifications. Possible values are: " + strings.Join(ValidSeverities, ", "),
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(SeverityValidator("Not a valid severity")),
				},
			},

			"should_send_push": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Should send Push notifications",
			},
			"delay_in_min_push": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Delay in minutes for Push notifications",
				Validators:          []validator.Int64{int64validator.Between(0, 60)},
			},
			"severities_push": schema.ListAttribute{
				Required:            true,
				MarkdownDescription: "Severities for Push notifications. Possible values are: " + strings.Join(ValidSeverities, ", "),
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(SeverityValidator("Not a valid severity")),
				},
			},

			"should_send_email": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Should send Email notifications",
			},
			"delay_in_min_email": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Delay in minutes for Email notifications",
				Validators:          []validator.Int64{int64validator.Between(0, 60)},
			},
			"severities_email": schema.ListAttribute{
				Required:            true,
				MarkdownDescription: "Severities for Email notifications. Possible values are: " + strings.Join(ValidSeverities, ", "),
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(SeverityValidator("Not a valid severity")),
				},
			},
			"disabled_intents_email": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Disabled intents for Email notifications. Possible values are: " + strings.Join(ValidIntents, ", "),
				ElementType:         types.StringType,
				Validators:          []validator.List{listvalidator.ValueStringsAre(IntentValidator("Not a valid intent"))},
			},
			"disabled_intents_voice": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Disabled intents for Voice Call notifications. Possible values are: " + strings.Join(ValidIntents, ", "),
				ElementType:         types.StringType,
				Validators:          []validator.List{listvalidator.ValueStringsAre(IntentValidator("Not a valid intent"))},
			},
			"disabled_intents_push": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Disabled intents for Push notifications. Possible values are: " + strings.Join(ValidIntents, ", "),
				ElementType:         types.StringType,
				Validators:          []validator.List{listvalidator.ValueStringsAre(IntentValidator("Not a valid intent"))},
			},
			"disabled_intents_sms": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Disabled intents for SMS notifications. Possible values are: " + strings.Join(ValidIntents, ", "),
				ElementType:         types.StringType,
				Validators:          []validator.List{listvalidator.ValueStringsAre(IntentValidator("Not a valid intent"))},
			},
		},
	}
}

func (r *UserIncidentNotificationSettings) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*AllQuietAPIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *AllQuietAPIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *UserIncidentNotificationSettings) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UserIncidentNotificationSettingsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.CreateUserIncidentNotificationSettingsResource(ctx, data.UserId.ValueString(), &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create user incident notification settings resource, got error: %s", err))
		return
	}
	mapUserIncidentNotificationSettingsResponseToModel(ctx, response, &data)

	tflog.Trace(ctx, "created user incident notification settings resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserIncidentNotificationSettings) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UserIncidentNotificationSettingsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.GetUserIncidentNotificationSettingsResource(ctx, data.UserId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read user incident notification settings resource, got error: %s", err))
		return
	}

	if response == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	mapUserIncidentNotificationSettingsResponseToModel(ctx, response, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserIncidentNotificationSettings) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data UserIncidentNotificationSettingsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.UpdateUserIncidentNotificationSettingsResource(ctx, data.UserId.ValueString(), &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update user incident notification settings resource, got error: %s", err))
		return
	}

	mapUserIncidentNotificationSettingsResponseToModel(ctx, response, &data)

	tflog.Trace(ctx, "updated user incident notification settings resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserIncidentNotificationSettings) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data UserIncidentNotificationSettingsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteUserIncidentNotificationSettingsResource(ctx, data.UserId.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete user incident notification settings resource, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted user incident notification settings resource")
}

func (r *UserIncidentNotificationSettings) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("user_id"), req, resp)
}

func mapUserIncidentNotificationSettingsResponseToModel(ctx context.Context, response *userIncidentNotificationSettingsResponse, data *UserIncidentNotificationSettingsModel) {
	data.Id = types.StringValue(response.UserId)
	data.UserId = types.StringValue(response.UserId)

	data.PhoneNumber = types.StringPointerValue(response.PhoneNumber)

	data.ShouldSendSMS = types.BoolPointerValue(response.ShouldSendSMS)
	data.DelayInMinSMS = types.Int64PointerValue(response.DelayInMinSMS)
	data.SeveritiesSMS = MapNullableList(ctx, response.SeveritiesSMS)

	data.ShouldCallVoice = types.BoolPointerValue(response.ShouldCallVoice)
	data.DelayInMinVoice = types.Int64PointerValue(response.DelayInMinVoice)
	data.SeveritiesVoice = MapNullableList(ctx, response.SeveritiesVoice)

	data.ShouldSendPush = types.BoolPointerValue(response.ShouldSendPush)
	data.DelayInMinPush = types.Int64PointerValue(response.DelayInMinPush)
	data.SeveritiesPush = MapNullableList(ctx, response.SeveritiesPush)

	data.ShouldSendEmail = types.BoolPointerValue(response.ShouldSendEmail)
	data.DelayInMinEmail = types.Int64PointerValue(response.DelayInMinEmail)
	data.SeveritiesEmail = MapNullableList(ctx, response.SeveritiesEmail)

	data.DisabledIntentsEmail = MapNullableList(ctx, response.DisabledIntentsEmail)
	data.DisabledIntentsVoice = MapNullableList(ctx, response.DisabledIntentsVoice)
	data.DisabledIntentsPush = MapNullableList(ctx, response.DisabledIntentsPush)
	data.DisabledIntentsSMS = MapNullableList(ctx, response.DisabledIntentsSMS)
}

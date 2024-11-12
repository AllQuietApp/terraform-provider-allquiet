// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &User{}
var _ resource.ResourceWithImportState = &User{}

func NewUser() resource.Resource {
	return &User{}
}

type User struct {
	client *AllQuietAPIClient
}

type UserModel struct {
	Id                           types.String                       `tfsdk:"id"`
	DisplayName                  types.String                       `tfsdk:"display_name"`
	Email                        types.String                       `tfsdk:"email"`
	PhoneNumber                  types.String                       `tfsdk:"phone_number"`
	TimeZoneId                   types.String                       `tfsdk:"time_zone_id"`
	IncidentNotificationSettings *IncidentNotificationSettingsModel `tfsdk:"incident_notification_settings"`
}

type IncidentNotificationSettingsModel struct {
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
}

func (r *User) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *User) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The user resource represents a user in All Quiet. Users can be members of users and receive notifications for incidents.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the user",
				Required:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email of the user",
				Required:            true,
				Validators: []validator.String{stringvalidator.RegexMatches(
					regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
					"must contain email matching the pattern '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$'",
				)},
			},
			"phone_number": schema.StringAttribute{
				MarkdownDescription: "The phone number of the user",
				Optional:            true,
				Validators: []validator.String{stringvalidator.RegexMatches(
					regexp.MustCompile(`^\+\d+$`),
					"must contain phone number in internatiol format matching the pattern '^\\+\\d+$'",
				)},
			},
			"time_zone_id": schema.StringAttribute{
				MarkdownDescription: "The timezone id, defaults to 'UTC' if not provided. Find all timezone ids [here](https://allquiet.app/api/public/v1/timezone)",
				Optional:            true,
				Default:             stringdefault.StaticString("UTC"),
				Computed:            true,
			},
			"incident_notification_settings": schema.SingleNestedAttribute{
				MarkdownDescription: "Settings which channels to use for incident notifications",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
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
						Optional:            true,
						MarkdownDescription: "Severities for SMS notifications",
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
						MarkdownDescription: "Severities for Voice Call notifications",
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
						MarkdownDescription: "Severities for Push notifications",
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
						MarkdownDescription: "Severities for Email notifications",
						ElementType:         types.StringType,
						Validators: []validator.List{
							listvalidator.ValueStringsAre(SeverityValidator("Not a valid severity")),
						},
					},
				},
			},
		},
	}
}

func (r *User) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *User) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UserModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	userResponse, err := r.client.CreateUserResource(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create user resource, got error: %s", err))
		return
	}
	mapUserResponseToModel(ctx, userResponse, &data)

	tflog.Trace(ctx, "created user resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *User) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UserModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	userResponse, err := r.client.GetUserResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get user resource, got error: %s", err))
		return
	}

	if userResponse == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to get user resource, got nil response")
		return
	}

	mapUserResponseToModel(ctx, userResponse, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *User) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data UserModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	userResponse, err := r.client.UpdateUserResource(ctx, data.Id.ValueString(), &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update user resource, got error: %s", err))
		return
	}

	mapUserResponseToModel(ctx, userResponse, &data)

	tflog.Trace(ctx, "updated user resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *User) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data UserModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUserResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete user resource, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted user resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *User) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapUserResponseToModel(ctx context.Context, response *userResponse, data *UserModel) {
	data.Id = types.StringValue(response.Id)
	data.DisplayName = types.StringValue(response.DisplayName)
	data.Email = types.StringValue(response.Email)
	data.PhoneNumber = types.StringPointerValue(response.PhoneNumber)
	data.TimeZoneId = types.StringValue(response.TimeZoneId)

	if response.IncidentNotificationSettings != nil {
		data.IncidentNotificationSettings = &IncidentNotificationSettingsModel{
			ShouldSendSMS: types.BoolValue(response.IncidentNotificationSettings.ShouldSendSMS),
			DelayInMinSMS: types.Int64Value(response.IncidentNotificationSettings.DelayInMinSMS),
			SeveritiesSMS: MapNullableListWithEmpty(ctx, response.IncidentNotificationSettings.SeveritiesSMS),

			ShouldCallVoice: types.BoolValue(response.IncidentNotificationSettings.ShouldCallVoice),
			DelayInMinVoice: types.Int64Value(response.IncidentNotificationSettings.DelayInMinVoice),
			SeveritiesVoice: MapNullableListWithEmpty(ctx, response.IncidentNotificationSettings.SeveritiesVoice),

			ShouldSendPush: types.BoolValue(response.IncidentNotificationSettings.ShouldSendPush),
			DelayInMinPush: types.Int64Value(response.IncidentNotificationSettings.DelayInMinPush),
			SeveritiesPush: MapNullableListWithEmpty(ctx, response.IncidentNotificationSettings.SeveritiesPush),

			ShouldSendEmail: types.BoolValue(response.IncidentNotificationSettings.ShouldSendEmail),
			DelayInMinEmail: types.Int64Value(response.IncidentNotificationSettings.DelayInMinEmail),
			SeveritiesEmail: MapNullableListWithEmpty(ctx, response.IncidentNotificationSettings.SeveritiesEmail),
		}
	}

}

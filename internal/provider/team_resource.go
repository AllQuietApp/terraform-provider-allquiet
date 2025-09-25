// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"

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
var _ resource.Resource = &Team{}
var _ resource.ResourceWithImportState = &Team{}

func NewTeam() resource.Resource {
	return &Team{}
}

// Team defines the resource implementation.
type Team struct {
	client *AllQuietAPIClient
}

type IncidentEngagementReportSettingsModel struct {
	DayOfWeek types.String `tfsdk:"day_of_week"`
	Time      types.String `tfsdk:"time"`
}

// TeamModel describes the resource data model.
type TeamModel struct {
	Id                               types.String                           `tfsdk:"id"`
	DisplayName                      types.String                           `tfsdk:"display_name"`
	TimeZoneId                       types.String                           `tfsdk:"time_zone_id"`
	IncidentEngagementReportSettings *IncidentEngagementReportSettingsModel `tfsdk:"incident_engagement_report_settings"`
	Labels                           types.List                             `tfsdk:"labels"`
}

func (r *Team) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (r *Team) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The team resource represents a team in All Quiet. Teams are used to group members and define schedules, its tiers and rorations.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the team",
				Required:            true,
			},
			"time_zone_id": schema.StringAttribute{
				MarkdownDescription: "The timezone id, defaults to 'UTC' if not provided. Find all timezone ids [here](https://allquiet.app/api/public/v1/timezone)",
				Optional:            true,
				Default:             stringdefault.StaticString("UTC"),
				Computed:            true,
			},
			"incident_engagement_report_settings": schema.SingleNestedAttribute{
				MarkdownDescription: "Settings when to send the incident report for the team",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"day_of_week": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Which day of the week to send the report",
						Validators:          []validator.String{stringvalidator.OneOf([]string{"sun", "mon", "tue", "wed", "thu", "fri", "sat"}...)},
					},
					"time": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Time of the day to send the report",
						Validators: []validator.String{stringvalidator.RegexMatches(
							regexp.MustCompile(`^([01]\d|2[0-3]):([0-5]\d)$`),
							"must contain time matching the pattern '^([01]\\d|2[0-3]):([0-5]\\d)$'",
						)},
					},
				},
			},
			"labels": schema.ListAttribute{
				MarkdownDescription: "The labels of the team",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *Team) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *Team) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TeamModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	teamResponse, err := r.client.CreateTeamResource(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create team resource, got error: %s", err))
		return
	}
	mapTeamResponseToModel(ctx, teamResponse, &data)

	tflog.Trace(ctx, "created team resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Team) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TeamModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	teamResponse, err := r.client.GetTeamResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get team resource, got error: %s", err))
		return
	}

	if teamResponse == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to get team resource, got nil response")
		return
	}

	mapTeamResponseToModel(ctx, teamResponse, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Team) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TeamModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	teamResponse, err := r.client.UpdateTeamResource(ctx, data.Id.ValueString(), &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update team resource, got error: %s", err))
		return
	}

	mapTeamResponseToModel(ctx, teamResponse, &data)

	tflog.Trace(ctx, "updated team resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Team) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TeamModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTeamResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update team resource, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted team resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Team) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapTeamResponseToModel(ctx context.Context, response *teamResponse, data *TeamModel) {
	data.Id = types.StringValue(response.Id)
	data.DisplayName = types.StringValue(response.DisplayName)
	data.TimeZoneId = types.StringValue(response.TimeZoneId)
	data.Labels = MapNullableList(ctx, response.Labels)
	if response.IncidentEngagementReportSettings != nil {
		data.IncidentEngagementReportSettings = &IncidentEngagementReportSettingsModel{
			DayOfWeek: types.StringValue(response.IncidentEngagementReportSettings.DayOfWeek),
			Time:      types.StringValue(response.IncidentEngagementReportSettings.Time),
		}
	}

}

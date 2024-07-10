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

type TeamMemberModel struct {
	Email types.String `tfsdk:"email"`
	Role  types.String `tfsdk:"role"`
}

type TeamTierModel struct {
	AutoEscalationAfterMinutes types.Int64         `tfsdk:"auto_escalation_after_minutes"`
	Schedules                  []TeamScheduleModel `tfsdk:"schedules"`
}

type TeamScheduleModel struct {
	ScheduleSettings *TeamScheduleSettingsModel `tfsdk:"schedule_settings"`
	RotationSettings *TeamRotationSettingsModel `tfsdk:"rotation_settings"`
	Rotations        []TeamRotationModel        `tfsdk:"rotations"`
}

type TeamRotationModel struct {
	Members []TeamRotationMemberModel `tfsdk:"members"`
}

type TeamRotationMemberModel struct {
	Email types.String `tfsdk:"email"`
}

type TeamScheduleSettingsModel struct {
	Start        types.String `tfsdk:"start"`
	End          types.String `tfsdk:"end"`
	SelectedDays types.List   `tfsdk:"selected_days"`
}

type TeamRotationSettingsModel struct {
	Repeats             types.String `tfsdk:"repeats"`
	StartsOnDayOfWeek   types.String `tfsdk:"starts_on_day_of_week"`
	StartsOnDateOfMonth types.Int64  `tfsdk:"starts_on_date_of_month"`
	StartsOnTime        types.String `tfsdk:"starts_on_time"`
	CustomRepeatUnit    types.String `tfsdk:"custom_repeat_unit"`
	CustomRepeatValue   types.Int64  `tfsdk:"custom_repeat_value"`
}

// TeamModel describes the resource data model.
type TeamModel struct {
	Id                               types.String                           `tfsdk:"id"`
	DisplayName                      types.String                           `tfsdk:"display_name"`
	TimeZoneId                       types.String                           `tfsdk:"time_zone_id"`
	IncidentEngagementReportSettings *IncidentEngagementReportSettingsModel `tfsdk:"incident_engagement_report_settings"`
	Members                          []TeamMemberModel                      `tfsdk:"members"`
	Tiers                            []TeamTierModel                        `tfsdk:"tiers"`
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
			"members": schema.SetNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"email": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Email of the member",
							Validators: []validator.String{stringvalidator.RegexMatches(
								regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
								"must contain email matching the pattern '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$'",
							)},
						},
						"role": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Role of the member (either 'Member' or 'Administrator')",
							Validators:          []validator.String{stringvalidator.OneOf([]string{"Member", "Administrator"}...)},
						},
					},
				},
			},
			"tiers": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"auto_escalation_after_minutes": schema.Int64Attribute{
							Optional:            true,
							MarkdownDescription: "After how many minutes the incident should be escalated to the next tier.",
							Validators: []validator.Int64{
								int64validator.Between(0, 60*24*30),
							},
						},
						"schedules": schema.ListNestedAttribute{
							Required: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"rotations": schema.ListNestedAttribute{
										Required: true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"members": schema.ListNestedAttribute{
													Required: true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"email": schema.StringAttribute{
																Required:            true,
																MarkdownDescription: "Email of the member",
																Validators: []validator.String{stringvalidator.RegexMatches(
																	regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
																	"must contain email matching the pattern '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$'",
																)},
															},
														},
													},
												},
											},
										},
									},
									"schedule_settings": schema.SingleNestedAttribute{
										MarkdownDescription: "Settings for the schedule",
										Optional:            true,
										Attributes: map[string]schema.Attribute{
											"start": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "Start time of the schedule",
												Validators: []validator.String{stringvalidator.RegexMatches(
													regexp.MustCompile(`^([01]\d|2[0-3]):([0-5]\d)$`),
													"must contain time matching the pattern '^([01]\\d|2[0-3]):([0-5]\\d)$'",
												)},
											},
											"end": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "End time of the schedule",
												Validators: []validator.String{stringvalidator.RegexMatches(
													regexp.MustCompile(`^([01]\d|2[0-3]):([0-5]\d)$`),
													"must contain time matching the pattern '^([01]\\d|2[0-3]):([0-5]\\d)$'",
												)},
											},
											"selected_days": schema.ListAttribute{
												Optional:            true,
												MarkdownDescription: "Selected days of the week",
												ElementType:         types.StringType,
												Validators: []validator.List{
													listvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"sun", "mon", "tue", "wed", "thu", "fri", "sat"}...)),
												},
											},
										},
									},
									"rotation_settings": schema.SingleNestedAttribute{
										MarkdownDescription: "Settings for the rotation",
										Optional:            true,
										Attributes: map[string]schema.Attribute{
											"repeats": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "The rotation will repeat on the given interval",
												Validators:          []validator.String{stringvalidator.OneOf([]string{"daily", "weekly", "biweekly", "monthly", "custom"}...)},
											},
											"starts_on_day_of_week": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "Starts on day of the week. Needs to be set if 'repeats' is not 'monthly'",
												Validators:          []validator.String{stringvalidator.OneOf([]string{"sun", "mon", "tue", "wed", "thu", "fri", "sat"}...)},
											},
											"starts_on_date_of_month": schema.Int64Attribute{
												Optional:            true,
												MarkdownDescription: "If set, starts on date of the month. Needs to be set if 'repeats' is 'monthly'",
												Validators:          []validator.Int64{int64validator.Between(1, 31)},
											},
											"starts_on_time": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "If set, starts on time of day. Needs to be set if 'repeats' is 'custom' and 'custom_repeat_unit' is 'hours'",
												Validators: []validator.String{stringvalidator.RegexMatches(
													regexp.MustCompile(`^([01]\d|2[0-3]):([0-5]\d)$`),
													"must contain time matching the pattern '^([01]\\d|2[0-3]):([0-5]\\d)$'",
												)},
											},
											"custom_repeat_unit": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "In what interval unit the rotation should repeat. Needs to be set if 'repeats' is 'custom'",
												Validators:          []validator.String{stringvalidator.OneOf([]string{"months", "weeks", "days", "hours"}...)},
											},
											"custom_repeat_value": schema.Int64Attribute{
												Optional:            true,
												MarkdownDescription: "How often the rotation should repeat. Needs to be set if 'repeats' is 'custom'",
												Validators:          []validator.Int64{int64validator.Between(1, 365)},
											},
										},
									},
								},
							},
						},
					},
				},
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
	if response.IncidentEngagementReportSettings != nil {
		data.IncidentEngagementReportSettings = &IncidentEngagementReportSettingsModel{
			DayOfWeek: types.StringValue(response.IncidentEngagementReportSettings.DayOfWeek),
			Time:      types.StringValue(response.IncidentEngagementReportSettings.Time),
		}
	}
	data.Members = mapTeamMembersResponseToData(response.Members)
	data.Tiers = mapTeamTiersResponseToData(ctx, response.Tiers)

}

func mapTeamTiersResponseToData(ctx context.Context, data []teamTier) []TeamTierModel {
	var tiers []TeamTierModel
	for _, tier := range data {

		var autoEscalationAfterMinutes types.Int64
		if tier.AutoEscalationAfterMinutes != nil {
			autoEscalationAfterMinutes = types.Int64PointerValue(tier.AutoEscalationAfterMinutes)
		} else {
			autoEscalationAfterMinutes = types.Int64Null()
		}

		tiers = append(tiers, TeamTierModel{
			AutoEscalationAfterMinutes: autoEscalationAfterMinutes,
			Schedules:                  mapTeamSchedulesResponseToData(ctx, tier.Schedules),
		})
	}
	return tiers
}

func mapTeamSchedulesResponseToData(ctx context.Context, teamSchedule []teamSchedule) []TeamScheduleModel {
	var schedules []TeamScheduleModel
	for _, schedule := range teamSchedule {
		schedules = append(schedules, TeamScheduleModel{
			ScheduleSettings: mapTeamScheduleSettingsResponseToData(ctx, schedule.ScheduleSettings),
			RotationSettings: mapTeamRotationSettingsResponseToData(schedule.RotationSettings),
			Rotations:        mapTeamRotationsResponseToData(schedule.Rotations),
		})
	}
	return schedules
}

func mapTeamRotationsResponseToData(data []teamRotation) []TeamRotationModel {
	var rotations []TeamRotationModel
	for _, rotation := range data {
		rotations = append(rotations, TeamRotationModel{
			Members: mapTeamRotationMembersResponseToData(rotation.Members),
		})
	}
	return rotations
}

func mapTeamRotationMembersResponseToData(data []rotationMember) []TeamRotationMemberModel {
	var members []TeamRotationMemberModel
	for _, member := range data {
		members = append(members, TeamRotationMemberModel{
			Email: types.StringValue(member.Email),
		})
	}
	return members
}

func mapTeamRotationSettingsResponseToData(rotationSettings *rotationSettings) *TeamRotationSettingsModel {
	if rotationSettings == nil {
		return nil
	}
	return &TeamRotationSettingsModel{
		Repeats:             types.StringPointerValue(rotationSettings.Repeats),
		StartsOnDayOfWeek:   types.StringPointerValue(rotationSettings.StartsOnDayOfWeek),
		StartsOnDateOfMonth: types.Int64PointerValue(rotationSettings.StartsOnDateOfMonth),
		StartsOnTime:        types.StringPointerValue(rotationSettings.StartsOnTime),
		CustomRepeatUnit:    types.StringPointerValue(rotationSettings.CustomRepeatUnit),
		CustomRepeatValue:   types.Int64PointerValue(rotationSettings.CustomRepeatValue),
	}
}

func mapTeamScheduleSettingsResponseToData(ctx context.Context, scheduleSettings *scheduleSettings) *TeamScheduleSettingsModel {
	if scheduleSettings == nil {
		return nil
	}
	return &TeamScheduleSettingsModel{
		Start:        types.StringPointerValue(scheduleSettings.Start),
		End:          types.StringPointerValue(scheduleSettings.End),
		SelectedDays: MapNullableList(ctx, &scheduleSettings.SelectedDays),
	}
}

func mapTeamMembersResponseToData(data []teamMember) []TeamMemberModel {
	var members []TeamMemberModel
	for _, member := range data {
		members = append(members, TeamMemberModel{
			Email: types.StringValue(member.Email),
			Role:  types.StringValue(member.Role),
		})
	}
	return members
}

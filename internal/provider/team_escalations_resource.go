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

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TeamEscalations{}
var _ resource.ResourceWithImportState = &TeamEscalations{}

func NewTeamEscalations() resource.Resource {
	return &TeamEscalations{}
}

// TeamEscalations defines the resource implementation.
type TeamEscalations struct {
	client *AllQuietAPIClient
}

type TeamEscalationsTierModel struct {
	AutoEscalationEnabled        types.Bool                        `tfsdk:"auto_escalation_enabled"`
	AutoEscalationAfterMinutes   types.Int64                       `tfsdk:"auto_escalation_after_minutes"`
	AutoEscalationStopMode       types.String                      `tfsdk:"auto_escalation_stop_mode"`
	AutoEscalationSeverities     types.List                        `tfsdk:"auto_escalation_severities"`
	AutoEscalationTimeFilters    *[]TeamEscalationsTimeFilterModel `tfsdk:"auto_escalation_time_filters"`
	Repeats                      types.Int64                       `tfsdk:"repeats"`
	RepeatsAfterMinutes          types.Int64                       `tfsdk:"repeats_after_minutes"`
	RepeatsStopMode              types.String                      `tfsdk:"repeats_stop_mode"`
	AutoAssignToTeams            types.List                        `tfsdk:"auto_assign_to_teams"`
	AutoAssignToTeamsSeverities  types.List                        `tfsdk:"auto_assign_to_teams_severities"`
	AutoAssignToTeamsTimeFilters *[]TeamEscalationsTimeFilterModel `tfsdk:"auto_assign_to_teams_time_filters"`
	Schedules                    []TeamEscalationsScheduleModel    `tfsdk:"schedules"`
}

type TeamEscalationsTimeFilterModel struct {
	SelectedDays types.List   `tfsdk:"selected_days"`
	From         types.String `tfsdk:"from"`
	Until        types.String `tfsdk:"until"`
}

type TeamEscalationsScheduleModel struct {
	ScheduleSettings *TeamEscalationsScheduleSettingsModel `tfsdk:"schedule_settings"`
	RotationSettings *TeamEscalationsRotationSettingsModel `tfsdk:"rotation_settings"`
	Rotations        []TeamEscalationsRotationModel        `tfsdk:"rotations"`
}

type TeamEscalationsRotationModel struct {
	Members []TeamEscalationsRotationMemberModel `tfsdk:"members"`
}

type TeamEscalationsRotationMemberModel struct {
	TeamMembershipId types.String `tfsdk:"team_membership_id"`
}

type TeamEscalationsScheduleSettingsModel struct {
	Start           types.String                          `tfsdk:"start"`
	End             types.String                          `tfsdk:"end"`
	SelectedDays    types.List                            `tfsdk:"selected_days"`
	WeeklySchedules *[]TeamEscalationsWeeklyScheduleModel `tfsdk:"weekly_schedules"`
}

type TeamEscalationsWeeklyScheduleModel struct {
	SelectedDays types.List   `tfsdk:"selected_days"`
	From         types.String `tfsdk:"from"`
	Until        types.String `tfsdk:"until"`
}

type TeamEscalationsRotationSettingsModel struct {
	Repeats             types.String `tfsdk:"repeats"`
	StartsOnDayOfWeek   types.String `tfsdk:"starts_on_day_of_week"`
	StartsOnDateOfMonth types.Int64  `tfsdk:"starts_on_date_of_month"`
	StartsOnTime        types.String `tfsdk:"starts_on_time"`
	CustomRepeatUnit    types.String `tfsdk:"custom_repeat_unit"`
	CustomRepeatValue   types.Int64  `tfsdk:"custom_repeat_value"`
	EffectiveFrom       types.String `tfsdk:"effective_from"`
	RotationMode        types.String `tfsdk:"rotation_mode"`
	AutoRotationSize    types.Int64  `tfsdk:"auto_rotation_size"`
}

// TeamEscalationsModel describes the resource data model.
type TeamEscalationsModel struct {
	Id              types.String               `tfsdk:"id"`
	TeamId          types.String               `tfsdk:"team_id"`
	EscalationTiers []TeamEscalationsTierModel `tfsdk:"escalation_tiers"`
	TierSettings    *TierSettingsModel         `tfsdk:"tier_settings"`
}

type TierSettingsModel struct {
	Repeats             types.Int64  `tfsdk:"repeats"`
	RepeatsAfterMinutes types.Int64  `tfsdk:"repeats_after_minutes"`
	RepeatsStopMode     types.String `tfsdk:"repeats_stop_mode"`
}

func (r *TeamEscalations) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_escalations"
}

func (r *TeamEscalations) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The Team Escalations resource represents a Team's Escalation Tiers in All Quiet. Escalation Tiers are used to group members and define schedules, its tiers and rorations.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"team_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Id of the associated team",
			},
			"tier_settings": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"repeats": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "How many times all tiers should repeat.",
						Validators: []validator.Int64{
							int64validator.Between(0, 16),
						},
					},
					"repeats_after_minutes": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "How many minutes after the last tier all tiers should repeat.",
						Validators: []validator.Int64{
							int64validator.Between(1, OneMonthInSeconds),
						},
					},
					"repeats_stop_mode": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "When all tiers should stop repeating. Possible values are: " + strings.Join(ValidEscalationModes, ", "),
						Validators:          []validator.String{stringvalidator.OneOf(ValidEscalationModes...)},
					},
				},
			},
			"escalation_tiers": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"auto_escalation_enabled": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Whether auto-escalation is enabled for this tier.",
						},
						"auto_escalation_after_minutes": schema.Int64Attribute{
							Optional:            true,
							MarkdownDescription: "Escalation cadence in minutes. After how many minutes the incident should be auto-escalated to the next tier or auto-assigned to teams.",
							Validators: []validator.Int64{
								int64validator.Between(0, 60*24*30),
							},
						},
						"auto_escalation_stop_mode": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "When the escalation should be stopped. Possible values are: " + strings.Join(ValidEscalationModes, ", "),
							Validators:          []validator.String{stringvalidator.OneOf(ValidEscalationModes...)},
						},
						"auto_escalation_severities": schema.ListAttribute{
							Optional:            true,
							MarkdownDescription: "Severities that should trigger auto-escalation. Possible values are: " + strings.Join(ValidSeverities, ", "),
							ElementType:         types.StringType,
							Validators: []validator.List{
								listvalidator.ValueStringsAre(SeverityValidator("Not a valid severity")),
							},
						},
						"auto_escalation_time_filters": schema.ListNestedAttribute{
							Optional:            true,
							MarkdownDescription: "Time filters in which auto-escalation should be triggered.",
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
								},
							},
						},
						"auto_assign_to_teams": schema.ListAttribute{
							Optional:            true,
							MarkdownDescription: "Team IDs that should be auto-assigned to.",
							ElementType:         types.StringType,
						},
						"auto_assign_to_teams_severities": schema.ListAttribute{
							Optional:            true,
							MarkdownDescription: "Severities that should trigger auto-assign to teams. Possible values are: " + strings.Join(ValidSeverities, ", "),
							ElementType:         types.StringType,
							Validators: []validator.List{
								listvalidator.ValueStringsAre(SeverityValidator("Not a valid severity")),
							},
						},
						"auto_assign_to_teams_time_filters": schema.ListNestedAttribute{
							Optional:            true,
							MarkdownDescription: "Time filters in which auto-assign to teams should be triggered.",
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
								},
							},
						},
						"repeats": schema.Int64Attribute{
							Optional:            true,
							MarkdownDescription: "How many times the rotation should repeat.",
							Validators: []validator.Int64{
								int64validator.Between(0, 16),
							},
						},
						"repeats_after_minutes": schema.Int64Attribute{
							Optional:            true,
							MarkdownDescription: "How many minutes after the rotation should repeat.",
							Validators: []validator.Int64{
								int64validator.Between(1, OneMonthInSeconds),
							},
						},
						"repeats_stop_mode": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "When this tier should stop being repeated. Possible values are: " + strings.Join(ValidEscalationModes, ", "),
							Validators:          []validator.String{stringvalidator.OneOf(ValidEscalationModes...)},
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
															"team_membership_id": schema.StringAttribute{
																Required:            true,
																MarkdownDescription: "Id of the team membership",
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
												MarkdownDescription: "Start time of the schedule. Format: HH:mm",
												DeprecationMessage:  "Use weekly_schedules instead",
												Validators:          []validator.String{TimeValidator("Not a valid time")},
											},
											"end": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "End time of the schedule. Format: HH:mm",
												DeprecationMessage:  "Use weekly_schedules instead",
												Validators:          []validator.String{TimeValidator("Not a valid time")},
											},
											"selected_days": schema.ListAttribute{
												Optional:            true,
												DeprecationMessage:  "Use weekly_schedules instead",
												MarkdownDescription: "Selected days of the week. Possible values are: " + strings.Join(ValidDaysOfWeek, ", "),
												ElementType:         types.StringType,
												Validators: []validator.List{
													listvalidator.ValueStringsAre(DaysOfWeekValidator("Not a valid day of week")),
												},
											},
											"weekly_schedules": schema.ListNestedAttribute{
												Optional:            true,
												MarkdownDescription: "Weekly schedules",
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
													},
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
												MarkdownDescription: "The rotation will repeat on the given interval. Possible values are: " + strings.Join(ValidRotationRepeats, ", "),
												Validators:          []validator.String{stringvalidator.OneOf(ValidRotationRepeats...)},
											},
											"starts_on_day_of_week": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "Starts on day of the week. Needs to be set if 'repeats' is not 'monthly'. Possible values are: " + strings.Join(ValidDaysOfWeek, ", "),
												Validators:          []validator.String{DaysOfWeekValidator("Not a valid day of week")},
											},
											"starts_on_date_of_month": schema.Int64Attribute{
												Optional:            true,
												MarkdownDescription: "If set, starts on date of the month. Needs to be set if 'repeats' is 'monthly'",
												Validators:          []validator.Int64{int64validator.Between(1, 31)},
											},
											"starts_on_time": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "If set, starts on time of day. Needs to be set if 'repeats' is 'custom' and 'custom_repeat_unit' is 'hours'. Format: HH:mm",
												Validators:          []validator.String{TimeValidator("Not a valid time")},
											},
											"custom_repeat_unit": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "In what interval unit the rotation should repeat. Needs to be set if 'repeats' is 'custom'. Possible values are: " + strings.Join(ValidCustomRepeatUnits, ", "),
												Validators:          []validator.String{stringvalidator.OneOf(ValidCustomRepeatUnits...)},
											},
											"custom_repeat_value": schema.Int64Attribute{
												Optional:            true,
												MarkdownDescription: "How often the rotation should repeat. Needs to be set if 'repeats' is 'custom'",
												Validators:          []validator.Int64{int64validator.Between(1, 365)},
											},
											"effective_from": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "If sets, the rotation will be effective from the given date in ISO 8601 format",
												Validators: []validator.String{stringvalidator.RegexMatches(
													regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`),
													"must contain ISO date matching the pattern '^\\d{4}-\\d{2}-\\d{2}$'",
												)},
											},
											"rotation_mode": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "The mode of the rotation. Possible values are: " + strings.Join(ValidRotationModes, ", "),
												Validators:          []validator.String{stringvalidator.OneOf(ValidRotationModes...)},
											},
											"auto_rotation_size": schema.Int64Attribute{
												Optional:            true,
												MarkdownDescription: "The size of the rotation",
												Validators:          []validator.Int64{int64validator.Between(1, 500)},
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

func (r *TeamEscalations) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TeamEscalations) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TeamEscalationsModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	teamEscalationsResponse, err := r.client.CreateTeamEscalationsResource(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create teamEscalations resource, got error: %s", err))
		return
	}
	mapTeamEscalationsResponseToModel(ctx, teamEscalationsResponse, &data)

	tflog.Trace(ctx, "created teamEscalations resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamEscalations) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TeamEscalationsModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	teamEscalationsResponse, err := r.client.GetTeamEscalationsResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get teamEscalations resource, got error: %s", err))
		return
	}

	if teamEscalationsResponse == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to get teamEscalations resource, got nil response")
		return
	}

	mapTeamEscalationsResponseToModel(ctx, teamEscalationsResponse, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamEscalations) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TeamEscalationsModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	teamEscalationsResponse, err := r.client.UpdateTeamEscalationsResource(ctx, data.Id.ValueString(), &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update teamEscalations resource, got error: %s", err))
		return
	}

	mapTeamEscalationsResponseToModel(ctx, teamEscalationsResponse, &data)

	tflog.Trace(ctx, "updated teamEscalations resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamEscalations) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TeamEscalationsModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTeamEscalationsResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update teamEscalations resource, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted teamEscalations resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamEscalations) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapTeamEscalationsResponseToModel(ctx context.Context, response *teamEscalationsResponse, data *TeamEscalationsModel) {
	data.Id = types.StringValue(response.Id)
	data.TeamId = types.StringValue(response.TeamId)
	data.EscalationTiers = mapTeamEscalationsTiersResponseToData(ctx, response.EscalationTiers)
}

func mapTeamEscalationsTiersResponseToData(ctx context.Context, data []teamEscalationsTier) []TeamEscalationsTierModel {
	var tiers []TeamEscalationsTierModel
	for _, tier := range data {

		var autoEscalationAfterMinutes types.Int64
		if tier.AutoEscalationAfterMinutes != nil {
			autoEscalationAfterMinutes = types.Int64PointerValue(tier.AutoEscalationAfterMinutes)
		} else {
			autoEscalationAfterMinutes = types.Int64Null()
		}

		tiers = append(tiers, TeamEscalationsTierModel{
			AutoEscalationEnabled:        types.BoolPointerValue(tier.AutoEscalationEnabled),
			AutoEscalationAfterMinutes:   autoEscalationAfterMinutes,
			AutoEscalationStopMode:       types.StringPointerValue(tier.AutoEscalationStopMode),
			AutoEscalationSeverities:     MapNullableList(ctx, tier.AutoEscalationSeverities),
			AutoEscalationTimeFilters:    mapTeamEscalationsTimeFiltersToData(ctx, tier.AutoEscalationTimeFilters),
			AutoAssignToTeams:            MapNullableList(ctx, tier.AutoAssignToTeams),
			AutoAssignToTeamsSeverities:  MapNullableList(ctx, tier.AutoAssignToTeamsSeverities),
			AutoAssignToTeamsTimeFilters: mapTeamEscalationsTimeFiltersToData(ctx, tier.AutoAssignToTeamsTimeFilters),
			Repeats:                      types.Int64PointerValue(tier.Repeats),
			RepeatsAfterMinutes:          types.Int64PointerValue(tier.RepeatsAfterMinutes),
			RepeatsStopMode:              types.StringPointerValue(tier.RepeatsStopMode),
			Schedules:                    mapTeamEscalationsSchedulesResponseToData(ctx, tier.Schedules),
		})
	}
	return tiers
}

func mapTeamEscalationsTimeFiltersToData(ctx context.Context, timeFilters *[]teamEscalationsTimeFilter) *[]TeamEscalationsTimeFilterModel {
	var timeFiltersData []TeamEscalationsTimeFilterModel
	if timeFilters == nil {
		return nil
	}
	for _, timeFilter := range *timeFilters {
		timeFiltersData = append(timeFiltersData, TeamEscalationsTimeFilterModel{
			SelectedDays: MapNullableList(ctx, timeFilter.SelectedDays),
			From:         types.StringPointerValue(timeFilter.From),
			Until:        types.StringPointerValue(timeFilter.Until),
		})
	}
	return &timeFiltersData
}

func mapTeamEscalationsSchedulesResponseToData(ctx context.Context, teamEscalationsSchedule []teamEscalationsSchedule) []TeamEscalationsScheduleModel {
	var schedules []TeamEscalationsScheduleModel
	for _, schedule := range teamEscalationsSchedule {
		schedules = append(schedules, TeamEscalationsScheduleModel{
			ScheduleSettings: mapTeamEscalationsScheduleSettingsResponseToData(ctx, schedule.ScheduleSettings),
			RotationSettings: mapTeamEscalationsRotationSettingsResponseToData(schedule.RotationSettings),
			Rotations:        mapTeamEscalationsRotationsResponseToData(schedule.Rotations),
		})
	}
	return schedules
}

func mapTeamEscalationsRotationsResponseToData(data []teamEscalationsRotation) []TeamEscalationsRotationModel {
	var rotations []TeamEscalationsRotationModel
	for _, rotation := range data {
		rotations = append(rotations, TeamEscalationsRotationModel{
			Members: mapTeamEscalationsRotationMembersResponseToData(rotation.Members),
		})
	}
	return rotations
}

func mapTeamEscalationsRotationMembersResponseToData(data []teamEscalationsRotationMember) []TeamEscalationsRotationMemberModel {
	var members []TeamEscalationsRotationMemberModel
	for _, member := range data {
		members = append(members, TeamEscalationsRotationMemberModel{
			TeamMembershipId: types.StringValue(member.TeamMembershipId),
		})
	}
	return members
}

func mapTeamEscalationsRotationSettingsResponseToData(rotationSettings *rotationSettings) *TeamEscalationsRotationSettingsModel {
	if rotationSettings == nil {
		return nil
	}
	return &TeamEscalationsRotationSettingsModel{
		Repeats:             types.StringPointerValue(rotationSettings.Repeats),
		StartsOnDayOfWeek:   types.StringPointerValue(rotationSettings.StartsOnDayOfWeek),
		StartsOnDateOfMonth: types.Int64PointerValue(rotationSettings.StartsOnDateOfMonth),
		StartsOnTime:        types.StringPointerValue(rotationSettings.StartsOnTime),
		CustomRepeatUnit:    types.StringPointerValue(rotationSettings.CustomRepeatUnit),
		CustomRepeatValue:   types.Int64PointerValue(rotationSettings.CustomRepeatValue),
		EffectiveFrom:       types.StringPointerValue(rotationSettings.EffectiveFrom),
		RotationMode:        types.StringPointerValue(rotationSettings.RotationMode),
		AutoRotationSize:    types.Int64PointerValue(rotationSettings.AutoRotationSize),
	}
}

func mapTeamEscalationsScheduleSettingsResponseToData(ctx context.Context, scheduleSettings *scheduleSettings) *TeamEscalationsScheduleSettingsModel {
	if scheduleSettings == nil {
		return nil
	}
	return &TeamEscalationsScheduleSettingsModel{
		Start:           types.StringPointerValue(scheduleSettings.Start),
		End:             types.StringPointerValue(scheduleSettings.End),
		SelectedDays:    MapNullableList(ctx, scheduleSettings.SelectedDays),
		WeeklySchedules: mapTeamEscalationsWeeklySchedulesResponseToData(ctx, scheduleSettings.WeeklySchedules),
	}
}

func mapTeamEscalationsWeeklySchedulesResponseToData(ctx context.Context, weeklySchedules *[]weeklySchedule) *[]TeamEscalationsWeeklyScheduleModel {
	if weeklySchedules == nil {
		return nil
	}
	var weeklySchedulesData []TeamEscalationsWeeklyScheduleModel
	for _, weeklySchedule := range *weeklySchedules {
		weeklySchedulesData = append(weeklySchedulesData, TeamEscalationsWeeklyScheduleModel{
			SelectedDays: MapNullableList(ctx, weeklySchedule.SelectedDays),
			From:         types.StringPointerValue(weeklySchedule.From),
			Until:        types.StringPointerValue(weeklySchedule.Until),
		})
	}
	return &weeklySchedulesData
}

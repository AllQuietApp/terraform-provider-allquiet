// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &Routing{}
var _ resource.ResourceWithImportState = &Routing{}

func NewRouting() resource.Resource {
	return &Routing{}
}

// Routing defines the resource implementation.
type Routing struct {
	client *AllQuietAPIClient
}

// RoutingModel describes the resource data model.
type RoutingModel struct {
	Id                     types.String            `tfsdk:"id"`
	DisplayName            types.String            `tfsdk:"display_name"`
	TeamId                 types.String            `tfsdk:"team_id"`
	Rules                  []RoutingRuleModel      `tfsdk:"rules"`
	TeamConnectionSettings *TeamConnectionSettings `tfsdk:"team_connection_settings"`
}

type RoutingRuleModel struct {
	Conditions *RoutingRuleConditionsModel `tfsdk:"conditions"`
	Actions    *RoutingRuleActionsModel    `tfsdk:"actions"`
	Channels   *RoutingRuleChannelsModel   `tfsdk:"channels"`
}

type RoutingRuleConditionsModel struct {
	Statuses        types.List                            `tfsdk:"statuses"`
	Severities      types.List                            `tfsdk:"severities"`
	Integrations    types.List                            `tfsdk:"integrations"`
	Intents         types.List                            `tfsdk:"intents"`
	Attributes      []RoutingRuleConditionsAttributeModel `tfsdk:"attributes"`
	DateRestriction *DateRestrictionModel                 `tfsdk:"date_restriction"`
	Schedule        *ScheduleModel                        `tfsdk:"schedule"`
}

type RoutingRuleConditionsAttributeModel struct {
	Name     types.String `tfsdk:"name"`
	Operator types.String `tfsdk:"operator"`
	Value    types.String `tfsdk:"value"`
}

type RoutingRuleActionsModel struct {
	AssignToTeams                 types.List                             `tfsdk:"assign_to_teams"`
	Discard                       types.Bool                             `tfsdk:"discard"`
	ChangeSeverity                types.String                           `tfsdk:"change_severity"`
	AddInteraction                types.String                           `tfsdk:"add_interaction"`
	RuleFlowControl               types.String                           `tfsdk:"rule_flow_control"`
	DelayActionsInMinutes         types.Int64                            `tfsdk:"delay_actions_in_minutes"`
	AffectsServices               types.List                             `tfsdk:"affects_services"`
	ForwardToOutboundIntegrations types.List                             `tfsdk:"forward_to_outbound_integrations"`
	SetAttributes                 []RoutingRuleActionsSetAttributesModel `tfsdk:"set_attributes"`
	SnoozeForRelativeInMinutes    types.Int64                            `tfsdk:"snooze_for_relative_in_minutes"`
	SnoozeUntilAbsolute           types.String                           `tfsdk:"snooze_until_absolute"`
	SnoozeUntilWeekdayAbsolute    types.String                           `tfsdk:"snooze_until_weekday_absolute"`
}

type RoutingRuleActionsSetAttributesModel struct {
	Name           types.String `tfsdk:"name"`
	Value          types.String `tfsdk:"value"`
	IsImage        types.Bool   `tfsdk:"is_image"`
	HideInPreviews types.Bool   `tfsdk:"hide_in_previews"`
}

type RoutingRuleChannelsModel struct {
	OutboundIntegrations      types.List `tfsdk:"outbound_integrations"`
	OutboundIntegrationsMuted types.Bool `tfsdk:"outbound_integrations_muted"`
	NotificationChannels      types.List `tfsdk:"notification_channels"`
	NotificationChannelsMuted types.Bool `tfsdk:"notification_channels_muted"`
}

type DateRestrictionModel struct {
	From  types.String `tfsdk:"from"`
	Until types.String `tfsdk:"until"`
}

type ScheduleModel struct {
	After      types.String `tfsdk:"after"`
	Before     types.String `tfsdk:"before"`
	DaysOfWeek types.List   `tfsdk:"days_of_week"`
}

func (r *Routing) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_routing"
}

func (r *Routing) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The routing resource allows you to define routing rules for incidents based on various conditions and actions.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the routing",
				Required:            true,
			},
			"team_id": schema.StringAttribute{
				MarkdownDescription: "The team id of the routing",
				Required:            true,
			},
			"rules": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"conditions": schema.SingleNestedAttribute{
							MarkdownDescription: "Settings for the schedule",
							Required:            true,
							Attributes: map[string]schema.Attribute{
								"statuses": schema.ListAttribute{
									Optional:            true,
									MarkdownDescription: "Statuses. Possible values are: " + strings.Join(ValidStatuses, ", "),
									ElementType:         types.StringType,
									Validators: []validator.List{
										listvalidator.ValueStringsAre(StatusValidator("Not a valid status")),
									},
								},
								"severities": schema.ListAttribute{
									Optional:            true,
									MarkdownDescription: "Severeties. Possible values are: " + strings.Join(ValidSeverities, ", "),
									ElementType:         types.StringType,
									Validators: []validator.List{
										listvalidator.ValueStringsAre(SeverityValidator("Not a valid severity")),
									},
								},
								"integrations": schema.ListAttribute{
									Optional:            true,
									MarkdownDescription: "Integration IDs",
									ElementType:         types.StringType,
									Validators: []validator.List{
										listvalidator.ValueStringsAre(GuidValidator("Not a valid GUID")),
									},
								},
								"intents": schema.ListAttribute{
									Optional:            true,
									MarkdownDescription: "Intents. Possible values are: " + strings.Join(ValidIntents, ", "),
									ElementType:         types.StringType,
									Validators: []validator.List{
										listvalidator.ValueStringsAre(IntentValidator("Not a valid intent")),
									},
								},
								"attributes": schema.ListNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												MarkdownDescription: "The name of the attribute",
												Required:            true,
											},
											"operator": schema.StringAttribute{
												MarkdownDescription: "The operator. Possible values are: " + strings.Join(ValidOperators, ", "),
												Required:            true,
												Validators:          []validator.String{OperatorValidator("Not a valid operator")},
											},
											"value": schema.StringAttribute{
												MarkdownDescription: "The value of the attribute to match with the operator against",
												Optional:            true,
											},
										},
									},
								},
								"date_restriction": schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{
										"from": schema.StringAttribute{
											Optional:    true,
											Description: "Start date for the routing rule (RFC3339 format)",
											Validators:  []validator.String{DateTimeValidator("Not a valid date")},
										},
										"until": schema.StringAttribute{
											Optional:    true,
											Description: "End date for the routing rule (RFC3339 format)",
											Validators:  []validator.String{DateTimeValidator("Not a valid date")},
										},
									},
								}, "schedule": schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{
										"after": schema.StringAttribute{
											Optional:    true,
											Description: "Time after which the rule is active (HH:mm format)",
											Validators:  []validator.String{TimeValidator("Not a valid time")},
										},
										"before": schema.StringAttribute{
											Optional:    true,
											Description: "Time before which the rule is active (HH:mm format)",
											Validators:  []validator.String{TimeValidator("Not a valid time")},
										},
										"days_of_week": schema.ListAttribute{
											Optional:    true,
											ElementType: types.StringType,
											Description: "Days of the week when the rule is active. Possible values are: " + strings.Join(ValidDaysOfWeek, ", "),
											Validators: []validator.List{
												listvalidator.ValueStringsAre(DaysOfWeekValidator("Not a valid day of week")),
											},
										},
									},
								},
							},
						},
						"actions": schema.SingleNestedAttribute{
							MarkdownDescription: "Settings for the schedule",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"assign_to_teams": schema.ListAttribute{
									Optional:            true,
									MarkdownDescription: "Will assign the incident to the specified teams.",
									ElementType:         types.StringType,
									Validators: []validator.List{
										listvalidator.ValueStringsAre(GuidValidator("Not a valid GUID")),
									},
								},
								"discard": schema.BoolAttribute{
									Optional:            true,
									Default:             booldefault.StaticBool(false),
									Computed:            true,
									MarkdownDescription: "If true will discard and delete the incident",
								},
								"change_severity": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Will change the severity of the incident. Possible values are: " + strings.Join(ValidSeverities, ", "),
									Validators:          []validator.String{SeverityValidator("Not a valid severity")},
								},
								"add_interaction": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Will add an interaction. For instance, you can auto resolve an incident by adding an interaction of intent 'Resolved'. Possible values are: " + strings.Join(ValidIntents, ", "),
									Validators:          []validator.String{IntentValidator("Not a valid intent")},
								},
								"rule_flow_control": schema.StringAttribute{
									Optional:            true,
									Default:             stringdefault.StaticString("Continue"),
									Computed:            true,
									MarkdownDescription: "If 'Skip' will not evaluate further rules. Possible values are: " + strings.Join(ValidRuleFlowControl, ", "),
									Validators:          []validator.String{RuleFlowValidator("Not a valid rule flow value")},
								},
								"delay_actions_in_minutes": schema.Int64Attribute{
									Optional:            true,
									MarkdownDescription: "Delay actions in minutes",
								},
								"affects_services": schema.ListAttribute{
									Optional:            true,
									MarkdownDescription: "Will affect the specified services. Only with add_interaction 'Affects'.",
									ElementType:         types.StringType,
									Validators: []validator.List{
										listvalidator.ValueStringsAre(GuidValidator("Not a valid GUID")),
									},
								},
								"forward_to_outbound_integrations": schema.ListAttribute{
									Optional:            true,
									MarkdownDescription: "Will forward to the specified outbound integrations. Only with add_interaction 'Forwarded'.",
									ElementType:         types.StringType,
									Validators: []validator.List{
										listvalidator.ValueStringsAre(GuidValidator("Not a valid GUID")),
									},
								},
								"set_attributes": schema.ListNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												MarkdownDescription: "The name of the attribute",
												Required:            true,
											},
											"value": schema.StringAttribute{
												MarkdownDescription: "The value of the attribute",
												Required:            true,
											},
											"is_image": schema.BoolAttribute{
												Optional:            true,
												Default:             booldefault.StaticBool(false),
												Computed:            true,
												MarkdownDescription: "If true will display the value as an image if it's a URL",
											},
											"hide_in_previews": schema.BoolAttribute{
												Optional:            true,
												Default:             booldefault.StaticBool(false),
												Computed:            true,
												MarkdownDescription: "If true will hide the value in previews",
											},
										},
									},
								},
								"snooze_for_relative_in_minutes": schema.Int64Attribute{
									Optional:            true,
									MarkdownDescription: "Snooze for relative in minutes",
								},
								"snooze_until_absolute": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Snooze until absolute",
									Validators:          []validator.String{TimeValidator("Not a valid time")},
								},
								"snooze_until_weekday_absolute": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Snooze until weekday absolute. Possible values are: " + strings.Join(ValidDaysOfWeek, ", "),
									Validators:          []validator.String{DaysOfWeekValidator("Not a valid day of week")},
								},
							},
						},
						"channels": schema.SingleNestedAttribute{
							MarkdownDescription: "Settings for the schedule",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"outbound_integrations": schema.ListAttribute{
									Optional:            true,
									MarkdownDescription: "Outbound integrations",
									ElementType:         types.StringType,
								},
								"outbound_integrations_muted": schema.BoolAttribute{
									Optional:            true,
									Default:             booldefault.StaticBool(false),
									Computed:            true,
									MarkdownDescription: "If true will mute the outbound integrations",
								},
								"notification_channels": schema.ListAttribute{
									Optional:            true,
									MarkdownDescription: "Notification channels",
									ElementType:         types.StringType,
									Validators: []validator.List{
										listvalidator.ValueStringsAre(NotificationChannelValidator("Not a valid channel")),
									},
								},
								"notification_channels_muted": schema.BoolAttribute{
									Optional:            true,
									Default:             booldefault.StaticBool(false),
									Computed:            true,
									MarkdownDescription: "If true will mute the notification channels",
								},
							},
						},
					},
				},
			},
			"team_connection_settings": schema.SingleNestedAttribute{
				MarkdownDescription: "The team connection settings for the routing",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"team_connection_mode": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The team connection mode for the routing. Possible values are: " + strings.Join(ValidTeamConnectionModes, ", "),
						Validators:          []validator.String{stringvalidator.OneOf(ValidTeamConnectionModes...)},
					},
					"team_ids": schema.ListAttribute{
						MarkdownDescription: "The team ids for the routing. If not provided, team_connection_mode must be set to 'OrganizationTeams'.",
						Optional:            true,
						ElementType:         types.StringType,
						Validators: []validator.List{
							listvalidator.ValueStringsAre(GuidValidator("Not a valid GUID")),
						},
					},
				},
			},
		},
	}
}

func (r *Routing) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *Routing) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RoutingModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	routingResponse, err := r.client.CreateRoutingResource(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create routing resource, got error: %s", err))
		return
	}

	mapRoutingResponseToModel(ctx, routingResponse, &data)

	tflog.Trace(ctx, "created routing resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Routing) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RoutingModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	routingResponse, err := r.client.GetRoutingResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get routing resource, got error: %s", err))
		return
	}

	if routingResponse == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to get routing resource, got nil response")
		return
	}

	mapRoutingResponseToModel(ctx, routingResponse, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Routing) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RoutingModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	routingResponse, err := r.client.UpdateRoutingResource(ctx, data.Id.ValueString(), &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update routing resource, got error: %s", err))
		return
	}

	mapRoutingResponseToModel(ctx, routingResponse, &data)

	tflog.Trace(ctx, "updated routing resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Routing) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RoutingModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRoutingResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update routing resource, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted routing resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Routing) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapRoutingResponseToModel(ctx context.Context, response *routingResponse, data *RoutingModel) {

	data.Id = types.StringValue(response.Id)
	data.DisplayName = types.StringValue(response.DisplayName)
	data.TeamId = types.StringValue(response.TeamId)
	data.Rules = mapRoutingRuleResponseToModel(ctx, response.Rules)
	data.TeamConnectionSettings = MapTeamConnectionSettingsResponseToModel(ctx, response.TeamConnectionSettings)
}

func mapRoutingRuleResponseToModel(ctx context.Context, rules []routingRule) []RoutingRuleModel {
	var result []RoutingRuleModel

	for _, rule := range rules {
		result = append(result, RoutingRuleModel{
			Conditions: mapRoutingRuleConditionsResponseToModel(ctx, rule.Conditions),
			Actions:    mapRoutingRuleActionsResponseToModel(ctx, rule.Actions),
			Channels:   mapRoutingRuleChannelsResponseToModel(ctx, rule.Channels),
		})
	}

	return result
}

func mapRoutingRuleConditionsResponseToModel(ctx context.Context, conditions *routingRuleConditions) *RoutingRuleConditionsModel {
	if conditions == nil {
		return nil
	}

	return &RoutingRuleConditionsModel{
		Statuses:        MapNullableList(ctx, conditions.Statuses),
		Severities:      MapNullableList(ctx, conditions.Severities),
		Integrations:    MapNullableList(ctx, conditions.Integrations),
		Intents:         MapNullableList(ctx, conditions.Intents),
		Attributes:      mapRoutingRuleConditionsAttributeResponseToModel(conditions.Attributes),
		DateRestriction: mapRoutingRuleDateRestrictionResponseToModel(conditions.DateRestriction),
		Schedule:        mapRoutingRuleScheduleResponseToModel(ctx, conditions.Schedule),
	}
}

func mapRoutingRuleDateRestrictionResponseToModel(dateRestriction *routingRuleDateRestriction) *DateRestrictionModel {
	if dateRestriction == nil {
		return nil
	}

	return &DateRestrictionModel{
		From:  types.StringPointerValue(dateRestriction.From),
		Until: types.StringPointerValue(dateRestriction.Until),
	}
}

func mapRoutingRuleScheduleResponseToModel(ctx context.Context, schedule *routingRuleSchedule) *ScheduleModel {
	if schedule == nil {
		return nil
	}

	return &ScheduleModel{
		After:      types.StringPointerValue(schedule.After),
		Before:     types.StringPointerValue(schedule.Before),
		DaysOfWeek: MapNullableList(ctx, schedule.DaysOfWeek),
	}
}

func mapRoutingRuleConditionsAttributeResponseToModel(attributes []routingRuleAttribute) []RoutingRuleConditionsAttributeModel {
	var result []RoutingRuleConditionsAttributeModel

	for _, attribute := range attributes {
		result = append(result, RoutingRuleConditionsAttributeModel{
			Name:     types.StringValue(attribute.Name),
			Operator: types.StringValue(attribute.Operator),
			Value:    types.StringPointerValue(attribute.Value),
		})
	}

	return result
}

func mapRoutingRuleActionsResponseToModel(ctx context.Context, actions *routingRuleActions) *RoutingRuleActionsModel {

	if actions == nil {
		return nil
	}

	return &RoutingRuleActionsModel{
		AssignToTeams:                 MapNullableList(ctx, actions.AssignToTeams),
		Discard:                       types.BoolValue(actions.Discard),
		ChangeSeverity:                types.StringPointerValue(actions.ChangeSeverity),
		AddInteraction:                types.StringPointerValue(actions.AddInteraction),
		RuleFlowControl:               types.StringPointerValue(actions.RuleFlowControl),
		DelayActionsInMinutes:         types.Int64PointerValue(actions.DelayActionsInMinutes),
		AffectsServices:               MapNullableList(ctx, actions.AffectsServices),
		ForwardToOutboundIntegrations: MapNullableList(ctx, actions.ForwardToOutboundIntegrations),
		SetAttributes:                 mapRoutingRuleActionsSetAttributesResponseToModel(actions.SetAttributes),
		SnoozeForRelativeInMinutes:    types.Int64PointerValue(actions.SnoozeForRelativeInMinutes),
		SnoozeUntilAbsolute:           types.StringPointerValue(actions.SnoozeUntilAbsolute),
		SnoozeUntilWeekdayAbsolute:    types.StringPointerValue(actions.SnoozeUntilWeekdayAbsolute),
	}
}

func mapRoutingRuleActionsSetAttributesResponseToModel(attributes *[]routingRuleSetAttribute) []RoutingRuleActionsSetAttributesModel {
	var result []RoutingRuleActionsSetAttributesModel

	for _, attribute := range *attributes {
		result = append(result, RoutingRuleActionsSetAttributesModel{
			Name:           types.StringValue(attribute.Name),
			Value:          types.StringValue(attribute.Value),
			IsImage:        types.BoolValue(attribute.IsImage),
			HideInPreviews: types.BoolValue(attribute.HideInPreviews),
		})
	}

	return result
}

func mapRoutingRuleChannelsResponseToModel(ctx context.Context, channels *routingRuleChannels) *RoutingRuleChannelsModel {
	if channels == nil {
		return nil
	}

	return &RoutingRuleChannelsModel{
		OutboundIntegrations:      MapNullableList(ctx, channels.OutboundIntegrations),
		OutboundIntegrationsMuted: types.BoolValue(channels.OutboundIntegrationsMuted),
		NotificationChannels:      MapNullableList(ctx, channels.NotificationChannels),
		NotificationChannelsMuted: types.BoolValue(channels.NotificationChannelsMuted),
	}
}

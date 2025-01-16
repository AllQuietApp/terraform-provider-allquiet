// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	Id          types.String       `tfsdk:"id"`
	DisplayName types.String       `tfsdk:"display_name"`
	TeamId      types.String       `tfsdk:"team_id"`
	Rules       []RoutingRuleModel `tfsdk:"rules"`
}

type RoutingRuleModel struct {
	Conditions *RoutingRuleConditionsModel `tfsdk:"conditions"`
	Actions    *RoutingRuleActionsModel    `tfsdk:"actions"`
	Channels   *RoutingRuleChannelsModel   `tfsdk:"channels"`
}

type RoutingRuleConditionsModel struct {
	Statuses     types.List                            `tfsdk:"statuses"`
	Severities   types.List                            `tfsdk:"severities"`
	Integrations types.List                            `tfsdk:"integrations"`
	Intents      types.List                            `tfsdk:"intents"`
	Attributes   []RoutingRuleConditionsAttributeModel `tfsdk:"attributes"`
}

type RoutingRuleConditionsAttributeModel struct {
	Name     types.String `tfsdk:"name"`
	Operator types.String `tfsdk:"operator"`
	Value    types.String `tfsdk:"value"`
}

type RoutingRuleActionsModel struct {
	AssignToTeams         types.List   `tfsdk:"assign_to_teams"`
	Discard               types.Bool   `tfsdk:"discard"`
	ChangeSeverity        types.String `tfsdk:"change_severity"`
	AddInteraction        types.String `tfsdk:"add_interaction"`
	RuleFlowControl       types.String `tfsdk:"rule_flow_control"`
	DelayActionsInMinutes types.Int64  `tfsdk:"delay_actions_in_minutes"`
}

type RoutingRuleChannelsModel struct {
	OutboundIntegrations      types.List `tfsdk:"outbound_integrations"`
	OutboundIntegrationsMuted types.Bool `tfsdk:"outbound_integrations_muted"`
	NotificationChannels      types.List `tfsdk:"notification_channels"`
	NotificationChannelsMuted types.Bool `tfsdk:"notification_channels_muted"`
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
									MarkdownDescription: "Statuses. Possible values are: Open, Resolved",
									ElementType:         types.StringType,
									Validators: []validator.List{
										listvalidator.ValueStringsAre(StatusValidator("Not a valid status")),
									},
								},
								"severities": schema.ListAttribute{
									Optional:            true,
									MarkdownDescription: "Severeties. Possible values are: Critical, Warning, Minor",
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
												MarkdownDescription: "The operator",
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
							},
						},
						"actions": schema.SingleNestedAttribute{
							MarkdownDescription: "Settings for the schedule",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"assign_to_teams": schema.ListAttribute{
									Optional:            true,
									MarkdownDescription: "Will assign the incident to the specified teams",
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
									MarkdownDescription: "Will change the severity of the incident",
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
									MarkdownDescription: "If 'Skip' will not evaluate further rules",
									Validators:          []validator.String{RuleFlowValidator("Not a valid rule flow value")},
								},
								"delay_actions_in_minutes": schema.Int64Attribute{
									Optional:            true,
									MarkdownDescription: "Delay actions in minutes",
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
		Statuses:     MapNullableList(ctx, conditions.Statuses),
		Severities:   MapNullableList(ctx, conditions.Severities),
		Integrations: MapNullableList(ctx, conditions.Integrations),
		Intents:      MapNullableList(ctx, conditions.Intents),
		Attributes:   mapRoutingRuleConditionsAttributeResponseToModel(conditions.Attributes),
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
		AssignToTeams:   MapNullableList(ctx, actions.AssignToTeams),
		Discard:         types.BoolValue(actions.Discard),
		ChangeSeverity:  types.StringPointerValue(actions.ChangeSeverity),
		AddInteraction:  types.StringPointerValue(actions.AddInteraction),
		RuleFlowControl: types.StringPointerValue(actions.RuleFlowControl),
	}
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

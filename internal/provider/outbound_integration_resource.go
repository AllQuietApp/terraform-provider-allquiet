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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &OutboundIntegration{}
var _ resource.ResourceWithImportState = &OutboundIntegration{}

func NewOutboundIntegration() resource.Resource {
	return &OutboundIntegration{}
}

type OutboundIntegration struct {
	client *AllQuietAPIClient
}

type OutboundIntegrationModel struct {
	Id                          types.String            `tfsdk:"id"`
	DisplayName                 types.String            `tfsdk:"display_name"`
	TeamId                      types.String            `tfsdk:"team_id"`
	Type                        types.String            `tfsdk:"type"`
	TriggersOnlyOnForwarded     types.Bool              `tfsdk:"triggers_only_on_forwarded"`
	SkipUpdatingAfterForwarding types.Bool              `tfsdk:"skip_updating_after_forwarding"`
	TeamConnectionSettings      *TeamConnectionSettings `tfsdk:"team_connection_settings"`
	SlackSettings               *SlackSettings          `tfsdk:"slack_settings"`
	MattermostSettings          *MattermostSettings     `tfsdk:"mattermost_settings"`
}

func (r *OutboundIntegration) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_outbound_integration"
}

func (r *OutboundIntegration) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The `outbound_integration` resource represents an outbound integration in All Quiet. Outbound integrations are used to send alerts to external systems like Slack or Discord.",

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
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the integration. See all types here: https://allquiet.app/api/public/v1/outbound-integration/types",
				Required:            true,
			},
			"triggers_only_on_forwarded": schema.BoolAttribute{
				MarkdownDescription: "If true, the integration will only trigger once explicitly forwarded.",
				Optional:            true,
				Computed:            true,
			},
			"skip_updating_after_forwarding": schema.BoolAttribute{
				MarkdownDescription: "If true, the integration will not trigger on updates, once it has been forwarded.",
				Optional:            true,
				Computed:            true,
			},
			"team_connection_settings": schema.SingleNestedAttribute{
				MarkdownDescription: "The team connection settings for the integration",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"team_connection_mode": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The team connection mode for the integration. Possible values are: " + strings.Join(ValidTeamConnectionModes, ", "),
						Validators:          []validator.String{stringvalidator.OneOf(ValidTeamConnectionModes...)},
					},
					"team_ids": schema.ListAttribute{
						MarkdownDescription: "The team ids for the integration. If not provided, team_connection_mode must be set to 'OrganizationTeams'.",
						Optional:            true,
						ElementType:         types.StringType,
						Validators: []validator.List{
							listvalidator.ValueStringsAre(GuidValidator("Not a valid GUID")),
						},
					},
				},
			},
			"slack_settings": schema.SingleNestedAttribute{
				MarkdownDescription: "Slack-specific settings for the integration. Only applicable when type is 'Slack'.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"selected_channel_ids": schema.ListAttribute{
						MarkdownDescription: "List of Slack channel IDs to send notifications to. Either this or severity_based_channel_settings must be provided, but not both.",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"severity_based_channel_settings": schema.SingleNestedAttribute{
						MarkdownDescription: "Severity-based channel settings. Either this or selected_channel_ids must be provided, but not both.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"selected_channel_ids_minor": schema.ListAttribute{
								MarkdownDescription: "List of Slack channel IDs for minor severity notifications.",
								Optional:            true,
								ElementType:         types.StringType,
							},
							"selected_channel_ids_warning": schema.ListAttribute{
								MarkdownDescription: "List of Slack channel IDs for warning severity notifications.",
								Optional:            true,
								ElementType:         types.StringType,
							},
							"selected_channel_ids_critical": schema.ListAttribute{
								MarkdownDescription: "List of Slack channel IDs for critical severity notifications.",
								Optional:            true,
								ElementType:         types.StringType,
							},
						},
					},
					"on_call_reminder_schedule_settings": schema.SingleNestedAttribute{
						MarkdownDescription: "Schedule settings for on-call reminders.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"run_time": schema.StringAttribute{
								MarkdownDescription: "The time to run the reminder in HH:mm format (e.g., '09:00').",
								Optional:            true,
							},
							"days_of_week": schema.ListAttribute{
								MarkdownDescription: "List of days of the week for the reminder schedule (e.g., ['Monday', 'Wednesday']).",
								Optional:            true,
								ElementType:         types.StringType,
							},
						},
					},
					"on_call_reminder_channel_ids": schema.ListAttribute{
						MarkdownDescription: "List of Slack channel IDs for on-call reminders.",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"tag_on_call_members": schema.BoolAttribute{
						MarkdownDescription: "If true, tag on-call members in Slack notifications.",
						Optional:            true,
					},
					"is_slack_message_payload_read_only": schema.BoolAttribute{
						MarkdownDescription: "If true, the Slack message payload will be read-only.",
						Optional:            true,
					},
				},
			},
			"mattermost_settings": schema.SingleNestedAttribute{
				MarkdownDescription: "Mattermost-specific settings for the integration. Only applicable when type is 'Mattermost'.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"send_incidents_to_mattermost": schema.BoolAttribute{
						MarkdownDescription: "If true, incidents will be sent to Mattermost.",
						Required:            true,
					},
					"create_incidents_from_mattermost": schema.BoolAttribute{
						MarkdownDescription: "If true, incidents can be created from Mattermost.",
						Required:            true,
					},
					"base_url": schema.StringAttribute{
						MarkdownDescription: "The Mattermost server URL (e.g. https://your-mattermost-server.com).",
						Optional:            true,
					},
					"bot_token": schema.StringAttribute{
						MarkdownDescription: "The Mattermost bot token.",
						Optional:            true,
						Sensitive:           true,
					},
					"slash_command_token": schema.StringAttribute{
						MarkdownDescription: "The Mattermost slash command token.",
						Optional:            true,
						Sensitive:           true,
					},
					"selected_channel_ids": schema.ListAttribute{
						MarkdownDescription: "List of Mattermost channel IDs to send notifications to. Either this or severity_based_channel_settings must be provided, but not both.",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"severity_based_channel_settings": schema.SingleNestedAttribute{
						MarkdownDescription: "Severity-based channel settings. Either this or selected_channel_ids must be provided, but not both.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"selected_channel_ids_minor": schema.ListAttribute{
								MarkdownDescription: "List of Mattermost channel IDs for minor severity notifications.",
								Optional:            true,
								ElementType:         types.StringType,
							},
							"selected_channel_ids_warning": schema.ListAttribute{
								MarkdownDescription: "List of Mattermost channel IDs for warning severity notifications.",
								Optional:            true,
								ElementType:         types.StringType,
							},
							"selected_channel_ids_critical": schema.ListAttribute{
								MarkdownDescription: "List of Mattermost channel IDs for critical severity notifications.",
								Optional:            true,
								ElementType:         types.StringType,
							},
						},
					},
					"selected_team_id": schema.StringAttribute{
						MarkdownDescription: "The Mattermost team ID.",
						Optional:            true,
					},
					"is_message_read_only": schema.BoolAttribute{
						MarkdownDescription: "If true, the Mattermost message will be read-only.",
						Optional:            true,
					},
				},
			},
		},
	}
}

func (r *OutboundIntegration) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OutboundIntegration) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OutboundIntegrationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integrationResponse, err := r.client.CreateOutboundIntegrationResource(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create integration resource, got error: %s", err))
		return
	}

	mapOutboundIntegrationResponseToModel(ctx, integrationResponse, &data)

	tflog.Trace(ctx, "created integration resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OutboundIntegration) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OutboundIntegrationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	integrationResponse, err := r.client.GetOutboundIntegrationResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get integration resource, got error: %s", err))
		return
	}

	if integrationResponse == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to get integration resource, got nil response")
		return
	}

	mapOutboundIntegrationResponseToModel(ctx, integrationResponse, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OutboundIntegration) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data OutboundIntegrationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	integrationResponse, err := r.client.UpdateOutboundIntegrationResource(ctx, data.Id.ValueString(), &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update integration resource, got error: %s", err))
		return
	}

	mapOutboundIntegrationResponseToModel(ctx, integrationResponse, &data)

	tflog.Trace(ctx, "updated integration resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OutboundIntegration) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OutboundIntegrationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOutboundIntegrationResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update integration resource, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted integration resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OutboundIntegration) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapOutboundIntegrationResponseToModel(ctx context.Context, response *outboundIntegrationResponse, data *OutboundIntegrationModel) {

	data.Id = types.StringValue(response.Id)
	data.DisplayName = types.StringValue(response.DisplayName)
	data.TeamId = types.StringValue(response.TeamId)
	data.Type = types.StringValue(response.Type)
	data.TriggersOnlyOnForwarded = types.BoolPointerValue(response.TriggersOnlyOnForwarded)
	data.SkipUpdatingAfterForwarding = types.BoolPointerValue(response.SkipUpdatingAfterForwarding)

	if response.TeamConnectionSettings != nil {
		data.TeamConnectionSettings = &TeamConnectionSettings{
			TeamConnectionMode: types.StringValue(response.TeamConnectionSettings.TeamConnectionMode),
			TeamIds:            MapNullableList(ctx, response.TeamConnectionSettings.TeamIds),
		}
	}

	data.SlackSettings = MapSlackSettingsResponseToModel(ctx, response.SlackSettings)
	data.MattermostSettings = MapMattermostSettingsResponseToModel(ctx, response.MattermostSettings)
}

type SlackSettings struct {
	SelectedChannelIds             types.List                    `tfsdk:"selected_channel_ids"`
	SeverityBasedChannelSettings   *SeverityBasedChannelSettings `tfsdk:"severity_based_channel_settings"`
	OnCallReminderScheduleSettings *ReminderScheduleSettings     `tfsdk:"on_call_reminder_schedule_settings"`
	OnCallReminderChannelIds       types.List                    `tfsdk:"on_call_reminder_channel_ids"`
	TagOnCallMembers               types.Bool                    `tfsdk:"tag_on_call_members"`
	IsSlackMessagePayloadReadOnly  types.Bool                    `tfsdk:"is_slack_message_payload_read_only"`
}

type SeverityBasedChannelSettings struct {
	SelectedChannelIdsMinor    types.List `tfsdk:"selected_channel_ids_minor"`
	SelectedChannelIdsWarning  types.List `tfsdk:"selected_channel_ids_warning"`
	SelectedChannelIdsCritical types.List `tfsdk:"selected_channel_ids_critical"`
}

type ReminderScheduleSettings struct {
	RunTime    types.String `tfsdk:"run_time"`
	DaysOfWeek types.List   `tfsdk:"days_of_week"`
}

type MattermostSettings struct {
	SendIncidentsToMattermost     types.Bool                              `tfsdk:"send_incidents_to_mattermost"`
	CreateIncidentsFromMattermost types.Bool                              `tfsdk:"create_incidents_from_mattermost"`
	BaseUrl                       types.String                            `tfsdk:"base_url"`
	BotToken                      types.String                            `tfsdk:"bot_token"`
	SlashCommandToken             types.String                            `tfsdk:"slash_command_token"`
	SelectedChannelIds            types.List                              `tfsdk:"selected_channel_ids"`
	SeverityBasedChannelSettings  *MattermostSeverityBasedChannelSettings `tfsdk:"severity_based_channel_settings"`
	SelectedTeamId                types.String                            `tfsdk:"selected_team_id"`
	IsMessageReadOnly             types.Bool                              `tfsdk:"is_message_read_only"`
}

type MattermostSeverityBasedChannelSettings struct {
	SelectedChannelIdsMinor    types.List `tfsdk:"selected_channel_ids_minor"`
	SelectedChannelIdsWarning  types.List `tfsdk:"selected_channel_ids_warning"`
	SelectedChannelIdsCritical types.List `tfsdk:"selected_channel_ids_critical"`
}

type slackSettings struct {
	SelectedChannelIds             *[]string                     `json:"selectedChannelIds"`
	SeverityBasedChannelSettings   *severityBasedChannelSettings `json:"severityBasedChannelSettings"`
	OnCallReminderScheduleSettings *reminderScheduleSettings     `json:"onCallReminderScheduleSettings"`
	OnCallReminderChannelIds       *[]string                     `json:"onCallReminderChannelIds"`
	TagOnCallMembers               *bool                         `json:"tagOnCallMembers"`
	IsSlackMessagePayloadReadOnly  *bool                         `json:"isSlackMessagePayloadReadOnly"`
}

type severityBasedChannelSettings struct {
	SelectedChannelIdsMinor    *[]string `json:"selectedChannelIdsMinor"`
	SelectedChannelIdsWarning  *[]string `json:"selectedChannelIdsWarning"`
	SelectedChannelIdsCritical *[]string `json:"selectedChannelIdsCritical"`
}

type reminderScheduleSettings struct {
	RunTime    *string   `json:"runTime"`
	DaysOfWeek *[]string `json:"daysOfWeek"`
}

type mattermostSettings struct {
	SendIncidentsToMattermost     *bool                                   `json:"sendIncidentsToMattermost"`
	CreateIncidentsFromMattermost *bool                                   `json:"createIncidentsFromMattermost"`
	BaseUrl                       *string                                 `json:"baseUrl"`
	BotToken                      *string                                 `json:"botToken"`
	SlashCommandToken             *string                                 `json:"slashCommandToken"`
	SelectedChannelIds            *[]string                               `json:"selectedChannelIds"`
	SeverityBasedChannelSettings  *mattermostSeverityBasedChannelSettings `json:"severityBasedChannelSettings"`
	SelectedTeamId                *string                                 `json:"selectedTeamId"`
	IsMessageReadOnly             *bool                                   `json:"isMessageReadOnly"`
}

type mattermostSeverityBasedChannelSettings struct {
	SelectedChannelIdsMinor    *[]string `json:"selectedChannelIdsMinor"`
	SelectedChannelIdsWarning  *[]string `json:"selectedChannelIdsWarning"`
	SelectedChannelIdsCritical *[]string `json:"selectedChannelIdsCritical"`
}

func MapSlackSettingsToRequest(settings *SlackSettings) *slackSettings {
	if settings == nil {
		return nil
	}

	result := &slackSettings{
		SelectedChannelIds:            ListToStringArray(settings.SelectedChannelIds),
		OnCallReminderChannelIds:      ListToStringArray(settings.OnCallReminderChannelIds),
		TagOnCallMembers:              settings.TagOnCallMembers.ValueBoolPointer(),
		IsSlackMessagePayloadReadOnly: settings.IsSlackMessagePayloadReadOnly.ValueBoolPointer(),
	}

	if settings.SeverityBasedChannelSettings != nil {
		result.SeverityBasedChannelSettings = &severityBasedChannelSettings{
			SelectedChannelIdsMinor:    ListToStringArray(settings.SeverityBasedChannelSettings.SelectedChannelIdsMinor),
			SelectedChannelIdsWarning:  ListToStringArray(settings.SeverityBasedChannelSettings.SelectedChannelIdsWarning),
			SelectedChannelIdsCritical: ListToStringArray(settings.SeverityBasedChannelSettings.SelectedChannelIdsCritical),
		}
	}

	if settings.OnCallReminderScheduleSettings != nil {
		var runTime *string
		if !settings.OnCallReminderScheduleSettings.RunTime.IsNull() && !settings.OnCallReminderScheduleSettings.RunTime.IsUnknown() {
			runTimeStr := settings.OnCallReminderScheduleSettings.RunTime.ValueString()
			runTime = &runTimeStr
		}
		result.OnCallReminderScheduleSettings = &reminderScheduleSettings{
			RunTime:    runTime,
			DaysOfWeek: ListToStringArray(settings.OnCallReminderScheduleSettings.DaysOfWeek),
		}
	}

	return result
}

func MapMattermostSettingsToRequest(settings *MattermostSettings) *mattermostSettings {
	if settings == nil {
		return nil
	}

	result := &mattermostSettings{
		SendIncidentsToMattermost:     settings.SendIncidentsToMattermost.ValueBoolPointer(),
		CreateIncidentsFromMattermost: settings.CreateIncidentsFromMattermost.ValueBoolPointer(),
		IsMessageReadOnly:             settings.IsMessageReadOnly.ValueBoolPointer(),
	}

	if !settings.BaseUrl.IsNull() && !settings.BaseUrl.IsUnknown() {
		s := settings.BaseUrl.ValueString()
		result.BaseUrl = &s
	}
	if !settings.BotToken.IsNull() && !settings.BotToken.IsUnknown() {
		s := settings.BotToken.ValueString()
		result.BotToken = &s
	}
	if !settings.SlashCommandToken.IsNull() && !settings.SlashCommandToken.IsUnknown() {
		s := settings.SlashCommandToken.ValueString()
		result.SlashCommandToken = &s
	}
	if !settings.SelectedTeamId.IsNull() && !settings.SelectedTeamId.IsUnknown() {
		s := settings.SelectedTeamId.ValueString()
		result.SelectedTeamId = &s
	}

	result.SelectedChannelIds = ListToStringArray(settings.SelectedChannelIds)

	if settings.SeverityBasedChannelSettings != nil {
		result.SeverityBasedChannelSettings = &mattermostSeverityBasedChannelSettings{
			SelectedChannelIdsMinor:    ListToStringArray(settings.SeverityBasedChannelSettings.SelectedChannelIdsMinor),
			SelectedChannelIdsWarning:  ListToStringArray(settings.SeverityBasedChannelSettings.SelectedChannelIdsWarning),
			SelectedChannelIdsCritical: ListToStringArray(settings.SeverityBasedChannelSettings.SelectedChannelIdsCritical),
		}
	}

	return result
}

func MapMattermostSettingsResponseToModel(ctx context.Context, settings *mattermostSettings) *MattermostSettings {
	if settings == nil {
		return nil
	}

	// Only treat booleans as "has" when true, so API response { sendIncidentsToMattermost: false, ... } with nothing else set maps to nil (plan was null)
	hasSendIncidents := settings.SendIncidentsToMattermost != nil && *settings.SendIncidentsToMattermost
	hasCreateIncidents := settings.CreateIncidentsFromMattermost != nil && *settings.CreateIncidentsFromMattermost
	hasBaseUrl := settings.BaseUrl != nil && *settings.BaseUrl != ""
	hasBotToken := settings.BotToken != nil && *settings.BotToken != ""
	hasSlashCommandToken := settings.SlashCommandToken != nil && *settings.SlashCommandToken != ""
	hasSelectedChannelIds := settings.SelectedChannelIds != nil && len(*settings.SelectedChannelIds) > 0
	hasSeverityBased := settings.SeverityBasedChannelSettings != nil &&
		((settings.SeverityBasedChannelSettings.SelectedChannelIdsMinor != nil && len(*settings.SeverityBasedChannelSettings.SelectedChannelIdsMinor) > 0) ||
			(settings.SeverityBasedChannelSettings.SelectedChannelIdsWarning != nil && len(*settings.SeverityBasedChannelSettings.SelectedChannelIdsWarning) > 0) ||
			(settings.SeverityBasedChannelSettings.SelectedChannelIdsCritical != nil && len(*settings.SeverityBasedChannelSettings.SelectedChannelIdsCritical) > 0))
	hasSelectedTeamId := settings.SelectedTeamId != nil && *settings.SelectedTeamId != ""
	hasIsMessageReadOnly := settings.IsMessageReadOnly != nil

	if !hasSendIncidents && !hasCreateIncidents && !hasBaseUrl && !hasBotToken && !hasSlashCommandToken && !hasSelectedChannelIds && !hasSeverityBased && !hasSelectedTeamId && !hasIsMessageReadOnly {
		return nil
	}

	result := &MattermostSettings{
		SendIncidentsToMattermost:     types.BoolPointerValue(settings.SendIncidentsToMattermost),
		CreateIncidentsFromMattermost: types.BoolPointerValue(settings.CreateIncidentsFromMattermost),
		BaseUrl:                       types.StringPointerValue(settings.BaseUrl),
		BotToken:                      types.StringPointerValue(settings.BotToken),
		SlashCommandToken:             types.StringPointerValue(settings.SlashCommandToken),
		SelectedChannelIds:            MapNullableList(ctx, settings.SelectedChannelIds),
		SelectedTeamId:                types.StringPointerValue(settings.SelectedTeamId),
		IsMessageReadOnly:             types.BoolPointerValue(settings.IsMessageReadOnly),
	}

	if settings.SeverityBasedChannelSettings != nil {
		result.SeverityBasedChannelSettings = &MattermostSeverityBasedChannelSettings{
			SelectedChannelIdsMinor:    MapNullableList(ctx, settings.SeverityBasedChannelSettings.SelectedChannelIdsMinor),
			SelectedChannelIdsWarning:  MapNullableList(ctx, settings.SeverityBasedChannelSettings.SelectedChannelIdsWarning),
			SelectedChannelIdsCritical: MapNullableList(ctx, settings.SeverityBasedChannelSettings.SelectedChannelIdsCritical),
		}
	}

	return result
}

func MapSlackSettingsResponseToModel(ctx context.Context, settings *slackSettings) *SlackSettings {
	if settings == nil {
		return nil
	}

	hasSelectedChannelIds := settings.SelectedChannelIds != nil && len(*settings.SelectedChannelIds) > 0
	hasSeverityBased := settings.SeverityBasedChannelSettings != nil &&
		((settings.SeverityBasedChannelSettings.SelectedChannelIdsMinor != nil && len(*settings.SeverityBasedChannelSettings.SelectedChannelIdsMinor) > 0) ||
			(settings.SeverityBasedChannelSettings.SelectedChannelIdsWarning != nil && len(*settings.SeverityBasedChannelSettings.SelectedChannelIdsWarning) > 0) ||
			(settings.SeverityBasedChannelSettings.SelectedChannelIdsCritical != nil && len(*settings.SeverityBasedChannelSettings.SelectedChannelIdsCritical) > 0))
	hasOnCallReminderChannelIds := settings.OnCallReminderChannelIds != nil && len(*settings.OnCallReminderChannelIds) > 0
	hasOnCallReminderSchedule := settings.OnCallReminderScheduleSettings != nil &&
		(settings.OnCallReminderScheduleSettings.RunTime != nil || (settings.OnCallReminderScheduleSettings.DaysOfWeek != nil && len(*settings.OnCallReminderScheduleSettings.DaysOfWeek) > 0))
	hasTagOnCallMembers := settings.TagOnCallMembers != nil
	hasIsSlackMessagePayloadReadOnly := settings.IsSlackMessagePayloadReadOnly != nil

	if !hasSelectedChannelIds && !hasSeverityBased && !hasOnCallReminderChannelIds && !hasOnCallReminderSchedule && !hasTagOnCallMembers && !hasIsSlackMessagePayloadReadOnly {
		return nil
	}

	result := &SlackSettings{
		SelectedChannelIds:            MapNullableList(ctx, settings.SelectedChannelIds),
		OnCallReminderChannelIds:      MapNullableList(ctx, settings.OnCallReminderChannelIds),
		TagOnCallMembers:              types.BoolPointerValue(settings.TagOnCallMembers),
		IsSlackMessagePayloadReadOnly: types.BoolPointerValue(settings.IsSlackMessagePayloadReadOnly),
	}

	if settings.SeverityBasedChannelSettings != nil {
		result.SeverityBasedChannelSettings = &SeverityBasedChannelSettings{
			SelectedChannelIdsMinor:    MapNullableList(ctx, settings.SeverityBasedChannelSettings.SelectedChannelIdsMinor),
			SelectedChannelIdsWarning:  MapNullableList(ctx, settings.SeverityBasedChannelSettings.SelectedChannelIdsWarning),
			SelectedChannelIdsCritical: MapNullableList(ctx, settings.SeverityBasedChannelSettings.SelectedChannelIdsCritical),
		}
	}

	if settings.OnCallReminderScheduleSettings != nil {
		runTime := types.StringPointerValue(settings.OnCallReminderScheduleSettings.RunTime)
		result.OnCallReminderScheduleSettings = &ReminderScheduleSettings{
			RunTime:    runTime,
			DaysOfWeek: MapNullableList(ctx, settings.OnCallReminderScheduleSettings.DaysOfWeek),
		}
	}

	return result
}

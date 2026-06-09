package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	callRoutingMediaPurposeWelcome   = "welcome-audio"
	callRoutingMediaPurposeVoicemail = "voicemail-message"
	callRoutingMediaPurposeClosing   = "closing-message"
)

type TwilioSettingsModel struct {
	AccessSettings *TwilioAccessSettingsModel `tfsdk:"access_settings"`
	CallFlowConfig *CallFlowConfigModel       `tfsdk:"call_flow_config"`
}

type TwilioAccessSettingsModel struct {
	ApiBaseDomain  types.String `tfsdk:"api_base_domain"`
	AccountSid     types.String `tfsdk:"account_sid"`
	AuthToken      types.String `tfsdk:"auth_token"`
	PhoneNumber    types.String `tfsdk:"phone_number"`
	PhoneNumberSid types.String `tfsdk:"phone_number_sid"`
}

type CallFlowConfigModel struct {
	WelcomeMessage            *PromptModel      `tfsdk:"welcome_message"`
	CallDisplayMode           types.String      `tfsdk:"call_display_mode"`
	AttachVoicemailTranscript types.Bool        `tfsdk:"attach_voicemail_transcript"`
	CallerAllowlist           types.List        `tfsdk:"caller_allowlist"`
	CallerBlocklist           types.List        `tfsdk:"caller_blocklist"`
	Route                     *RouteConfigModel `tfsdk:"route"`
	Menu                      *MenuConfigModel  `tfsdk:"menu"`
}

type RouteConfigModel struct {
	TeamId                            types.String `tfsdk:"team_id"`
	IncidentSeverity                  types.String `tfsdk:"incident_severity"`
	DirectVoicemail                   types.Bool   `tfsdk:"direct_voicemail"`
	SendToVoicemailIfNoResponse       types.Bool   `tfsdk:"send_to_voicemail_if_no_response"`
	VoicemailMessage                  *PromptModel `tfsdk:"voicemail_message"`
	ClosingMessage                    *PromptModel `tfsdk:"closing_message"`
	AutoResolveOnAnswer               types.Bool   `tfsdk:"auto_resolve_on_answer"`
	SkipAlreadyDialedNumbers          types.Bool   `tfsdk:"skip_already_dialed_numbers"`
	MaxUsersToTry                     types.Int64  `tfsdk:"max_users_to_try"`
	DialTimeoutPerUserInSeconds       types.Int64  `tfsdk:"dial_timeout_per_user_in_seconds"`
	AcceptCallPrompt                  *PromptModel `tfsdk:"accept_call_prompt"`
	TargetConfirmGatherTimeoutSeconds types.Int64  `tfsdk:"target_confirm_gather_timeout_seconds"`
	TargetConfirmMaxAttempts          types.Int64  `tfsdk:"target_confirm_max_attempts"`
}

type MenuConfigModel struct {
	Options                  []MenuOptionModel `tfsdk:"options"`
	MenuMaxAttempts          types.Int64       `tfsdk:"menu_max_attempts"`
	MenuGatherTimeoutSeconds types.Int64       `tfsdk:"menu_gather_timeout_seconds"`
}

type MenuOptionModel struct {
	Key   types.String     `tfsdk:"key"`
	Route RouteConfigModel `tfsdk:"route"`
}

type PromptModel struct {
	Type      types.String    `tfsdk:"type"`
	Text      types.String    `tfsdk:"text"`
	Locale    types.String    `tfsdk:"locale"`
	Voice     types.String    `tfsdk:"voice"`
	AudioFile *AudioFileModel `tfsdk:"audio_file"`
}

type AudioFileModel struct {
	Source        types.String `tfsdk:"source"`
	ObjectKey     types.String `tfsdk:"object_key"`
	FileName      types.String `tfsdk:"file_name"`
	ContentType   types.String `tfsdk:"content_type"`
	FileSizeBytes types.Int64  `tfsdk:"file_size_bytes"`
	SourceHash    types.String `tfsdk:"source_hash"`
}

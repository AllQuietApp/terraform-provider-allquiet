package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func twilioIntegrationSettingsSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"access_settings": schema.SingleNestedAttribute{
			MarkdownDescription: "Twilio account credentials and inbound phone number",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"api_base_domain": schema.StringAttribute{
					MarkdownDescription: "Twilio API base domain",
					Optional:            true,
					Computed:            true,
				},
				"account_sid": schema.StringAttribute{
					MarkdownDescription: "Twilio account SID",
					Required:            true,
					Sensitive:           true,
				},
				"auth_token": schema.StringAttribute{
					MarkdownDescription: "Twilio auth token",
					Required:            true,
					Sensitive:           true,
				},
				"phone_number": schema.StringAttribute{
					MarkdownDescription: "Inbound phone number in E.164 format",
					Required:            true,
				},
				"phone_number_sid": schema.StringAttribute{
					MarkdownDescription: "Twilio phone number SID",
					Optional:            true,
					Computed:            true,
					Sensitive:           true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
		},
		"call_flow_config": schema.SingleNestedAttribute{
			MarkdownDescription: "Inbound call routing flow configuration",
			Optional:            true,
			Attributes:          callFlowConfigSchemaAttributes(),
		},
	}
}

func callFlowConfigSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"welcome_message": promptSchemaAttribute("Welcome message played at the start of the call", true),
		"call_display_mode": schema.StringAttribute{
			MarkdownDescription: "Caller ID display mode. Possible values: ShowCallerNumber, ShowInboundNumber",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("ShowInboundNumber"),
			Validators:          []validator.String{CallDisplayModeValidator("Not a valid call display mode")},
		},
		"attach_voicemail_transcript": schema.BoolAttribute{
			MarkdownDescription: "Attach voicemail transcript to incidents",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(true),
		},
		"caller_allowlist": schema.ListAttribute{
			MarkdownDescription: "Phone numbers or country prefixes allowed to call (E.164, e.g. +14155552671 or +1)",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"caller_blocklist": schema.ListAttribute{
			MarkdownDescription: "Phone numbers or country prefixes blocked from calling (E.164)",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"route": schema.SingleNestedAttribute{
			MarkdownDescription: "Default route when no menu is configured or after menu selection",
			Optional:            true,
			Attributes:          routeConfigSchemaAttributes(true),
		},
		"menu": schema.SingleNestedAttribute{
			MarkdownDescription: "IVR menu configuration",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"options": schema.ListNestedAttribute{
					MarkdownDescription: "Menu options (max 9)",
					Optional:            true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"key": schema.StringAttribute{
								MarkdownDescription: "DTMF key (0-9)",
								Required:            true,
							},
							"route": schema.SingleNestedAttribute{
								MarkdownDescription: "Route for this menu option",
								Required:            true,
								Attributes:          routeConfigSchemaAttributes(true),
							},
						},
					},
				},
				"menu_max_attempts": schema.Int64Attribute{
					MarkdownDescription: "Maximum menu input attempts (1-5)",
					Optional:            true,
					Computed:            true,
					Default:             int64default.StaticInt64(2),
					Validators:          []validator.Int64{int64validator.Between(1, 5)},
				},
				"menu_gather_timeout_seconds": schema.Int64Attribute{
					MarkdownDescription: "Seconds to wait for menu input (3-60)",
					Optional:            true,
					Computed:            true,
					Default:             int64default.StaticInt64(10),
					Validators:          []validator.Int64{int64validator.Between(3, 60)},
				},
			},
		},
	}
}

func routeConfigSchemaAttributes(includeAudioPrompts bool) map[string]schema.Attribute {
	attrs := map[string]schema.Attribute{
		"team_id": schema.StringAttribute{
			MarkdownDescription: "Team to route calls to",
			Optional:            true,
			Validators:          []validator.String{GuidValidator("Not a valid team id")},
		},
		"incident_severity": schema.StringAttribute{
			MarkdownDescription: "Incident severity for routed calls. Possible values: Critical, Warning, Minor",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("Critical"),
			Validators:          []validator.String{SeverityValidator("Not a valid severity")},
		},
		"direct_voicemail": schema.BoolAttribute{
			MarkdownDescription: "Send callers directly to voicemail",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"send_to_voicemail_if_no_response": schema.BoolAttribute{
			MarkdownDescription: "Send to voicemail if on-call users do not answer",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(true),
		},
		"auto_resolve_on_answer": schema.BoolAttribute{
			MarkdownDescription: "Auto-resolve incident when call is answered",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(true),
		},
		"skip_already_dialed_numbers": schema.BoolAttribute{
			MarkdownDescription: "Skip phone numbers already dialed in this session",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"max_users_to_try": schema.Int64Attribute{
			MarkdownDescription: "Maximum on-call users to try (1-10)",
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(3),
			Validators:          []validator.Int64{int64validator.Between(1, 10)},
		},
		"dial_timeout_per_user_in_seconds": schema.Int64Attribute{
			MarkdownDescription: "Dial timeout per user in seconds (5-120)",
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(20),
			Validators:          []validator.Int64{int64validator.Between(5, 120)},
		},
		"target_confirm_gather_timeout_seconds": schema.Int64Attribute{
			MarkdownDescription: "Seconds to wait for accept-call confirmation (3-60)",
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(10),
			Validators:          []validator.Int64{int64validator.Between(3, 60)},
		},
		"target_confirm_max_attempts": schema.Int64Attribute{
			MarkdownDescription: "Maximum accept-call confirmation attempts (1-5)",
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(1),
			Validators:          []validator.Int64{int64validator.Between(1, 5)},
		},
	}

	if includeAudioPrompts {
		attrs["voicemail_message"] = promptSchemaAttribute("Voicemail greeting message", true)
		attrs["closing_message"] = promptSchemaAttribute("Closing message after voicemail", true)
	}
	attrs["accept_call_prompt"] = promptSchemaAttribute("Prompt played when dialing on-call users (text only)", false)

	return attrs
}

func promptSchemaAttribute(description string, allowAudio bool) schema.SingleNestedAttribute {
	attrs := map[string]schema.Attribute{
		"type": schema.StringAttribute{
			MarkdownDescription: "Prompt type: Text or Audio",
			Required:            true,
			Validators:          []validator.String{PromptTypeValidator("Not a valid prompt type")},
		},
		"text": schema.StringAttribute{
			MarkdownDescription: "Text-to-speech content when type is Text",
			Optional:            true,
		},
		"locale": schema.StringAttribute{
			MarkdownDescription: "Voice locale for text prompts (see call-routing-options API)",
			Optional:            true,
		},
		"voice": schema.StringAttribute{
			MarkdownDescription: "Voice name for text prompts",
			Optional:            true,
		},
	}

	if allowAudio {
		attrs["audio_file"] = schema.SingleNestedAttribute{
			MarkdownDescription: "Uploaded audio file when type is Audio",
			Optional:            true,
			PlanModifiers: []planmodifier.Object{
				audioFilePlanModifier{},
			},
			Attributes: map[string]schema.Attribute{
				"source": schema.StringAttribute{
					MarkdownDescription: "Local path to MP3 or WAV file to upload",
					Required:            true,
				},
				"object_key": schema.StringAttribute{
					MarkdownDescription: "S3 object key for the uploaded audio",
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"file_name": schema.StringAttribute{
					MarkdownDescription: "Original file name",
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"content_type": schema.StringAttribute{
					MarkdownDescription: "Uploaded content type",
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"file_size_bytes": schema.Int64Attribute{
					MarkdownDescription: "Uploaded file size in bytes",
					Computed:            true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.UseStateForUnknown(),
					},
				},
				"source_hash": schema.StringAttribute{
					MarkdownDescription: "SHA-256 hash of the source file content",
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
		}
	}

	return schema.SingleNestedAttribute{
		MarkdownDescription: description,
		Optional:            true,
		Attributes:          attrs,
	}
}

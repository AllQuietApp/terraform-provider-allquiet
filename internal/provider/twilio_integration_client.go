package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type twilioSettingsResponse struct {
	AccessSettings *twilioAccessSettingsResponse `json:"accessSettings"`
	CallFlowConfig *callFlowConfigResponse       `json:"callFlowConfig"`
}

type twilioAccessSettingsResponse struct {
	ApiBaseDomain  string  `json:"apiBaseDomain"`
	AccountSid     *string `json:"accountSid"`
	AuthToken      *string `json:"authToken"`
	PhoneNumber    *string `json:"phoneNumber"`
	PhoneNumberSid *string `json:"phoneNumberSid"`
}

type callFlowConfigResponse struct {
	WelcomeMessage            *promptResponse      `json:"welcomeMessage"`
	CallDisplayMode           string               `json:"callDisplayMode"`
	AttachVoicemailTranscript *bool                `json:"attachVoicemailTranscript"`
	CallerAllowlist           *[]string            `json:"callerAllowlist"`
	CallerBlocklist           *[]string            `json:"callerBlocklist"`
	Route                     *routeConfigResponse `json:"route"`
	Menu                      *menuConfigResponse  `json:"menu"`
}

type routeConfigResponse struct {
	TeamId                            *string         `json:"teamId"`
	IncidentSeverity                  string          `json:"incidentSeverity"`
	DirectVoicemail                   *bool           `json:"directVoicemail"`
	SendToVoicemailIfNoResponse       *bool           `json:"sendToVoicemailIfNoResponse"`
	VoicemailMessage                  *promptResponse `json:"voicemailMessage"`
	ClosingMessage                    *promptResponse `json:"closingMessage"`
	AutoResolveOnAnswer               *bool           `json:"autoResolveOnAnswer"`
	SkipAlreadyDialedNumbers          *bool           `json:"skipAlreadyDialedNumbers"`
	MaxUsersToTry                     *int64          `json:"maxUsersToTry"`
	DialTimeoutPerUserInSeconds       *int64          `json:"dialTimeoutPerUserInSeconds"`
	AcceptCallPrompt                  *promptResponse `json:"acceptCallPrompt"`
	TargetConfirmGatherTimeoutSeconds *int64          `json:"targetConfirmGatherTimeoutSeconds"`
	TargetConfirmMaxAttempts          *int64          `json:"targetConfirmMaxAttempts"`
}

type menuConfigResponse struct {
	Options                  *[]menuOptionResponse `json:"options"`
	MenuMaxAttempts          *int64                `json:"menuMaxAttempts"`
	MenuGatherTimeoutSeconds *int64                `json:"menuGatherTimeoutSeconds"`
}

type menuOptionResponse struct {
	Key   string              `json:"key"`
	Route routeConfigResponse `json:"route"`
}

type promptResponse struct {
	Type      string                        `json:"type"`
	Text      *string                       `json:"text"`
	Locale    *string                       `json:"locale"`
	Voice     *string                       `json:"voice"`
	AudioFile *callRoutingAudioFileResponse `json:"audioFile"`
}

func mapTwilioCreateRequest(plan *TwilioSettingsModel) *twilioSettingsResponse {
	if plan == nil {
		return nil
	}

	return &twilioSettingsResponse{
		AccessSettings: mapTwilioAccessSettingsCreateRequest(plan.AccessSettings),
		CallFlowConfig: mapCallFlowConfigCreateRequest(plan.CallFlowConfig),
	}
}

func mapTwilioAccessSettingsCreateRequest(plan *TwilioAccessSettingsModel) *twilioAccessSettingsResponse {
	if plan == nil {
		return nil
	}

	return &twilioAccessSettingsResponse{
		ApiBaseDomain:  plan.ApiBaseDomain.ValueString(),
		AccountSid:     plan.AccountSid.ValueStringPointer(),
		AuthToken:      plan.AuthToken.ValueStringPointer(),
		PhoneNumber:    plan.PhoneNumber.ValueStringPointer(),
		PhoneNumberSid: plan.PhoneNumberSid.ValueStringPointer(),
	}
}

func mapCallFlowConfigCreateRequest(plan *CallFlowConfigModel) *callFlowConfigResponse {
	if plan == nil {
		return nil
	}

	response := &callFlowConfigResponse{
		WelcomeMessage:            mapPromptCreateRequest(plan.WelcomeMessage),
		CallDisplayMode:           plan.CallDisplayMode.ValueString(),
		AttachVoicemailTranscript: plan.AttachVoicemailTranscript.ValueBoolPointer(),
		CallerAllowlist:           stringListForAPI(plan.CallerAllowlist),
		CallerBlocklist:           stringListForAPI(plan.CallerBlocklist),
		Route:                     mapRouteConfigCreateRequest(plan.Route),
		Menu:                      mapMenuConfigCreateRequest(plan.Menu),
	}
	return response
}

func mapMenuConfigCreateRequest(plan *MenuConfigModel) *menuConfigResponse {
	if plan == nil {
		return nil
	}

	var options []menuOptionResponse
	for _, option := range plan.Options {
		route := mapRouteConfigCreateRequest(&option.Route)
		if route == nil {
			route = &routeConfigResponse{}
		}
		options = append(options, menuOptionResponse{
			Key:   option.Key.ValueString(),
			Route: *route,
		})
	}

	return &menuConfigResponse{
		Options:                  &options,
		MenuMaxAttempts:          plan.MenuMaxAttempts.ValueInt64Pointer(),
		MenuGatherTimeoutSeconds: plan.MenuGatherTimeoutSeconds.ValueInt64Pointer(),
	}
}

func mapRouteConfigCreateRequest(plan *RouteConfigModel) *routeConfigResponse {
	if plan == nil {
		return nil
	}

	return &routeConfigResponse{
		TeamId:                            plan.TeamId.ValueStringPointer(),
		IncidentSeverity:                  plan.IncidentSeverity.ValueString(),
		DirectVoicemail:                   plan.DirectVoicemail.ValueBoolPointer(),
		SendToVoicemailIfNoResponse:       plan.SendToVoicemailIfNoResponse.ValueBoolPointer(),
		VoicemailMessage:                  mapPromptCreateRequest(plan.VoicemailMessage),
		ClosingMessage:                    mapPromptCreateRequest(plan.ClosingMessage),
		AutoResolveOnAnswer:               plan.AutoResolveOnAnswer.ValueBoolPointer(),
		SkipAlreadyDialedNumbers:          plan.SkipAlreadyDialedNumbers.ValueBoolPointer(),
		MaxUsersToTry:                     plan.MaxUsersToTry.ValueInt64Pointer(),
		DialTimeoutPerUserInSeconds:       plan.DialTimeoutPerUserInSeconds.ValueInt64Pointer(),
		AcceptCallPrompt:                  mapPromptCreateRequest(plan.AcceptCallPrompt),
		TargetConfirmGatherTimeoutSeconds: plan.TargetConfirmGatherTimeoutSeconds.ValueInt64Pointer(),
		TargetConfirmMaxAttempts:          plan.TargetConfirmMaxAttempts.ValueInt64Pointer(),
	}
}

func mapPromptCreateRequest(plan *PromptModel) *promptResponse {
	if plan == nil {
		return nil
	}

	return &promptResponse{
		Type:      plan.Type.ValueString(),
		Text:      plan.Text.ValueStringPointer(),
		Locale:    plan.Locale.ValueStringPointer(),
		Voice:     plan.Voice.ValueStringPointer(),
		AudioFile: mapAudioFileCreateRequest(plan.AudioFile),
	}
}

func mapAudioFileCreateRequest(plan *AudioFileModel) *callRoutingAudioFileResponse {
	if plan == nil {
		return nil
	}
	if plan.ObjectKey.IsNull() || plan.ObjectKey.IsUnknown() || plan.ObjectKey.ValueString() == "" {
		return nil
	}

	return &callRoutingAudioFileResponse{
		ObjectKey:     plan.ObjectKey.ValueString(),
		FileName:      plan.FileName.ValueStringPointer(),
		ContentType:   plan.ContentType.ValueStringPointer(),
		FileSizeBytes: plan.FileSizeBytes.ValueInt64Pointer(),
	}
}

func mapTwilioResponseToModel(ctx context.Context, response *twilioSettingsResponse, prior *TwilioSettingsModel) *TwilioSettingsModel {
	if response == nil {
		return nil
	}

	var priorAccess *TwilioAccessSettingsModel
	var priorCallFlow *CallFlowConfigModel
	if prior != nil {
		priorAccess = prior.AccessSettings
		priorCallFlow = prior.CallFlowConfig
	}

	return &TwilioSettingsModel{
		AccessSettings: mapTwilioAccessSettingsResponseToModel(response.AccessSettings, priorAccess),
		CallFlowConfig: mapCallFlowConfigResponseToModel(ctx, response.CallFlowConfig, priorCallFlow),
	}
}

func mapTwilioAccessSettingsResponseToModel(response *twilioAccessSettingsResponse, prior *TwilioAccessSettingsModel) *TwilioAccessSettingsModel {
	if response == nil {
		return nil
	}

	model := &TwilioAccessSettingsModel{
		ApiBaseDomain:  types.StringValue(response.ApiBaseDomain),
		AccountSid:     types.StringPointerValue(response.AccountSid),
		PhoneNumber:    types.StringPointerValue(response.PhoneNumber),
		PhoneNumberSid: types.StringPointerValue(response.PhoneNumberSid),
	}

	if prior != nil {
		if !prior.AccountSid.IsNull() && !prior.AccountSid.IsUnknown() {
			model.AccountSid = prior.AccountSid
		}
		if !prior.AuthToken.IsNull() && !prior.AuthToken.IsUnknown() {
			model.AuthToken = prior.AuthToken
		} else if response.AuthToken != nil {
			model.AuthToken = types.StringValue(*response.AuthToken)
		}
		model.PhoneNumberSid = mergeStringFromPrior(model.PhoneNumberSid, prior.PhoneNumberSid)
	} else if response.AuthToken != nil {
		model.AuthToken = types.StringValue(*response.AuthToken)
	}

	return model
}

func mapCallFlowConfigResponseToModel(ctx context.Context, response *callFlowConfigResponse, prior *CallFlowConfigModel) *CallFlowConfigModel {
	if response == nil {
		return nil
	}

	var priorWelcome *PromptModel
	var priorRoute *RouteConfigModel
	var priorMenu *MenuConfigModel
	if prior != nil {
		priorWelcome = prior.WelcomeMessage
		priorRoute = prior.Route
		priorMenu = prior.Menu
	}

	model := &CallFlowConfigModel{
		WelcomeMessage:            mapPromptResponseToModel(response.WelcomeMessage, priorWelcome),
		CallDisplayMode:           types.StringValue(response.CallDisplayMode),
		AttachVoicemailTranscript: BoolPointerWithDefaultTrue(response.AttachVoicemailTranscript),
		CallerAllowlist:           MapNullableList(ctx, response.CallerAllowlist),
		CallerBlocklist:           MapNullableList(ctx, response.CallerBlocklist),
		Route:                     mapRouteConfigResponseToModel(response.Route, priorRoute),
		Menu:                      mapMenuConfigResponseToModel(response.Menu, priorMenu),
	}

	if prior != nil {
		model.CallDisplayMode = mergeStringFromPrior(model.CallDisplayMode, prior.CallDisplayMode)
		model.CallerAllowlist = mergeStringListFromPrior(model.CallerAllowlist, prior.CallerAllowlist)
		model.CallerBlocklist = mergeStringListFromPrior(model.CallerBlocklist, prior.CallerBlocklist)
		model.AttachVoicemailTranscript = mergeBoolFromPrior(model.AttachVoicemailTranscript, prior.AttachVoicemailTranscript)
	}

	return model
}

func mapMenuConfigResponseToModel(response *menuConfigResponse, prior *MenuConfigModel) *MenuConfigModel {
	if response == nil {
		return nil
	}

	var options []MenuOptionModel
	if response.Options != nil {
		for i, option := range *response.Options {
			var priorRoute *RouteConfigModel
			if prior != nil && i < len(prior.Options) {
				priorRoute = &prior.Options[i].Route
			}
			options = append(options, MenuOptionModel{
				Key:   types.StringValue(option.Key),
				Route: *mapRouteConfigResponseToModel(&option.Route, priorRoute),
			})
		}
	}

	return &MenuConfigModel{
		Options:                  options,
		MenuMaxAttempts:          types.Int64PointerValue(response.MenuMaxAttempts),
		MenuGatherTimeoutSeconds: types.Int64PointerValue(response.MenuGatherTimeoutSeconds),
	}
}

func mergeRouteConfigFromPrior(model *RouteConfigModel, prior *RouteConfigModel) {
	model.IncidentSeverity = mergeStringFromPrior(model.IncidentSeverity, prior.IncidentSeverity)
	model.DirectVoicemail = mergeBoolFromPrior(model.DirectVoicemail, prior.DirectVoicemail)
	model.SendToVoicemailIfNoResponse = mergeBoolFromPrior(model.SendToVoicemailIfNoResponse, prior.SendToVoicemailIfNoResponse)
	model.AutoResolveOnAnswer = mergeBoolFromPrior(model.AutoResolveOnAnswer, prior.AutoResolveOnAnswer)
	model.SkipAlreadyDialedNumbers = mergeBoolFromPrior(model.SkipAlreadyDialedNumbers, prior.SkipAlreadyDialedNumbers)
	model.MaxUsersToTry = mergeInt64FromPrior(model.MaxUsersToTry, prior.MaxUsersToTry)
	model.DialTimeoutPerUserInSeconds = mergeInt64FromPrior(model.DialTimeoutPerUserInSeconds, prior.DialTimeoutPerUserInSeconds)
	model.TargetConfirmGatherTimeoutSeconds = mergeInt64FromPrior(model.TargetConfirmGatherTimeoutSeconds, prior.TargetConfirmGatherTimeoutSeconds)
	model.TargetConfirmMaxAttempts = mergeInt64FromPrior(model.TargetConfirmMaxAttempts, prior.TargetConfirmMaxAttempts)
}

func mergeStringFromPrior(mapped, prior types.String) types.String {
	if !prior.IsNull() && !prior.IsUnknown() {
		if mapped.IsNull() || mapped.IsUnknown() || mapped.ValueString() == "" {
			return prior
		}
	}
	return mapped
}

func mergeBoolFromPrior(mapped, prior types.Bool) types.Bool {
	if mapped.IsNull() && !prior.IsNull() && !prior.IsUnknown() {
		return prior
	}
	return mapped
}

func mergeInt64FromPrior(mapped, prior types.Int64) types.Int64 {
	if mapped.IsNull() && !prior.IsNull() && !prior.IsUnknown() {
		return prior
	}
	return mapped
}

func mergeStringListFromPrior(mapped, prior types.List) types.List {
	if prior.IsNull() || prior.IsUnknown() {
		if mapped.IsNull() || mapped.IsUnknown() {
			return prior
		}
		if len(mapped.Elements()) == 0 {
			return prior
		}
	}
	return mapped
}

func mapRouteConfigResponseToModel(response *routeConfigResponse, prior *RouteConfigModel) *RouteConfigModel {
	if response == nil {
		return nil
	}

	var priorVoicemail *PromptModel
	var priorClosing *PromptModel
	var priorAccept *PromptModel
	if prior != nil {
		priorVoicemail = prior.VoicemailMessage
		priorClosing = prior.ClosingMessage
		priorAccept = prior.AcceptCallPrompt
	}

	model := &RouteConfigModel{
		TeamId:                            types.StringPointerValue(response.TeamId),
		IncidentSeverity:                  types.StringValue(response.IncidentSeverity),
		DirectVoicemail:                   types.BoolPointerValue(response.DirectVoicemail),
		SendToVoicemailIfNoResponse:       types.BoolPointerValue(response.SendToVoicemailIfNoResponse),
		VoicemailMessage:                  mapPromptResponseToModel(response.VoicemailMessage, priorVoicemail),
		ClosingMessage:                    mapPromptResponseToModel(response.ClosingMessage, priorClosing),
		AutoResolveOnAnswer:               types.BoolPointerValue(response.AutoResolveOnAnswer),
		SkipAlreadyDialedNumbers:          types.BoolPointerValue(response.SkipAlreadyDialedNumbers),
		MaxUsersToTry:                     types.Int64PointerValue(response.MaxUsersToTry),
		DialTimeoutPerUserInSeconds:       types.Int64PointerValue(response.DialTimeoutPerUserInSeconds),
		AcceptCallPrompt:                  mapPromptResponseToModel(response.AcceptCallPrompt, priorAccept),
		TargetConfirmGatherTimeoutSeconds: types.Int64PointerValue(response.TargetConfirmGatherTimeoutSeconds),
		TargetConfirmMaxAttempts:          types.Int64PointerValue(response.TargetConfirmMaxAttempts),
	}

	if prior != nil {
		mergeRouteConfigFromPrior(model, prior)
	}

	return model
}

func mapPromptResponseToModel(response *promptResponse, prior *PromptModel) *PromptModel {
	if response == nil {
		return nil
	}

	return &PromptModel{
		Type:      types.StringValue(response.Type),
		Text:      types.StringPointerValue(response.Text),
		Locale:    types.StringPointerValue(response.Locale),
		Voice:     types.StringPointerValue(response.Voice),
		AudioFile: mapAudioFileResponseToModel(response.AudioFile, priorAudioFile(prior)),
	}
}

func priorAudioFile(prior *PromptModel) *AudioFileModel {
	if prior == nil {
		return nil
	}
	return prior.AudioFile
}

func mapAudioFileResponseToModel(response *callRoutingAudioFileResponse, prior *AudioFileModel) *AudioFileModel {
	if response == nil {
		return nil
	}

	model := &AudioFileModel{
		ObjectKey:     types.StringValue(response.ObjectKey),
		FileName:      types.StringPointerValue(response.FileName),
		ContentType:   types.StringPointerValue(response.ContentType),
		FileSizeBytes: types.Int64PointerValue(response.FileSizeBytes),
	}

	if prior != nil {
		if !prior.Source.IsNull() && !prior.Source.IsUnknown() {
			model.Source = prior.Source
		}
		if !prior.SourceHash.IsNull() && !prior.SourceHash.IsUnknown() {
			model.SourceHash = prior.SourceHash
		}
	}

	return model
}

func stringListForAPI(list types.List) *[]string {
	if list.IsNull() || list.IsUnknown() {
		empty := []string{}
		return &empty
	}

	return ListToStringArray(list)
}

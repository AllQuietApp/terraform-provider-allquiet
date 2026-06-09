package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMapCallFlowConfigResponseToModel_PreservesNullCallerListsFromPrior(t *testing.T) {
	ctx := context.Background()
	prior := &CallFlowConfigModel{
		CallerAllowlist: types.ListNull(types.StringType),
		CallerBlocklist: types.ListNull(types.StringType),
		CallDisplayMode: types.StringValue("ShowInboundNumber"),
	}

	response := &callFlowConfigResponse{
		CallDisplayMode: "ShowInboundNumber",
		CallerAllowlist: &[]string{},
		CallerBlocklist: &[]string{},
	}

	model := mapCallFlowConfigResponseToModel(ctx, response, prior)

	if !model.CallerAllowlist.IsNull() {
		t.Fatal("expected caller allowlist to remain null")
	}
	if !model.CallerBlocklist.IsNull() {
		t.Fatal("expected caller blocklist to remain null")
	}
}

func TestMapRouteConfigResponseToModel_PreservesSchemaDefaultsFromPrior(t *testing.T) {
	prior := &RouteConfigModel{
		IncidentSeverity:            types.StringValue("Critical"),
		DirectVoicemail:             types.BoolValue(false),
		SendToVoicemailIfNoResponse: types.BoolValue(true),
		MaxUsersToTry:               types.Int64Value(3),
	}

	response := &routeConfigResponse{
		TeamId: stringPtr("11111111-1111-1111-1111-111111111111"),
	}

	model := mapRouteConfigResponseToModel(response, prior)

	if got, want := model.IncidentSeverity.ValueString(), "Critical"; got != want {
		t.Fatalf("IncidentSeverity = %q, want %q", got, want)
	}
	if model.DirectVoicemail.ValueBool() {
		t.Fatal("expected DirectVoicemail to remain false")
	}
	if !model.SendToVoicemailIfNoResponse.ValueBool() {
		t.Fatal("expected SendToVoicemailIfNoResponse to remain true")
	}
	if got, want := model.MaxUsersToTry.ValueInt64(), int64(3); got != want {
		t.Fatalf("MaxUsersToTry = %d, want %d", got, want)
	}
}

func TestMapTwilioAccessSettingsResponseToModel_PreservesSensitiveValuesFromPrior(t *testing.T) {
	prior := &TwilioAccessSettingsModel{
		AccountSid: types.StringValue("AC123"),
		AuthToken:  types.StringValue("secret"),
	}

	response := &twilioAccessSettingsResponse{
		ApiBaseDomain: "api.twilio.com",
		AccountSid:    stringPtr("AC999"),
		PhoneNumber:   stringPtr("+15551234567"),
	}

	model := mapTwilioAccessSettingsResponseToModel(response, prior)

	if got, want := model.AccountSid.ValueString(), "AC123"; got != want {
		t.Fatalf("AccountSid = %q, want %q", got, want)
	}
	if got, want := model.AuthToken.ValueString(), "secret"; got != want {
		t.Fatalf("AuthToken = %q, want %q", got, want)
	}
}

func stringPtr(value string) *string {
	return &value
}

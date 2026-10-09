package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMapIntegrationCreateRequest_NoiseLimitInSeconds(t *testing.T) {
	limit := int64(120)
	plan := &IntegrationModel{
		DisplayName:         types.StringValue("Datadog"),
		TeamId:              types.StringValue("team-id"),
		Labels:              types.ListNull(types.StringType),
		IsMuted:             types.BoolValue(false),
		IsInMaintenance:     types.BoolValue(false),
		Type:                types.StringValue("Datadog"),
		NoiseLimitInSeconds: types.Int64Value(limit),
	}

	request := mapIntegrationCreateRequest(plan)
	if request.NoiseLimitInSeconds == nil || *request.NoiseLimitInSeconds != limit {
		t.Fatalf("expected noise limit %d, got %#v", limit, request.NoiseLimitInSeconds)
	}
}

func TestMapIntegrationCreateRequest_OmitsUnsetNoiseLimit(t *testing.T) {
	plan := &IntegrationModel{
		DisplayName:         types.StringValue("Datadog"),
		TeamId:              types.StringValue("team-id"),
		Labels:              types.ListNull(types.StringType),
		IsMuted:             types.BoolValue(false),
		IsInMaintenance:     types.BoolValue(false),
		Type:                types.StringValue("Datadog"),
		NoiseLimitInSeconds: types.Int64Null(),
	}

	request := mapIntegrationCreateRequest(plan)
	if request.NoiseLimitInSeconds != nil {
		t.Fatalf("expected nil noise limit, got %#v", request.NoiseLimitInSeconds)
	}
}

func TestMapIntegrationResponseToModel_NoiseLimitInSeconds(t *testing.T) {
	ctx := context.Background()
	limit := int64(90)
	response := &integrationResponse{
		Id:                  "integration-id",
		DisplayName:         "Datadog",
		TeamId:              "team-id",
		Type:                "Datadog",
		NoiseLimitInSeconds: &limit,
	}
	data := &IntegrationModel{}

	mapIntegrationResponseToModel(ctx, response, data, nil)

	if data.NoiseLimitInSeconds.IsNull() || data.NoiseLimitInSeconds.ValueInt64() != limit {
		t.Fatalf("expected noise limit %d, got %#v", limit, data.NoiseLimitInSeconds)
	}
}

func TestMapIntegrationResponseToModel_NullNoiseLimit(t *testing.T) {
	ctx := context.Background()
	response := &integrationResponse{
		Id:          "integration-id",
		DisplayName: "Datadog",
		TeamId:      "team-id",
		Type:        "Datadog",
	}
	data := &IntegrationModel{
		NoiseLimitInSeconds: types.Int64Value(120),
	}

	mapIntegrationResponseToModel(ctx, response, data, nil)

	if !data.NoiseLimitInSeconds.IsNull() {
		t.Fatalf("expected null noise limit, got %#v", data.NoiseLimitInSeconds)
	}
}

func TestMapHttpMonitoringCreateRequest_DisableFollowRedirect(t *testing.T) {
	plan := &HttpMonitoringModel{
		Headers:                     types.MapNull(types.StringType),
		OverrideAcceptedStatusCodes: types.ListNull(types.Int64Type),
		DisableFollowRedirect:       types.BoolValue(true),
	}

	request := mapHttpMonitoringCreateRequest(plan)
	if request == nil || !request.DisableFollowRedirect {
		t.Fatalf("expected disableFollowRedirect true, got %#v", request)
	}

	body, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["disableFollowRedirect"] != true {
		t.Fatalf("expected disableFollowRedirect true in JSON, got %#v", payload["disableFollowRedirect"])
	}
}

func TestMapWebhookAuthenticationUpdateRequest_OmitsWhenPlanAndPriorAbsent(t *testing.T) {
	got := mapWebhookAuthenticationUpdateRequest(nil, nil)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
}

func TestMapWebhookAuthenticationUpdateRequest_ClearsWhenRemovedFromPlan(t *testing.T) {
	prior := &WebhookAuthenticationModel{
		Type: types.StringValue("bearer"),
		Bearer: &BearerModel{
			Token: types.StringValue("secret"),
		},
	}
	got := mapWebhookAuthenticationUpdateRequest(nil, prior)
	if got == nil || got.Type != "none" {
		t.Fatalf("expected type none, got %#v", got)
	}
	if got.Bearer != nil || got.PrivateIpFilter != nil {
		t.Fatalf("expected no bearer or ip filter, got %#v", got)
	}
}

func TestMapWebhookAuthenticationUpdateRequest_SendsPlannedBearer(t *testing.T) {
	plan := &WebhookAuthenticationModel{
		Type: types.StringValue("bearer"),
		Bearer: &BearerModel{
			Token: types.StringValue("my-token"),
		},
	}
	got := mapWebhookAuthenticationUpdateRequest(plan, nil)
	if got == nil || got.Type != "bearer" || got.Bearer == nil || got.Bearer.Token != "my-token" {
		t.Fatalf("expected bearer token, got %#v", got)
	}
}

func TestMapIntegrationUpdateRequest_OmitsWebhookAuthWhenNeverConfigured(t *testing.T) {
	plan := &IntegrationModel{
		DisplayName:     types.StringValue("Datadog"),
		TeamId:          types.StringValue("team-id"),
		Labels:          types.ListNull(types.StringType),
		IsMuted:         types.BoolValue(false),
		IsInMaintenance: types.BoolValue(false),
		Type:            types.StringValue("Datadog"),
	}
	prior := &IntegrationModel{
		DisplayName:     types.StringValue("Datadog"),
		TeamId:          types.StringValue("team-id"),
		IsMuted:         types.BoolValue(false),
		IsInMaintenance: types.BoolValue(false),
		Type:            types.StringValue("Datadog"),
	}

	request := mapIntegrationUpdateRequest(plan, prior)
	if request.WebhookAuthentication != nil {
		t.Fatalf("expected nil webhook authentication, got %#v", request.WebhookAuthentication)
	}
}

func TestNormalizeWebhookAuthenticationInState_PlanOmitted(t *testing.T) {
	mapped := &WebhookAuthenticationModel{
		Type: types.StringValue("bearer"),
		Bearer: &BearerModel{
			Token: types.StringValue("secret"),
		},
	}
	got := normalizeWebhookAuthenticationInState(mapped, mapped, true)
	if got != nil {
		t.Fatalf("expected nil when plan omitted block, got %#v", got)
	}
}

func TestNormalizeWebhookAuthenticationInState_ClearedApiWithNullPrior(t *testing.T) {
	mapped := mapWebhookAuthenticationResponseToModel(&webhookAuthenticationResponse{Type: "none"})
	got := normalizeWebhookAuthenticationInState(mapped, nil, false)
	if got != nil {
		t.Fatalf("expected nil for cleared api shape, got %#v", got)
	}
}

func TestNormalizeWebhookAuthenticationInState_KeepsBearerFromApi(t *testing.T) {
	mapped := mapWebhookAuthenticationResponseToModel(&webhookAuthenticationResponse{
		Type: "bearer",
		Bearer: &webhookAuthenticationBearerResponse{
			Token: "secret",
		},
	})
	got := normalizeWebhookAuthenticationInState(mapped, nil, false)
	if got == nil || got.Bearer == nil || got.Bearer.Token.ValueString() != "secret" {
		t.Fatalf("expected bearer in state, got %#v", got)
	}
}

func TestNormalizeWebhookAuthenticationInState_KeepsExplicitNoneInPrior(t *testing.T) {
	prior := &WebhookAuthenticationModel{
		Type: types.StringValue("none"),
	}
	mapped := mapWebhookAuthenticationResponseToModel(&webhookAuthenticationResponse{Type: "none"})
	got := normalizeWebhookAuthenticationInState(mapped, prior, false)
	if got == nil || got.Type.ValueString() != "none" {
		t.Fatalf("expected explicit none in state, got %#v", got)
	}
}

func TestWebhookAuthenticationResponseIsAbsentEffectively(t *testing.T) {
	if !webhookAuthenticationResponseIsAbsentEffectively(nil) {
		t.Fatal("expected nil response to be absent")
	}
	if !webhookAuthenticationResponseIsAbsentEffectively(&webhookAuthenticationResponse{Type: "none"}) {
		t.Fatal("expected type none to be absent")
	}
	ip := "203.0.113.1"
	if webhookAuthenticationResponseIsAbsentEffectively(&webhookAuthenticationResponse{Type: "none", PrivateIpFilter: &ip}) {
		t.Fatal("expected ip filter to count as present")
	}
}

func TestMapHttpMonitoringResponseToModel_DisableFollowRedirect(t *testing.T) {
	ctx := context.Background()
	response := &httpMonitoringResponse{
		Url:                   "https://example.com",
		Method:                "GET",
		DisableFollowRedirect: true,
	}

	model := mapHttpMonitoringResponseToModel(ctx, response)
	if model == nil || model.DisableFollowRedirect.IsNull() || !model.DisableFollowRedirect.ValueBool() {
		t.Fatalf("expected disable_follow_redirect true, got %#v", model)
	}
}

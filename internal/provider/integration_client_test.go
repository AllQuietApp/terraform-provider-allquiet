package provider

import (
	"context"
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

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

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &OutboundIntegration{}
var _ resource.ResourceWithImportState = &OutboundIntegration{}

func NewOutboundIntegration() resource.Resource {
	return &OutboundIntegration{}
}

// OutboundIntegration defines the resource implementation.
type OutboundIntegration struct {
	client *AllQuietAPIClient
}

// OutboundIntegrationModel describes the resource data model.
type OutboundIntegrationModel struct {
	Id                      types.String `tfsdk:"id"`
	DisplayName             types.String `tfsdk:"display_name"`
	TeamId                  types.String `tfsdk:"team_id"`
	Type                    types.String `tfsdk:"type"`
	TriggersOnlyOnForwarded types.Bool   `tfsdk:"triggers_only_on_forwarded"`

	TeamConnectionSettings *OutboundIntegrationTeamConnectionSettings `tfsdk:"team_connection_settings"`
}

type OutboundIntegrationTeamConnectionSettings struct {
	TeamConnectionMode types.String `tfsdk:"team_connection_mode"`
	TeamIds            types.List   `tfsdk:"team_ids"`
}

func (r *OutboundIntegration) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_outbound_integration"
}

func (r *OutboundIntegration) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
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
				MarkdownDescription: "If true, the integration will only trigger when explicitly forwarded",
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
		},
	}
}

func (r *OutboundIntegration) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OutboundIntegration) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OutboundIntegrationModel

	// Read Terraform plan data into the model
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

	// Read Terraform prior state data into the model
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

	// Read Terraform plan data into the model
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

	// Read Terraform prior state data into the model
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

	if response.TeamConnectionSettings != nil {
		data.TeamConnectionSettings = &OutboundIntegrationTeamConnectionSettings{
			TeamConnectionMode: types.StringValue(response.TeamConnectionSettings.TeamConnectionMode),
			TeamIds:            MapNullableList(ctx, response.TeamConnectionSettings.TeamIds),
		}
	}
}

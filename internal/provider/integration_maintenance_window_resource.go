// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

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
var _ resource.Resource = &IntegrationMaintenanceWindow{}
var _ resource.ResourceWithImportState = &IntegrationMaintenanceWindow{}

func NewIntegrationMaintenanceWindow() resource.Resource {
	return &IntegrationMaintenanceWindow{}
}

// IntegrationMaintenanceWindow defines the resource implementation.
type IntegrationMaintenanceWindow struct {
	client *AllQuietAPIClient
}

// IntegrationMaintenanceWindowModel describes the resource data model.
type IntegrationMaintenanceWindowModel struct {
	Id            types.String `tfsdk:"id"`
	IntegrationId types.String `tfsdk:"integration_id"`
	Start         types.String `tfsdk:"start"`
	End           types.String `tfsdk:"end"`
	Description   types.String `tfsdk:"description"`
	Type          types.String `tfsdk:"type"`
}

func (r *IntegrationMaintenanceWindow) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_maintenance_window"
}

func (r *IntegrationMaintenanceWindow) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The `integration_maintenance_window` resource represents an integration maintenance window in All Quiet. Integration maintenance windows are used to define maintenance windows for an integration.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"integration_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Id of the associated integration",
			},
			"start": schema.StringAttribute{
				Optional:    true,
				Description: "Start of the maintenance window (RFC3339 format)",
				Validators:  []validator.String{DateTimeValidator("Not a valid date")},
			},
			"end": schema.StringAttribute{
				Optional:    true,
				Description: "End of the maintenance window (RFC3339 format)",
				Validators:  []validator.String{DateTimeValidator("Not a valid date")},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the maintenance window",
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of the maintenance window. Possible values are: " + strings.Join(ValidMaintenanceWindowTypes, ", "),
				Validators: []validator.String{
					stringvalidator.OneOf(ValidMaintenanceWindowTypes...),
				},
			},
		},
	}
}

func (r *IntegrationMaintenanceWindow) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*AllQuietAPIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *AllQuietAPIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *IntegrationMaintenanceWindow) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IntegrationMaintenanceWindowModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integrationResponse, err := r.client.CreateIntegrationMaintenanceWindowResource(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create integration maintenance window resource, got error: %s", err))
		return
	}

	mapIntegrationMaintenanceWindowResponseToModel(integrationResponse, &data)

	tflog.Trace(ctx, "created integration maintenance window resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IntegrationMaintenanceWindow) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IntegrationMaintenanceWindowModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	integrationResponse, err := r.client.GetIntegrationMaintenanceWindowResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get integration maintenance window resource, got error: %s", err))
		return
	}

	if integrationResponse == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to get integration maintenance window resource, got nil response")
		return
	}

	mapIntegrationMaintenanceWindowResponseToModel(integrationResponse, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IntegrationMaintenanceWindow) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IntegrationMaintenanceWindowModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	integrationResponse, err := r.client.UpdateIntegrationMaintenanceWindowResource(ctx, data.Id.ValueString(), &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update integration maintenance window resource, got error: %s", err))
		return
	}

	mapIntegrationMaintenanceWindowResponseToModel(integrationResponse, &data)

	tflog.Trace(ctx, "updated integration maintenance window resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IntegrationMaintenanceWindow) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IntegrationMaintenanceWindowModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteIntegrationMaintenanceWindowResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete integration maintenance window resource, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted integration maintenance window resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IntegrationMaintenanceWindow) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapIntegrationMaintenanceWindowResponseToModel(response *integrationMaintenanceWindowResponse, data *IntegrationMaintenanceWindowModel) {
	data.Id = types.StringValue(response.Id)
	data.IntegrationId = types.StringValue(response.IntegrationId)
	data.Start = types.StringPointerValue(response.Start)
	data.End = types.StringPointerValue(response.End)
	data.Description = types.StringPointerValue(response.Description)
	data.Type = types.StringValue(response.Type)
}

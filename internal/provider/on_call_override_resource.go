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
var _ resource.Resource = &OnCallOverride{}
var _ resource.ResourceWithImportState = &OnCallOverride{}

func NewOnCallOverride() resource.Resource {
	return &OnCallOverride{}
}

type OnCallOverride struct {
	client *AllQuietAPIClient
}

type OnCallOverrideModel struct {
	Id                 types.String `tfsdk:"id"`
	UserId             types.String `tfsdk:"user_id"`
	Type               types.String `tfsdk:"type"`
	Start              types.String `tfsdk:"start"`
	End                types.String `tfsdk:"end"`
	ReplacementUserIds types.List   `tfsdk:"replacement_user_ids"`
}

var ValidOnCallOverrideTypes = []string{"online", "offline"}

func OnCallOverrideTypeValidator(message string) validator.String {
	return stringvalidator.OneOf(ValidOnCallOverrideTypes...)
}

func (r *OnCallOverride) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_on_call_override"
}

func (r *OnCallOverride) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The user resource represents a user in All Quiet. Users can be members of users and receive notifications for incidents.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.StringAttribute{
				MarkdownDescription: "The user id of the user",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the override. Possible values are: " + strings.Join(ValidOnCallOverrideTypes, ", "),
				Required:            true,
				Validators:          []validator.String{OnCallOverrideTypeValidator("Invalid on call override type")},
			},
			"start": schema.StringAttribute{
				Required:    true,
				Description: "Start date / time for the override (RFC3339 format)",
				Validators:  []validator.String{DateTimeValidator("Not a valid date / time")},
			},
			"end": schema.StringAttribute{
				Required:    true,
				Description: "End date / time for the override (RFC3339 format)",
				Validators:  []validator.String{DateTimeValidator("Not a valid date / time")},
			},
			"replacement_user_ids": schema.ListAttribute{
				MarkdownDescription: "Replacement user ids",
				Optional:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(GuidValidator("Not a valid GUID")),
				},
			},
		},
	}
}

func (r *OnCallOverride) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OnCallOverride) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OnCallOverrideModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	userResponse, err := r.client.CreateOnCallOverrideResource(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create user resource, got error: %s", err))
		return
	}
	mapOnCallOverrideResponseToModel(ctx, userResponse, &data)

	tflog.Trace(ctx, "created user resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OnCallOverride) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OnCallOverrideModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	userResponse, err := r.client.GetOnCallOverrideResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get user resource, got error: %s", err))
		return
	}

	if userResponse == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to get user resource, got nil response")
		return
	}

	mapOnCallOverrideResponseToModel(ctx, userResponse, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OnCallOverride) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data OnCallOverrideModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	userResponse, err := r.client.UpdateOnCallOverrideResource(ctx, data.Id.ValueString(), &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update override resource, got error: %s", err))
		return
	}

	mapOnCallOverrideResponseToModel(ctx, userResponse, &data)

	tflog.Trace(ctx, "updated override resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OnCallOverride) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OnCallOverrideModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOnCallOverrideResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete override resource, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted override resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OnCallOverride) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapOnCallOverrideResponseToModel(ctx context.Context, response *onCallOverrideResponse, data *OnCallOverrideModel) {
	data.Id = types.StringValue(response.Id)
	data.UserId = types.StringValue(response.UserId)
	data.Type = types.StringValue(response.Type)
	data.Start = types.StringValue(response.Start)
	data.End = types.StringValue(response.End)
	data.ReplacementUserIds = MapNullableList(ctx, response.ReplacementUserIds)
}

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
var _ resource.Resource = &TeamMembership{}
var _ resource.ResourceWithImportState = &TeamMembership{}

func NewTeamMembership() resource.Resource {
	return &TeamMembership{}
}

type TeamMembership struct {
	client *AllQuietAPIClient
}

type TeamMembershipModel struct {
	Id     types.String `tfsdk:"id"`
	TeamId types.String `tfsdk:"team_id"`
	UserId types.String `tfsdk:"user_id"`
	Role   types.String `tfsdk:"role"`
}

func (r *TeamMembership) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_membership"
}

func (r *TeamMembership) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"team_id": schema.StringAttribute{
				MarkdownDescription: "The team id that the user is a member of",
				Required:            true,
			},
			"user_id": schema.StringAttribute{
				MarkdownDescription: "The user id of the user",
				Required:            true,
			},
			"role": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Role of the member. Possible values are: " + strings.Join(ValidTeamMembershipRoles, ", "),
				Validators:          []validator.String{stringvalidator.OneOf(ValidTeamMembershipRoles...)},
			},
		},
	}
}

func (r *TeamMembership) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TeamMembership) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TeamMembershipModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	userResponse, err := r.client.CreateTeamMembershipResource(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create user resource, got error: %s", err))
		return
	}
	mapTeamMembershipResponseToModel(userResponse, &data)

	tflog.Trace(ctx, "created user resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamMembership) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TeamMembershipModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	userResponse, err := r.client.GetTeamMembershipResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get user resource, got error: %s", err))
		return
	}

	if userResponse == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to get user resource, got nil response")
		return
	}

	mapTeamMembershipResponseToModel(userResponse, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamMembership) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TeamMembershipModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	userResponse, err := r.client.UpdateTeamMembershipResource(ctx, data.Id.ValueString(), &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update user resource, got error: %s", err))
		return
	}

	mapTeamMembershipResponseToModel(userResponse, &data)

	tflog.Trace(ctx, "updated user resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamMembership) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TeamMembershipModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTeamMembershipResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete user resource, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted user resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TeamMembership) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapTeamMembershipResponseToModel(response *teamMembershipResponse, data *TeamMembershipModel) {
	data.Id = types.StringValue(response.Id)
	data.UserId = types.StringValue(response.UserId)
	data.TeamId = types.StringValue(response.TeamId)
	data.Role = types.StringValue(response.Role)
}

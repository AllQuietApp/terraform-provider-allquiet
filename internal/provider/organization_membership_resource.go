// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

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
var _ resource.Resource = &OrganizationMembership{}
var _ resource.ResourceWithImportState = &OrganizationMembership{}

func NewOrganizationMembership() resource.Resource {
	return &OrganizationMembership{}
}

type OrganizationMembership struct {
	client *AllQuietAPIClient
}

type OrganizationMembershipModel struct {
	Id     types.String `tfsdk:"id"`
	UserId types.String `tfsdk:"user_id"`
	Role   types.String `tfsdk:"role"`
}

func (r *OrganizationMembership) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_membership"
}

func (r *OrganizationMembership) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"role": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Role of the member (either 'Owner' or 'Administrator')",
				Validators:          []validator.String{stringvalidator.OneOf([]string{"Owner", "Administrator"}...)},
			},
		},
	}
}

func (r *OrganizationMembership) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationMembership) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OrganizationMembershipModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	organizationMembershipResponse, err := r.client.CreateOrganizationMembershipResource(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create organization membership resource, got error: %s", err))
		return
	}
	mapOrganizationMembershipResponseToModel(organizationMembershipResponse, &data)

	tflog.Trace(ctx, "created organization membership resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationMembership) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OrganizationMembershipModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	organizationMembershipResponse, err := r.client.GetOrganizationMembershipResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get organization membership resource, got error: %s", err))
		return
	}

	if organizationMembershipResponse == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to get organization membership resource, got nil response")
		return
	}

	mapOrganizationMembershipResponseToModel(organizationMembershipResponse, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationMembership) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data OrganizationMembershipModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	organizationMembershipResponse, err := r.client.UpdateOrganizationMembershipResource(ctx, data.Id.ValueString(), &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update organization membership resource, got error: %s", err))
		return
	}

	mapOrganizationMembershipResponseToModel(organizationMembershipResponse, &data)

	tflog.Trace(ctx, "updated organization membership resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationMembership) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OrganizationMembershipModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOrganizationMembershipResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete organization membership resource, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted organization membership resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationMembership) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapOrganizationMembershipResponseToModel(response *organizationMembershipResponse, data *OrganizationMembershipModel) {
	data.Id = types.StringValue(response.Id)
	data.UserId = types.StringValue(response.UserId)
	data.Role = types.StringValue(response.Role)
}

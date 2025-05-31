// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TeamMembershipDataSource{}

func NewTeamMembershipDataSource() datasource.DataSource {
	return &TeamMembershipDataSource{}
}

// TeamMembershipDataSource defines the data source implementation.
type TeamMembershipDataSource struct {
	client *AllQuietAPIClient
}

// TeamMembershipDataSourceModel describes the data source data model.
type TeamMembershipDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	UserId      types.String `tfsdk:"user_id"`
	TeamId      types.String `tfsdk:"team_id"`
	Role        types.String `tfsdk:"role"`
	ActivatedAt types.String `tfsdk:"activated_at"`
}

func (d *TeamMembershipDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_membership"
}

func (d *TeamMembershipDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Team membership data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Team membership ID",
				Computed:            true,
			},
			"user_id": schema.StringAttribute{
				MarkdownDescription: "ID of the user to look up",
				Required:            true,
			},
			"team_id": schema.StringAttribute{
				MarkdownDescription: "ID of the team to look up",
				Required:            true,
			},
			"role": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Role of the user in the team. Possible values are: " + strings.Join(ValidTeamMembershipRoles, ", "),
				Validators:          []validator.String{stringvalidator.OneOf(ValidTeamMembershipRoles...)},
			},
			"activated_at": schema.StringAttribute{
				MarkdownDescription: "Date and time if the membership was activated / accepted by the user",
				Computed:            true,
			},
		},
	}
}

func (d *TeamMembershipDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

func (d *TeamMembershipDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TeamMembershipDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	teamMembershipResponse, err := d.client.GetTeamMembershipDataSource(ctx, &data, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get team membership resource, got error: %s", err))
		return
	}

	if teamMembershipResponse == nil {
		resp.Diagnostics.AddError("Client Error", "Did not find a team membership with the provided id, user_id, or team_id")
		return
	}

	mapTeamMembershipResponseToDataSourceModel(teamMembershipResponse, &data)

	tflog.Trace(ctx, "read a data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapTeamMembershipResponseToDataSourceModel(teamMembershipResponse *teamMembershipDataSourceResponse, data *TeamMembershipDataSourceModel) {
	data.Id = types.StringValue(teamMembershipResponse.Id)
	data.UserId = types.StringValue(teamMembershipResponse.UserId)
	data.TeamId = types.StringValue(teamMembershipResponse.TeamId)
	data.Role = types.StringValue(teamMembershipResponse.Role)
	data.ActivatedAt = types.StringValue(teamMembershipResponse.ActivatedAt)
}

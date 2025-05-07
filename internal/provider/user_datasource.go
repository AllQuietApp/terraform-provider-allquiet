// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &UserDataSource{}

func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

// UserDataSource defines the data source implementation.
type UserDataSource struct {
	client *AllQuietAPIClient
}

// UserDataSourceModel describes the data source data model.
type UserDataSourceModel struct {
	Id             types.String `tfsdk:"id"`
	Email          types.String `tfsdk:"email"`
	DisplayName    types.String `tfsdk:"display_name"`
	ScimExternalId types.String `tfsdk:"scim_external_id"`
}

func (d *UserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *UserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "User data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "User ID",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Email address of the user to look up",
				Optional:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the user to look up",
				Optional:            true,
			},
			"scim_external_id": schema.StringAttribute{
				MarkdownDescription: "If the user was provisioned by SCIM, this is the SCIM external ID of the user to look up",
				Optional:            true,
			},
		},
	}
}

func (d *UserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	userResponse, err := d.client.GetUserDataSource(ctx, &data, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get user resource, got error: %s", err))
		return
	}

	if userResponse == nil {
		resp.Diagnostics.AddError("Client Error", "Did not find a user with the provided id, email, display name, or scim external id")
		return
	}

	mapUserResponseToDataSourceModel(userResponse, &data)

	tflog.Trace(ctx, "read a data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapUserResponseToDataSourceModel(userResponse *userDataSourceResponse, data *UserDataSourceModel) {
	data.Id = types.StringValue(userResponse.Id)
	data.Email = types.StringValue(userResponse.Email)
	data.DisplayName = types.StringValue(userResponse.DisplayName)
	data.ScimExternalId = types.StringValue(userResponse.ScimExternalId)
}

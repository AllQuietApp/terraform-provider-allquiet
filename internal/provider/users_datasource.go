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
var _ datasource.DataSource = &UsersDataSource{}

func NewUsersDataSource() datasource.DataSource {
	return &UsersDataSource{}
}

// UsersDataSource defines the data source implementation.
type UsersDataSource struct {
	client *AllQuietAPIClient
}

// UsersDataSourceModel describes the data source data model.
type UsersDataSourceModel struct {
	Email       types.String          `tfsdk:"email"`
	DisplayName types.String          `tfsdk:"display_name"`
	Users       []UserDataSourceModel `tfsdk:"users"`
}

func (d *UsersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

func (d *UsersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Users data source",
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				MarkdownDescription: "Email address of the user to look up",
				Optional:            true,
				Sensitive:           true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the user to look up",
				Optional:            true,
				Sensitive:           true,
			},
			"users": schema.ListNestedAttribute{
				MarkdownDescription: "List of users",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "User ID",
							Computed:            true,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "Email address of the user",
							Computed:            true,
							Sensitive:           true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "Display name of the user",
							Computed:            true,
							Sensitive:           true,
						},
						"scim_external_id": schema.StringAttribute{
							MarkdownDescription: "If the user was provisioned by SCIM, this is the SCIM external ID of the user",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *UsersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UsersDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	usersResponse, err := d.client.GetUsersDataSource(ctx, &data, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get user resource, got error: %s", err))
		return
	}

	if usersResponse == nil {
		resp.Diagnostics.AddError("Client Error", "Did not find a user with the provided id, email, display name, or scim external id")
		return
	}

	mapUsersResponseToDataSourceModel(usersResponse, &data)

	tflog.Trace(ctx, "read a data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapUsersResponseToDataSourceModel(usersResponse *usersDataSourceResponse, data *UsersDataSourceModel) {
	if usersResponse.Users == nil {
		data.Users = make([]UserDataSourceModel, 0, len(usersResponse.Users))
		return
	}

	data.Users = make([]UserDataSourceModel, 0, len(usersResponse.Users))

	for _, user := range usersResponse.Users {
		data.Users = append(data.Users, UserDataSourceModel{
			Id:             types.StringValue(user.Id),
			Email:          types.StringValue(user.Email),
			DisplayName:    types.StringValue(user.DisplayName),
			ScimExternalId: types.StringValue(user.ScimExternalId),
		})
	}
}

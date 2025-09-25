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
var _ datasource.DataSource = &TeamDataSource{}

func NewTeamDataSource() datasource.DataSource {
	return &TeamDataSource{}
}

// TeamDataSource defines the data source implementation.
type TeamDataSource struct {
	client *AllQuietAPIClient
}

// TeamDataSourceModel describes the data source data model.
type TeamDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	TimeZoneId  types.String `tfsdk:"time_zone_id"`
	Labels      types.List   `tfsdk:"labels"`
}

func (d *TeamDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (d *TeamDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Team data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Team ID",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the team to look up",
				Required:            true,
			},
			"time_zone_id": schema.StringAttribute{
				MarkdownDescription: "The timezone id. Find all timezone ids [here](https://allquiet.app/api/public/v1/timezone)",
				Computed:            true,
			},
			"labels": schema.ListAttribute{
				MarkdownDescription: "Labels of the team",
				Computed:            true,
				ElementType:         types.StringType,
				Optional:            true,
			},
		},
	}
}

func (d *TeamDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TeamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TeamDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	teamResponse, err := d.client.GetTeamDataSource(ctx, &data, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get team resource, got error: %s", err))
		return
	}

	if teamResponse == nil {
		resp.Diagnostics.AddError("Client Error", "Did not find a team with the provided id or display name")
		return
	}

	mapTeamResponseToDataSourceModel(teamResponse, &data)

	tflog.Trace(ctx, "read a data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapTeamResponseToDataSourceModel(teamResponse *teamDataSourceResponse, data *TeamDataSourceModel) {
	data.Id = types.StringValue(teamResponse.Id)
	data.DisplayName = types.StringValue(teamResponse.DisplayName)
}

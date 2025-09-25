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
var _ datasource.DataSource = &TeamsDataSource{}

func NewTeamsDataSource() datasource.DataSource {
	return &TeamsDataSource{}
}

// TeamsDataSource defines the data source implementation.
type TeamsDataSource struct {
	client *AllQuietAPIClient
}

// TeamsDataSourceModel describes the data source data model.
type TeamsDataSourceModel struct {
	DisplayName types.String          `tfsdk:"display_name"`
	Teams       []TeamDataSourceModel `tfsdk:"teams"`
}

func (d *TeamsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_teams"
}

func (d *TeamsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Teams data source",
		Attributes: map[string]schema.Attribute{
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the team to look up",
				Optional:            true,
			},
			"teams": schema.ListNestedAttribute{
				MarkdownDescription: "List of teams",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Team ID",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "Display name of the team",
							Computed:            true,
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
				},
			},
		},
	}
}

func (d *TeamsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TeamsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TeamsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	teamsResponse, err := d.client.GetTeamsDataSource(ctx, &data, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get team resource, got error: %s", err))
		return
	}

	if teamsResponse == nil {
		resp.Diagnostics.AddError("Client Error", "Did not find a team with the provided id or display name")
		return
	}

	mapTeamsResponseToDataSourceModel(ctx, teamsResponse, &data)

	tflog.Trace(ctx, "read a data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapTeamsResponseToDataSourceModel(ctx context.Context, teamsResponse *teamsDataSourceResponse, data *TeamsDataSourceModel) {
	if teamsResponse.Teams == nil {
		data.Teams = make([]TeamDataSourceModel, 0, len(teamsResponse.Teams))
		return
	}

	data.Teams = make([]TeamDataSourceModel, 0, len(teamsResponse.Teams))

	for _, team := range teamsResponse.Teams {
		data.Teams = append(data.Teams, TeamDataSourceModel{
			DisplayName: types.StringValue(team.DisplayName),
			TimeZoneId:  types.StringValue(team.TimeZoneId),
			Id:          types.StringValue(team.Id),
			Labels:      MapNullableList(ctx, team.Labels),
		})
	}
}

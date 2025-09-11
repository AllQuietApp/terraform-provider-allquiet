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
var _ datasource.DataSource = &OnCallOverridesDataSource{}

func NewOnCallOverridesDataSource() datasource.DataSource {
	return &OnCallOverridesDataSource{}
}

// OnCallOverridesDataSource defines the data source implementation.
type OnCallOverridesDataSource struct {
	client *AllQuietAPIClient
}

// OnCallOverridesDataSourceModel describes the data source data model.
type OnCallOverridesDataSourceModel struct {
	UserId          types.String                    `tfsdk:"user_id"`
	OnCallOverrides []OnCallOverrideDataSourceModel `tfsdk:"on_call_overrides"`
}

type OnCallOverrideDataSourceModel struct {
	Id                 types.String `tfsdk:"id"`
	UserId             types.String `tfsdk:"user_id"`
	Type               types.String `tfsdk:"type"`
	Start              types.String `tfsdk:"start"`
	End                types.String `tfsdk:"end"`
	ReplacementUserIds types.List   `tfsdk:"replacement_user_ids"`
}

func (d *OnCallOverridesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_on_call_overrides"
}

func (d *OnCallOverridesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "On call overrides data source",
		Attributes: map[string]schema.Attribute{
			"user_id": schema.StringAttribute{
				MarkdownDescription: "ID of the user to filter by",
				Optional:            true,
			},
			"on_call_overrides": schema.ListNestedAttribute{
				MarkdownDescription: "List of on call overrides",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "On call override ID",
							Computed:            true,
						},
						"user_id": schema.StringAttribute{
							MarkdownDescription: "User ID",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the override",
							Computed:            true,
						},
						"start": schema.StringAttribute{
							MarkdownDescription: "Start date and time of the override",
							Computed:            true,
						},
						"end": schema.StringAttribute{
							MarkdownDescription: "End date and time of the override",
							Computed:            true,
						},
						"replacement_user_ids": schema.ListAttribute{
							MarkdownDescription: "Replacement user IDs",
							Computed:            true,
							ElementType:         types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *OnCallOverridesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OnCallOverridesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OnCallOverridesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	onCallOverridesResponse, err := d.client.GetOnCallOverridesDataSource(ctx, &data, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get team membership resource, got error: %s", err))
		return
	}

	if onCallOverridesResponse == nil {
		resp.Diagnostics.AddError("Client Error", "Did not find a team membership with the provided id, user_id, or team_id")
		return
	}

	mapOnCallOverridesResponseToDataSourceModel(ctx, onCallOverridesResponse, &data)

	tflog.Trace(ctx, "read a data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapOnCallOverridesResponseToDataSourceModel(ctx context.Context, onCallOverridesResponse *onCallOverridesDataSourceResponse, data *OnCallOverridesDataSourceModel) {
	data.OnCallOverrides = []OnCallOverrideDataSourceModel{}
	for _, onCallOverride := range onCallOverridesResponse.OnCallOverrides {
		data.OnCallOverrides = append(data.OnCallOverrides, OnCallOverrideDataSourceModel{
			Id:                 types.StringValue(onCallOverride.Id),
			UserId:             types.StringValue(onCallOverride.UserId),
			Type:               types.StringValue(onCallOverride.Type),
			Start:              types.StringValue(onCallOverride.Start),
			End:                types.StringValue(onCallOverride.End),
			ReplacementUserIds: MapNullableList(ctx, onCallOverride.ReplacementUserIds),
		})
	}
}

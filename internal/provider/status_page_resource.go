// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
var _ resource.Resource = &StatusPage{}
var _ resource.ResourceWithImportState = &StatusPage{}

func NewStatusPage() resource.Resource {
	return &StatusPage{}
}

// StatusPage defines the resource implementation.
type StatusPage struct {
	client *AllQuietAPIClient
}

// StatusPageModel describes the resource data model.
type StatusPageModel struct {
	Id                            types.String        `tfsdk:"id"`
	DisplayName                   types.String        `tfsdk:"display_name"`
	PublicTitle                   types.String        `tfsdk:"public_title"`
	PublicDescription             types.String        `tfsdk:"public_description"`
	Slug                          types.String        `tfsdk:"slug"`
	Services                      types.List          `tfsdk:"services"`
	PublicCompanyUrl              types.String        `tfsdk:"public_company_url"`
	PublicCompanyName             types.String        `tfsdk:"public_company_name"`
	PublicSupportUrl              types.String        `tfsdk:"public_support_url"`
	PublicSupportEmail            types.String        `tfsdk:"public_support_email"`
	HistoryInDays                 types.Int64         `tfsdk:"history_in_days"`
	TimeZoneId                    types.String        `tfsdk:"time_zone_id"`
	DisablePublicSubscription     types.Bool          `tfsdk:"disable_public_subscription"`
	PublicSeverityMappingMinor    types.String        `tfsdk:"public_severity_mapping_minor"`
	PublicSeverityMappingWarning  types.String        `tfsdk:"public_severity_mapping_warning"`
	PublicSeverityMappingCritical types.String        `tfsdk:"public_severity_mapping_critical"`
	BannerBackgroundColor         types.String        `tfsdk:"banner_background_color"`
	BannerBackgroundColorDarkMode types.String        `tfsdk:"banner_background_color_dark_mode"`
	BannerTextColor               types.String        `tfsdk:"banner_text_color"`
	BannerTextColorDarkMode       types.String        `tfsdk:"banner_text_color_dark_mode"`
	CustomHostSettings            *CustomHostSettings `tfsdk:"custom_host_settings"`
}

type CustomHostSettings struct {
	Host types.String `tfsdk:"host"`
}

func (r *StatusPage) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_status_page"
}

func (r *StatusPage) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The `status_page` resource represents a status page in All Quiet.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the status page",
				Required:            true,
			},
			"public_title": schema.StringAttribute{
				MarkdownDescription: "The public title of the status page",
				Required:            true,
			},
			"public_description": schema.StringAttribute{
				MarkdownDescription: "The public description of the status page",
				Optional:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "The slug of the status page. Provide slug or custom host settings.",
				Optional:            true,
			},
			"custom_host_settings": schema.SingleNestedAttribute{
				MarkdownDescription: "The custom host settings of the status page (CNAME). Provide slug or custom host settings.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"host": schema.StringAttribute{
						MarkdownDescription: "The host of the status page",
						Required:            true,
					},
				},
			},
			"services": schema.ListAttribute{
				Optional:            true,
				MarkdownDescription: "The service ids of the status page",
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(GuidValidator("Not a valid GUID")),
				},
			},
			"public_company_url": schema.StringAttribute{
				MarkdownDescription: "The public company url of the status page",
				Optional:            true,
			},
			"public_company_name": schema.StringAttribute{
				MarkdownDescription: "The public company name of the status page",
				Optional:            true,
			},
			"public_support_url": schema.StringAttribute{
				MarkdownDescription: "The public support url of the status page",
				Optional:            true,
			},
			"public_support_email": schema.StringAttribute{
				MarkdownDescription: "The public support email of the status page",
				Optional:            true,
			},
			"history_in_days": schema.Int64Attribute{
				MarkdownDescription: "The history in days of the status page",
				Required:            true,
			},
			"time_zone_id": schema.StringAttribute{
				MarkdownDescription: "The time zone id of the status page",
				Optional:            true,
			},
			"disable_public_subscription": schema.BoolAttribute{
				MarkdownDescription: "The disable public subscription of the status page",
				Required:            true,
			},
			"public_severity_mapping_minor": schema.StringAttribute{
				MarkdownDescription: "The public severity mapping minor of the status page",
				Optional:            true,
			},
			"public_severity_mapping_warning": schema.StringAttribute{
				MarkdownDescription: "The public severity mapping warning of the status page",
				Optional:            true,
			},
			"public_severity_mapping_critical": schema.StringAttribute{
				MarkdownDescription: "The public severity mapping critical of the status page",
				Optional:            true,
			},
			"banner_background_color": schema.StringAttribute{
				MarkdownDescription: "The banner background color of the status page. Must be a valid hex color.",
				Optional:            true,
				Validators: []validator.String{
					HexColorValidator("Not a valid hex color"),
				},
			},
			"banner_background_color_dark_mode": schema.StringAttribute{
				MarkdownDescription: "The banner background color dark mode of the status page. Must be a valid hex color.",
				Optional:            true,
				Validators: []validator.String{
					HexColorValidator("Not a valid hex color"),
				},
			},
			"banner_text_color": schema.StringAttribute{
				MarkdownDescription: "The banner text color of the status page. Must be a valid hex color.",
				Optional:            true,
				Validators: []validator.String{
					HexColorValidator("Not a valid hex color"),
				},
			},
			"banner_text_color_dark_mode": schema.StringAttribute{
				MarkdownDescription: "The banner text color dark mode of the status page. Must be a valid hex color.",
				Optional:            true,
				Validators: []validator.String{
					HexColorValidator("Not a valid hex color"),
				},
			},
		},
	}
}

func (r *StatusPage) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *StatusPage) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data StatusPageModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	statusPageResponse, err := r.client.CreateStatusPageResource(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create status page resource, got error: %s", err))
		return
	}

	mapStatusPageResponseToModel(ctx, statusPageResponse, &data)

	tflog.Trace(ctx, "created status page resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StatusPage) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data StatusPageModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	statusPageResponse, err := r.client.GetStatusPageResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get status page resource, got error: %s", err))
		return
	}

	if statusPageResponse == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to get status page resource, got nil response")
		return
	}

	mapStatusPageResponseToModel(ctx, statusPageResponse, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StatusPage) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data StatusPageModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	statusPageResponse, err := r.client.UpdateStatusPageResource(ctx, data.Id.ValueString(), &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update status page resource, got error: %s", err))
		return
	}

	mapStatusPageResponseToModel(ctx, statusPageResponse, &data)

	tflog.Trace(ctx, "updated status page resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StatusPage) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data StatusPageModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteStatusPageResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete status page resource, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted status page resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StatusPage) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapStatusPageResponseToModel(ctx context.Context, response *statusPageResponse, data *StatusPageModel) {

	data.Id = types.StringValue(response.Id)
	data.DisplayName = types.StringValue(response.DisplayName)
	data.PublicTitle = types.StringValue(response.PublicTitle)
	data.PublicDescription = types.StringPointerValue(response.PublicDescription)
	data.Slug = types.StringPointerValue(response.Slug)
	data.Services = MapNullableList(ctx, response.ServiceIds)
	data.PublicCompanyUrl = types.StringPointerValue(response.PublicCompanyUrl)
	data.PublicCompanyName = types.StringPointerValue(response.PublicCompanyName)
	data.PublicSupportUrl = types.StringPointerValue(response.PublicSupportUrl)
	data.PublicSupportEmail = types.StringPointerValue(response.PublicSupportEmail)
	data.HistoryInDays = types.Int64Value(response.HistoryInDays)
	data.TimeZoneId = types.StringPointerValue(response.TimeZoneId)
	data.DisablePublicSubscription = types.BoolValue(response.DisablePublicSubscription)
	data.PublicSeverityMappingMinor = types.StringPointerValue(response.PublicSeverityMappingMinor)
	data.PublicSeverityMappingWarning = types.StringPointerValue(response.PublicSeverityMappingWarning)
	data.PublicSeverityMappingCritical = types.StringPointerValue(response.PublicSeverityMappingCritical)
	data.BannerBackgroundColor = types.StringPointerValue(response.BannerBackgroundColor)
	data.BannerBackgroundColorDarkMode = types.StringPointerValue(response.BannerBackgroundColorDarkMode)
	data.BannerTextColor = types.StringPointerValue(response.BannerTextColor)
	data.BannerTextColorDarkMode = types.StringPointerValue(response.BannerTextColorDarkMode)
	data.CustomHostSettings = mapCustomHostSettingsResponseToModel(response.CustomHostSettings)
}

func mapCustomHostSettingsResponseToModel(response *customHostSettingsResponse) *CustomHostSettings {
	if response == nil {
		return nil
	}

	return &CustomHostSettings{
		Host: types.StringValue(response.Host),
	}
}

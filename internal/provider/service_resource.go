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
var _ resource.Resource = &Service{}
var _ resource.ResourceWithImportState = &Service{}

func NewService() resource.Resource {
	return &Service{}
}

// Service defines the resource implementation.
type Service struct {
	client *AllQuietAPIClient
}

// ServiceModel describes the resource data model.
type ServiceModel struct {
	Id                     types.String            `tfsdk:"id"`
	DisplayName            types.String            `tfsdk:"display_name"`
	PublicTitle            types.String            `tfsdk:"public_title"`
	PublicDescription      types.String            `tfsdk:"public_description"`
	Templates              *[]ServiceTemplateModel `tfsdk:"templates"`
	TeamConnectionSettings *TeamConnectionSettings `tfsdk:"team_connection_settings"`
}

type ServiceTemplateModel struct {
	Id          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Message     types.String `tfsdk:"message"`
}

func (r *Service) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

func (r *Service) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The `service` resource represents a service in All Quiet.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the service",
				Required:            true,
			},
			"public_title": schema.StringAttribute{
				MarkdownDescription: "The public title of the service",
				Required:            true,
			},
			"public_description": schema.StringAttribute{
				MarkdownDescription: "The public description of the service",
				Optional:            true,
			},
			"templates": schema.ListNestedAttribute{
				MarkdownDescription: "The templates of the service",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Id",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "The display name of the template",
							Required:            true,
						},
						"message": schema.StringAttribute{
							MarkdownDescription: "The message of the template",
							Required:            true,
						},
					},
				},
			},
			"team_connection_settings": schema.SingleNestedAttribute{
				MarkdownDescription: "The team connection settings for the integration",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"team_connection_mode": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The team connection mode for the integration. Possible values are: " + strings.Join(ValidTeamConnectionModes, ", "),
						Validators:          []validator.String{stringvalidator.OneOf(ValidTeamConnectionModes...)},
					},
					"team_ids": schema.ListAttribute{
						MarkdownDescription: "The team ids for the integration. If not provided, team_connection_mode must be set to 'OrganizationTeams'.",
						Optional:            true,
						ElementType:         types.StringType,
						Validators: []validator.List{
							listvalidator.ValueStringsAre(GuidValidator("Not a valid GUID")),
						},
					},
				},
			},
		},
	}
}

func (r *Service) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *Service) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ServiceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceResponse, err := r.client.CreateServiceResource(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create service resource, got error: %s", err))
		return
	}

	mapServiceResponseToModel(ctx, serviceResponse, &data)

	tflog.Trace(ctx, "created service resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Service) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ServiceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	serviceResponse, err := r.client.GetServiceResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get service resource, got error: %s", err))
		return
	}

	if serviceResponse == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to get service resource, got nil response")
		return
	}

	mapServiceResponseToModel(ctx, serviceResponse, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Service) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ServiceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	serviceResponse, err := r.client.UpdateServiceResource(ctx, data.Id.ValueString(), &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update service resource, got error: %s", err))
		return
	}

	mapServiceResponseToModel(ctx, serviceResponse, &data)

	tflog.Trace(ctx, "updated service resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Service) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ServiceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteServiceResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete service resource, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted service resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Service) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapServiceResponseToModel(ctx context.Context, response *serviceResponse, data *ServiceModel) {

	data.Id = types.StringValue(response.Id)
	data.DisplayName = types.StringValue(response.DisplayName)
	data.PublicTitle = types.StringValue(response.PublicTitle)
	data.PublicDescription = types.StringPointerValue(response.PublicDescription)
	data.Templates = mapServiceTemplateResponseToModel(response.Templates)
	data.TeamConnectionSettings = MapTeamConnectionSettingsResponseToModel(ctx, response.TeamConnectionSettings)
}

func mapServiceTemplateResponseToModel(templates *[]serviceTemplate) *[]ServiceTemplateModel {
	if templates == nil {
		return nil
	}

	var result []ServiceTemplateModel

	for _, template := range *templates {
		result = append(result, ServiceTemplateModel{
			Id:          types.StringPointerValue(template.Id),
			DisplayName: types.StringValue(template.DisplayName),
			Message:     types.StringValue(template.Message),
		})
	}

	return &result
}

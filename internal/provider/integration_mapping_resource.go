// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IntegrationMapping{}
var _ resource.ResourceWithImportState = &IntegrationMapping{}

func NewIntegrationMapping() resource.Resource {
	return &IntegrationMapping{}
}

// IntegrationMapping defines the resource implementation.
type IntegrationMapping struct {
	client *AllQuietAPIClient
}

// IntegrationMappingModel describes the resource data model.
type IntegrationMappingModel struct {
	Id                types.String                              `tfsdk:"id"`
	IntegrationId     types.String                              `tfsdk:"integration_id"`
	AttributesMapping *IntegrationMappingAttributesMappingModel `tfsdk:"attributes_mapping"`
}

type IntegrationMappingAttributesMappingModel struct {
	Attributes              []IntegrationMappingAttributeModel `tfsdk:"attributes"`
	GroupingWindowInSeconds types.Int64                        `tfsdk:"grouping_window_in_seconds"`
}

type IntegrationMappingAttributeModel struct {
	Name types.String `tfsdk:"name"`

	IsImage        types.Bool `tfsdk:"is_image"`
	HideInPreviews types.Bool `tfsdk:"hide_in_previews"`
	IsGroupingKey  types.Bool `tfsdk:"is_grouping_key"`

	Mappings []IntegrationMappingMappingModel `tfsdk:"mappings"`
}

type IntegrationMappingMappingModel struct {
	XPath    types.String `tfsdk:"xpath"`
	JSONPath types.String `tfsdk:"json_path"`
	Regex    types.String `tfsdk:"regex"`
	Replace  types.String `tfsdk:"replace"`
	Map      types.String `tfsdk:"map"`
	Static   types.String `tfsdk:"static"`
}

func (r *IntegrationMapping) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_mapping"
}

func (r *IntegrationMapping) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The `integration_mapping` resource represents an integration mapping in All Quiet. Integration mappings are used to map attributes from an external system to All Quiet attributes. This allows you to map attributes like severity or status from an external system to All Quiet attributes. \n\n" +
			"An integration mapping is associated with an integration and is optional. You only need to define an integration mapping if you want to override the default mapping of the integration.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"integration_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Id of the associated integration",
			},
			"attributes_mapping": schema.SingleNestedAttribute{
				MarkdownDescription: "The attributes mapping of the integration",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"grouping_window_in_seconds": schema.Int64Attribute{
						MarkdownDescription: "The grouping window in seconds",
						Optional:            true,
					},
					"attributes": schema.ListNestedAttribute{
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									MarkdownDescription: "The name of the attribute",
									Required:            true,
								},
								"is_image": schema.BoolAttribute{
									MarkdownDescription: "Whether the attribute is an image",
									Optional:            true,
								},
								"hide_in_previews": schema.BoolAttribute{
									MarkdownDescription: "Whether the attribute is hidden in previews",
									Optional:            true,
								},
								"is_grouping_key": schema.BoolAttribute{
									MarkdownDescription: "Whether the attribute is a grouping key",
									Optional:            true,
								},
								"mappings": schema.ListNestedAttribute{
									Required:            true,
									MarkdownDescription: "The attribute's mappings",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"xpath": schema.StringAttribute{
												MarkdownDescription: "A XPath expression to map HTML or XML. ( [w3schools](https://www.w3schools.com/xml/xpath_intro.asp))",
												Optional:            true,
											},
											"json_path": schema.StringAttribute{
												MarkdownDescription: "A JSONPath expression to map JSON ([goessner.net/articles/JsonPath](https://goessner.net/articles/JsonPath/))",
												Optional:            true,
											},
											"regex": schema.StringAttribute{
												MarkdownDescription: "A regular expression to extract parts of text. The regex is evaluated with the .NET/C# flavor. If groups are matched, the named group 'result' is returned. If no group is named 'result' the last group is returned. If no groups are found the whole match is returned. ( regex101.com)",
												Optional:            true,
											},
											"replace": schema.StringAttribute{
												MarkdownDescription: "Works together with the regex. Example: you could use the regex '(\\d+) and the replace value 'https://sentry.io/issues/$1/' to create a link to a Sentry issue.",
												Optional:            true,
											},
											"map": schema.StringAttribute{
												MarkdownDescription: "A simple map expression mapping values from A to B. The expression A->1,B->2,->3 will map the value 'A' to '1' and 'B' to '2' and fallback to '3' if no match is found. You can also omit the fallback. The result will then evaluate to the original value.",
												Optional:            true,
											},
											"static": schema.StringAttribute{
												MarkdownDescription: "A static string. The result will always be this string.",
												Optional:            true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *IntegrationMapping) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IntegrationMapping) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IntegrationMappingModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integrationResponse, err := r.client.CreateIntegrationMappingResource(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create integration resource, got error: %s", err))
		return
	}

	mapIntegrationMappingResponseToModel(integrationResponse, &data)

	tflog.Trace(ctx, "created integration resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IntegrationMapping) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IntegrationMappingModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	integrationResponse, err := r.client.GetIntegrationMappingResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get integration resource, got error: %s", err))
		return
	}

	if integrationResponse == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to get integration resource, got nil response")
		return
	}

	mapIntegrationMappingResponseToModel(integrationResponse, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IntegrationMapping) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IntegrationMappingModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	integrationResponse, err := r.client.UpdateIntegrationMappingResource(ctx, data.Id.ValueString(), &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update integration resource, got error: %s", err))
		return
	}

	mapIntegrationMappingResponseToModel(integrationResponse, &data)

	tflog.Trace(ctx, "updated integration resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IntegrationMapping) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IntegrationMappingModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteIntegrationMappingResource(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update integration resource, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted integration resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IntegrationMapping) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapIntegrationMappingResponseToModel(response *integrationMappingResponse, data *IntegrationMappingModel) {
	data.Id = types.StringValue(response.Id)
	data.IntegrationId = types.StringValue(response.IntegrationId)
	data.AttributesMapping = &IntegrationMappingAttributesMappingModel{
		GroupingWindowInSeconds: types.Int64PointerValue(response.AttributesMapping.GroupingWindowInSeconds),
		Attributes:              make([]IntegrationMappingAttributeModel, len(response.AttributesMapping.Attributes)),
	}

	for i, attribute := range response.AttributesMapping.Attributes {
		data.AttributesMapping.Attributes[i].Name = types.StringValue(attribute.Name)
		data.AttributesMapping.Attributes[i].IsImage = types.BoolPointerValue(attribute.IsImage)
		data.AttributesMapping.Attributes[i].HideInPreviews = types.BoolPointerValue(attribute.HideInPreviews)
		data.AttributesMapping.Attributes[i].IsGroupingKey = types.BoolPointerValue(attribute.IsGroupingKey)

		data.AttributesMapping.Attributes[i].Mappings = make([]IntegrationMappingMappingModel, len(attribute.Mappings))
		for j, mapping := range attribute.Mappings {
			data.AttributesMapping.Attributes[i].Mappings[j] = IntegrationMappingMappingModel{
				XPath:    types.StringPointerValue(mapping.XPath),
				JSONPath: types.StringPointerValue(mapping.JSONPath),
				Regex:    types.StringPointerValue(mapping.Regex),
				Replace:  types.StringPointerValue(mapping.Replace),
				Map:      types.StringPointerValue(mapping.Map),
				Static:   types.StringPointerValue(mapping.Static),
			}
		}
	}
}

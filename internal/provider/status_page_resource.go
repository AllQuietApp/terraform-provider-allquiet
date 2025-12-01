// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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
	Id                                types.String                   `tfsdk:"id"`
	DisplayName                       types.String                   `tfsdk:"display_name"`
	PublicTitle                       types.String                   `tfsdk:"public_title"`
	PublicDescription                 types.String                   `tfsdk:"public_description"`
	Slug                              types.String                   `tfsdk:"slug"`
	Services                          types.List                     `tfsdk:"services"`
	ServiceGroups                     *[]StatusPageServiceGroupModel `tfsdk:"service_groups"`
	PublicCompanyUrl                  types.String                   `tfsdk:"public_company_url"`
	PublicCompanyName                 types.String                   `tfsdk:"public_company_name"`
	PublicSupportUrl                  types.String                   `tfsdk:"public_support_url"`
	PublicSupportEmail                types.String                   `tfsdk:"public_support_email"`
	HistoryInDays                     types.Int64                    `tfsdk:"history_in_days"`
	TimeZoneId                        types.String                   `tfsdk:"time_zone_id"`
	DisablePublicSubscription         types.Bool                     `tfsdk:"disable_public_subscription"`
	PublicSeverityMappingMinor        types.String                   `tfsdk:"public_severity_mapping_minor"`
	PublicSeverityMappingWarning      types.String                   `tfsdk:"public_severity_mapping_warning"`
	PublicSeverityMappingCritical     types.String                   `tfsdk:"public_severity_mapping_critical"`
	BannerBackgroundColor             types.String                   `tfsdk:"banner_background_color"`
	BannerBackgroundColorDarkMode     types.String                   `tfsdk:"banner_background_color_dark_mode"`
	BannerTextColor                   types.String                   `tfsdk:"banner_text_color"`
	BannerTextColorDarkMode           types.String                   `tfsdk:"banner_text_color_dark_mode"`
	CustomHostSettings                *CustomHostSettings            `tfsdk:"custom_host_settings"`
	DisablePublicPage                 types.Bool                     `tfsdk:"disable_public_page"`
	DisablePublicJson                 types.Bool                     `tfsdk:"disable_public_json"`
	PrivateIpFilter                   types.String                   `tfsdk:"private_ip_filter"`
	PrivateUserAuthenticationRequired types.Bool                     `tfsdk:"private_user_authentication_required"`
	EnableSMSSubscription             types.Bool                     `tfsdk:"enable_sms_subscription"`
	BodyBackgroundColor               types.String                   `tfsdk:"body_background_color"`
	BodyBackgroundColorDarkMode       types.String                   `tfsdk:"body_background_color_dark_mode"`
	SecondaryBackgroundColor          types.String                   `tfsdk:"secondary_background_color"`
	SecondaryBackgroundColorDarkMode  types.String                   `tfsdk:"secondary_background_color_dark_mode"`
	PrimaryTextColor                  types.String                   `tfsdk:"primary_text_color"`
	PrimaryTextColorDarkMode          types.String                   `tfsdk:"primary_text_color_dark_mode"`
	SecondaryTextColor                types.String                   `tfsdk:"secondary_text_color"`
	SecondaryTextColorDarkMode        types.String                   `tfsdk:"secondary_text_color_dark_mode"`
	ButtonBackgroundColor             types.String                   `tfsdk:"button_background_color"`
	ButtonBackgroundColorDarkMode     types.String                   `tfsdk:"button_background_color_dark_mode"`
	ButtonTextColor                   types.String                   `tfsdk:"button_text_color"`
	ButtonTextColorDarkMode           types.String                   `tfsdk:"button_text_color_dark_mode"`
	DecimalPlaces                     types.Int64                    `tfsdk:"decimal_places"`
}

type StatusPageServiceGroupModel struct {
	Id                types.String `tfsdk:"id"`
	PublicDisplayName types.String `tfsdk:"public_display_name"`
	PublicDescription types.String `tfsdk:"public_description"`
	Services          types.List   `tfsdk:"services"`
}

type CustomHostSettings struct {
	Host                                   types.String `tfsdk:"host"`
	CloudFlareCreateCustomHostNameResponse types.Object `tfsdk:"cloudflare_create_custom_hostname_response"`
}

type CloudFlareCreateCustomHostNameResponse struct {
	Errors   types.List                          `tfsdk:"errors"`
	Messages types.List                          `tfsdk:"messages"`
	Success  types.Bool                          `tfsdk:"success"`
	Result   *CreateCustomHostNameResponseResult `tfsdk:"result"`
}

type CloudFlareResponseInfo struct {
	Code    types.Int64  `tfsdk:"code"`
	Message types.String `tfsdk:"message"`
}

type CreateCustomHostNameResponseResult struct {
	Id                        types.String                                 `tfsdk:"id"`
	Hostname                  types.String                                 `tfsdk:"hostname"`
	Status                    types.String                                 `tfsdk:"status"`
	OwnershipVerification     *CustomHostSettingsOwnershipVerification     `tfsdk:"ownership_verification"`
	OwnershipVerificationHttp *CustomHostSettingsOwnershipVerificationHttp `tfsdk:"ownership_verification_http"`
	VerificationErrors        types.List                                   `tfsdk:"verification_errors"`
	Ssl                       *CustomHostSettingsSsl                       `tfsdk:"ssl"`
}

type CustomHostSettingsOwnershipVerification struct {
	Type  types.String `tfsdk:"type"`
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type CustomHostSettingsOwnershipVerificationHttp struct {
	HttpBody types.String `tfsdk:"http_body"`
	HttpUrl  types.String `tfsdk:"http_url"`
}

type CustomHostSettingsSsl struct {
	Id                types.String `tfsdk:"id"`
	Method            types.String `tfsdk:"method"`
	Status            types.String `tfsdk:"status"`
	ValidationErrors  types.List   `tfsdk:"validation_errors"`
	ValidationRecords types.List   `tfsdk:"validation_records"`
}

type CustomHostSettingsSslValidationRecord struct {
	Emails   types.List   `tfsdk:"emails"`
	HttpBody types.String `tfsdk:"http_body"`
	HttpUrl  types.String `tfsdk:"http_url"`
	TxtName  types.String `tfsdk:"txt_name"`
	TxtValue types.String `tfsdk:"txt_value"`
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
					"cloudflare_create_custom_hostname_response": schema.SingleNestedAttribute{
						MarkdownDescription: "The CloudFlare custom hostname response containing verification and SSL details",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"errors": schema.ListNestedAttribute{
								MarkdownDescription: "List of errors from CloudFlare",
								Computed:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"code": schema.Int64Attribute{
											MarkdownDescription: "Error code",
											Computed:            true,
										},
										"message": schema.StringAttribute{
											MarkdownDescription: "Error message",
											Computed:            true,
										},
									},
								},
							},
							"messages": schema.ListNestedAttribute{
								MarkdownDescription: "List of messages from CloudFlare",
								Computed:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"code": schema.Int64Attribute{
											MarkdownDescription: "Message code",
											Computed:            true,
										},
										"message": schema.StringAttribute{
											MarkdownDescription: "Message text",
											Computed:            true,
										},
									},
								},
							},
							"success": schema.BoolAttribute{
								MarkdownDescription: "Whether the request was successful",
								Computed:            true,
							},
							"result": schema.SingleNestedAttribute{
								MarkdownDescription: "The result containing custom hostname details",
								Computed:            true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										MarkdownDescription: "The ID of the custom hostname",
										Computed:            true,
									},
									"hostname": schema.StringAttribute{
										MarkdownDescription: "The hostname of the custom hostname",
										Computed:            true,
									},
									"status": schema.StringAttribute{
										MarkdownDescription: "The status of the custom hostname",
										Computed:            true,
									},
									"ownership_verification": schema.SingleNestedAttribute{
										MarkdownDescription: "The ownership verification details for the custom hostname",
										Computed:            true,
										Attributes: map[string]schema.Attribute{
											"type": schema.StringAttribute{
												MarkdownDescription: "The type of ownership verification",
												Computed:            true,
											},
											"name": schema.StringAttribute{
												MarkdownDescription: "The name for ownership verification",
												Computed:            true,
											},
											"value": schema.StringAttribute{
												MarkdownDescription: "The value for ownership verification",
												Computed:            true,
											},
										},
									},
									"ownership_verification_http": schema.SingleNestedAttribute{
										MarkdownDescription: "The HTTP ownership verification details for the custom hostname",
										Computed:            true,
										Attributes: map[string]schema.Attribute{
											"http_body": schema.StringAttribute{
												MarkdownDescription: "The HTTP body for ownership verification",
												Computed:            true,
											},
											"http_url": schema.StringAttribute{
												MarkdownDescription: "The HTTP URL for ownership verification",
												Computed:            true,
											},
										},
									},
									"verification_errors": schema.ListAttribute{
										MarkdownDescription: "List of verification errors for the custom hostname",
										Computed:            true,
										ElementType:         types.StringType,
									},
									"ssl": schema.SingleNestedAttribute{
										MarkdownDescription: "The SSL configuration for the custom hostname",
										Computed:            true,
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												MarkdownDescription: "The SSL ID",
												Computed:            true,
											},
											"method": schema.StringAttribute{
												MarkdownDescription: "The SSL method",
												Computed:            true,
											},
											"status": schema.StringAttribute{
												MarkdownDescription: "The SSL status",
												Computed:            true,
											},
											"validation_errors": schema.ListAttribute{
												MarkdownDescription: "List of SSL validation errors",
												Computed:            true,
												ElementType:         types.StringType,
											},
											"validation_records": schema.ListNestedAttribute{
												MarkdownDescription: "List of SSL validation records",
												Computed:            true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"emails": schema.ListAttribute{
															MarkdownDescription: "List of emails for validation",
															Computed:            true,
															ElementType:         types.StringType,
														},
														"http_body": schema.StringAttribute{
															MarkdownDescription: "The HTTP body for validation",
															Computed:            true,
														},
														"http_url": schema.StringAttribute{
															MarkdownDescription: "The HTTP URL for validation",
															Computed:            true,
														},
														"txt_name": schema.StringAttribute{
															MarkdownDescription: "The TXT record name for validation",
															Computed:            true,
														},
														"txt_value": schema.StringAttribute{
															MarkdownDescription: "The TXT record value for validation",
															Computed:            true,
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
				},
			},
			"services": schema.ListAttribute{
				Optional:            true,
				DeprecationMessage:  "Use service_groups instead",
				MarkdownDescription: "The service ids of the status page",
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(GuidValidator("Not a valid GUID")),
				},
			},
			"service_groups": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The service groups of the status page",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Internal id of the service group",
						},
						"public_display_name": schema.StringAttribute{
							MarkdownDescription: "The public display name of the service group",
							Required:            true,
						},
						"public_description": schema.StringAttribute{
							MarkdownDescription: "The public description of the service group",
							Optional:            true,
						},
						"services": schema.ListAttribute{
							Required:            true,
							MarkdownDescription: "The service ids of the service group",
							ElementType:         types.StringType,
							Validators: []validator.List{
								listvalidator.ValueStringsAre(GuidValidator("Not a valid GUID")),
							},
						},
					},
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
			"disable_public_page": schema.BoolAttribute{
				MarkdownDescription: "Disable public access to the status page. When enabled, the status page will not be publicly accessible.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"disable_public_json": schema.BoolAttribute{
				MarkdownDescription: "Disable public access to the status page JSON API. When enabled, the JSON API endpoint will not be publicly accessible.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"private_ip_filter": schema.StringAttribute{
				MarkdownDescription: "Private IP filter (CIDR format) to restrict access to the status page. Only IPs matching the filter will be able to access the page.",
				Optional:            true,
			},
			"private_user_authentication_required": schema.BoolAttribute{
				MarkdownDescription: "Require user authentication to access the status page. When enabled, users must be authenticated All Quiet users of your organization to view the status page. Private user authentication is not allowed for custom host settings (CNAME).",
				Optional:            true,
			},
			"enable_sms_subscription": schema.BoolAttribute{
				MarkdownDescription: "Enable SMS subscription for status page updates. Allows users to subscribe to status updates via SMS.",
				Optional:            true,
			},
			"body_background_color": schema.StringAttribute{
				MarkdownDescription: "The body background color of the status page. Must be a valid hex color.",
				Optional:            true,
				Validators: []validator.String{
					HexColorValidator("Not a valid hex color"),
				},
			},
			"body_background_color_dark_mode": schema.StringAttribute{
				MarkdownDescription: "The body background color dark mode of the status page. Must be a valid hex color.",
				Optional:            true,
				Validators: []validator.String{
					HexColorValidator("Not a valid hex color"),
				},
			},
			"secondary_background_color": schema.StringAttribute{
				MarkdownDescription: "The secondary background color of the status page. Must be a valid hex color.",
				Optional:            true,
				Validators: []validator.String{
					HexColorValidator("Not a valid hex color"),
				},
			},
			"secondary_background_color_dark_mode": schema.StringAttribute{
				MarkdownDescription: "The secondary background color dark mode of the status page. Must be a valid hex color.",
				Optional:            true,
				Validators: []validator.String{
					HexColorValidator("Not a valid hex color"),
				},
			},
			"primary_text_color": schema.StringAttribute{
				MarkdownDescription: "The primary text color of the status page. Must be a valid hex color.",
				Optional:            true,
				Validators: []validator.String{
					HexColorValidator("Not a valid hex color"),
				},
			},
			"primary_text_color_dark_mode": schema.StringAttribute{
				MarkdownDescription: "The primary text color dark mode of the status page. Must be a valid hex color.",
				Optional:            true,
				Validators: []validator.String{
					HexColorValidator("Not a valid hex color"),
				},
			},
			"secondary_text_color": schema.StringAttribute{
				MarkdownDescription: "The secondary text color of the status page. Must be a valid hex color.",
				Optional:            true,
				Validators: []validator.String{
					HexColorValidator("Not a valid hex color"),
				},
			},
			"secondary_text_color_dark_mode": schema.StringAttribute{
				MarkdownDescription: "The secondary text color dark mode of the status page. Must be a valid hex color.",
				Optional:            true,
				Validators: []validator.String{
					HexColorValidator("Not a valid hex color"),
				},
			},
			"button_background_color": schema.StringAttribute{
				MarkdownDescription: "The button background color of the status page. Must be a valid hex color.",
				Optional:            true,
				Validators: []validator.String{
					HexColorValidator("Not a valid hex color"),
				},
			},
			"button_background_color_dark_mode": schema.StringAttribute{
				MarkdownDescription: "The button background color dark mode of the status page. Must be a valid hex color.",
				Optional:            true,
				Validators: []validator.String{
					HexColorValidator("Not a valid hex color"),
				},
			},
			"button_text_color": schema.StringAttribute{
				MarkdownDescription: "The button text color of the status page. Must be a valid hex color.",
				Optional:            true,
				Validators: []validator.String{
					HexColorValidator("Not a valid hex color"),
				},
			},
			"button_text_color_dark_mode": schema.StringAttribute{
				MarkdownDescription: "The button text color dark mode of the status page. Must be a valid hex color.",
				Optional:            true,
				Validators: []validator.String{
					HexColorValidator("Not a valid hex color"),
				},
			},
			"decimal_places": schema.Int64Attribute{
				MarkdownDescription: "The number of decimal places to display on the status page. Must be between 0 and 8.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 8),
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
	data.CustomHostSettings = mapCustomHostSettingsResponseToModel(ctx, response.CustomHostSettings)
	data.ServiceGroups = mapStatusPageServiceGroupsResponseToModel(ctx, response.ServiceGroups)
	data.DisablePublicPage = types.BoolPointerValue(response.DisablePublicPage)
	data.DisablePublicJson = types.BoolPointerValue(response.DisablePublicJson)
	data.PrivateUserAuthenticationRequired = types.BoolPointerValue(response.PrivateUserAuthenticationRequired)
	data.EnableSMSSubscription = types.BoolPointerValue(response.EnableSMSSubscription)
	data.PrivateIpFilter = types.StringPointerValue(response.PrivateIpFilter)
	data.BodyBackgroundColor = types.StringPointerValue(response.BodyBackgroundColor)
	data.BodyBackgroundColorDarkMode = types.StringPointerValue(response.BodyBackgroundColorDarkMode)
	data.SecondaryBackgroundColor = types.StringPointerValue(response.SecondaryBackgroundColor)
	data.SecondaryBackgroundColorDarkMode = types.StringPointerValue(response.SecondaryBackgroundColorDarkMode)
	data.PrimaryTextColor = types.StringPointerValue(response.PrimaryTextColor)
	data.PrimaryTextColorDarkMode = types.StringPointerValue(response.PrimaryTextColorDarkMode)
	data.SecondaryTextColor = types.StringPointerValue(response.SecondaryTextColor)
	data.SecondaryTextColorDarkMode = types.StringPointerValue(response.SecondaryTextColorDarkMode)
	data.ButtonBackgroundColor = types.StringPointerValue(response.ButtonBackgroundColor)
	data.ButtonBackgroundColorDarkMode = types.StringPointerValue(response.ButtonBackgroundColorDarkMode)
	data.ButtonTextColor = types.StringPointerValue(response.ButtonTextColor)
	data.ButtonTextColorDarkMode = types.StringPointerValue(response.ButtonTextColorDarkMode)
	data.DecimalPlaces = types.Int64PointerValue(response.DecimalPlaces)
}

func mapStatusPageServiceGroupsResponseToModel(ctx context.Context, response *[]statusPageServiceGroupResponse) *[]StatusPageServiceGroupModel {
	if response == nil {
		return nil
	}

	serviceGroups := make([]StatusPageServiceGroupModel, len(*response))
	for i, serviceGroup := range *response {
		serviceGroups[i] = StatusPageServiceGroupModel{
			Id:                types.StringValue(serviceGroup.Id),
			PublicDisplayName: types.StringValue(serviceGroup.PublicDisplayName),
			PublicDescription: types.StringPointerValue(serviceGroup.PublicDescription),
			Services:          MapNullableList(ctx, serviceGroup.ServiceIds),
		}
	}
	return &serviceGroups
}

func mapCustomHostSettingsResponseToModel(ctx context.Context, response *customHostSettingsResponse) *CustomHostSettings {
	if response == nil {
		return nil
	}

	result := &CustomHostSettings{
		Host: types.StringValue(response.Host),
	}

	if response.CloudFlareCreateCustomHostNameResponse != nil {
		result.CloudFlareCreateCustomHostNameResponse = mapCloudFlareCreateCustomHostNameResponseToModel(ctx, response.CloudFlareCreateCustomHostNameResponse)
	} else {
		// Set to null if not present
		result.CloudFlareCreateCustomHostNameResponse = types.ObjectNull(getCloudFlareCreateCustomHostNameResponseObjectType().AttrTypes)
	}

	return result
}

func getCloudFlareCreateCustomHostNameResponseObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"errors":   types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"code": types.Int64Type, "message": types.StringType}}},
			"messages": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"code": types.Int64Type, "message": types.StringType}}},
			"success":  types.BoolType,
			"result": types.ObjectType{AttrTypes: map[string]attr.Type{
				"id":                          types.StringType,
				"hostname":                    types.StringType,
				"status":                      types.StringType,
				"ownership_verification":      types.ObjectType{AttrTypes: map[string]attr.Type{"type": types.StringType, "name": types.StringType, "value": types.StringType}},
				"ownership_verification_http": types.ObjectType{AttrTypes: map[string]attr.Type{"http_body": types.StringType, "http_url": types.StringType}},
				"verification_errors":         types.ListType{ElemType: types.StringType},
				"ssl": types.ObjectType{AttrTypes: map[string]attr.Type{
					"id":                types.StringType,
					"method":            types.StringType,
					"status":            types.StringType,
					"validation_errors": types.ListType{ElemType: types.StringType},
					"validation_records": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
						"emails":    types.ListType{ElemType: types.StringType},
						"http_body": types.StringType,
						"http_url":  types.StringType,
						"txt_name":  types.StringType,
						"txt_value": types.StringType,
					}}},
				}},
			}},
		},
	}
}

func mapCloudFlareCreateCustomHostNameResponseToModel(ctx context.Context, response *cloudFlareCreateCustomHostNameResponse) types.Object {
	objectType := getCloudFlareCreateCustomHostNameResponseObjectType()
	attrValues := make(map[string]attr.Value)

	// Map Errors
	errorObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"code":    types.Int64Type,
			"message": types.StringType,
		},
	}
	if response.Errors != nil {
		errorObjects := make([]types.Object, len(*response.Errors))
		for i, err := range *response.Errors {
			obj, _ := types.ObjectValue(errorObjectType.AttrTypes, map[string]attr.Value{
				"code":    types.Int64Value(int64(err.Code)),
				"message": types.StringPointerValue(err.Message),
			})
			errorObjects[i] = obj
		}
		errorList, _ := types.ListValueFrom(ctx, errorObjectType, errorObjects)
		attrValues["errors"] = errorList
	} else {
		attrValues["errors"] = types.ListNull(errorObjectType)
	}

	// Map Messages
	if response.Messages != nil {
		messageObjects := make([]types.Object, len(*response.Messages))
		for i, msg := range *response.Messages {
			obj, _ := types.ObjectValue(errorObjectType.AttrTypes, map[string]attr.Value{
				"code":    types.Int64Value(int64(msg.Code)),
				"message": types.StringPointerValue(msg.Message),
			})
			messageObjects[i] = obj
		}
		messageList, _ := types.ListValueFrom(ctx, errorObjectType, messageObjects)
		attrValues["messages"] = messageList
	} else {
		attrValues["messages"] = types.ListNull(errorObjectType)
	}

	// Map Success
	attrValues["success"] = types.BoolValue(response.Success)

	// Map Result
	if response.Result != nil {
		resultObj := mapCreateCustomHostNameResponseResultToObject(ctx, response.Result)
		attrValues["result"] = resultObj
	} else {
		attrValues["result"] = types.ObjectNull(map[string]attr.Type{
			"id":                          types.StringType,
			"hostname":                    types.StringType,
			"status":                      types.StringType,
			"ownership_verification":      types.ObjectType{AttrTypes: map[string]attr.Type{"type": types.StringType, "name": types.StringType, "value": types.StringType}},
			"ownership_verification_http": types.ObjectType{AttrTypes: map[string]attr.Type{"http_body": types.StringType, "http_url": types.StringType}},
			"verification_errors":         types.ListType{ElemType: types.StringType},
			"ssl": types.ObjectType{AttrTypes: map[string]attr.Type{
				"id":                types.StringType,
				"method":            types.StringType,
				"status":            types.StringType,
				"validation_errors": types.ListType{ElemType: types.StringType},
				"validation_records": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
					"emails":    types.ListType{ElemType: types.StringType},
					"http_body": types.StringType,
					"http_url":  types.StringType,
					"txt_name":  types.StringType,
					"txt_value": types.StringType,
				}}},
			}},
		})
	}

	obj, _ := types.ObjectValue(objectType.AttrTypes, attrValues)
	return obj
}

func mapCreateCustomHostNameResponseResultToObject(ctx context.Context, result *createCustomHostNameResponseResult) types.Object {
	resultAttrTypes := map[string]attr.Type{
		"id":                          types.StringType,
		"hostname":                    types.StringType,
		"status":                      types.StringType,
		"ownership_verification":      types.ObjectType{AttrTypes: map[string]attr.Type{"type": types.StringType, "name": types.StringType, "value": types.StringType}},
		"ownership_verification_http": types.ObjectType{AttrTypes: map[string]attr.Type{"http_body": types.StringType, "http_url": types.StringType}},
		"verification_errors":         types.ListType{ElemType: types.StringType},
		"ssl": types.ObjectType{AttrTypes: map[string]attr.Type{
			"id":                types.StringType,
			"method":            types.StringType,
			"status":            types.StringType,
			"validation_errors": types.ListType{ElemType: types.StringType},
			"validation_records": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
				"emails":    types.ListType{ElemType: types.StringType},
				"http_body": types.StringType,
				"http_url":  types.StringType,
				"txt_name":  types.StringType,
				"txt_value": types.StringType,
			}}},
		}},
	}
	resultAttrValues := make(map[string]attr.Value)

	resultAttrValues["id"] = types.StringValue(result.Id)
	resultAttrValues["hostname"] = types.StringValue(result.Hostname)
	resultAttrValues["status"] = types.StringValue(result.Status)

	// Ownership verification
	if result.OwnershipVerification != nil {
		ownershipVerificationObj, _ := types.ObjectValue(map[string]attr.Type{"type": types.StringType, "name": types.StringType, "value": types.StringType}, map[string]attr.Value{
			"type":  types.StringPointerValue(result.OwnershipVerification.Type),
			"name":  types.StringPointerValue(result.OwnershipVerification.Name),
			"value": types.StringPointerValue(result.OwnershipVerification.Value),
		})
		resultAttrValues["ownership_verification"] = ownershipVerificationObj
	} else {
		resultAttrValues["ownership_verification"] = types.ObjectNull(map[string]attr.Type{"type": types.StringType, "name": types.StringType, "value": types.StringType})
	}

	// Ownership verification HTTP
	if result.OwnershipVerificationHttp != nil {
		ownershipVerificationHttpObj, _ := types.ObjectValue(map[string]attr.Type{"http_body": types.StringType, "http_url": types.StringType}, map[string]attr.Value{
			"http_body": types.StringPointerValue(result.OwnershipVerificationHttp.HttpBody),
			"http_url":  types.StringPointerValue(result.OwnershipVerificationHttp.HttpUrl),
		})
		resultAttrValues["ownership_verification_http"] = ownershipVerificationHttpObj
	} else {
		resultAttrValues["ownership_verification_http"] = types.ObjectNull(map[string]attr.Type{"http_body": types.StringType, "http_url": types.StringType})
	}

	// Verification errors
	if result.VerificationErrors != nil {
		resultAttrValues["verification_errors"] = MapNullableList(ctx, result.VerificationErrors)
	} else {
		resultAttrValues["verification_errors"] = types.ListNull(types.StringType)
	}

	// SSL
	if result.Ssl != nil {
		sslObj := mapSslToObject(ctx, result.Ssl)
		resultAttrValues["ssl"] = sslObj
	} else {
		resultAttrValues["ssl"] = types.ObjectNull(map[string]attr.Type{
			"id":                types.StringType,
			"method":            types.StringType,
			"status":            types.StringType,
			"validation_errors": types.ListType{ElemType: types.StringType},
			"validation_records": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
				"emails":    types.ListType{ElemType: types.StringType},
				"http_body": types.StringType,
				"http_url":  types.StringType,
				"txt_name":  types.StringType,
				"txt_value": types.StringType,
			}}},
		})
	}

	resultObj, _ := types.ObjectValue(resultAttrTypes, resultAttrValues)
	return resultObj
}

func mapSslToObject(ctx context.Context, ssl *customHostSettingsSsl) types.Object {
	sslAttrTypes := map[string]attr.Type{
		"id":                types.StringType,
		"method":            types.StringType,
		"status":            types.StringType,
		"validation_errors": types.ListType{ElemType: types.StringType},
		"validation_records": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"emails":    types.ListType{ElemType: types.StringType},
			"http_body": types.StringType,
			"http_url":  types.StringType,
			"txt_name":  types.StringType,
			"txt_value": types.StringType,
		}}},
	}
	sslAttrValues := make(map[string]attr.Value)

	sslAttrValues["id"] = types.StringPointerValue(ssl.Id)
	sslAttrValues["method"] = types.StringPointerValue(ssl.Method)
	sslAttrValues["status"] = types.StringPointerValue(ssl.Status)

	if ssl.ValidationErrors != nil {
		sslAttrValues["validation_errors"] = MapNullableList(ctx, ssl.ValidationErrors)
	} else {
		sslAttrValues["validation_errors"] = types.ListNull(types.StringType)
	}

	if ssl.ValidationRecords != nil {
		validationRecordObjectType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"emails":    types.ListType{ElemType: types.StringType},
				"http_body": types.StringType,
				"http_url":  types.StringType,
				"txt_name":  types.StringType,
				"txt_value": types.StringType,
			},
		}
		validationRecordObjects := make([]types.Object, len(*ssl.ValidationRecords))
		for i, record := range *ssl.ValidationRecords {
			recordObj, _ := types.ObjectValue(validationRecordObjectType.AttrTypes, map[string]attr.Value{
				"emails":    MapNullableList(ctx, record.Emails),
				"http_body": types.StringPointerValue(record.HttpBody),
				"http_url":  types.StringPointerValue(record.HttpUrl),
				"txt_name":  types.StringPointerValue(record.TxtName),
				"txt_value": types.StringPointerValue(record.TxtValue),
			})
			validationRecordObjects[i] = recordObj
		}
		validationRecordsList, _ := types.ListValueFrom(ctx, validationRecordObjectType, validationRecordObjects)
		sslAttrValues["validation_records"] = validationRecordsList
	} else {
		sslAttrValues["validation_records"] = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"emails":    types.ListType{ElemType: types.StringType},
				"http_body": types.StringType,
				"http_url":  types.StringType,
				"txt_name":  types.StringType,
				"txt_value": types.StringType,
			},
		})
	}

	sslObj, _ := types.ObjectValue(sslAttrTypes, sslAttrValues)
	return sslObj
}

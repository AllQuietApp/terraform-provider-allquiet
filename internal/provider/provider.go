// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure AllQuietProvider satisfies various provider interfaces.
var _ provider.Provider = &AllQuietProvider{}
var _ provider.ProviderWithFunctions = &AllQuietProvider{}

// AllQuietProvider defines the provider implementation.
type AllQuietProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// AllQuietProviderModel describes the provider data model.
type AllQuietProviderModel struct {
	ApiKey types.String `tfsdk:"api_key"`
	Region types.String `tfsdk:"api_region"`
}

func (p *AllQuietProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "allquiet"
	resp.Version = p.version
}

func (p *AllQuietProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "All Quiet's API key. If not provided explicitly, make sure to provide it via the `ALLQUIET_API_KEY` environment variable",
				Optional:            true,
			},
			"api_region": schema.StringAttribute{
				MarkdownDescription: "All Quiet's API key. US or EU.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"us", "eu"}...),
				},
			},
		},
	}
}

func (p *AllQuietProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config AllQuietProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown All Quiet API Key",
			"The provider cannot create the All Quiet API client as there is an unknown configuration value for the All Quiet API Key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ALLQUIET_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	endpoint := os.Getenv("ALLQUIET_ENDPOINT")
	apiKey := os.Getenv("ALLQUIET_API_KEY")

	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if endpoint == "" {
		if config.Region.ValueString() == "eu" {
			endpoint = "https://allquiet.eu/api/public/v1"
		} else {
			endpoint = "https://allquiet.app/api/public/v1"
		}
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing All Quiet API API Key",
			"The provider cannot create the All Quiet API client as there is a missing or empty value for the All Quiet API api_key. "+
				"Set the username value in the configuration or use the ALLQUIET_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	basicAuthUsername := os.Getenv("ALLQUIET_BASIC_AUTH_USERNAME")
	basicAuthPassword := os.Getenv("ALLQUIET_BASIC_AUTH_PASSWORD")

	var basicAuth *BasicAuth
	if basicAuthUsername != "" && basicAuthPassword != "" {
		basicAuth = &BasicAuth{
			Username: basicAuthUsername,
			Password: basicAuthPassword,
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client := NewAllQuietAPIClient(apiKey, endpoint, basicAuth)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *AllQuietProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTeam,
		NewUser,
		NewTeamMembership,
		NewTeamEscalations,
		NewIntegration,
		NewIntegrationMapping,
		NewOutboundIntegration,
		NewRouting,
		NewService,
		NewStatusPage,
	}
}

func (p *AllQuietProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *AllQuietProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AllQuietProvider{
			version: version,
		}
	}
}

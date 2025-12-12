package validators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type emailIntegrationSettingsValidator struct{}

func (v emailIntegrationSettingsValidator) Description(_ context.Context) string {
	return "When integration type is 'Email', integration_settings.email must be specified"
}

func (v emailIntegrationSettingsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v emailIntegrationSettingsValidator) ValidateObject(ctx context.Context, request validator.ObjectRequest, response *validator.ObjectResponse) {
	// Get the type attribute from the parent resource.
	var integrationType types.String
	diags := request.Config.GetAttribute(ctx, path.Root("type"), &integrationType)
	if diags.HasError() || integrationType.IsNull() || integrationType.IsUnknown() {
		// If we can't get the type, skip validation.
		return
	}

	if integrationType.ValueString() != "Email" {
		return
	}

	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Missing Required Attribute",
			"When integration type is 'Email', an empty `email` block must be specified in the `integration_settings` block: `integration_settings = { email = {} }`",
		)
		return
	}

	configValue := request.ConfigValue
	emailAttr, exists := configValue.Attributes()["email"]
	if !exists {
		response.Diagnostics.AddAttributeError(
			request.Path.AtName("email"),
			"Missing Required Attribute",
			"When integration type is 'Email', an empty `integration_settings.email` must be specified",
		)
		return
	}

	// Check if email is null.
	if emailAttr.IsNull() {
		response.Diagnostics.AddAttributeError(
			request.Path.AtName("email"),
			"Missing Required Attribute",
			"When integration type is 'Email', an empty `integration_settings.email` must be specified",
		)
		return
	}
}

// EmailIntegrationSettings returns a validator that ensures integration_settings.email
// is specified when the integration type is "Email".
func EmailIntegrationSettings() emailIntegrationSettingsValidator {
	return emailIntegrationSettingsValidator{}
}

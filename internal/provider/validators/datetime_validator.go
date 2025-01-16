package validators

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type dateTimeValidator struct {
	message string
}

func (v dateTimeValidator) Description(_ context.Context) string {
	return v.message
}

func (v dateTimeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v dateTimeValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	_, err := time.Parse("2006-01-02T15:04:05Z", request.ConfigValue.ValueString())
	if err != nil {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Date Format",
			fmt.Sprintf("%s: %s", v.message, err.Error()),
		)
	}
}

func DateTime(message string) dateTimeValidator {
	return dateTimeValidator{
		message: message,
	}
}

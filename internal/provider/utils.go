package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NonNullableArrayToStringArray(array *[]string) []string {
	if array == nil {
		return []string{}
	}
	return *array
}

func ListToStringArray(list types.List) *[]string {
	if list.IsNull() {
		return nil
	}

	result := make([]string, len(list.Elements()))

	for i, item := range list.Elements() {
		if item.IsUnknown() || item.IsNull() {
			continue
		}
		strValue, ok := item.(types.String)
		if !ok { // type assertion failed
			return nil
		}

		result[i] = strValue.ValueString()
	}
	return &result
}

func MapNullableList(ctx context.Context, stringArray *[]string) types.List {
	if stringArray == nil {
		return types.ListNull(types.StringType)
	}
	var stringList []types.String
	for _, s := range *stringArray {
		stringList = append(stringList, types.StringValue(s))
	}

	listValue, diags := types.ListValueFrom(ctx, types.StringType, stringList)
	if diags.HasError() {
		return types.List{}
	}

	return listValue
}

type badRequestResponse struct {
	Errors map[string][]string `json:"errors"`
}

func logErrorResponse(resp *http.Response) error {

	err := fmt.Errorf("%s %s: %d", resp.Request.Method, resp.Request.URL.RequestURI(), resp.StatusCode)

	if resp.StatusCode == 400 {

		var badRequestResponse badRequestResponse
		errDecode := json.NewDecoder(resp.Body).Decode(&badRequestResponse)
		if errDecode != nil {
			err = fmt.Errorf("%s\ncould not decode 400 response: %s", err, errDecode)
		}

		for key, value := range badRequestResponse.Errors {
			err = fmt.Errorf("%s\n%s: %s", err, key, value)
		}
	}

	return err
}

func IntentValidator(message string) validator.String {
	return stringvalidator.OneOf([]string{"Investigated", "Commented", "Escalated", "Resolved", "Unresolved", "Created", "Deleted", "Updated"}...)
}

func SeverityValidator(message string) validator.String {
	return stringvalidator.OneOf([]string{"Critical", "Warning", "Minor"}...)
}

func StatusValidator(message string) validator.String {
	return stringvalidator.OneOf([]string{"Open", "Resolved"}...)
}

func RuleFlowValidator(message string) validator.String {
	return stringvalidator.OneOf([]string{"Continue", "Skip"}...)
}

func NotificationChannelValidator(message string) validator.String {
	return stringvalidator.OneOf([]string{"Email", "VoiceCall", "SMS", "Push"}...)
}

func OperatorValidator(message string) validator.String {
	return stringvalidator.OneOf([]string{"=", "!=", "contains"}...)
}

func GuidValidator(message string) validator.String {
	return stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`), message)
}

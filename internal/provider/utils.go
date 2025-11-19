package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/AllQuietApp/terraform-provider-internal/internal/provider/validators"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ValidIntents = []string{
	"Resolved",
	"Investigated",
	"Escalated",
	"Commented",
	"Unresolved",
	"Assigned",
	"Affects",
	"Forwarded",
	"Archived",
	"Unarchived",
	"Created",
	"Deleted",
	"Updated",
	"Snoozed",
	"Unsnoozed",
}

func NonNullableArrayToStringArray(array *[]string) []string {
	if array == nil {
		return []string{}
	}
	return *array
}

func ListToNonNullableStringArray(list types.List) *[]string {
	if list.IsNull() {
		return &[]string{}
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

func mapNullableListWithEmpty(ctx context.Context, stringArray *[]string) types.List {
	if stringArray == nil {
		return types.ListNull(types.StringType)
	}

	if (len(*stringArray)) == 0 {
		listValue, _ := types.ListValueFrom(ctx, types.StringType, []string{})
		return listValue
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

func MapNullableList(ctx context.Context, stringArray *[]string) types.List {
	return mapNullableListWithEmpty(ctx, stringArray)
}

type badRequestResponse struct {
	Errors map[string][]string `json:"errors"`
}

type badRequestResultResponse struct {
	Succeeded bool                            `json:"succeeded"`
	Errors    []badRequestResultErrorResponse `json:"errors"`
}

type badRequestResultErrorResponse struct {
	Description string `json:"description"`
	Field       string `json:"field"`
}

func handleBadRequestResponse(data []byte) (error, error) {
	var badRequestResponse badRequestResponse
	reader := bytes.NewReader(data)
	errDecode := json.NewDecoder(reader).Decode(&badRequestResponse)

	if errDecode != nil {
		return errDecode, nil
	}

	err := fmt.Errorf("")
	for key, value := range badRequestResponse.Errors {
		err = fmt.Errorf("%s\n%s: %s", err, key, value)
	}

	return nil, err
}

func handleBadRequestResultResponse(data []byte) (error, error) {
	var badRequestResponse badRequestResultResponse

	reader := bytes.NewReader(data)
	errDecode := json.NewDecoder(reader).Decode(&badRequestResponse)

	if errDecode != nil {
		return errDecode, nil
	}

	err := fmt.Errorf("")
	for _, value := range badRequestResponse.Errors {
		if value.Field == "" {
			err = fmt.Errorf("%s\n%s", err, value.Description)
			continue
		}

		err = fmt.Errorf("%s\n%s: %s", err, value.Field, value.Description)
	}

	return nil, err
}

var logErrorRequest = false

func logErrorResponse(resp *http.Response, req interface{}) error {
	err := fmt.Errorf("%s %s: %d", resp.Request.Method, resp.Request.URL.RequestURI(), resp.StatusCode)

	if resp.StatusCode == 400 || resp.StatusCode == 401 || resp.StatusCode == 403 {
		data, _ := io.ReadAll(resp.Body)

		errDecode, errBadRequest := handleBadRequestResponse(data)
		if errDecode != nil {
			errDecode, errBadRequest = handleBadRequestResultResponse(data)

			if errDecode != nil {
				err = fmt.Errorf("%s\ncould not decode %d response: %s", err, resp.StatusCode, errDecode)
			}
		}

		if errBadRequest != nil {
			err = fmt.Errorf("%s\n%s", err, errBadRequest)
		}
	}

	if logErrorRequest && req != nil {
		b, _ := json.Marshal(req)
		err = fmt.Errorf("%s\nrequest: %s", err, string(b))
	}

	return err
}

var ValidTimes = []string{"00:00", "00:15", "00:30", "00:45", "01:00",
	"01:15", "01:30", "01:45", "02:00", "02:15", "02:30", "02:45", "03:00",
	"03:15", "03:30", "03:45", "04:00", "04:15", "04:30", "04:45", "05:00",
	"05:15", "05:30", "05:45", "06:00", "06:15", "06:30", "06:45", "07:00",
	"07:15", "07:30", "07:45", "08:00", "08:15", "08:30", "08:45", "09:00",
	"09:15", "09:30", "09:45", "10:00", "10:15", "10:30", "10:45", "11:00",
	"11:15", "11:30", "11:45", "12:00", "12:15", "12:30", "12:45", "13:00",
	"13:15", "13:30", "13:45", "14:00", "14:15", "14:30", "14:45", "15:00",
	"15:15", "15:30", "15:45", "16:00", "16:15", "16:30", "16:45", "17:00",
	"17:15", "17:30", "17:45", "18:00", "18:15", "18:30", "18:45", "19:00",
	"19:15", "19:30", "19:45", "20:00", "20:15", "20:30", "20:45", "21:00",
	"21:15", "21:30", "21:45", "22:00", "22:15", "22:30", "22:45", "23:00",
	"23:15", "23:30", "23:45"}

func DateTimeValidator(message string) validator.String {
	return validators.DateTime(message)
}

func TimeValidator(message string) validator.String {
	return stringvalidator.OneOf(ValidTimes...)
}

func IntentValidator(message string) validator.String {
	return stringvalidator.OneOf(ValidIntents...)
}

var ValidDaysOfWeek = []string{"sun", "mon", "tue", "wed", "thu", "fri", "sat"}

func DaysOfWeekValidator(message string) validator.String {
	return stringvalidator.OneOf(ValidDaysOfWeek...)
}

var ValidSeverities = []string{"Critical", "Warning", "Minor"}

func SeverityValidator(message string) validator.String {
	return stringvalidator.OneOf(ValidSeverities...)
}

var ValidStatuses = []string{"Open", "Resolved"}

func StatusValidator(message string) validator.String {
	return stringvalidator.OneOf(ValidStatuses...)
}

var ValidRuleFlowControl = []string{"Continue", "Skip"}

func RuleFlowValidator(message string) validator.String {
	return stringvalidator.OneOf(ValidRuleFlowControl...)
}

var ValidNotificationChannels = []string{"Email", "VoiceCall", "SMS", "Push"}

func NotificationChannelValidator(message string) validator.String {
	return stringvalidator.OneOf(ValidNotificationChannels...)
}

var ValidMaintenanceWindowTypes = []string{"maintenance", "muted"}

var ValidOperators = []string{"=", "!=", "contains", "!contains", ">", ">=", "<", "<="}

var ValidAttributesMatchTypes = []string{"all", "any"}

var ValidRotationModes = []string{"explicit", "auto"}

var ValidCustomRepeatUnits = []string{"months", "weeks", "days", "hours"}

var ValidRotationRepeats = []string{"daily", "weekly", "biweekly", "monthly", "custom"}

var ValidTeamConnectionModes = []string{"OrganizationTeams", "SelectedTeams"}

var ValidTeamMembershipRoles = []string{"Member", "Administrator"}

var ValidOrganizationMembershipRoles = []string{"Member", "Owner", "Administrator"}

var ValidEscalationModes = []string{"resolved", "acknowledged"}

func OperatorValidator(message string) validator.String {
	return stringvalidator.OneOf(ValidOperators...)
}

func GuidValidator(message string) validator.String {
	return stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`), message)
}

func RandomizeExample(example string) string {
	return strings.Replace(example, "@allquiet.app", "+"+uuid.New().String()+"@allquiet.app", -1)
}

func HexColorValidator(message string) validator.String {
	return stringvalidator.RegexMatches(regexp.MustCompile(`^#([0-9a-fA-F]{6})$`), message)
}

const OneMonthInSeconds = 2629746

var ValidWebhookAuthenticationTypes = []string{"bearer"}

func WebhookAuthenticationTypeValidator(message string) validator.String {
	return stringvalidator.OneOf(ValidWebhookAuthenticationTypes...)
}

var ValidIntervalsInSeconds = []int64{30 * 1, 60 * 1, 60 * 2, 60 * 5, 60 * 10, 60 * 15, 60 * 30, 60 * 60, 60 * 1440}
var ValidIntervalsInSecondsAsString = convertInt64ArrayToStringArray(ValidIntervalsInSeconds)

func convertInt64ArrayToStringArray(array []int64) []string {
	result := make([]string, len(array))
	for i, value := range array {
		result[i] = strconv.FormatInt(value, 10)
	}
	return result
}

func IntervalInSecondsValidator(message string) validator.Int64 {
	return int64validator.OneOf(ValidIntervalsInSeconds...)
}

var ValidHttpMonitoringAuthenticationTypes = []string{"Basic", "Bearer", "None"}

func HttpMonitoringAuthenticationTypeValidator(message string) validator.String {
	return stringvalidator.OneOf(ValidHttpMonitoringAuthenticationTypes...)
}

var ValidHttpMonitoringMethods = []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"}

func HttpMonitoringMethodValidator(message string) validator.String {
	return stringvalidator.OneOf(ValidHttpMonitoringMethods...)
}

var ValidTimeoutsHttpMonitoringInMilliseconds = []int64{50, 100, 200, 500, 1000, 2000, 5000, 10000, 30000, 60000}

func ValidTimeoutsHttpMonitoringInMillisecondsValidator(message string) validator.Int64 {
	return int64validator.OneOf(ValidTimeoutsHttpMonitoringInMilliseconds...)
}

var ValidTimeoutsPingMonitorInMilliseconds = []int64{50, 100, 200, 500, 1000, 2000, 5000}

func ValidTimeoutsPingMonitorInMillisecondsValidator(message string) validator.Int64 {
	return int64validator.OneOf(ValidTimeoutsPingMonitorInMilliseconds...)
}

func AddQueryParam(currentUrl string, key string, value string) string {

	if strings.Contains(currentUrl, "?") {
		return fmt.Sprintf("%s&%s=%s", currentUrl, key, url.QueryEscape(value))
	}
	return fmt.Sprintf("%s?%s=%s", currentUrl, key, url.QueryEscape(value))
}

func GetAccTestEnv() string {
	endpoint := os.Getenv("ALLQUIET_ENDPOINT")

	if endpoint == "" || strings.Contains(endpoint, "https://allquiet.app") {
		return "prod"
	}
	if strings.Contains(endpoint, "https://allquiet-test.app") {
		return "test"
	}
	return "local"
}

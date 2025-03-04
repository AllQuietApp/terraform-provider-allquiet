package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type teamEscalationsResponse struct {
	Id              string
	TeamId          string
	EscalationTiers []teamEscalationsTier
}

type teamEscalationsCreateRequest struct {
	TeamId          string                `json:"teamId"`
	EscalationTiers []teamEscalationsTier `json:"escalationTiers"`
}

type teamEscalationsTier struct {
	AutoEscalationEnabled        *bool                        `json:"autoEscalationEnabled"`
	AutoEscalationAfterMinutes   *int64                       `json:"autoEscalationAfterMinutes"`
	AutoEscalationSeverities     *[]string                    `json:"autoEscalationSeverities"`
	AutoEscalationTimeFilters    *[]teamEscalationsTimeFilter `json:"autoEscalationTimeFilters"`
	AutoAssignToTeams            *[]string                    `json:"autoAssignToTeams"`
	AutoAssignToTeamsSeverities  *[]string                    `json:"autoAssignToTeamsSeverities"`
	AutoAssignToTeamsTimeFilters *[]teamEscalationsTimeFilter `json:"autoAssignToTeamsTimeFilters"`
	Repeats                      *int64                       `json:"repeats"`
	RepeatsAfterMinutes          *int64                       `json:"repeatsAfterMinutes"`
	Schedules                    []teamEscalationsSchedule    `json:"schedules"`
}

type teamEscalationsTimeFilter struct {
	SelectedDays *[]string `json:"selectedDays"`
	From         *string   `json:"from"`
	Until        *string   `json:"until"`
}

type teamEscalationsSchedule struct {
	ScheduleSettings *scheduleSettings         `json:"scheduleSettings"`
	RotationSettings *rotationSettings         `json:"rotationSettings"`
	Rotations        []teamEscalationsRotation `json:"rotations"`
}

type teamEscalationsRotation struct {
	Members []teamEscalationsRotationMember `json:"members"`
}

type teamEscalationsRotationMember struct {
	TeamMembershipId string `json:"teamMembershipId"`
}

type scheduleSettings struct {
	Start           *string           `json:"start"`
	End             *string           `json:"end"`
	SelectedDays    *[]string         `json:"selectedDays"`
	WeeklySchedules *[]weeklySchedule `json:"weeklySchedules"`
}

type weeklySchedule struct {
	SelectedDays *[]string `json:"selectedDays"`
	From         *string   `json:"from"`
	Until        *string   `json:"until"`
}

type rotationSettings struct {
	Repeats             *string `json:"repeats"`
	StartsOnDayOfWeek   *string `json:"startsOnDayOfWeek"`
	StartsOnDateOfMonth *int64  `json:"startsOnDateOfMonth"`
	StartsOnTime        *string `json:"startsOnTime"`
	CustomRepeatUnit    *string `json:"customRepeatUnit"`
	CustomRepeatValue   *int64  `json:"customRepeatValue"`
	EffectiveFrom       *string `json:"effectiveFrom"`
	RotationMode        *string `json:"rotationMode"`
	AutoRotationSize    *int64  `json:"autoRotationSize"`
}

func mapTeamEscalationsCreateRequest(plan *TeamEscalationsModel) *teamEscalationsCreateRequest {
	tiers := make([]teamEscalationsTier, len(plan.EscalationTiers))
	for i, tier := range plan.EscalationTiers {
		mappedTier := mapTier(tier)
		tiers[i] = *mappedTier
	}

	return &teamEscalationsCreateRequest{
		TeamId:          plan.TeamId.ValueString(),
		EscalationTiers: tiers,
	}
}

func mapTier(tier TeamEscalationsTierModel) *teamEscalationsTier {

	schedules := make([]teamEscalationsSchedule, len(tier.Schedules))
	for i, schedule := range tier.Schedules {

		schedules[i] = teamEscalationsSchedule{}

		if schedule.ScheduleSettings != nil {
			selectedDays := ListToStringArray(schedule.ScheduleSettings.SelectedDays)

			schedules[i].ScheduleSettings = &scheduleSettings{
				Start:           schedule.ScheduleSettings.Start.ValueStringPointer(),
				End:             schedule.ScheduleSettings.End.ValueStringPointer(),
				SelectedDays:    selectedDays,
				WeeklySchedules: mapTeamEscalationsWeeklySchedulesToRequest(schedule.ScheduleSettings.WeeklySchedules),
			}
		}

		if schedule.RotationSettings != nil {
			schedules[i].RotationSettings = &rotationSettings{
				Repeats:             schedule.RotationSettings.Repeats.ValueStringPointer(),
				StartsOnDayOfWeek:   schedule.RotationSettings.StartsOnDayOfWeek.ValueStringPointer(),
				StartsOnDateOfMonth: schedule.RotationSettings.StartsOnDateOfMonth.ValueInt64Pointer(),
				StartsOnTime:        schedule.RotationSettings.StartsOnTime.ValueStringPointer(),
				CustomRepeatUnit:    schedule.RotationSettings.CustomRepeatUnit.ValueStringPointer(),
				CustomRepeatValue:   schedule.RotationSettings.CustomRepeatValue.ValueInt64Pointer(),
				EffectiveFrom:       schedule.RotationSettings.EffectiveFrom.ValueStringPointer(),
				RotationMode:        schedule.RotationSettings.RotationMode.ValueStringPointer(),
				AutoRotationSize:    schedule.RotationSettings.AutoRotationSize.ValueInt64Pointer(),
			}
		}

		rotations := make([]teamEscalationsRotation, len(schedule.Rotations))
		for j, rotation := range schedule.Rotations {
			members := make([]teamEscalationsRotationMember, len(rotation.Members))
			for k, member := range rotation.Members {
				members[k] = teamEscalationsRotationMember{
					TeamMembershipId: member.TeamMembershipId.ValueString(),
				}
			}
			rotations[j] = teamEscalationsRotation{
				Members: members,
			}
		}
		schedules[i].Rotations = rotations
	}

	return &teamEscalationsTier{
		AutoEscalationEnabled:        tier.AutoEscalationEnabled.ValueBoolPointer(),
		AutoEscalationAfterMinutes:   tier.AutoEscalationAfterMinutes.ValueInt64Pointer(),
		AutoEscalationSeverities:     ListToStringArray(tier.AutoEscalationSeverities),
		AutoEscalationTimeFilters:    mapTeamEscalationsTimeFiltersToRequest(tier.AutoEscalationTimeFilters),
		AutoAssignToTeams:            ListToStringArray(tier.AutoAssignToTeams),
		AutoAssignToTeamsSeverities:  ListToStringArray(tier.AutoAssignToTeamsSeverities),
		AutoAssignToTeamsTimeFilters: mapTeamEscalationsTimeFiltersToRequest(tier.AutoAssignToTeamsTimeFilters),
		Repeats:                      tier.Repeats.ValueInt64Pointer(),
		RepeatsAfterMinutes:          tier.RepeatsAfterMinutes.ValueInt64Pointer(),
		Schedules:                    schedules,
	}
}

func mapTeamEscalationsTimeFiltersToRequest(timeFilters *[]TeamEscalationsTimeFilterModel) *[]teamEscalationsTimeFilter {
	if timeFilters == nil {
		return nil
	}

	requestTimeFilters := make([]teamEscalationsTimeFilter, len(*timeFilters))
	for i, timeFilter := range *timeFilters {
		requestTimeFilters[i] = teamEscalationsTimeFilter{
			SelectedDays: ListToStringArray(timeFilter.SelectedDays),
			From:         timeFilter.From.ValueStringPointer(),
			Until:        timeFilter.Until.ValueStringPointer(),
		}
	}

	return &requestTimeFilters
}

func mapTeamEscalationsWeeklySchedulesToRequest(weeklySchedules *[]TeamEscalationsWeeklyScheduleModel) *[]weeklySchedule {
	if weeklySchedules == nil {
		return nil
	}

	requestWeeklySchedules := make([]weeklySchedule, len(*weeklySchedules))
	for i, ws := range *weeklySchedules {
		requestWeeklySchedules[i] = weeklySchedule{
			SelectedDays: ListToStringArray(ws.SelectedDays),
			From:         ws.From.ValueStringPointer(),
			Until:        ws.Until.ValueStringPointer(),
		}
	}

	return &requestWeeklySchedules
}

func (c *AllQuietAPIClient) CreateTeamEscalationsResource(ctx context.Context, data *TeamEscalationsModel) (*teamEscalationsResponse, error) {
	reqBody := mapTeamEscalationsCreateRequest(data)

	url := "/team-escalations"
	httpResp, err := c.post(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result teamEscalationsResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) DeleteTeamEscalationsResource(ctx context.Context, id string) error {
	url := fmt.Sprintf("/team-escalations/%s", url.PathEscape(id))
	httpResp, err := c.delete(ctx, url)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return logErrorResponse(httpResp, nil)
	}

	return nil
}

func (c *AllQuietAPIClient) UpdateTeamEscalationsResource(ctx context.Context, id string, data *TeamEscalationsModel) (*teamEscalationsResponse, error) {
	reqBody := mapTeamEscalationsCreateRequest(data)

	url := fmt.Sprintf("/team-escalations/%s", url.PathEscape(id))
	httpResp, err := c.put(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result teamEscalationsResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) GetTeamEscalationsResource(ctx context.Context, id string) (*teamEscalationsResponse, error) {
	url := fmt.Sprintf("/team-escalations/%s", url.PathEscape(id))
	httpResp, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, logErrorResponse(httpResp, nil)
	}

	var result teamEscalationsResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

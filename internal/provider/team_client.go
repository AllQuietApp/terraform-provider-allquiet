package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type teamResponse struct {
	Id                               string
	DisplayName                      string
	TimeZoneId                       string
	IncidentEngagementReportSettings *incidentEngagementReportSettings
	Members                          []teamMember
	Tiers                            []teamTier
}

type incidentEngagementReportSettings struct {
	DayOfWeek string
	Time      string
}

type teamCreateRequest struct {
	DisplayName                      string                            `json:"displayName"`
	TimeZoneId                       string                            `json:"timeZoneId"`
	Members                          []teamMember                      `json:"members"`
	IncidentEngagementReportSettings *incidentEngagementReportSettings `json:"incidentEngagementReportSettings"`
	Tiers                            []teamTier                        `json:"tiers"`
}

type teamMember struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type teamTier struct {
	AutoEscalationAfterMinutes *int64         `json:"autoEscalationAfterMinutes"`
	Schedules                  []teamSchedule `json:"schedules"`
}

type teamSchedule struct {
	ScheduleSettings *scheduleSettings `json:"scheduleSettings"`
	RotationSettings *rotationSettings `json:"rotationSettings"`
	Rotations        []teamRotation    `json:"rotations"`
}

type teamRotation struct {
	Members []rotationMember `json:"members"`
}

type rotationMember struct {
	Email string `json:"email"`
}

type scheduleSettings struct {
	Start        *string   `json:"start"`
	End          *string   `json:"end"`
	SelectedDays *[]string `json:"selectedDays"`
}

type rotationSettings struct {
	Repeats             *string `json:"repeats"`
	StartsOnDayOfWeek   *string `json:"startsOnDayOfWeek"`
	StartsOnDateOfMonth *int64  `json:"startsOnDateOfMonth"`
}

func mapTeamCreateRequest(plan *TeamModel) *teamCreateRequest {
	members := make([]teamMember, len(plan.Members))
	for i, member := range plan.Members {
		members[i] = teamMember{
			Email: member.Email.ValueString(),
			Role:  member.Role.ValueString(),
		}
	}

	var settings *incidentEngagementReportSettings

	if plan.IncidentEngagementReportSettings != nil {
		settings = &incidentEngagementReportSettings{
			DayOfWeek: *plan.IncidentEngagementReportSettings.DayOfWeek.ValueStringPointer(),
			Time:      plan.IncidentEngagementReportSettings.Time.ValueString(),
		}
	}

	tiers := make([]teamTier, len(plan.Tiers))
	for i, tier := range plan.Tiers {
		mappedTier := mapTier(tier)
		tiers[i] = *mappedTier
	}

	return &teamCreateRequest{
		DisplayName:                      plan.DisplayName.ValueString(),
		TimeZoneId:                       plan.TimeZoneId.ValueString(),
		Members:                          members,
		IncidentEngagementReportSettings: settings,
		Tiers:                            tiers,
	}
}

func mapTier(tier TeamTierModel) *teamTier {

	schedules := make([]teamSchedule, len(tier.Schedules))
	for i, schedule := range tier.Schedules {

		schedules[i] = teamSchedule{}

		if schedule.ScheduleSettings != nil {
			selectedDays := ListToStringArray(schedule.ScheduleSettings.SelectedDays)

			schedules[i].ScheduleSettings = &scheduleSettings{
				Start:        schedule.ScheduleSettings.Start.ValueStringPointer(),
				End:          schedule.ScheduleSettings.End.ValueStringPointer(),
				SelectedDays: selectedDays,
			}
		}

		if schedule.RotationSettings != nil {
			schedules[i].RotationSettings = &rotationSettings{
				Repeats:             schedule.RotationSettings.Repeats.ValueStringPointer(),
				StartsOnDayOfWeek:   schedule.RotationSettings.StartsOnDayOfWeek.ValueStringPointer(),
				StartsOnDateOfMonth: schedule.RotationSettings.StartsOnDateOfMonth.ValueInt64Pointer(),
			}
		}

		rotations := make([]teamRotation, len(schedule.Rotations))
		for j, rotation := range schedule.Rotations {
			members := make([]rotationMember, len(rotation.Members))
			for k, member := range rotation.Members {
				members[k] = rotationMember{
					Email: member.Email.ValueString(),
				}
			}
			rotations[j] = teamRotation{
				Members: members,
			}
		}
		schedules[i].Rotations = rotations
	}

	return &teamTier{
		AutoEscalationAfterMinutes: tier.AutoEscalationAfterMinutes.ValueInt64Pointer(),
		Schedules:                  schedules,
	}
}

func (c *AllQuietAPIClient) CreateTeamResource(ctx context.Context, data *TeamModel) (*teamResponse, error) {
	reqBody := mapTeamCreateRequest(data)

	url := "/team"
	httpResp, err := c.post(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		logErrorResponse(httpResp)
		return nil, fmt.Errorf("non-200 response from API for POST %s: %d", url, httpResp.StatusCode)
	}

	var result teamResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) DeleteTeamResource(ctx context.Context, id string) error {
	url := fmt.Sprintf("/team/%s", url.PathEscape(id))
	httpResp, err := c.delete(ctx, url)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		logErrorResponse(httpResp)
		return fmt.Errorf("non-200 response from API for DELETE %s: %d", url, httpResp.StatusCode)
	}

	return nil
}

func (c *AllQuietAPIClient) UpdateTeamResource(ctx context.Context, id string, data *TeamModel) (*teamResponse, error) {
	reqBody := mapTeamCreateRequest(data)

	url := fmt.Sprintf("/team/%s", url.PathEscape(id))
	httpResp, err := c.put(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		logErrorResponse(httpResp)
		return nil, fmt.Errorf("non-200 response from API for PUT %s: %d", url, httpResp.StatusCode)
	}

	var result teamResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) GetTeamResource(ctx context.Context, id string) (*teamResponse, error) {
	url := fmt.Sprintf("/team/%s", url.PathEscape(id))
	httpResp, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		logErrorResponse(httpResp)
		return nil, fmt.Errorf("non-200 response from API for GET %s: %d", url, httpResp.StatusCode)
	}

	var result teamResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

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
	Labels                           *[]string
}

type incidentEngagementReportSettings struct {
	DayOfWeek string
	Time      string
}

type teamCreateRequest struct {
	DisplayName                      string                            `json:"displayName"`
	TimeZoneId                       string                            `json:"timeZoneId"`
	IncidentEngagementReportSettings *incidentEngagementReportSettings `json:"incidentEngagementReportSettings"`
	Labels                           *[]string                         `json:"labels"`
}

func mapTeamCreateRequest(plan *TeamModel) *teamCreateRequest {
	var settings *incidentEngagementReportSettings

	if plan.IncidentEngagementReportSettings != nil {
		settings = &incidentEngagementReportSettings{
			DayOfWeek: *plan.IncidentEngagementReportSettings.DayOfWeek.ValueStringPointer(),
			Time:      plan.IncidentEngagementReportSettings.Time.ValueString(),
		}
	}

	return &teamCreateRequest{
		DisplayName:                      plan.DisplayName.ValueString(),
		TimeZoneId:                       plan.TimeZoneId.ValueString(),
		IncidentEngagementReportSettings: settings,
		Labels:                           ListToStringArray(plan.Labels),
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
		return nil, logErrorResponse(httpResp, nil)
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
		return logErrorResponse(httpResp, nil)
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
		return nil, logErrorResponse(httpResp, nil)
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
		return nil, logErrorResponse(httpResp, nil)
	}

	var result teamResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

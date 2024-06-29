package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type routingResponse struct {
	Id          string        `json:"id"`
	DisplayName string        `json:"displayName"`
	TeamId      string        `json:"teamId"`
	Rules       []routingRule `json:"rules"`
}

type routingRule struct {
	Conditions *routingRuleConditions `json:"conditions"`
	Actions    *routingRuleActions    `json:"actions"`
	Channels   *routingRuleChannels   `json:"channels"`
}

type routingRuleConditions struct {
	Statuses     *[]string              `json:"statuses"`
	Severities   *[]string              `json:"severities"`
	Integrations *[]string              `json:"integrations"`
	Intents      *[]string              `json:"intents"`
	Attributes   []routingRuleAttribute `json:"attributes"`
}

type routingRuleAttribute struct {
	Name     string  `json:"name"`
	Operator string  `json:"operator"`
	Value    *string `json:"value"`
}

type routingRuleActions struct {
	RouteToTeams    *[]string `json:"routeToTeams"`
	ChangeSeverity  *string   `json:"changeSeverity"`
	AddInteraction  *string   `json:"addInteraction"`
	RuleFlowControl *string   `json:"ruleFlowControl"`
	Discard         bool      `json:"discard"`
}

type routingRuleChannels struct {
	OutboundIntegrations      *[]string `json:"outboundIntegrations"`
	OutboundIntegrationsMuted bool      `json:"outboundIntegrationsMuted"`
	NotificationChannels      *[]string `json:"notificationChannels"`
	NotificationChannelsMuted bool      `json:"notificationChannelsMuted"`
}

type routingCreateRequest struct {
	DisplayName string        `json:"displayName"`
	TeamId      string        `json:"teamId"`
	Rules       []routingRule `json:"rules"`
}

func mapRoutingCreateRequest(plan *RoutingModel) *routingCreateRequest {
	return &routingCreateRequest{
		DisplayName: plan.DisplayName.ValueString(),
		TeamId:      plan.TeamId.ValueString(),
		Rules:       mapRoutingRules(plan.Rules),
	}
}

func mapRoutingRules(rules []RoutingRuleModel) []routingRule {
	result := make([]routingRule, len(rules))
	for i, rule := range rules {
		result[i] = routingRule{
			Conditions: mapRoutingRuleConditions(rule.Conditions),
			Actions:    mapRoutingRuleActions(rule.Actions),
			Channels:   mapRoutingRuleChannels(rule.Channels),
		}
	}
	return result
}

func mapRoutingRuleConditions(conditions *RoutingRuleConditionsModel) *routingRuleConditions {
	if conditions == nil {
		return nil
	}

	return &routingRuleConditions{
		Statuses:     ListToStringArray(conditions.Statuses),
		Severities:   ListToStringArray(conditions.Severities),
		Integrations: ListToStringArray(conditions.Integrations),
		Intents:      ListToStringArray(conditions.Intents),
		Attributes:   mapRoutingRuleAttributes(conditions.Attributes),
	}
}

func mapRoutingRuleAttributes(attributes []RoutingRuleConditionsAttributeModel) []routingRuleAttribute {
	result := make([]routingRuleAttribute, len(attributes))
	for i, attribute := range attributes {
		result[i] = routingRuleAttribute{
			Name:     attribute.Name.ValueString(),
			Operator: attribute.Operator.ValueString(),
			Value:    attribute.Value.ValueStringPointer(),
		}
	}
	return result
}

func mapRoutingRuleActions(actions *RoutingRuleActionsModel) *routingRuleActions {
	if actions == nil {
		return nil
	}

	return &routingRuleActions{
		RouteToTeams:    ListToStringArray(actions.RouteToTeams),
		ChangeSeverity:  actions.ChangeSeverity.ValueStringPointer(),
		AddInteraction:  actions.AddInteraction.ValueStringPointer(),
		RuleFlowControl: actions.RuleFlowControl.ValueStringPointer(),
		Discard:         actions.Discard.ValueBool(),
	}
}

func mapRoutingRuleChannels(channels *RoutingRuleChannelsModel) *routingRuleChannels {
	if channels == nil {
		return nil
	}

	return &routingRuleChannels{
		OutboundIntegrations:      ListToStringArray(channels.OutboundIntegrations),
		OutboundIntegrationsMuted: channels.OutboundIntegrationsMuted.ValueBool(),
		NotificationChannels:      ListToStringArray(channels.NotificationChannels),
		NotificationChannelsMuted: channels.NotificationChannelsMuted.ValueBool(),
	}
}

func (c *AllQuietAPIClient) CreateRoutingResource(ctx context.Context, data *RoutingModel) (*routingResponse, error) {
	reqBody := mapRoutingCreateRequest(data)

	url := "/routing"
	httpResp, err := c.post(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		logErrorResponse(httpResp)
		return nil, fmt.Errorf("non-200 response from API for POST %s: %d", url, httpResp.StatusCode)
	}

	var result routingResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) DeleteRoutingResource(ctx context.Context, id string) error {
	url := fmt.Sprintf("/routing/%s", url.PathEscape(id))
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

func (c *AllQuietAPIClient) UpdateRoutingResource(ctx context.Context, id string, data *RoutingModel) (*routingResponse, error) {
	reqBody := mapRoutingCreateRequest(data)

	url := fmt.Sprintf("/routing/%s", url.PathEscape(id))
	httpResp, err := c.put(ctx, url, reqBody)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		logErrorResponse(httpResp)
		return nil, fmt.Errorf("non-200 response from API for PUT %s: %d", url, httpResp.StatusCode)
	}

	var result routingResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *AllQuietAPIClient) GetRoutingResource(ctx context.Context, id string) (*routingResponse, error) {
	url := fmt.Sprintf("/routing/%s", url.PathEscape(id))
	httpResp, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		logErrorResponse(httpResp)
		return nil, fmt.Errorf("non-200 response from API for GET %s: %d", url, httpResp.StatusCode)
	}

	var result routingResponse
	err = json.NewDecoder(httpResp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

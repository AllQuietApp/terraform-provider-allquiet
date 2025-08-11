// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TeamConnectionSettings is a generic model for team connection settings
// that can be reused across different resources like outbound integrations and routing.
type TeamConnectionSettings struct {
	TeamConnectionMode types.String `tfsdk:"team_connection_mode"`
	TeamIds            types.List   `tfsdk:"team_ids"`
}

// teamConnectionSettings is the generic client-side response type for team connection settings
// that can be reused across different resources like outbound integrations and routing.
type teamConnectionSettings struct {
	TeamConnectionMode string    `json:"teamConnectionMode"`
	TeamIds            *[]string `json:"teamIds"`
}

// MapTeamConnectionSettingsToRequest maps the generic TeamConnectionSettings model
// to the generic teamConnectionSettings request type.
func MapTeamConnectionSettingsToRequest(settings *TeamConnectionSettings) *teamConnectionSettings {
	if settings == nil {
		return nil
	}

	return &teamConnectionSettings{
		TeamConnectionMode: settings.TeamConnectionMode.ValueString(),
		TeamIds:            ListToStringArray(settings.TeamIds),
	}
}

func MapTeamConnectionSettingsResponseToModel(ctx context.Context, settings *teamConnectionSettings) *TeamConnectionSettings {
	if settings != nil {
		return &TeamConnectionSettings{
			TeamConnectionMode: types.StringValue(settings.TeamConnectionMode),
			TeamIds:            MapNullableList(ctx, settings.TeamIds),
		}
	}

	return nil
}

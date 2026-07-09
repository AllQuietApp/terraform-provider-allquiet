package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMapStatusPageServiceGroupsRequestToModel_EmptyPublicDisplayName(t *testing.T) {
	serviceIds, _ := types.ListValueFrom(nil, types.StringType, []string{
		"11111111-1111-1111-1111-111111111111",
	})
	plan := []StatusPageServiceGroupModel{
		{
			PublicDisplayName: types.StringValue(""),
			Services:          serviceIds,
		},
	}

	result := mapStatusPageServiceGroupsRequestToModel(&plan)

	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(*result) != 1 {
		t.Fatalf("len(*result) = %d, want 1", len(*result))
	}
	if got, want := (*result)[0].PublicDisplayName, ""; got != want {
		t.Fatalf("PublicDisplayName = %q, want %q", got, want)
	}
}

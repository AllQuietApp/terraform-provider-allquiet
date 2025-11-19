// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOnCallOverridesDataSource(t *testing.T) {
	uid := uuid.New().String()
	email := fmt.Sprintf("acceptance-tests+millie+%s@allquiet.app", uid)
	uid2 := uuid.New().String()
	email2 := fmt.Sprintf("acceptance-tests+millie+%s@allquiet.app", uid2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOnCallOverridesDataSourceConfig("Millie Bobby Brown", email, "Miley Cyrus", email2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.allquiet_on_call_overrides.test1", "on_call_overrides.#", "2"),
					resource.TestCheckResourceAttr("data.allquiet_on_call_overrides.test2", "on_call_overrides.#", "1"),
				),
			},
		},
	})
}

func testAccOnCallOverridesDataSourceConfig(displayName, email, displayName2, email2 string) string {
	return fmt.Sprintf(`

		resource "allquiet_user" "user1" {
			display_name = %[1]q
			email        = %[2]q
		}

		resource "allquiet_user" "user2" {
			display_name = %[3]q
			email        = %[4]q
		}

		resource "allquiet_on_call_override" "user1_override1" {
			user_id = allquiet_user.user1.id
			type = "online"
			start = "2025-09-11T00:00:00Z"
			end = "2025-09-11T00:20:00Z"
		}

		resource "allquiet_on_call_override" "user1_override2" {
			user_id = allquiet_user.user1.id
			type = "offline"
			start = "2025-10-01T00:00:00Z"
			end = "2025-11-01T00:00:00Z"
		}

		resource "allquiet_on_call_override" "user2_override1" {
			user_id = allquiet_user.user2.id
			type = "online"
			start = "2025-11-01T00:00:00Z"
			end = "2025-12-01T00:00:00Z"
			replacement_user_ids = [allquiet_user.user1.id]
		}

		data "allquiet_on_call_overrides" "test1" {
			user_id = allquiet_user.user1.id
			depends_on = [allquiet_on_call_override.user1_override1, allquiet_on_call_override.user1_override2, allquiet_user.user1]
		}

		data "allquiet_on_call_overrides" "test2" {
			user_id = allquiet_user.user2.id
			depends_on = [allquiet_on_call_override.user2_override1, allquiet_user.user1, allquiet_user.user2]
		}
	`, displayName, email, displayName2, email2)
}

func TestAccOnCallOverridesDataSourceExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccOnCallOverridesDataSourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.allquiet_on_call_overrides.example1", "on_call_overrides.#", "2"),
				),
			},
		},
	})
}

func testAccOnCallOverridesDataSourceExample() string {
	absPath, _ := filepath.Abs("../../examples/data-sources/allquiet_on_call_overrides/data-source.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return RandomizeExample(string(dat))
}

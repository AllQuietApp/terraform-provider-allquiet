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

func TestAccOnCallOverrideResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccOnCallOverrideResourceConfig("online"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("allquiet_on_call_override.test", "user_id"),
					resource.TestCheckResourceAttr("allquiet_on_call_override.test", "type", "online"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_on_call_override.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccOnCallOverrideResourceConfig("offline"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("allquiet_on_call_override.test", "user_id"),
					resource.TestCheckResourceAttr("allquiet_on_call_override.test", "type", "offline"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccOnCallOverrideResourceExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccOnCallOverrideResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("allquiet_on_call_override.millie_brown_override1", "user_id"),
					resource.TestCheckResourceAttrSet("allquiet_on_call_override.millie_brown_override2", "user_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_on_call_override.millie_brown_override1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// ImportState testing
			{
				ResourceName:      "allquiet_on_call_override.millie_brown_override2",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccOnCallOverrideResourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("allquiet_on_call_override.millie_brown_override1", "user_id"),
					resource.TestCheckResourceAttrSet("allquiet_on_call_override.millie_brown_override2", "user_id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccOnCallOverrideResourceConfig(override_type string) string {
	return fmt.Sprintf(`

resource "allquiet_user" "millie_brown" {
  display_name =  "Millie Bobby Brown"
  email = "acceptance-tests+millie+%s@allquiet.app"
}

resource "allquiet_user" "taylor_swift" {
  display_name =  "Taylor Swift"
  email = "acceptance-tests+taylor+%s@allquiet.app"
}

resource "allquiet_on_call_override" "test" {
	user_id = allquiet_user.millie_brown.id
	type = "%s"
	start = "2025-09-11T00:00:00Z"
	end = "2025-09-11T00:20:00Z"
	replacement_user_ids = [allquiet_user.taylor_swift.id]
}

`, uuid.New().String(), uuid.New().String(), override_type)

}

func testAccOnCallOverrideResourceExample() string {
	absPath, _ := filepath.Abs("../../examples/resources/allquiet_on_call_override/resource.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return RandomizeExample(string(dat))
}

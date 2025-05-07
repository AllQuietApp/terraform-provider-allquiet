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

func TestAccUserDataSource(t *testing.T) {
	uid := uuid.New().String()
	email := fmt.Sprintf("acceptance-tests+millie+%s@allquiet.app", uid)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig("Millie Bobby Brown", email),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.allquiet_user.test_by_email", "email", email),
					resource.TestCheckResourceAttr("data.allquiet_user.test_by_display_name", "display_name", "Millie Bobby Brown"),
				),
			},
		},
	})
}

func TestAccUserDataSourceExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserDataSourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.allquiet_user.test_by_email", "display_name", "Millie Bobby Brown"),
					resource.TestCheckResourceAttr("data.allquiet_user.test_by_display_name", "display_name", "Millie Bobby Brown"),
				),
			},
		},
	})
}

func testAccUserDataSourceConfig(displayName, email string) string {
	return fmt.Sprintf(`

		resource "allquiet_user" "test" {
			display_name = %[1]q
			email        = %[2]q
			phone_number = "+12035479055"
		}

		data "allquiet_user" "test_by_email" {
			email = %[2]q
			depends_on = [allquiet_user.test]
		}

		data "allquiet_user" "test_by_display_name" {
			display_name = %[1]q
			depends_on = [allquiet_user.test]
		}

	`, displayName, email)
}

func testAccUserDataSourceExample() string {
	absPath, _ := filepath.Abs("../../examples/data-sources/allquiet_user/data-source.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	return RandomizeExample(string(dat))
}

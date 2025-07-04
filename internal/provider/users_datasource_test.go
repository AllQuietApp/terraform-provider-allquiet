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

func TestAccUsersDataSource(t *testing.T) {
	uid := uuid.New().String()
	emailPrefix := fmt.Sprintf("acceptance-tests+millie+ds+%s", uid)
	displayName := fmt.Sprintf("Millie Bobby Brown %s", uid)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUsersDataSourceConfig(displayName, emailPrefix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.allquiet_users.test_by_email", "users.#", "3"),
				),
			},
		},
	})
}

func testAccUsersDataSourceConfig(displayName, emailPrefix string) string {
	return fmt.Sprintf(`

		resource "allquiet_user" "test1" {
			display_name = "%[1]s 1"
			email        = "%[2]s+1@allquiet.app"
			phone_number = "+12035479051"
		}

		resource "allquiet_user" "test2" {
			display_name = "%[1]s 2"
			email        = "%[2]s+2@allquiet.app"
			phone_number = "+12035479052"
		}

		resource "allquiet_user" "test3" {
			display_name = "%[1]s 3"
			email        = "%[2]s+3@allquiet.app"
			phone_number = "+12035479053"
		}

		data "allquiet_users" "test_by_email" {
			email        = "%[2]s"
			depends_on = [allquiet_user.test1, allquiet_user.test2, allquiet_user.test3]
		}
	`, displayName, emailPrefix)
}

func TestAccUsersDataSourceExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUsersDataSourceExample(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.allquiet_users.users_by_email", "users.#", "2"),
					resource.TestCheckResourceAttrSet("data.allquiet_users.users_by_email", "users.#"),
				),
			},
		},
	})
}
func testAccUsersDataSourceExample() string {
	absPath, _ := filepath.Abs("../../examples/data-sources/allquiet_users/data-source.tf")

	dat, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}

	example := string(dat)
	return RandomizeExample(example)
}

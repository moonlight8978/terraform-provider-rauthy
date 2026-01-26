// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package passwordpolicy_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/moonlight8978/terraform-provider-rauthy/internal/provider/acctest"
)

func TestAccPasswordPolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPasswordPolicyResourceConfig(10, 100, 3),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"rauthy_password_policy.default",
						tfjsonpath.New("length_min"),
						knownvalue.Int64Exact(10),
					),
					statecheck.ExpectKnownValue(
						"rauthy_password_policy.default",
						tfjsonpath.New("length_max"),
						knownvalue.Int64Exact(100),
					),
					statecheck.ExpectKnownValue(
						"rauthy_password_policy.default",
						tfjsonpath.New("include_digits"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"rauthy_password_policy.default",
						tfjsonpath.New("include_lower_case"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"rauthy_password_policy.default",
						tfjsonpath.New("include_upper_case"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"rauthy_password_policy.default",
						tfjsonpath.New("include_special"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"rauthy_password_policy.default",
						tfjsonpath.New("not_recently_used"),
						knownvalue.Int64Exact(0),
					),
					statecheck.ExpectKnownValue(
						"rauthy_password_policy.default",
						tfjsonpath.New("valid_days"),
						knownvalue.Int64Exact(3),
					),
				},
			},
			{
				Config: testAccPasswordPolicyResourceConfig(20, 200, 0),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"rauthy_password_policy.default",
						tfjsonpath.New("length_min"),
						knownvalue.Int64Exact(20),
					),
					statecheck.ExpectKnownValue(
						"rauthy_password_policy.default",
						tfjsonpath.New("length_max"),
						knownvalue.Int64Exact(200),
					),
					statecheck.ExpectKnownValue(
						"rauthy_password_policy.default",
						tfjsonpath.New("include_digits"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"rauthy_password_policy.default",
						tfjsonpath.New("include_lower_case"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"rauthy_password_policy.default",
						tfjsonpath.New("include_upper_case"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"rauthy_password_policy.default",
						tfjsonpath.New("include_special"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"rauthy_password_policy.default",
						tfjsonpath.New("not_recently_used"),
						knownvalue.Int64Exact(0),
					),
					statecheck.ExpectKnownValue(
						"rauthy_password_policy.default",
						tfjsonpath.New("valid_days"),
						knownvalue.Int64Exact(0),
					),
				},
			},
		},
	})
}

func testAccPasswordPolicyResourceConfig(lengthMin int64, lengthMax int64, validDays int64) string {
	return fmt.Sprintf(`
resource "rauthy_password_policy" "default" {
  length_min = %[1]d
  length_max = %[2]d
  include_digits = 1
  include_lower_case = 1
  include_upper_case = 1
  include_special = 1
  not_recently_used = 0
  valid_days = %[3]d
}
`, lengthMin, lengthMax, validDays)
}

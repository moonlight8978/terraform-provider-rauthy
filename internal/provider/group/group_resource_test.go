package group_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/moonlight8978/terraform-provider-rauthy/internal/provider/acctest"
)

func TestAccGroupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupResourceConfig("test-group"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"rauthy_group.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("test-group"),
					),
					statecheck.ExpectKnownValue(
						"rauthy_group.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
				},
			},
			{
				ResourceName:      "rauthy_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGroupResourceConfig("test-group-updated"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"rauthy_group.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("test-group-updated"),
					),
				},
			},
		},
	})
}

func testAccGroupResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "rauthy_group" "test" {
	name = %[1]q
}
`, name)
}

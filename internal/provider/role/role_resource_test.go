package role_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/moonlight8978/terraform-provider-rauthy/internal/provider/acctest"
)

func TestAccRoleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleResourceConfig("test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rauthy_role.test", "name", "test"),
					resource.TestCheckResourceAttrSet("rauthy_role.test", "id"),
				),
			},
			{
				ResourceName:      "rauthy_role.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRoleResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "rauthy_role" "test" {
	name = %[1]q
}
`, name)
}

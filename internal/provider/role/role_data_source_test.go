package role_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/moonlight8978/terraform-provider-rauthy/internal/provider/acctest"
)

func TestAccRoleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleDataSourceConfig("test-role-ds"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.rauthy_role.test", "name", "test-role-ds"),
					resource.TestCheckResourceAttrSet("data.rauthy_role.test", "id"),
				),
			},
		},
	})
}

func testAccRoleDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "rauthy_role" "test" {
	name = %[1]q
}

data "rauthy_role" "test" {
	name = rauthy_role.test.name
}
`, name)
}

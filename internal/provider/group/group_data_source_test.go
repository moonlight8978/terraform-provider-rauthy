package group_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/moonlight8978/terraform-provider-rauthy/internal/provider/acctest"
)

func TestAccGroupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupDataSourceConfig("test-group-ds"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.rauthy_group.test", "name", "test-group-ds"),
					resource.TestCheckResourceAttrSet("data.rauthy_group.test", "id"),
				),
			},
		},
	})
}

func testAccGroupDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "rauthy_group" "test" {
	name = %[1]q
}

data "rauthy_group" "test" {
	name = rauthy_group.test.name
}
`, name)
}

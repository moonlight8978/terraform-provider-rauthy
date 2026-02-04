package oidc_client_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/moonlight8978/terraform-provider-rauthy/internal/provider/acctest"
)

func TestAccOidcClientDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOidcClientDataSourceConfig("test-client-ds"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.rauthy_client.test", "name", "rauthy_client.test", "name"),
					resource.TestCheckResourceAttrSet("data.rauthy_client.test", "id"),
				),
			},
		},
	})
}

func testAccOidcClientDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "rauthy_client" "test" {
	id = %[1]q
	name = %[1]q
	redirect_uris = ["http://localhost:8080"]
}

data "rauthy_client" "test" {
	id = rauthy_client.test.id
}
`, name)
}

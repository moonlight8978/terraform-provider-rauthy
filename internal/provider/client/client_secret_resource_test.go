package client_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/moonlight8978/terraform-provider-rauthy/internal/provider/acctest"
)

func TestAccClientSecretResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClientSecretResourceConfig("testsecret"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"rauthy_client_secret.test",
						tfjsonpath.New("client_id"),
						knownvalue.StringExact("testsecret"),
					),
					statecheck.ExpectKnownValue(
						"rauthy_client_secret.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("one"),
					),
					statecheck.ExpectKnownValue(
						"rauthy_client_secret.test",
						tfjsonpath.New("secret"),
						knownvalue.NotNull(),
					),
				},
			},
			{
				ResourceName:      "rauthy_client_secret.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccClientSecretResourceConfig(clientId string) string {
	return fmt.Sprintf(`
resource "rauthy_client" "test" {
	id = %[1]q
	name = %[1]q
	confidential = true
	redirect_uris = ["http://localhost/callback"]
}

resource "rauthy_client_secret" "test" {
	id = "one"
	client_id = rauthy_client.test.id
}
`, clientId)
}

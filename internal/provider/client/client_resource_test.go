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

func TestAccClientResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClientResourceConfig("google", "Google", true, []string{"http://localhost/callback"}),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"rauthy_client.google",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Google"),
					),
					statecheck.ExpectKnownValue(
						"rauthy_client.google",
						tfjsonpath.New("confidential"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"rauthy_client.google",
						tfjsonpath.New("id"),
						knownvalue.StringExact("google"),
					),
					statecheck.ExpectKnownValue(
						"rauthy_client.google",
						tfjsonpath.New("redirect_uris"),
						knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("http://localhost/callback")}),
					),
				},
			},
			{
				ResourceName:      "rauthy_client.google",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccClientResourceConfig(id string, name string, confidential bool, redirectURIs []string) string {
	return fmt.Sprintf(`
resource "rauthy_client" "google" {
	id = %[1]q
	name = %[2]q
	confidential = %[3]t
	redirect_uris = %[4]q
	enabled = true
	flows_enabled = ["authorization_code"]
	access_token_alg = "EdDSA"
	id_token_alg = "EdDSA"
	auth_code_lifetime = 10
	access_token_lifetime = 10
	scopes = ["openid"]
	challenges = ["S256"]
	force_mfa = false
}
`, id, name, confidential, redirectURIs)
}

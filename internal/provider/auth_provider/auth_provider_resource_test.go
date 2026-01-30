package auth_provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/moonlight8978/terraform-provider-rauthy/internal/provider/acctest"
)

func TestAccAuthProviderResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthProviderResourceConfig("google", "Google"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"rauthy_auth_provider.google",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Google"),
					),
					statecheck.ExpectKnownValue(
						"rauthy_auth_provider.google",
						tfjsonpath.New("id"),
						knownvalue.StringExact("google"),
					),
					statecheck.ExpectKnownValue(
						"rauthy_auth_provider.google",
						tfjsonpath.New("issuer"),
						knownvalue.StringExact("https://accounts.google.com"),
					),
				},
			},
			{
				ResourceName:      "rauthy_auth_provider.google",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAuthProviderResourceConfig("google", "Google 2"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"rauthy_auth_provider.google",
						tfjsonpath.New("name"),
						knownvalue.StringExact("Google 2"),
					),
				},
			},
		},
	})
}

func testAccAuthProviderResourceConfig(id string, name string) string {
	return fmt.Sprintf(`
resource "rauthy_auth_provider" "google" {
	id = %[1]q
	name = %[2]q
	typ = "google"
	issuer = "https://accounts.google.com"
	client_id = "google-client-id"
	client_secret = "google-client-secret"
	authorization_endpoint = "https://accounts.google.com/o/oauth2/v2/auth"
	token_endpoint = "https://oauth2.googleapis.com/token"
	userinfo_endpoint = "https://openidconnect.googleapis.com/v1/userinfo"
	jwks_endpoint = "https://www.googleapis.com/oauth2/v3/certs"
	scope = "openid profile email"
	enabled = true
}
`, id, name)
}

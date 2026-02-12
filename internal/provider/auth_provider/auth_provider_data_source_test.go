package auth_provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/moonlight8978/terraform-provider-rauthy/internal/provider/acctest"
)

func TestAccAuthProviderDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthProviderDataSourceConfig("google-ds", "Google DS"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.rauthy_auth_provider.test", "name", "rauthy_auth_provider.test", "name"),
					resource.TestCheckResourceAttrSet("data.rauthy_auth_provider.test", "id"),
					resource.TestCheckResourceAttr("data.rauthy_auth_provider.test", "issuer", "https://accounts.google.com"),
					resource.TestCheckResourceAttr("data.rauthy_auth_provider.test", "client_id", "google-client-id"),
				),
			},
		},
	})
}

func testAccAuthProviderDataSourceConfig(id string, name string) string {
	return fmt.Sprintf(`
resource "rauthy_auth_provider" "test" {
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

data "rauthy_auth_provider" "test" {
	id = rauthy_auth_provider.test.id
}
`, id, name)
}

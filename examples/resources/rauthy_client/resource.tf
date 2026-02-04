resource "rauthy_client" "google" {
  id                        = "google"
  name                      = "Google"
  confidential              = true
  redirect_uris             = ["http://localhost/callback"]
  post_logout_redirect_uris = ["http://localhost/logout"]
  enabled                   = true
  flows_enabled             = ["authorization_code"]
  access_token_alg          = "EdDSA"
  id_token_alg              = "EdDSA"
  auth_code_lifetime        = 60
  access_token_lifetime     = 1800
  scopes                    = ["openid"]
  challenges                = ["S256"]
  force_mfa                 = false
}

resource "rauthy_auth_provider" "default" {
  name                   = "Google"
  typ                    = "google"
  issuer                 = "https://accounts.google.com"
  client_id              = "google-client-id"
  client_secret          = "google-client-secret"
  authorization_endpoint = "https://accounts.google.com/o/oauth2/v2/auth"
  token_endpoint         = "https://oauth2.googleapis.com/token"
  userinfo_endpoint      = "https://openidconnect.googleapis.com/v1/userinfo"
  scope                  = "openid profile email"
  use_pkce               = true

  auto_link       = true
  auto_onboarding = false

  enabled = true
}

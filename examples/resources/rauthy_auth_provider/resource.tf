resource "rauthy_provider" "default" {
  id                     = "google"
  name                   = "Google"
  typ                    = "google"
  issuer                 = "https://accounts.google.com"
  client_id              = "google-client-id"
  client_secret          = "google-client-secret"
  authorization_endpoint = "https://accounts.google.com/o/oauth2/v2/auth"
  token_endpoint         = "https://oauth2.googleapis.com/token"
  userinfo_endpoint      = "https://openidconnect.googleapis.com/v1/userinfo"
  jwks_endpoint          = "https://www.googleapis.com/oauth2/v3/certs"
  scope                  = "openid profile email"
  enabled                = true
}

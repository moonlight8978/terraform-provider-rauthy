resource "rauthy_auth_provider" "github" {
  enabled = true

  auto_onboarding = true
  auto_link       = true

  name                = "GitHub"
  typ                 = "github"
  issuer              = "github.com"
  client_id           = "github-client-id"
  client_secret       = "github-client-secret"
  client_secret_basic = true
  client_secret_post  = false
  scope               = "user:email"
  use_pkce            = false

  authorization_endpoint = "https://github.com/login/oauth/authorize"
  token_endpoint         = "https://github.com/login/oauth/access_token"
  userinfo_endpoint      = "https://api.github.com/user"

  mfa_claim_path  = "$.two_factor_authentication"
  mfa_claim_value = "true"
}

resource "rauthy_auth_provider" "google" {
  enabled = true

  auto_onboarding = true
  auto_link       = true

  name                = "Google"
  typ                 = "google"
  issuer              = "https://accounts.google.com"
  client_id           = "google-client-id"
  client_secret       = "google-client-secret"
  client_secret_basic = true
  client_secret_post  = false
  scope               = "openid profile email"
  use_pkce            = true

  authorization_endpoint = "https://accounts.google.com/o/oauth2/v2/auth"
  token_endpoint         = "https://oauth2.googleapis.com/token"
  userinfo_endpoint      = "https://openidconnect.googleapis.com/v1/userinfo"
}

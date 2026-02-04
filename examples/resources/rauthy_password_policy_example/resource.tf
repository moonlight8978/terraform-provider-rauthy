resource "rauthy_password_policy" "default" {
  include_upper_case = 1
  include_lower_case = 1
  include_digits     = 1
  include_special    = 0
  length_min         = 14
  length_max         = 128
  not_recently_used  = 3
  valid_days         = 180
}

resource "casdoor_organization" "acme" {
  name          = "acme"
  display_name  = "Acme"
  website_url   = "https://acme.example.com"
  password_type = "bcrypt"
  country_codes = ["US"]
  languages     = ["en"]
}

resource "casdoor_webhook" "audit" {
  name         = "webhook-acme-audit"
  organization = casdoor_organization.acme.name
  url          = "https://audit.acme.example.com/casdoor"
  method       = "POST"
  content_type = "application/json"
  events       = ["signup", "login", "logout", "update-user"]
  is_enabled   = true

  headers = [
    {
      name  = "Authorization"
      value = "Bearer ${var.audit_token}"
    },
  ]
}

resource "casdoor_provider" "github" {
  name          = "provider-acme-github"
  display_name  = "GitHub"
  category      = "OAuth"
  type          = "GitHub"
  client_id     = var.github_client_id
  client_secret = var.github_client_secret
}

resource "casdoor_provider" "smtp" {
  owner         = casdoor_organization.acme.name
  name          = "provider-acme-smtp"
  category      = "Email"
  type          = "SMTP"
  host          = "smtp.example.com"
  port          = 465
  client_id     = "noreply@acme.example.com"
  client_secret = var.smtp_password
  title         = "Your verification code"
  content       = "Your code is %s, valid for 5 minutes."
}

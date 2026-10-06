resource "casdoor_application" "portal" {
  name          = "app-acme-portal"
  organization  = casdoor_organization.acme.name
  display_name  = "Acme Portal"
  cert          = casdoor_cert.acme.name
  redirect_uris = ["https://portal.acme.example.com/callback"]
  grant_types   = ["authorization_code", "refresh_token"]

  enable_password = true
  enable_sign_up  = true

  providers = [
    {
      name        = casdoor_provider.github.name
      owner       = "admin"
      can_sign_in = true
      can_sign_up = true
      rule        = "None"
    },
  ]
}

output "portal_client_id" {
  value = casdoor_application.portal.client_id
}

output "portal_client_secret" {
  value     = casdoor_application.portal.client_secret
  sensitive = true
}

resource "casdoor_invitation" "beta" {
  owner        = casdoor_organization.acme.name
  name         = "beta"
  display_name = "Beta testers"
  code         = "acme-beta-2026"
  quota        = 100
  application  = casdoor_application.portal.name
  state        = "Active"
}

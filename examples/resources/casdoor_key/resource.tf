resource "casdoor_key" "ci" {
  owner        = casdoor_organization.acme.name
  name         = "ci"
  display_name = "CI key"
  type         = "Organization"
  organization = casdoor_organization.acme.name
  state        = "Active"
}

resource "casdoor_role" "admins" {
  owner        = casdoor_organization.acme.name
  name         = "admins"
  display_name = "Admins"
  users        = ["acme/alice"]
  is_enabled   = true
}

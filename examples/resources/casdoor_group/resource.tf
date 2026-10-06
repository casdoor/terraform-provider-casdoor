resource "casdoor_group" "engineering" {
  owner        = casdoor_organization.acme.name
  name         = "engineering"
  display_name = "Engineering"
  type         = "Virtual"
  parent_id    = casdoor_organization.acme.name
  is_top_group = true
  is_enabled   = true
}

resource "casdoor_adapter" "policies" {
  owner       = casdoor_organization.acme.name
  name        = "policies"
  table       = "acme_policy"
  use_same_db = true
}

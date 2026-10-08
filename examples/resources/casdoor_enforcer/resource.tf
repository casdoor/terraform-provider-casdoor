resource "casdoor_enforcer" "api" {
  owner        = casdoor_organization.acme.name
  name         = "api"
  display_name = "API enforcer"
  model        = "${casdoor_model.rbac.owner}/${casdoor_model.rbac.name}"
  adapter      = "${casdoor_adapter.policies.owner}/${casdoor_adapter.policies.name}"
}

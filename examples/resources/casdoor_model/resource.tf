resource "casdoor_model" "rbac" {
  owner        = casdoor_organization.acme.name
  name         = "rbac"
  display_name = "RBAC"
  model_text   = <<-EOT
    [request_definition]
    r = sub, obj, act

    [policy_definition]
    p = sub, obj, act

    [role_definition]
    g = _, _

    [policy_effect]
    e = some(where (p.eft == allow))

    [matchers]
    m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
  EOT
}

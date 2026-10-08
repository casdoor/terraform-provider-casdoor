resource "casdoor_site" "portal" {
  owner        = casdoor_organization.acme.name
  name         = "portal"
  display_name = "Acme Portal"
  domain       = "portal.acme.example.com"
  host         = "http://10.0.0.5:8080"
  ssl_mode     = "HTTPS Only"
  rules        = ["${casdoor_rule.block_bots.owner}/${casdoor_rule.block_bots.name}"]
}

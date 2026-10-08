resource "casdoor_rule" "block_bots" {
  owner  = casdoor_organization.acme.name
  name   = "block-bots"
  type   = "User-Agent"
  action = "Block"
  reason = "Your request is blocked."

  expressions = [
    {
      name     = "bots"
      operator = "contains"
      value    = "BadBot"
    },
  ]
}

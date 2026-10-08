resource "casdoor_coupon" "launch" {
  owner              = casdoor_organization.acme.name
  name               = "launch"
  display_name       = "Launch discount"
  code               = "LAUNCH20"
  discount_type      = "percentage"
  discount           = 20
  scope              = "universal"
  quantity           = 100
  max_usage_per_user = 1
  currency           = "USD"
  state              = "Active"
}

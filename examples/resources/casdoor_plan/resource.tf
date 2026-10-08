resource "casdoor_plan" "pro_monthly" {
  owner        = casdoor_organization.acme.name
  name         = "pro-monthly"
  display_name = "Pro (monthly)"
  price        = 19
  currency     = "USD"
  period       = "Monthly"
  role         = "${casdoor_role.admins.owner}/${casdoor_role.admins.name}"
  is_enabled   = true
}

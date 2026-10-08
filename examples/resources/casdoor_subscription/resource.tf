resource "casdoor_subscription" "alice" {
  owner        = casdoor_organization.acme.name
  name         = "sub-alice"
  display_name = "Alice - Pro"
  user         = casdoor_user.alice.name
  pricing      = casdoor_pricing.default.name
  plan         = casdoor_plan.pro_monthly.name
  period       = "Monthly"
  start_time   = "2026-01-01T00:00:00Z"
  end_time     = "2027-01-01T00:00:00Z"
  state        = "Active"
}

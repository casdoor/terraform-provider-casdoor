resource "casdoor_pricing" "default" {
  owner          = casdoor_organization.acme.name
  name           = "default"
  display_name   = "Acme pricing"
  plans          = ["${casdoor_plan.pro_monthly.owner}/${casdoor_plan.pro_monthly.name}"]
  application    = casdoor_application.portal.name
  trial_duration = 7
  is_enabled     = true
}

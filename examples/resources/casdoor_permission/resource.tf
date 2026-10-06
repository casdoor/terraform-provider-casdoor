resource "casdoor_permission" "portal_admin" {
  owner         = casdoor_organization.acme.name
  name          = "portal-admin"
  display_name  = "Portal admin"
  roles         = ["acme/admins"]
  resource_type = "Application"
  resources     = [casdoor_application.portal.name]
  actions       = ["Read", "Write"]
  effect        = "Allow"
  is_enabled    = true
}

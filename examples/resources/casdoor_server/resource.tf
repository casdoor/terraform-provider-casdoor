resource "casdoor_server" "tools" {
  owner        = casdoor_organization.acme.name
  name         = "tools"
  display_name = "Acme tools"
  url          = "https://mcp.acme.example.com/mcp"
  token        = var.mcp_token
  application  = casdoor_application.portal.name
}

resource "casdoor_agent" "assistant" {
  owner        = casdoor_organization.acme.name
  name         = "assistant"
  display_name = "Assistant"
  url          = "https://agent.acme.example.com"
  token        = var.agent_token
  application  = casdoor_application.portal.name
}

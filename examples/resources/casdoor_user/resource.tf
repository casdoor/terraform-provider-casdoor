resource "casdoor_user" "alice" {
  owner        = casdoor_organization.acme.name
  name         = "alice"
  display_name = "Alice"
  email        = "alice@acme.example.com"
  password     = var.alice_password
  groups       = ["acme/engineering"]

  properties = {
    team = "platform"
  }
}

variable "alice_password" {
  type      = string
  sensitive = true
}

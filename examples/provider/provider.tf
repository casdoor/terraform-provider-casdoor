terraform {
  required_providers {
    casdoor = {
      source = "casdoor/casdoor"
    }
  }
}

# The client ID and secret of an application in the built-in organization,
# e.g. app-built-in, can manage every organization.
provider "casdoor" {
  endpoint      = "https://door.casdoor.com"
  client_id     = var.casdoor_client_id
  client_secret = var.casdoor_client_secret
}

variable "casdoor_client_id" {
  type = string
}

variable "casdoor_client_secret" {
  type      = string
  sensitive = true
}

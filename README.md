# Terraform Provider for Casdoor

[![CI](https://github.com/casdoor/terraform-provider-casdoor/actions/workflows/ci.yml/badge.svg)](https://github.com/casdoor/terraform-provider-casdoor/actions/workflows/ci.yml)
[![Terraform Registry](https://img.shields.io/badge/terraform-casdoor%2Fcasdoor-844FBA?logo=terraform)](https://registry.terraform.io/providers/casdoor/casdoor)

The official Terraform provider for [Casdoor](https://casdoor.ai), an open-source IAM / SSO platform. It manages organizations, applications, users, providers, certs, roles, permissions and groups as code.

- Provider documentation: [registry.terraform.io/providers/casdoor/casdoor](https://registry.terraform.io/providers/casdoor/casdoor/latest/docs)
- Guide: [Terraform on casdoor.ai](https://casdoor.ai/docs/deployment/terraform)

## Usage

```terraform
terraform {
  required_providers {
    casdoor = {
      source = "casdoor/casdoor"
    }
  }
}

provider "casdoor" {
  endpoint      = "https://door.casdoor.com"
  client_id     = var.casdoor_client_id
  client_secret = var.casdoor_client_secret
}

resource "casdoor_organization" "acme" {
  name         = "acme"
  display_name = "Acme"
}

resource "casdoor_application" "portal" {
  name          = "app-acme-portal"
  organization  = casdoor_organization.acme.name
  redirect_uris = ["https://portal.acme.example.com/callback"]
}
```

The provider calls the Casdoor API with the client ID and client secret of an application. Use an application of the `built-in` organization, e.g. `app-built-in`, to manage every organization. The settings can also be passed with the `CASDOOR_ENDPOINT`, `CASDOOR_CLIENT_ID` and `CASDOOR_CLIENT_SECRET` environment variables.

Only the attributes set in the configuration are written to Casdoor. The other fields of an object keep the values they have in Casdoor, so a resource can manage part of an object that is also edited in the Casdoor UI.

Existing objects can be imported: organizations and applications by name, the other objects by `owner/name`.

```shell
terraform import casdoor_organization.acme acme
terraform import casdoor_user.alice acme/alice
```

## Resources

| Resource | Casdoor object |
|---|---|
| `casdoor_organization` | Organization |
| `casdoor_application` | Application (OAuth 2.0 / OIDC / SAML client), with its providers |
| `casdoor_user` | User |
| `casdoor_provider` | Provider (OAuth, SAML, Email, SMS, Storage, Payment, Captcha, ...) |
| `casdoor_cert` | Cert |
| `casdoor_role` | Role |
| `casdoor_permission` | Permission |
| `casdoor_group` | Group |

## Development

Requirements: Go 1.25+ and Terraform 1.0+.

```shell
go build ./...
go generate ./...
```

`go generate` regenerates `docs/` from the schema and `examples/`.

To try a local build, point Terraform to it in `~/.terraformrc` (`%APPDATA%/terraform.rc` on Windows):

```hcl
provider_installation {
  dev_overrides {
    "casdoor/casdoor" = "/path/to/go/bin"
  }
  direct {}
}
```

then run `go install .` and use `terraform plan` / `terraform apply` directly, without `terraform init`.

## Release

Every push to `master` runs semantic-release on the commit messages (`feat:` for a minor version, `fix:` for a patch). When there is a new version, it pushes the `vX.Y.Z` tag, and GoReleaser builds, signs and publishes the GitHub release, which the Terraform Registry picks up.

## License

[Apache-2.0](LICENSE)

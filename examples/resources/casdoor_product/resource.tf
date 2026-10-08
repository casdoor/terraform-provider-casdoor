resource "casdoor_product" "pro" {
  owner        = casdoor_organization.acme.name
  name         = "pro"
  display_name = "Acme Pro"
  currency     = "USD"
  price        = 19
  quantity     = 1000
  providers    = [casdoor_provider.stripe.name]
  state        = "Published"
}

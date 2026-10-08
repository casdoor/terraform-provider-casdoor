resource "casdoor_syncer" "hr" {
  name          = "syncer-acme-hr"
  organization  = casdoor_organization.acme.name
  type          = "Database"
  database_type = "mysql"
  host          = "db.acme.example.com"
  port          = 3306
  user          = "casdoor"
  password      = var.hr_db_password
  database      = "hr"
  table         = "employee"
  sync_interval = 60
  is_read_only  = true
  is_enabled    = true

  table_columns = [
    {
      name         = "id"
      type         = "string"
      casdoor_name = "Id"
      is_key       = true
    },
    {
      name         = "login"
      type         = "string"
      casdoor_name = "Name"
    },
    {
      name         = "email"
      type         = "string"
      casdoor_name = "Email"
    },
  ]
}

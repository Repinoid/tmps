// https://registry.terraform.io/providers/yandex-cloud/yandex/latest/docs/resources/ydb_table
// Создание строковых таблиц в базе данных
// https://yandex.cloud/ru/docs/ydb/terraform/row-tables
resource "yandex_ydb_table" "metric_table" {
  path = "accounts"
  connection_string = var.dbEndpoint
  
column {
      name = "username"
      type = "Utf8"
      not_null = true
    }
    column {
      name = "password"
      type = "Utf8"
      not_null = true
    }
    column {
      name = "token"
      type = "Utf8"
      not_null = false
    }
    column {
      name = "created_at"
      type = "Timestamp"
      not_null = true
    }

  primary_key = ["username"]
  
}


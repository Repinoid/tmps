terraform {
  required_providers {
    yandex = {
      source = "yandex-cloud/yandex"
    }
  }
  required_version = ">= 0.13"
}
provider "yandex" {
  token     = var.token
  cloud_id  = var.cloud_id
  folder_id = var.folder_id
  zone      = var.compute-default-zone
}
///////////////////////////////////////////////////////////////////
// Создает архив. Пакует папку source_dir в архив с именем output_path
data "archive_file" "lambda" {
  type        = "zip"
  source_dir  = "./tozip/"  // путь в папке с кодом
  output_path = "goim.zip"  // zip пакуется сюда
}


# output "hash" {
#   value = data.archive_file.lambda.output_base64sha256
# }

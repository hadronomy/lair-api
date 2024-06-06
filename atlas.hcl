data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./internal/models",
    "--dialect", "postgres",
  ]
}

locals {
  db_host = urlescape(getenv("DB_HOST"))
  db_port = urlescape(getenv("DB_PORT"))
  db_user = urlescape(getenv("DB_USERNAME"))
  db_password = urlescape(getenv("DB_PASSWORD"))
  db_name = urlescape(getenv("DB_DATABASE"))
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "postgres://${local.db_user}:${local.db_password}@${local.db_host}:${local.db_port}/${local.db_name}?sslmode=disable"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

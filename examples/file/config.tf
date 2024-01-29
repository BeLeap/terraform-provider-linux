terraform {
  required_providers {
    linux = {
      source = "beleap/linux"
    }
  }
}

provider "linux" {
  host     = "localhost"
  username = "root"
  password = "root"
}

data "linux_file" "test" {
  path = "/etc"
  type = "directory"
}

output "test" {
  value = data.linux_file.test
}

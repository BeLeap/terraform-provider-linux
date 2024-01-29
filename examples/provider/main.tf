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

data "linux_user" "root" {
  username = "root"
}

output "root" {
  value = data.linux_user.root
}

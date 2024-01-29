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

resource "linux_user" "changseo_jang" {
  username = "testuser-changseo-jang"
  gid      = 2000
}

output "root" {
  value = data.linux_user.root
}

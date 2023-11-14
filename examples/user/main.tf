terraform {
  required_providers {
    linux = {
      source = "beleap/linux"
    }
  }
}

provider "linux" {
  host        = "node01.titanv.exp.riiid.cloud"
  username    = "root"
  private_key = file("../../ssh-keys/id_rsa")
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

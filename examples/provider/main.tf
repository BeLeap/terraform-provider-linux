terraform {
  required_providers {
    linux = {
      source = "beleap/linux"
    }
  }
}

provider "linux" {
  host        = "test-node.fox-deneb.ts.net"
  username    = "root"
  private_key = file("../../ssh-keys/id_rsa")
}

data "linux_user" "root" {
  username = "root"
}

output "root" {
  value = data.linux_user.root
}

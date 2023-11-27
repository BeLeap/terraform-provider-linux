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

data "linux_file" "test" {
  path = "/etc"
  type = "directory"
}

output "test" {
  value = data.linux_file.test
}

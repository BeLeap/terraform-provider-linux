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
  private_key = file("~/.ssh/gpu_id_rsa")
}

data "linux_dummy" "dummy" {}


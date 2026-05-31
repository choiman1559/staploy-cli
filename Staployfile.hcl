# THIS IS Example of "staploy-cli file" action
configure {
  address = "127.0.0.1"
  port    = 18090
  use_name = true
}

manage "testbin" {
  create {
     description = "test"
  }
  upload {
     path = "./out/testbin_0.1.0.tar"
  }
#  delete {
#     versions = ["0.1.0"]
#  }
}

target "deploy_test" {
  workers = ["0b8a272b-1d85-4c1e-ba2a-102cb0a14147"]

  deploy "testbin" {
    version     = "0.1.0"
    post_deploy = "testbin"
  }
}

target "push_test" {
  workers = ["0b8a272b-1d85-4c1e-ba2a-102cb0a14147"]

  push_only "testbin" {
    version = "0.1.0"
  }

  set_active "testbin" {
    version     = "0.1.0"
    post_deploy = "testbin"
  }
}

target "remove_test" {
  workers = ["0b8a272b-1d85-4c1e-ba2a-102cb0a14147"]

  remove "testbin" {
    version = "0.1.0"
    autoremove = false
  }
}



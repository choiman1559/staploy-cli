alias {
  app "staploy-build" {
    name    = "staploy-cli"
    version = "1.0.0"
  }
}

build "alias:staploy-build" {
  output_dir  = "out"
  executable = ["staploy-cli"]
  lib_version = "shell:go version"

  i386 { path = "out/386" }
  x86_64 { path = "out/amd64" }
  arm { path = "out/arm" }
  aarch64 { path = "out/arm64" }
  riscv64 { path = "out/riscv64" }
  mipsel { path = "out/mipsle" }
  mips64el { path = "out/mips64le" }
  mips { path = "out/mips" }
  mips64 { path = "out/mips64" }
  ppc64 { path = "out/ppc64" }
  ppc64le { path = "out/ppc64le" }
  s390x { path = "out/s390x" }
  loong64 { path = "out/loong64" }
}

configure {
  address      = "shell:echo $STAPLOY_HOST_ADDR"
  port         = "shell:echo $STAPLOY_HOST_PORT"
  enforce_uuid = false
}

manage "alias:staploy-build" {
  upload {}
}

target "staploy-deploy" {
  workers = ["group:all"]

  where {
    arch = ["!arm", "!aarch64", "!mips"]
  }

  # First unlink current activated version
  unset "alias:staploy-build" {
  }

  # Then push & set desire version
  deploy "alias:staploy-build" {
  }
}
pkg_name=bldr_package_sync
pkg_origin=indellient
pkg_version="0.1.0"
pkg_scaffolding="core/scaffolding-go"
pkg_bin_dirs=(bin)
scaffolding_go_build_deps=(
  "github.com/sirupsen/logrus"
  "github.com/BurntSushi/toml"
  "github.com/urfave/cli"
  "github.com/hashicorp/go-version"
)
pkg_build_deps=(
  "core/git"
  "core/busybox-static"
)
pkg_svc_user="root"
pkg_svc_group=$pkg_svc_user
export GOFLAGS=-mod=vendor

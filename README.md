# Go bindings for DPDK framework.
[![Documentation](https://godoc.org/github.com/tbarbette/go-dpdk?status.svg)](http://godoc.org/github.com/tbarbette/go-dpdk) [![Build Status](https://github.com/tbarbette/go-dpdk/actions/workflows/unit.yml/badge.svg)](https://github.com/tbarbette/go-dpdk/actions/workflows/unit.yml) [![codecov](https://codecov.io/gh/tbarbette/go-dpdk/branch/master/graph/badge.svg?token=1XW04KL02S)](https://codecov.io/gh/tbarbette/go-dpdk)

# Building apps

Starting from DPDK 21.05, `pkg-config` becomes the only official way to build DPDK apps. Because of it `go-dpdk` uses `#cgo pkg-config` directive to link against your DPDK distribution.

Go compiler may fail to accept some C compiler flags. You can fix it by submitting those flags to environment:
```
export CGO_CFLAGS_ALLOW="-mrtm"
export CGO_LDFLAGS_ALLOW="-Wl,--(?:no-)?whole-archive"
```

## Non-system DPDK installation

If DPDK is not installed globally, point `pkg-config` to your installation:
```sh
export DPDK_PATH=/path/to/dpdk/install
export PKG_CONFIG_PATH=$DPDK_PATH/lib/x86_64-linux-gnu/pkgconfig:$PKG_CONFIG_PATH
```

This applies both when building this repo directly and when using it as a dependency in another Go module — `go build` calls `pkg-config` using your current environment.

## Cache invalidation when switching DPDK versions

Go caches CGo builds keyed on `CGO_CFLAGS`/`CGO_LDFLAGS`. If you switch between DPDK installations, export these variables explicitly so Go detects the change and invalidates its cache automatically:
```sh
export CGO_CFLAGS=$(pkg-config --cflags libdpdk)
```
Without this, `go build` may silently reuse a binary linked against a different DPDK version. In this repo, `make build` handles this automatically via a `.dpdk-version` stamp file.

# Caveats
* Only dynamic linking is viable at this point.
* If you isolate CPU cores with `isolcpus` kernel parameter then `GOMAXPROCS` should be manually specified to reflect the actual number of logical cores in CPU mask. E.g. if `isolcpus=12-95` on a 96-core machine then default value for `GOMAXPROCS` would be 12 but it should be at least 84.

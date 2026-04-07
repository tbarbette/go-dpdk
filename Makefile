# Capture compile flags from the currently resolved libdpdk.
# When DPDK_PATH / PKG_CONFIG_PATH changes, CGO_CFLAGS changes too,
# which invalidates Go's build cache automatically.
# Linking is still handled by "#cgo pkg-config: libdpdk" in source files.
CGO_CFLAGS  := $(shell pkg-config --cflags libdpdk)
CGO_LDFLAGS := $(shell pkg-config --libs   libdpdk 2>/dev/null)
DPDK_VERSION := $(shell pkg-config --modversion libdpdk)
DPDK_LIBDIR  := $(shell pkg-config --variable=libdir libdpdk)

export CGO_CFLAGS
export CGO_LDFLAGS

APPS := $(wildcard app/*)
APP_BINS := $(patsubst app/%,bin/%,$(APPS))
STAMP := .dpdk-version

.PHONY: build libs test clean $(APPS)

# Regenerate stamp only when DPDK version or libdir actually changes.
$(STAMP): FORCE
	@echo "$(DPDK_VERSION) $(DPDK_LIBDIR)" | cmp -s - $@ || \
	  (echo "DPDK changed to $(DPDK_VERSION) ($(DPDK_LIBDIR)), rebuilding..."; \
	   echo "$(DPDK_VERSION) $(DPDK_LIBDIR)" > $@; \
	   go clean -cache)

FORCE:

build: $(STAMP) libs $(APP_BINS)

libs: $(STAMP)
	go build ./...

bin/%: app/% $(STAMP)
	@mkdir -p bin
	go build -o $@ ./$<

test:
	go test ./...

test-bin:
	@mkdir -p bin
	@for pkg in $(shell go list ./... | grep -v '/app/'); do \
	  name=$$(echo $$pkg | sed 's|.*/||'); \
	  go test -c -o bin/$$name.test $$pkg && echo "built bin/$$name.test" || true; \
	done

clean:
	go clean -cache
	rm -rf bin $(STAMP)

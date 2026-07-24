MODULE := $(shell grep ^module go.mod | awk '{print $$2}')
CLI_TOOL := /home/cuj1559/Documents/Staploy/staploy-cli/out/temp/go_build_staploy_cli
SRC_DIR := protobuf
PROTO_OUT_DIR := app/proto
BUILD_OUT_DIR := out

GOPATH_DEFAULT := $(shell go env GOPATH)
ifeq ($(GOPATH_DEFAULT),)
    GOPATH_DEFAULT := $(HOME)/go
endif

GOPATH_BIN := $(GOPATH_DEFAULT)/bin
PROTO_FILES := $(shell find $(SRC_DIR) -name "*.proto")
M_ARGS := $(foreach file,$(PROTO_FILES),--go_opt=M$(subst $(SRC_DIR)/,,$(file))=$(MODULE)/$(PROTO_OUT_DIR))

.PHONY: all proto buildPkg clean buildAll $(ARCHES)
all: proto buildAll buildPkg clean

ARCHES := 386 amd64 arm arm64 riscv64 mipsle mips64le
buildAll: $(ARCHES)
createAll: buildAll buildPkg

$(ARCHES):
	CGO_ENABLED=0 GOOS=linux GOARCH=$@ go build -ldflags="-s -w" -o $(BUILD_OUT_DIR)/$@/staploy-cli staploy-cli

buildPkg:
	$(CLI_TOOL) file -f ./build_pkg.hcl -v

proto:
	@mkdir -p $(PROTO_OUT_DIR)
	protoc -I=$(SRC_DIR) \
		--plugin=protoc-gen-go=$(GOPATH_BIN)/protoc-gen-go \
		--go_out=$(PROTO_OUT_DIR) \
		--go_opt=paths=source_relative \
		$(M_ARGS) \
		$(PROTO_FILES)

clean:
	rm -rf $(PROTO_OUT_DIR)

cleanBuilds:
	rm -rf $(BUILD_OUT_DIR)

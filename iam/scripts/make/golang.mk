
GO ?= go
# Use regexp to match.
GO_SUPPORT_VERSION ?= 1.13|1.14|1.15|1.16|1.17|1.18|1.19|1.20

# 把Git信息传递到程序里的
GO_VERSION += -X $(VERSION_PACKAGE).GitVersion=$(VERSION) \
	-X $(VERSION_PACKAGE).GitCommit=$(GIT_COMMIT) \
	-X $(VERSION_PACKAGE).GitTreeState=$(GIT_TREE_STATE) \
	-X $(VERSION_PACKAGE).BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GO_LDFLAGS += $(GO_VERSION)

DLV := $(shell which dlv)
ifneq ($(DLV),)
		# 如果二进制文件在编译时没有关闭优化功能，可能很难正确地调试它。
		# 请考虑在Go 1.10或更高版本上用-gcflags="all=-N -l "编译调试二进制文件，
		# 在Go的早期版本上用-gcflags="-N -l"。
		GO_BUILD_FLAGS += -gcflags "all=-N -l"
    	LDFLAGS = ""
endif

GO_BUILD_FLAGS += -ldflags "$(GO_LDFLAGS)"

ifeq ($(GOOS),windows)
	GO_OUT_EXT := .exe
endif

ifeq ($(ROOT_PACKAGE),)
	$(error the variable ROOT_PACKAGE must be set prior to including golang.mk)
endif

GOPATH := $(shell go env GOPATH)
GOBIN ?= $(GOPATH)/bin

# 从cmd文件夹取出命令集合，排除md文件
COMMANDS ?= $(filter-out %.md, $(wildcard ${ROOT_DIR}/cmd/*))
# 从COMMANDS各个路径里拿到不是文件夹的。
BINS ?= $(foreach cmd,${COMMANDS},$(notdir ${cmd}))

ifeq (${COMMANDS},)
  $(error Could not determine COMMANDS, set ROOT_DIR or run in source dir)
endif
ifeq (${BINS},)
  $(error Could not determine BINS, set ROOT_DIR or run in source dir)
endif

.PHONY: go.build.check
go.build.check:
ifneq ($(shell $(GO) version | grep -q -E '\bgo($(GO_SUPPORTED_VERSIONS))\b' && echo 0 || echo 1), 0)
	$(error unsupported go version. Please make install one of the following supported version: '$(GO_SUPPORTED_VERSIONS)')
endif

# Input: make go.build.iam-apiserver.linux_amd64
.PHONY: go.build.%
go.build.%:
	$(eval COMMAND := $(word 1, $(subst ., ,$*)))
	$(eval PLATFORM := $(word 2, $(subst ., ,$*)))
	$(eval OS := $(word 1, $(subst ., ,PLATFORM)))
	$(eval ARCH := $(word 2, $(subst ., ,PLATFORM)))
	@echo "$(COMMON_ARROW) Building binary $(COMMAND) $(VERSION) for $(OS) $(ARCH)."
	@mkdir -p $(OUTPUT_DIR)/platforms/$(OS)/$(ARCH)
	@CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) $(GO) build $(GO_BUILD_FLAGS) -o $(OUTPUT_DIR)/platforms/$(OS)/$(ARCH)/$(COMMAND)$(GO_OUT_EXT) $(ROOT_PACKAGE)/cmd/$(COMMAND)


.PHONY: go.build
go.build: go.build.check $(addprefix go.build., $(addprefix $(BINS)., $(PLATFORM)))

.PHONY: go.build.multi-platform
go.build.multi-platform: go.build.check $(foreach plt, $(PLATFORMS), $(addprefix go.build., $(addprefix $(BINS)., $(plt))))

.PHONY: go.lint
go.lint: tools.check.golangci-lint
	@echo "$(COMMON_ARROW) Run golangci to lint source codes"
	@golangci-lint run -c $(ROOT_DIR)/../.golangci.yaml $(ROOT_DIR)/...

.PHONY: go.clean
go.clean:
	@echo "$(COMMON_ARROW) Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)

.PHONY: go.test
go.test: tools.verify.go-junit-report
	@echo "===========> Run unit test"
	@set -o pipefail;$(GO) test -race -cover -coverprofile=$(OUTPUT_DIR)/coverage.out \
		-timeout=10m -shuffle=on -short -v `go list ./...|\
		egrep -v $(subst $(SPACE),'|',$(sort $(EXCLUDE_TESTS)))` 2>&1 | \
		tee >(go-junit-report --set-exit-code >$(OUTPUT_DIR)/report.xml)
	@sed -i '/mock_.*.go/d' $(OUTPUT_DIR)/coverage.out # remove mock_.*.go files from test coverage
	@$(GO) tool cover -html=$(OUTPUT_DIR)/coverage.out -o $(OUTPUT_DIR)/coverage.html

.PHONY: go.test.cover
go.test.cover: go.test
	@$(GO) tool cover -func=$(OUTPUT_DIR)/coverage.out | \
		awk -v target=$(COVERAGE) -f $(ROOT_DIR)/scripts/coverage.awk

.PHONY: go.updates
go.updates: tools.verify.go-mod-outdated
	@$(GO) list -u -m -json all | go-mod-outdated -update -direct
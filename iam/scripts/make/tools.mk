


.PHONY: tools.install.%
tools.install.%:
	@echo "$(COMMON_ARROW) Installing $*"
	@$(MAKE) install.$*

.PHONY: install.golangci-lint
install.golangci-lint:
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2
	@#golangci-lint completion bash > $(HOME)/.golangci-lint.bash
	@#if ! grep -q .golangci-lint.bash $(HOME)/.bashrc; then echo "source \$$HOME/.golangci-lint.bash" >> $(HOME)/.bashrc; fi
	@echo "golangci-lint installed."

.PHONY: install.go-mod-outdated
install.go-mod-outdated:
	@$(GO) install github.com/psampaz/go-mod-outdated@latest

.PHONY: tools.check.%
tools.check.%:
	@if ! which $* &>/dev/null; then $(MAKE) tools.install.$*; fi


.PHONY: install.go-junit-report
install.go-junit-report:
	@$(GO) install github.com/jstemmer/go-junit-report@latest
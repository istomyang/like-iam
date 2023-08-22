

.DEFAULT_GOAL = all


SHELL := /usr/bin/env bash

.PHONY: all
all:
	@echo "this is test"
	@echo $(MAKEFILE_LIST)

.PHONY: test1.%
test1.%: all
	echo $*
	echo $@
	echo $<
	echo $%


### tidy: run go mod tidy.
.PHONY: tidy
tidy:
	@$(GO) mod tidy


### help: Show this help info.
.PHONY: help
help: Makefile
	@printf "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:\n"
	@# 1 筛选注释，2 替换冒号为空行（表格），3 每一行加空格
	@sed -n 's/^###//p' $< | column -t -s ':' | sed -e 's/^/ /'
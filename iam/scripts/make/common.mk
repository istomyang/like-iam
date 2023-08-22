
SHELL := /usr/bin/env bash

ROOT_PACKAGE=istomyang.github.com/like-iam/iam
# 用于把版本信息录入到程序里面
VERSION_PACKAGE=istomyang.github.com/like-iam/component-base/version

# is common.mk file dir
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
ROOT_DIR ?= $(abspath $(shell cd $(COMMON_SELF_DIR)/../.. && pwd -P))
OUTPUT_DIR ?= $(ROOT_DIR)/_output
$(shell mkdir -p $(OUTPUT_DIR))
TOOLS_DIR ?= $(OUTPUT_DIR)/tools
$(shell mkdir -p $(TOOLS_DIR))
TMP_DIR ?= $(OUTPUT_DIR)/tmp
$(shell mkdir -p $(TMP_DIR))

# set the version number. you should not need to do this
# for the majority of scenarios.
ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --tags --always --match='v*')
endif
# Check if the tree is dirty.  default to dirty
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
	GIT_TREE_STATE="clean"
endif
GIT_COMMIT:=$(shell git rev-parse HEAD)

# Minimum test coverage
ifeq ($(origin COVERAGE),undefined)
COVERAGE := 60
endif

# The OS must be linux when building docker images
PLATFORMS ?= linux_amd64 linux_arm64
# The OS can be linux/windows/darwin when building binaries
# PLATFORMS ?= darwin_amd64 windows_amd64 linux_amd64 linux_arm64

ifeq($(origin $(PLATFORM)), undefined)
	GOOS ?= $(shell go env GOOS)
	GOARCH ?= $(shell go env GOARCH)
	PLATFORM := $(GOOS)_$(GOARCH)
	# Build images for container should be in linux.
	IMAGE_PLATFORM = linux_$(GOARCH)
else
	GOOS := $(word 1, $(subst _, ,$(PLATFORM)))
	GOARCH := $(word 2, $(subst _, ,$(PLATFORM)))
	IMAGE_PLATFORM := $(PLATFORM)
endif


# Makefile settings
ifndef V
MAKEFLAGS += --no-print-directory
endif

COMMON_ARROW := =========>

COMMA := ,
SPACE :=
SPACE +=
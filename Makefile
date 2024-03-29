# Copyright 2018 The Kubeflow Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Bump these on release
VERSION_MAJOR ?= 0
VERSION_MINOR ?= 1
VERSION_BUILD ?= 0

VERSION ?= v$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_BUILD)

GOOS ?= $(shell go env GOOS)
GOARCH = amd64
BUILD_DIR ?= ./out
ORG := github.com/kubeflow
PROJECT := fairing
REPOPATH ?= $(ORG)/$(PROJECT)
RELEASE_BUCKET ?= $(PROJECT)

SUPPORTED_PLATFORMS := linux-$(GOARCH) darwin-$(GOARCH) windows-$(GOARCH).exe
BUILD_PACKAGE = $(REPOPATH)/cmd/$(PROJECT)
VERSION_PACKAGE = $(REPOPATH)/pkg/$(PROJECT)/version


GO_LDFLAGS :="
GO_LDFLAGS += -X $(VERSION_PACKAGE).version=$(VERSION)
GO_LDFLAGS += -X $(VERSION_PACKAGE).buildDate=$(shell date +'%Y-%m-%dT%H:%M:%SZ')
GO_LDFLAGS += -X $(VERSION_PACKAGE).gitCommit=$(shell git rev-parse HEAD)
GO_LDFLAGS += -X $(VERSION_PACKAGE).gitTreeState=$(if $(shell git status --porcelain),dirty,clean)
GO_LDFLAGS +="

$(BUILD_DIR)/$(PROJECT): $(BUILD_DIR)/$(PROJECT)-$(GOOS)-$(GOARCH)
	cp $(BUILD_DIR)/$(PROJECT)-$(GOOS)-$(GOARCH) $@

$(BUILD_DIR)/$(PROJECT)-%-$(GOARCH): $(GO_FILES) $(BUILD_DIR)
	GOOS=$* GOARCH=$(GOARCH) CGO_ENABLED=0 go build -ldflags $(GO_LDFLAGS) -o $@ $(BUILD_PACKAGE)

%.sha256: %
	shasum -a 256 $< &> $@

%.exe: %
	mv $< $@

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

.PRECIOUS: $(foreach platform, $(SUPPORTED_PLATFORMS), $(BUILD_DIR)/$(PROJECT)-$(platform))

.PHONY: cross
cross: $(foreach platform, $(SUPPORTED_PLATFORMS), $(BUILD_DIR)/$(PROJECT)-$(platform).sha256)

.PHONY: test
test: $(BUILD_DIR)/$(PROJECT)
	@ ./test.sh

.PHONY: install
install: $(GO_FILES) $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go install -ldflags $(GO_LDFLAGS) $(BUILD_PACKAGE)

.PHONY: release
release: cross
	gsutil cp $(BUILD_DIR)/$(PROJECT)-* gs://$(RELEASE_BUCKET)/$(VERSION)/
	gsutil cp $(BUILD_DIR)/$(PROJECT)-* gs://$(RELEASE_BUCKET)/latest/

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

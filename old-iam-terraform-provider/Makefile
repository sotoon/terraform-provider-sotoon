NAME=sotoon
BINARY=terraform-provider-$(NAME)
BINARY_PATH=~/go/bin
LATEST_VERSION=$(shell cat ./VERSION)

.PHONY: build
build:
	mkdir -p $(BINARY_PATH)
	go mod tidy
	go build -o $(BINARY_PATH)/$(BINARY)
	
.PHONY: safebuild
safebuild:
	go install mvdan.cc/garble@latest
	$(BINARY_PATH)/garble -tiny -literals build -o $(BINARY_PATH)/$(BINARY)

.PHONY: test
test:
	TF_ACC=1 go test -mod=vendor -v ./... -count=1 -coverprofile cover.out
	go tool cover -func=cover.out
	
.PHONY: coverage-serve
coverage-serve:
	go tool cover -html=cover.out

.PHONY: docs
docs:
	go generate ./...
	sed -i 's/SOTOON_TERRAFORM_PROVIDER_REGISTRY/$(LATEST_VERSION)/' ./docs/*.md
	sed -i 's/# sotoon Provider/# Sotoon Provider/' ./docs/*.md
	sed -i 's/page_title: "sotoon Provider"/page_title: "Sotoon Provider"/' ./docs/*.md
	mv ./docs/index.md ./docs/provider.md
	mv ./docs/home.md ./docs/index.md
	mkdocs build

.PHONY: deploy-docs
deploy-docs:
	aws s3 --endpoint-url=https://s3.thr1.sotoon.ir rm s3://terraform-docs/ --recursive
	aws s3 --endpoint-url=https://s3.thr1.sotoon.ir sync ./site/ s3://terraform-docs/

.PHONY: doclint
doclint:
	terraform fmt -check -recursive ./examples

.PHONY: fixdoclint
fixdoclint:
	terraform fmt -recursive ./examples

.PHONY: release
release:
	goreleaser release --clean

.PHONY: push-provider
push-provider:
	bash scripts/registry-push.bash

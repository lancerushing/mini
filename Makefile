SHELL := bash

export BUILD_VERSION := $(shell date -u +'%Y-%m-%dT%H.%M.%SZ')
export GO_PATH := $(shell go env GOPATH)


########################################
.PHONY: tidy
tidy:
	go mod tidy -v

########################################
.PHONY: build
build: tidy fmt lint
	go build -race -ldflags "-s -w" -v .

########################################
.PHONY: package
package: build
	$(eval PACKAGE_DIR := /tmp/mini-$(BUILD_VERSION))
	mkdir $(PACKAGE_DIR)
	cp mini $(PACKAGE_DIR)/
	cp mini.yaml $(PACKAGE_DIR)/

	mkdir $(PACKAGE_DIR)/server
	rsync -a server/templates $(PACKAGE_DIR)/server

	cd /tmp && \
	tar czvf mini-$(BUILD_VERSION).tar.gz mini-$(BUILD_VERSION)

	echo "Package File Ready: /tmp/mini-$(BUILD_VERSION).tar.gz"


########################################
.PHONY: dev
dev:
	$(GO_PATH)/bin/modd

########################################
.PHONY: test
test:
	go test -race ./...


########################################
.PHONY: fmt
fmt:
	gofmt -s -l -w .

########################################
.PHONY: lint
lint:
	$(GO_PATH)/bin/golint -set_exit_status ./...

########################################
.PHONY: check
check: tidy fmt lint test
	# slow check
	$(GO_PATH)/bin/golangci-lint run --enable-all

########################################
.PHONY: create-db
create-db:
	sudo -u postgres psql postgres -c "CREATE DATABASE mini;"
	sudo -u postgres psql mini < db.sql
	# Test connection
	PGPASSWORD=mini_app psql --host localhost --username mini_app mini -c "SELECT 'success' as result;"



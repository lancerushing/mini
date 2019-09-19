SHELL := bash

export BUILD_VERSION := $(shell date -u +'%Y-%m-%dT%H.%M.%SZ')


########################################
.PHONY: tidy
tidy:
	go mod tidy -v

########################################
.PHONY: build
build: tidy
	go build -race-ldflags "-s -w" -v .

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
	gin  --port 4000 --appPort 4001 --immediate --build . --buildArgs "-v" --all run

########################################
.PHONY: test
test:
	go test -race ./...


########################################
.PHONY: fmt
fmt:
	go fmt *.go
	go fmt ./server/

########################################
.PHONY: create-db
create-db:
	sudo -u postgres psql postgres -c "CREATE DATABASE mini;"
	sudo -u postgres psql mini < db.sql
	# Test connection
	PGPASSWORD=mini_app psql --host localhost --username mini_app mini -c "SELECT 'success' as result;"




# Variables
# UNAME   := $(shell uname -s)

COMMIT_ID := `git log -1 --format=%H`
COMMIT_DATE := `git log -1 --format=%aI`
VERSION := $${CI_COMMIT_TAG:-SNAPSHOT-$(COMMIT_ID)}

.PHONY: ensure
ensure:
  @dep ensure
#
#.PHONY: test
#test: ensure
# go test -v -coverprofile=coverage.out $$(go list ./... | grep -v '/vendor/') && go tool cover -func=coverage.out

.PHONY: build
build: ensure
  @GOOS=linux GOARCH=arm64 go build -ldflags "-s -w -X github.com/tnosaj/weather-station-exporter/version.Version=$(VERSION) -X github.com/tnosaj/weather-station-exporter/version.Commit=$(COMMIT_ID) -X github.com/tnosaj/weather-station-exporter/version.Date=$(COMMIT_DATE)"

PLUGIN_FILE := plugins.txt
BUILD_TAG := $(tag)
BUILD_TIME := `date +%FT%T%z`
COMMIT_SHA1 := `git rev-parse HEAD`
TAG := `git describe --tags`
PLUGIN_LIST := $(shell cat ${PLUGIN_FILE})
PLUGIN_OUTPUT_PATH := plugins
all: so
	go get -d -v; \
	go build -ldflags " \
		-X github.com/cool2645/kotori-ng/version.BuildTag=${BUILD_TAG} \
		-X 'github.com/cool2645/kotori-ng/version.BuildTime=${BUILD_TIME}' \
		-X github.com/cool2645/kotori-ng/version.GitCommitSHA1=${COMMIT_SHA1} \
		-X github.com/cool2645/kotori-ng/version.GitTag=${TAG} \
		"
so:
	for plugin in $(PLUGIN_LIST); do \
		go get -d -v $$plugin; \
		cd $(shell go env GOPATH)/src/$$plugin; \
		go build -ldflags " \
			-X $$plugin.BuildTag=${BUILD_TAG} \
			-X '$$plugin.BuildTime=${BUILD_TIME}' \
			-X $$plugin.GitCommitSHA1=${COMMIT_SHA1} \
			-X $$plugin.GitTag=${TAG} \
			" --buildmode=plugin; \
		cd -; \
		cp $(shell go env GOPATH)/src/$$plugin/*.so $(PLUGIN_OUTPUT_PATH); \
	done

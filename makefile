file := plugins.txt
path := plugins
pluginlist := $(shell cat ${file})
all: so
	 go get -d -v; \
	 go build
so:
	 for plugin in $(pluginlist); do \
		 go get -d -v $$plugin; \
		 cd $(shell go env GOPATH)/src/$$plugin; \
		 go build --buildmode=plugin; \
		 cd -; \
		 cp $(shell go env GOPATH)/src/$$plugin/*.so plugins; \
	 done

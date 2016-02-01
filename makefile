ifndef TAVIS_HOST
  TAVIS_HOST:=$(warning Please define TAVIS_HOST environment variable, using default "localhost")localhost
endif

ifndef TAVIS_PORT
  TAVIS_PORT:=$(warning Please define TAVIS_PORT environment variable, using "8080")8080
endif

ifndef TAVIS_CONFIG
  TAVIS_CONFIG:=$(warning Please define TAVIS_CONFIG environment variable, using "config/config.json.sample")config/config.json.sample
endif

all:
	@echo "Please specify a target!"

install/linux:
	install.sh linux

install/darwin:
	install.sh darwin

build:
	protoc --go_out=. payload/*.proto
	go build -o tivan

run: build
	./tivan -config $(TAVIS_CONFIG) -verbose

profile/:
	mkdir -p $@

profile/cpu: profile/
	go tool pprof --pdf tivan $(TAVIS_HOST):$(TAVIS_PORT)/debug/pprof/profile > $@.pdf

profile/goroutine: profile/
	go tool pprof --pdf tivan $(TAVIS_HOST):$(TAVIS_PORT)/debug/pprof/goroutine > $@.pdf

profile/heap: profile/
	go tool pprof --pdf tivan $(TAVIS_HOST):$(TAVIS_PORT)/debug/pprof/heap > $@.pdf

tests/stress/single:
	wrk -t5 -c10 -d60s -s tests/stress-test/single.lua "$(TAVIS_HOST):$(TAVIS_PORT)/api/v1/track"

tests/stress/multi:
	wrk -t5 -c10 -d1m -s tests/stress-test/multi.lua "$(TAVIS_HOST):$(TAVIS_PORT)/api/v1/track"

ifndef THOST
  THOST:=$(warning Please define THOST environment variable, using default "http://localhost")http://localhost
endif

ifndef TPORT
  TPORT:=$(warning Please define TPORT environment variable, using "8080")8080
endif

ifndef TCONFIG
  TCONFIG:=$(warning Please define TCONFIG environment variable, using "config/config.json.sample")config/config.json.sample
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
	./tivan -config $(TCONFIG) -verbose

profile/:
	mkdir -p $@

profile/cpu: profile/
	go tool pprof --pdf tivan $(THOST):$(TPORT)/debug/pprof/profile > $@.pdf

profile/goroutine: profile/
	go tool pprof --pdf tivan $(THOST):$(TPORT)/debug/pprof/goroutine > $@.pdf

profile/heap: profile/
	go tool pprof --pdf tivan $(THOST):$(TPORT)/debug/pprof/heap > $@.pdf

tests/stress/single:
	wrk -t5 -c10 -d60s -s tests/stress-test/single.lua "$(THOST):$(TPORT)/api/v1/track"

tests/stress/multi:
	wrk -t5 -c10 -d1m -s tests/stress-test/multi.lua "$(THOST):$(TPORT)/api/v1/track"

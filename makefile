ifndef TAVIS_SCHEMA
  TAVIS_SCHEMA:=http://
endif

ifndef TAVIS_HOST
  TAVIS_HOST:=$(warning Please define TAVIS_HOST environment variable, using default localhost")localhost
endif

ifndef TAVIS_PORT
  TAVIS_PORT:=$(warning Please define TAVIS_PORT environment variable, using "8080")8080
endif

ifndef TAVIS_CONFIG
  TAVIS_CONFIG:=$(warning Please define TAVIS_CONFIG environment variable, using "config/config.json.sample")config/config.json.sample
endif

all:
	@echo "Please specify a target!"

install/%:
	./install.sh $*

build:
	protoc --go_out=. payload/*.proto
	protoc --python_out=utils payload/*.proto
	go build -o tivan

run: build
	./tivan -config $(TAVIS_CONFIG) -verbose

profile/:
	mkdir -p $@

profile/cpu: profile/
	go tool pprof --pdf tivan $(TAVIS_SCHEMA)$(TAVIS_HOST):$(TAVIS_PORT)/debug/pprof/profile > $@.pdf

profile/goroutine: profile/
	go tool pprof --pdf tivan $(TAVIS_SCHEMA)$(TAVIS_HOST):$(TAVIS_PORT)/debug/pprof/goroutine > $@.pdf

profile/heap: profile/
	go tool pprof --pdf tivan $(TAVIS_SCHEMA)$(TAVIS_HOST):$(TAVIS_PORT)/debug/pprof/heap > $@.pdf

tests/send:
	python utils/send_single_message.py $(TAVIS_CONFIG) "$(TAVIS_SCHEMA)$(TAVIS_HOST):$(TAVIS_PORT)/api/v1/track"

tests/stress/1-messages/:
	python utils/generate_stress_messages.py $(TAVIS_CONFIG) $@

tests/stress/n-messages/:
	python utils/generate_stress_messages.py -r $(TAVIS_CONFIG) $@

tests/stress/%: tests/stress/%-messages/
	cd $< && wrk -t5 -c10 -d1m -s ../run.lua "$(TAVIS_SCHEMA)$(TAVIS_HOST):$(TAVIS_PORT)/api/v1/track"

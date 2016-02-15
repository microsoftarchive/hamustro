ifndef HAMUSTRO_SCHEMA
  HAMUSTRO_SCHEMA:=http://
endif

ifndef HAMUSTRO_HOST
  HAMUSTRO_HOST:=$(warning Please define HAMUSTRO_HOST environment variable, using default localhost")localhost
endif

ifndef HAMUSTRO_PORT
  HAMUSTRO_PORT:=$(warning Please define HAMUSTRO_PORT environment variable, using "8080")8080
endif

ifndef HAMUSTRO_CONFIG
  HAMUSTRO_CONFIG:=$(warning Please define HAMUSTRO_CONFIG environment variable, using "config/config.json.sample")config/config.json.sample
endif

PYC:=python

all:
	@echo "Please specify a target!"

install/%:
	./utils/installer/_$*.sh

src/%.go:
src/%/%.go:
src/%/%/%.go:

hamustro: src/payload/ src/*.go src/*/*.go src/*/*/*.go
	go build -o $@ src/*.go

src/payload/:
	protoc --go_out=. proto/*.proto
	mkdir -p $@ && mv proto/*.go src/payload/

utils/payload/:
	protoc --python_out=. proto/*.proto
	mkdir -p $@ && mv proto/*.py utils/payload/
	echo "from payload_pb2 import *" > $@/__init__.py

dev: hamustro utils/payload/
	./$< -config $(HAMUSTRO_CONFIG) -verbose

server: hamustro
	./$< -config $(HAMUSTRO_CONFIG)

profile/:
	mkdir -p $@

profile/cpu: profile/
	go tool pprof --pdf hamustro $(HAMUSTRO_SCHEMA)$(HAMUSTRO_HOST):$(HAMUSTRO_PORT)/debug/pprof/profile > $@.pdf

profile/goroutine: profile/
	go tool pprof --pdf hamustro $(HAMUSTRO_SCHEMA)$(HAMUSTRO_HOST):$(HAMUSTRO_PORT)/debug/pprof/goroutine > $@.pdf

profile/heap: profile/
	go tool pprof --pdf hamustro $(HAMUSTRO_SCHEMA)$(HAMUSTRO_HOST):$(HAMUSTRO_PORT)/debug/pprof/heap > $@.pdf

tests/run:
	go test -v ./...

tests/send:
	$(PYC) utils/send_single_message.py $(HAMUSTRO_CONFIG) "$(HAMUSTRO_SCHEMA)$(HAMUSTRO_HOST):$(HAMUSTRO_PORT)/api/v1/track"

tests/stress/1-messages/:
	$(PYC) utils/generate_stress_messages.py $(HAMUSTRO_CONFIG) $@

tests/stress/n-messages/:
	$(PYC) utils/generate_stress_messages.py -r $(HAMUSTRO_CONFIG) $@

tests/stress/%: tests/stress/%-messages/
	cd $< && wrk -t5 -c10 -d1m -s ../run.lua "$(HAMUSTRO_SCHEMA)$(HAMUSTRO_HOST):$(HAMUSTRO_PORT)/api/v1/track"

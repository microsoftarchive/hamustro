ifndef HAMUSTRO_SCHEMA
  HAMUSTRO_SCHEMA:=http://
endif

ifndef HAMUSTRO_HOST
  HAMUSTRO_HOST:=$(warning Please define HAMUSTRO_HOST environment variable, using default "localhost")localhost
endif

ifndef HAMUSTRO_PORT
  HAMUSTRO_PORT:=$(warning Please define HAMUSTRO_PORT environment variable, using "8080")8080
endif

ifndef HAMUSTRO_CONFIG
  HAMUSTRO_CONFIG:=$(warning Please define HAMUSTRO_CONFIG environment variable, using "config/config.json.sample")config/config.json.sample
endif

ifndef PYC
	PYC:=python
endif

all:
	@echo "Please specify a target!"

install/%:
	./utils/installer/_$*.sh

setup:
	$(PYC) utils/setup.py

src/%.go:
src/%/%.go:
src/%/%/%.go:

hamustro: src/payload/payload.pb.go \
          src/*.go \
          src/*/*.go \
          src/*/*/*.go
	go build -o $@ src/*.go

hamustro_linux: src/payload/payload.pb.go \
          src/*.go \
          src/*/*.go \
          src/*/*/*.go
	CGO_ENABLED=0 GOOS=linux go build -o $@ src/*.go

src/payload/payload.pb.go:
	protoc --go_out=. proto/*.proto
	mkdir -p $(dir $@) && mv proto/*.go $@

utils/payload/payload_pb2.py:
	protoc --python_out=. proto/*.proto
	mkdir -p $(dir $@) && mv proto/*.py $@

dev: hamustro utils/payload/payload_pb2.py
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

profile/lines:
	git ls-files | xargs cat | wc -l

tests/run:
	go test -v -cover ./...

tests/send/%:
	$(PYC) utils/send_single_message.py --format "$*" $(HAMUSTRO_CONFIG) "$(HAMUSTRO_SCHEMA)$(HAMUSTRO_HOST):$(HAMUSTRO_PORT)/api/v1/track"

tests/1-messages/:
	$(PYC) utils/generate_stress_messages.py $(HAMUSTRO_CONFIG) $@

tests/n-messages/:
	$(PYC) utils/generate_stress_messages.py -r $(HAMUSTRO_CONFIG) $@

tests/protobuf/%: tests/%-messages/
	cd $< && wrk -t5 -c10 -d1m -s ../protobuf.lua "$(HAMUSTRO_SCHEMA)$(HAMUSTRO_HOST):$(HAMUSTRO_PORT)/api/v1/track"

tests/json/%: tests/%-messages/
	cd $< && wrk -t5 -c10 -d1m -s ../json.lua "$(HAMUSTRO_SCHEMA)$(HAMUSTRO_HOST):$(HAMUSTRO_PORT)/api/v1/track"

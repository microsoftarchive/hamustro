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

all:
	@echo "Please specify a target!"

install/%:
	./install.sh $*

build:
	protoc --go_out=. payload/*.proto
	protoc --python_out=utils payload/*.proto
	go build -o hamustro

run: build
	./hamustro -config $(HAMUSTRO_CONFIG) -verbose

profile/:
	mkdir -p $@

profile/cpu: profile/
	go tool pprof --pdf tivan $(HAMUSTRO_SCHEMA)$(HAMUSTRO_HOST):$(HAMUSTRO_PORT)/debug/pprof/profile > $@.pdf

profile/goroutine: profile/
	go tool pprof --pdf tivan $(HAMUSTRO_SCHEMA)$(HAMUSTRO_HOST):$(HAMUSTRO_PORT)/debug/pprof/goroutine > $@.pdf

profile/heap: profile/
	go tool pprof --pdf tivan $(HAMUSTRO_SCHEMA)$(HAMUSTRO_HOST):$(HAMUSTRO_PORT)/debug/pprof/heap > $@.pdf

tests/send:
	python utils/send_single_message.py $(HAMUSTRO_CONFIG) "$(HAMUSTRO_SCHEMA)$(HAMUSTRO_HOST):$(HAMUSTRO_PORT)/api/v1/track"

tests/stress/1-messages/:
	python utils/generate_stress_messages.py $(HAMUSTRO_CONFIG) $@

tests/stress/n-messages/:
	python utils/generate_stress_messages.py -r $(HAMUSTRO_CONFIG) $@

tests/stress/%: tests/stress/%-messages/
	cd $< && wrk -t5 -c10 -d1m -s ../run.lua "$(HAMUSTRO_SCHEMA)$(HAMUSTRO_HOST):$(HAMUSTRO_PORT)/api/v1/track"

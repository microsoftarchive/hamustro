
# ![Hamustro](docs/logo.png) - the collector of events

[![Travis](https://travis-ci.org/sub-ninja/hamustro.svg)](https://travis-ci.org/sub-ninja/hamustro)

## Overview

This collector meant to be a highly available [RESTful web service](https://github.com/sub-ninja/tivan/blob/master/main.go#L200) that receives events from client devices and secures them agnostic of cloud targets.

The collector is implemented in Go, runs on Ubuntu and OSX.

Events are sent in [Protobuf](https://github.com/sub-ninja/tivan/blob/master/payload/payload.proto) messages.

Currently supported cloud targets are (tested throughput on a c3.xlarge computer with 4vCPU in AWS):

* __Amazon Web Services Simple Notification Service__: 59k events/minute, 70 multi payload requests/s
* __Amazon Web Services Simple Storage Service__: 6.2M events/minute, 8k payload requests/s
* __Microsoft Azure Blob Storage__: 6.05M events/minute, 7.8k multi payload requests/s
* __Microsoft Azure Queue Storage__: 5k events/minute, 5 multi payload requests/s

6Wunderkinder used a similar node.js based service that secured messages in AWS SNS. Based on experiences we've rewritten the app in Go that can handle 20x more requests on equal hardware resources.

Inspired by UNIX philosophy (do one thing and do it well) and [Marcio Castilho's approach](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/).

## Clients

No official client is available at the moment. If you want to write your own please check out our [pseudo client specification](docs/pseudo-client.md).

## Installation

Please install [Go 1.5+](https://golang.org/dl/) and [Python 2.7 or 3.3+](https://www.python.org/downloads/).

```bash
$ sudo make install/go && source ~/.profile # you can install golang with this on OSX/Ubuntu if you need it
$ sudo make install/protobuf # initialize communication format
$ make install/pkg # golang dependencies
$ sudo make install/wrk # install http benchmarking tool
$ make install/utils # utils for development
```

After the package installation, please create your configuration file based on the [sample configuration](config/config.json.sample).

Set up your environment variables.

```bash
export HAMUSTRO_CONFIG="config/yourconfig.json"
export HAMUSTRO_HOST="localhost"
export HAMUSTRO_PORT="8080"
```

## Start collector

You can start the server for development with the following command:
```bash
$ make dev
```

In the _development_ mode it provides useful messages to track what's happening within the collector. Furthermore it notifies the clients with JSON responses on error.

To turn off the notifications and run the collector for production, please use the following command:

```bash
$ make server
```

## Tests

You can run the tests with
```bash
$ make tests/run
```

If you want to start a stress test, please use

```bash
$ make tests/stress/n # 1-25 payloads/request
$ make tests/stress/1 # 1 payload/request
```

During the stress test, you can profile the heap/cpu/goroutine usage easily in _development_ mode, just type

```bash
$ make profile/heap
$ make profile/cpu
$ make profile/goroutine
```

## License

Copyright © 2016, Bence Faludi.

Distributed under the MIT License.

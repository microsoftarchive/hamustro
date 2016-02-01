# Tivan - the collector of events

![](docs/tivan.jpg)

## Overview

This collector meant to be a highly available [RESTful web service](https://github.com/sub-ninja/tivan/blob/master/main.go#L200) that receives events from client devices and secures them agnostic of cloud targets.

The collector is implemented in Go, runs on Ubuntu and OSX.

Events are sent in [Protobuf](https://github.com/sub-ninja/tivan/blob/master/payload/payload.proto) messages.

Currently supported cloud targets are (tested throughput on a c3.xlarge computer with 4vCPU in AWS):

* Amazon Web Services Simple Notification Service (up to 59k events/minute, 70 multi payload requests/s),
* Microsoft Azure Blob Store (up to 2.6M events/minute, 3k multi payload requests/s),
* Microsoft Azure Queue Storage (up to 5k events/minute, 5 multi payload requests/s).

6Wunderkinder used a similar node.js based service that secured messages in AWS SNS. Based on experiences we've rewritten the app in Go that can handle 100x more requests on equal hardware resources.

Inspired by UNIX philosophy (do one thing and do it well) and [Marcio Castilho's approach](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/).

## Clients

No official client is available at the moment. If you want to write your own please check out our [pseudo client specification](docs/pseudo-client.md).

## License

Copyright © 2016, Bence Faludi.

Distributed under the MIT License.
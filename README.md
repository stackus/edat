![](https://github.com/stackus/edat/workflows/CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/stackus/edat)](https://goreportcard.com/report/github.com/stackus/edat)
[![](https://godoc.org/github.com/stackus/edat?status.svg)](https://pkg.go.dev/github.com/stackus/edat)

# edat - Event-Driven Architecture Toolkit

edat is an event-driven architecture library for Go.

## Installation

    go get -u github.com/stackus/edat

## Prerequisites

Go 1.15

## Features

edat provides opinionated plumbing to help with many aspects of the development of an event-driven application.

- Basic pubsub for events
- Asynchronous command and reply messaging
- Event sourcing
- Entity change publication
- Orchestrated sagas
- Transactional Outbox

## Examples

[FTGOGO](https://github.com/stackus/ftgogo) A golang rewrite of the FTGO Eventuate demonstration application using edat.

## TODOs

- Documentation
- Wiki Examples & Quickstart
- Tests, tests, and more tests

## Support Libraries

### Stores

- [edat-pgx](https://github.com/stackus/edat-pgx) Postgres

### Event Streams

- [edat-stan](https://github.com/stackus/edat-stan) NATS Streaming
- [edat-pgx](https://github.com/stackus/edat-pgx) Postgres (outbox store and message producer)

### Marshallers

- [edat-msgpack](https://github.com/stackus/edat-msgpack) MessagePack

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

MIT

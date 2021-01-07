# edat - Event-Driven Architecture Toolkit

edat is an event-driven architecture library for Go.

## Installation

    go get -u github.com/stackus/edat

## Prerequisities

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

[FTOGOGO](https://github.com/stackus/ftgogo) A golang rewrite of the FTGO Eventuate demonstration application using edat.

## TODOs

- Documentation
- Tests, tests, and more tests

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://github.com/stackus/edat/LICENSE)
package msg

import (
	"github.com/stackus/edat/log"
)

// PublisherOption options for PublisherPublisher
type PublisherOption func(*Publisher)

// WithPublisherLogger is an option to set the log.Logger of the Publisher
func WithPublisherLogger(logger log.Logger) PublisherOption {
	return func(publisher *Publisher) {
		publisher.logger = logger
	}
}

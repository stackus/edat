package outbox

import (
	"time"

	"github.com/stackus/edat/log"
	"github.com/stackus/edat/retry"
)

type PollingProcessorOption func(*PollingProcessor)

func WithPollingProcessorMessagesPerPolling(messagesPerPolling int) PollingProcessorOption {
	return func(processor *PollingProcessor) {
		processor.messagesPerPolling = messagesPerPolling
	}
}

func WithPollingProcessorPollingInterval(pollingInterval time.Duration) PollingProcessorOption {
	return func(processor *PollingProcessor) {
		processor.pollingInterval = pollingInterval
	}
}

func WithPollingProcessorRetryer(retryer retry.Retryer) PollingProcessorOption {
	return func(processor *PollingProcessor) {
		processor.retryer = retryer
	}
}

func WithPollingProcessorPurgeOlderThan(purgeOtherThan time.Duration) PollingProcessorOption {
	return func(processor *PollingProcessor) {
		processor.purgeOlderThan = purgeOtherThan
	}
}

func WithPollingProcessorPurgeInterval(purgeInterval time.Duration) PollingProcessorOption {
	return func(processor *PollingProcessor) {
		processor.purgeInterval = purgeInterval
	}
}

func WithPollingProcessorLogger(logger log.Logger) PollingProcessorOption {
	return func(processor *PollingProcessor) {
		processor.logger = logger
	}
}

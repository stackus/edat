package outbox

import (
	"time"

	"github.com/stackus/edat/log"
	"github.com/stackus/edat/retry"
)

// PollingProcessorOption options for PollingProcessor
type PollingProcessorOption func(*PollingProcessor)

// WithPollingProcessorMessagesPerPolling sets the number of messages to fetch for PollingProcessor
func WithPollingProcessorMessagesPerPolling(messagesPerPolling int) PollingProcessorOption {
	return func(processor *PollingProcessor) {
		processor.messagesPerPolling = messagesPerPolling
	}
}

// WithPollingProcessorPollingInterval sets the interval between attempts to fetch new messages for PollingProcessor
func WithPollingProcessorPollingInterval(pollingInterval time.Duration) PollingProcessorOption {
	return func(processor *PollingProcessor) {
		processor.pollingInterval = pollingInterval
	}
}

// WithPollingProcessorRetryer sets the retry strategy for failed calls for PollingProcessor
func WithPollingProcessorRetryer(retryer retry.Retryer) PollingProcessorOption {
	return func(processor *PollingProcessor) {
		processor.retryer = retryer
	}
}

// WithPollingProcessorPurgeOlderThan sets the max age of published messages to purge for PollingProcessor
func WithPollingProcessorPurgeOlderThan(purgeOtherThan time.Duration) PollingProcessorOption {
	return func(processor *PollingProcessor) {
		processor.purgeOlderThan = purgeOtherThan
	}
}

// WithPollingProcessorPurgeInterval sets the interval between attempts to purge published messages for PollingProcessor
func WithPollingProcessorPurgeInterval(purgeInterval time.Duration) PollingProcessorOption {
	return func(processor *PollingProcessor) {
		processor.purgeInterval = purgeInterval
	}
}

// WithPollingProcessorLogger sets the log.Logger for PollingProcessor
func WithPollingProcessorLogger(logger log.Logger) PollingProcessorOption {
	return func(processor *PollingProcessor) {
		processor.logger = logger
	}
}

package outbox

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/stackus/edat/log"
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/retry"
)

type PollingProcessor struct {
	in                 MessageStore
	out                msg.MessagePublisher
	messagesPerPolling int
	pollingInterval    time.Duration
	purgeOlderThan     time.Duration
	purgeInterval      time.Duration
	retryer            retry.Retryer
	logger             log.Logger
	stopping           chan struct{}
	close              sync.Once
}

func NewPollingProcessor(in MessageStore, out msg.MessagePublisher, options ...PollingProcessorOption) *PollingProcessor {
	p := &PollingProcessor{
		in:                 in,
		out:                out,
		messagesPerPolling: DefaultMessagesPerPolling,
		pollingInterval:    DefaultPollingInterval,
		purgeOlderThan:     DefaultPurgeOlderThan,
		purgeInterval:      DefaultPurgeInterval,
		retryer:            DefaultRetryer,
		logger:             log.DefaultLogger,
		stopping:           make(chan struct{}),
	}

	for _, option := range options {
		option(p)
	}

	p.logger.Trace("outbox.PollingProcessor constructed")

	return p
}

func (p *PollingProcessor) Start(ctx context.Context) error {
	cCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	group, gCtx := errgroup.WithContext(cCtx)

	group.Go(func() error {
		if <-p.stopping; true {
			cancel()
		}

		return nil
	})

	group.Go(func() error {
		return p.processMessages(gCtx)
	})

	group.Go(func() error {
		return p.purgePublished(gCtx)
	})

	p.logger.Trace("processor started")

	return group.Wait()
}

func (p *PollingProcessor) Stop(ctx context.Context) (err error) {
	p.close.Do(func() {
		close(p.stopping)

		done := make(chan struct{})
		go func() {
			// anything to wait for?
			close(done)
		}()

		select {
		case <-done:
			p.logger.Trace("done with internal cleanup")
		case <-ctx.Done():
			p.logger.Warn("timed out waiting for internal cleanup to complete")
		}
	})

	return
}

func (p *PollingProcessor) processMessages(ctx context.Context) error {
	pollingTimer := time.NewTimer(0)

	for {
		var err error
		var messages []Message

		err = p.retryer.Retry(ctx, func() error {
			messages, err = p.in.Fetch(ctx, p.messagesPerPolling)
			return err
		})
		if err != nil {
			p.logger.Error("error fetching messages", log.Error(err))
			return err
		}

		if len(messages) > 0 {
			p.logger.Trace("processing messages", log.Int("MessageCount", len(messages)))
			ids := make([]string, 0, len(messages))
			for _, message := range messages {
				var outgoingMsg msg.Message

				logger := p.logger.Sub(
					log.String("MessageID", message.MessageID),
					log.String("DestinationChannel", message.Destination),
				)

				outgoingMsg, err = message.ToMessage()
				if err != nil {
					logger.Error("error with transforming stored message", log.Error(err))
					// TODO this has potential to halt processing; systems need to be in place to fix or address
					return err
				}
				err = p.out.Publish(ctx, outgoingMsg)
				if err != nil {
					logger.Error("error publishing message", log.Error(err))
					// TODO this has potential to halt processing; systems need to be in place to fix or address
					return err
				}
				ids = append(ids, message.MessageID)
			}

			err = p.retryer.Retry(ctx, func() error {
				return p.in.MarkPublished(ctx, ids)
			})
			if err != nil {
				return err
			}

			continue
		}

		if !pollingTimer.Stop() {
			select {
			case <-pollingTimer.C:
			default:
			}
		}

		pollingTimer.Reset(p.pollingInterval)

		select {
		case <-ctx.Done():
			return nil
		case <-pollingTimer.C:
		}
	}
}

func (p *PollingProcessor) purgePublished(ctx context.Context) error {
	purgeTimer := time.NewTimer(0)

	for {
		err := p.retryer.Retry(ctx, func() error {
			return p.in.PurgePublished(ctx, p.purgeOlderThan)
		})
		if err != nil {
			p.logger.Error("error purging published messages", log.Error(err))
			return err
		}

		if !purgeTimer.Stop() {
			select {
			case <-purgeTimer.C:
			default:
			}
		}

		purgeTimer.Reset(p.purgeInterval)

		select {
		case <-ctx.Done():
			return nil
		case <-purgeTimer.C:
		}
	}
}

package saga

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/stackus/edat/core"
	"github.com/stackus/edat/log"
	"github.com/stackus/edat/msg"
)

// Orchestrator orchestrates local and distributed processes
type Orchestrator struct {
	definition    Definition
	instanceStore InstanceStore
	publisher     msg.CommandMessagePublisher
	logger        log.Logger
}

const sagaNotStarted = -1

var _ msg.MessageReceiver = (*Orchestrator)(nil)

// NewOrchestrator constructs a new Orchestrator
func NewOrchestrator(definition Definition, store InstanceStore, publisher msg.CommandMessagePublisher, options ...OrchestratorOption) *Orchestrator {
	o := &Orchestrator{
		definition:    definition,
		instanceStore: store,
		publisher:     publisher,
		logger:        log.DefaultLogger,
	}

	for _, option := range options {
		option(o)
	}

	o.logger.Trace("saga.Orchestrator constructed", log.String("SagaName", definition.SagaName()))

	return o
}

// Start creates a new instance of the saga and begins execution
func (o *Orchestrator) Start(ctx context.Context, sagaData core.SagaData) (*Instance, error) {
	instance := &Instance{
		sagaID:   uuid.New().String(),
		sagaName: o.definition.SagaName(),
		sagaData: sagaData,
	}

	err := o.instanceStore.Save(ctx, instance)
	if err != nil {
		return nil, err
	}

	logger := o.logger.Sub(
		log.String("SagaName", o.definition.SagaName()),
		log.String("SagaID", instance.sagaID),
	)

	logger.Trace("executing saga starting hook")
	o.definition.OnHook(SagaStarting, instance)

	results := o.executeNextStep(ctx, stepContext{step: sagaNotStarted}, sagaData)
	if results.failure != nil {
		logger.Error("error while starting saga orchestration", log.Error(results.failure))
		return nil, err
	}

	err = o.processResults(ctx, instance, results)
	if err != nil {
		logger.Error("error while processing results", log.Error(err))
		return nil, err
	}

	return instance, err
}

// ReplyChannel returns the channel replies are to be received from msg.Subscribers
func (o *Orchestrator) ReplyChannel() string {
	return o.definition.ReplyChannel()
}

// ReceiveMessage implements msg.MessageReceiver.ReceiveMessage
func (o *Orchestrator) ReceiveMessage(ctx context.Context, message msg.Message) error {
	replyName, sagaID, sagaName, err := o.replyMessageInfo(message)
	if err != nil {
		return nil
	}

	if sagaID == "" || (sagaName == "" || sagaName != o.definition.SagaName()) {
		o.logger.Error("cannot process message", log.String("NameValue", sagaName), log.String("IDValue", sagaID))
		return nil
	}

	logger := o.logger.Sub(
		log.String("ReplyName", replyName),
		log.String("SagaName", sagaName),
		log.String("SagaID", sagaID),
		log.String("MessageID", message.ID()),
	)

	logger.Debug("received saga reply message")

	reply, err := core.DeserializeReply(replyName, message.Payload())
	if err != nil {
		// sagas should not be receiving any replies that have not already been registered
		logger.Error("error decoding reply message payload", log.Error(err))
		return nil
	}

	replyMsg := msg.NewReply(reply, message.Headers())

	instance, err := o.instanceStore.Find(ctx, sagaName, sagaID)
	if err != nil {
		logger.Error("failed to locate saga instance data", log.Error(err))
		return nil
	}

	stepCtx := instance.getStepContext()

	results, err := o.handleReply(ctx, stepCtx, instance.SagaData(), replyMsg)
	if err != nil {
		logger.Error("saga reply handler returned an error", log.Error(err))
		return err
	}

	err = o.processResults(ctx, instance, results)
	if err != nil {
		logger.Error("error while processing results", log.Error(err))
		return err
	}

	return nil
}

func (o *Orchestrator) replyMessageInfo(message msg.Message) (string, string, string, error) {
	var err error
	var replyName, sagaID, sagaName string

	replyName, err = message.Headers().GetRequired(msg.MessageReplyName)
	if err != nil {
		o.logger.Error("error reading reply name", log.Error(err))
		return "", "", "", err
	}

	sagaID, err = message.Headers().GetRequired(MessageReplySagaID)
	if err != nil {
		o.logger.Error("error reading saga id", log.Error(err))
		return "", "", "", err
	}

	sagaName, err = message.Headers().GetRequired(MessageReplySagaName)
	if err != nil {
		o.logger.Error("error reading saga name", log.Error(err))
		return "", "", "", err
	}

	return replyName, sagaID, sagaName, nil
}

func (o *Orchestrator) processResults(ctx context.Context, instance *Instance, results *stepResults) error {
	var err error

	logger := o.logger.Sub(
		log.String("SagaName", o.definition.SagaName()),
		log.String("SagaID", instance.sagaID),
	)

	for {
		if results.failure != nil {
			logger.Trace("handling local failure result")
			results, err = o.handleReply(ctx, results.updatedStepContext, results.updatedSagaData, msg.WithFailure())
			if err != nil {
				logger.Error("error handling local failure result", log.Error(err))
				return err
			}
		} else {
			for _, command := range results.commands {
				err = o.publisher.PublishCommand(ctx, o.definition.ReplyChannel(), command, WithSagaInfo(instance))
				if err != nil {
					return err
				}
			}

			instance.updateStepContext(results.updatedStepContext)

			if results.updatedSagaData != nil {
				instance.sagaData = results.updatedSagaData
			}

			if results.updatedStepContext.ended {
				o.processEnd(instance)
			}

			err = o.instanceStore.Update(ctx, instance)
			if err != nil {
				logger.Error("error saving saga instance", log.Error(err))
				return err
			}

			if !results.local {
				logger.Trace("exiting step loop")
				break
			}

			// handle a local success outcome and kick off the next step
			logger.Trace("handling local success result")
			results, err = o.handleReply(ctx, results.updatedStepContext, results.updatedSagaData, msg.WithSuccess())
			if err != nil {
				logger.Error("error handling local success result", log.Error(err))
				return err
			}
		}
	}

	return nil
}

func (o *Orchestrator) processEnd(instance *Instance) {
	logger := o.logger.Sub(
		log.String("SagaName", o.definition.SagaName()),
		log.String("SagaID", instance.sagaID),
	)

	if instance.compensating {
		logger.Trace("executing saga compensated hook")
		o.definition.OnHook(SagaCompensated, instance)
	} else {
		logger.Trace("executing saga completed hook")
		o.definition.OnHook(SagaCompleted, instance)
	}
	logger.Trace("saga has finished all steps")
}

func (o *Orchestrator) handleReply(ctx context.Context, stepCtx stepContext, sagaData core.SagaData, message msg.Reply) (*stepResults, error) {
	replyName := message.Reply().ReplyName()

	logger := o.logger.Sub(
		log.String("SagaName", o.definition.SagaName()),
		log.String("SagaID", message.Headers().Get(MessageReplySagaID)),
		log.String("ReplyName", replyName),
	)

	if stepCtx.step >= len(o.definition.Steps()) || stepCtx.step < 0 {
		logger.Error("current step is out of bounds", log.Int("Step", stepCtx.step))
		return nil, fmt.Errorf("current step is out of bounds: 0-%d, got %d", len(o.definition.Steps()), stepCtx.step)
	}
	step := o.definition.Steps()[stepCtx.step]

	// handle specific replies
	if handler := step.getReplyHandler(replyName, stepCtx.compensating); handler != nil {
		logger.Trace("saga reply handler found")
		err := handler(ctx, sagaData, message.Reply())
		if err != nil {
			logger.Error("saga reply handler returned an error", log.Error(err))
			return nil, err
		}
	}

	outcome, err := message.Headers().GetRequired(msg.MessageReplyOutcome)
	if err != nil {
		logger.Error("error reading reply outcome", log.Error(err))
		return nil, err
	}

	logger.Trace("reply outcome", log.String("Outcome", outcome))

	success := outcome == msg.ReplyOutcomeSuccess

	switch {
	case success:
		logger.Trace("advancing to next step")
		return o.executeNextStep(ctx, stepCtx, sagaData), nil
	case stepCtx.compensating:
		// we're already failing, we can't fail any more
		logger.Error("received a failure outcome while compensating", log.Int("Step", stepCtx.step))
		return nil, fmt.Errorf("received failure outcome while compensating")
	default:
		logger.Trace("compensating to previous step")
		return o.executeNextStep(ctx, stepCtx.compensate(), sagaData), nil
	}
}

func (o *Orchestrator) executeNextStep(ctx context.Context, stepCtx stepContext, sagaData core.SagaData) *stepResults {
	var stepDelta = 1
	var direction = 1
	var step Step

	if stepCtx.compensating {
		direction = -1
	}

	sagaSteps := o.definition.Steps()

	for i := stepCtx.step + direction; i >= 0 && i < len(sagaSteps); i += direction {
		if step = sagaSteps[i]; step.hasInvocableAction(ctx, sagaData, stepCtx.compensating) {
			break
		}

		// Skips steps that have no action for the direction (compensating or not compensating)
		stepDelta++
	}

	// if no step to execute exists the saga is done
	if step == nil {
		return &stepResults{updatedStepContext: stepCtx.end()}
	}

	nextCtx := stepCtx.next(stepDelta)

	results := &stepResults{
		updatedSagaData:    sagaData,
		updatedStepContext: nextCtx,
	}

	step.execute(ctx, sagaData, stepCtx.compensating)(results)

	return results
}

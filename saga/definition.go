package saga

// Definition interface
type Definition interface {
	SagaName() string
	ReplyChannel() string
	Steps() []Step
	OnHook(hook LifecycleHook, instance *Instance)
}

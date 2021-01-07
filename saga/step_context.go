package saga

type stepContext struct {
	step         int
	compensating bool
	ended        bool
}

func (s *stepContext) next(stepIndex int) stepContext {
	if s.compensating {
		return stepContext{step: s.step - stepIndex, compensating: s.compensating}
	}

	return stepContext{step: s.step + stepIndex, compensating: s.compensating}
}

func (s *stepContext) compensate() stepContext {
	return stepContext{step: s.step, compensating: true, ended: s.ended}
}

func (s *stepContext) end() stepContext {
	return stepContext{step: s.step, compensating: s.compensating, ended: true}
}

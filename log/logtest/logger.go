package logtest

import (
	"github.com/stackus/edat/log/logmocks"
)

func MockLogger(setup func(m *logmocks.Logger)) *logmocks.Logger {
	m := &logmocks.Logger{}
	setup(m)
	return m
}

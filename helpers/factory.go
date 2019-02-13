package helpers

import "github.com/cloudfoundry/bosh-utils/logger"

// Logger ...
type Logger struct {
	Logger logger.Logger
	Tag    string
}

// NewLogger ...
func NewLogger(logger logger.Logger, tag string) *Logger {
	return &Logger{
		Logger: logger,
		Tag:    tag,
	}
}

// Debug ...
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.Logger.Debug(l.Tag, msg, args...)
}

// Error ...
func (l *Logger) Error(msg string, args ...interface{}) {
	l.Logger.Error(l.Tag, msg, args...)
}

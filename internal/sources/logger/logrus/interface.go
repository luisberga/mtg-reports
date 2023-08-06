package logrus

type Logger interface {
	Info(args ...interface{})
	WithFields(fields Fields) CustomEntry
	WithError(err error) CustomEntry
	Error(args ...interface{})
	Warn(args ...interface{})
}

type CustomEntry interface {
	Fatal(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Info(args ...interface{})
}

type Fields map[string]interface{}

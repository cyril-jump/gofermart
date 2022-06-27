package logger

type Logger interface {
	Info(msg string)
	Warn(msg string, err error)
	Fatal(msg string, err error)
	Close()
}

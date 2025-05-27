package cclogger

type Level uint8

const (
	Debug Level = iota
	Info
	Warn
	Error
)

type LogField struct {
	Key   string
	Value any
}

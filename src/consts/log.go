package consts

type LogLevel string

const SHOW_LOG = true

const TIME_FORMAT_MS = "2006-01-02 15:04:05.000"

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
	FATAL LogLevel = "FATAL"
)

const (
	NULL_CAN    = -1
	NULL_OFFSET = -1
)

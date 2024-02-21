package logstream

import "fmt"

type LogStream struct {
	Message    string
	ClientId   string
	LoggerName string
}

func (ls *LogStream) String() string {
	return fmt.Sprintf(
		"id: %s, loggername: %s, msg: %s",
		ls.ClientId, ls.LoggerName, ls.Message,
	)
}

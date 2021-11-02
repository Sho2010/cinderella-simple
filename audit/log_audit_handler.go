package audit

import (
	"fmt"
	"io"
)

type LogHandler struct {
	AuditEventHandler
	LogWriter io.Writer
}

func NewLogHandler(w io.Writer) *LogHandler {
	h := LogHandler{}
	h.LogWriter = w
	return &h
}

func (h *LogHandler) Start(ch <-chan AuditEvent) {
	fmt.Println("Start LogHandler")

	for e := range ch {
		fmt.Println("[LogHandler] Audit event received and write data")
		fmt.Fprintf(h.LogWriter, "Dump Message:%s\n", e.GetMessage())
	}
}

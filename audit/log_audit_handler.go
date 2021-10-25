package audit

import (
	"fmt"
	"io"
)

type LogHandler struct {
	LogWriter io.Writer
}

func (h *LogHandler) Start(event <-chan AuditEvent) {
	fmt.Println("Start handler")

	for e := range event {
		fmt.Println("Audit event received and dump event data")
		fmt.Printf("Dump description:%s\n", e.GetMessage())
	}
}

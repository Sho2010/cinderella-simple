package k8s

import (
	"time"

	"github.com/Sho2010/cinderella-simple/audit"
)

type CleanupEvent struct {
	deleteTime  time.Time
	description string
}

func (e *CleanupEvent) GetMessage() string {
	return e.description
}

// func (e *CleanupEvent) GetType() {
// 	return e.description
// }

func (e *CleanupEvent) GetTime() time.Time {
	return e.deleteTime
}

func RaiseCleanupEvent(description string) {

	e := &CleanupEvent{
		deleteTime:  time.Now(),
		description: description,
	}

	audit.AuditCh <- e
}

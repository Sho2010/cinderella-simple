package k8s

import (
	"time"

	"github.com/Sho2010/cinderella-simple/audit"
)

type CleanupEvent struct {
	deleteTime  time.Time
	description string
}

// Verify interface compliance
var _ audit.AuditEvent = (*CleanupEvent)(nil)

func (e *CleanupEvent) GetMessage() string {
	return e.description
}

// func (e *CleanupEvent) GetType() {
// 	return e.description
// }

func (e *CleanupEvent) EventAt() time.Time {
	return e.deleteTime
}

func RaiseCleanupEvent(description string) {

	e := &CleanupEvent{
		deleteTime:  time.Now(),
		description: description,
	}

	audit.PublishEvent(e)
}

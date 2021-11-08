package k8s

import (
	"time"

	"github.com/Sho2010/cinderella-simple/audit"
)

type ResourceCreateEvent struct {
	createAt time.Time
	message  string
}

// Verify interface compliance
var _ audit.AuditEvent = (*ResourceCreateEvent)(nil)

func (e *ResourceCreateEvent) GetMessage() string {
	return e.message
}

func GetType() string {
	return "ResourceCreate"
}

func (e *ResourceCreateEvent) EventAt() time.Time {
	return e.createAt
}

func RaiseResourceCreateEvent(message string) {
	e := &ResourceCreateEvent{
		createAt: time.Time{},
		message:  message,
	}
	audit.PublishEvent(e)
}

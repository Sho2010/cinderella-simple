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
var _ audit.AuditEvent = &ResourceCreateEvent{}

func (e *ResourceCreateEvent) GetMessage() string {
	return e.message
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

func (e *ResourceCreateEvent) GetType() audit.AuditType {
	return audit.AuditTypeCreateKubernetesResource
}

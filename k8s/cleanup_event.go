package k8s

import (
	"fmt"
	"strings"
	"time"

	"github.com/Sho2010/cinderella-simple/audit"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CleanupEvent struct {
	cleanupAt      time.Time
	deletedObjects []metav1.Object
	errs           deleteErrors
}

// Verify interface compliance
var _ audit.AuditEvent = &CleanupEvent{}

func (e *CleanupEvent) GetMessage() string {

	var sb strings.Builder
	if len(e.deletedObjects) > 0 {
		fmt.Fprintf(&sb, "[%s] Deleted expired %d objects:", e.GetType(), len(e.deletedObjects))
		for _, obj := range e.deletedObjects {
			fmt.Fprintf(&sb, " %s/%s", obj.GetNamespace(), obj.GetName())
		}
		fmt.Fprint(&sb, "\n")
	}

	if len(e.errs) > 0 {
		fmt.Fprintf(&sb, "[%s] Failed to delete expired objects:", e.GetType())
		for _, err := range e.errs {
			fmt.Fprintf(&sb, " %s/%s: %s", err.target.GetNamespace(), err.target.GetName(), err.err.Error())
		}
		fmt.Fprint(&sb, "\n")
	}

	return sb.String()
}

func (e *CleanupEvent) GetType() audit.AuditType {
	return audit.AuditTypeCleanup
}

func (e *CleanupEvent) EventAt() time.Time {
	return e.cleanupAt
}

func (e *CleanupEvent) String() string {
	return e.GetMessage()
}

func newCleanupEvent(deletedObjects []metav1.Object, errs deleteErrors) *CleanupEvent {
	fmt.Printf("%#v", deletedObjects)

	return &CleanupEvent{
		cleanupAt:      time.Now(),
		deletedObjects: deletedObjects,
		errs:           errs,
	}
}

func publishCleanupEvent(deletedObjects []metav1.Object, errs deleteErrors) {
	if len(deletedObjects) > 0 || len(errs) > 0 {
		audit.PublishEvent(newCleanupEvent(deletedObjects, errs))
	}
}

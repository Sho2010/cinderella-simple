package claim

import (
	"time"
)

type ClaimRegisterEvent struct {
	message string
	eventAt time.Time
}

func (e *ClaimRegisterEvent) GetMessage() string {
	return e.message
}
func (e *ClaimRegisterEvent) EventAt() time.Time {
	return e.eventAt
}

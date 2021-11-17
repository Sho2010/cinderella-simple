package model

import (
	"time"

	"github.com/Sho2010/cinderella-simple/audit"
)

type ClaimRegisterEvent struct {
	message string
	eventAt time.Time
	subject string
}

// Verify interface compliance
var _ audit.AuditEvent = &ClaimRegisterEvent{}

func (e *ClaimRegisterEvent) GetMessage() string {
	return e.message
}

func (e *ClaimRegisterEvent) EventAt() time.Time {
	return e.eventAt
}

func (e *ClaimRegisterEvent) GetType() audit.AuditType {
	return audit.AuditTypeRegisterClaim
}

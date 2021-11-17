package event

import "time"

type ClaimEventType string

var _claimEventCh = make(chan ClaimEvent)

const (
	ClaimEventCreated  ClaimEventType = "Created"
	ClaimEventAceepted ClaimEventType = "Accepted"
	ClaimEventRejected ClaimEventType = "Rejected"
	ClaimEventExpired  ClaimEventType = "Expired"
)

type ClaimEvent interface {
	Subject() string
	EventAt() time.Time
	EventType() ClaimEventType
}

type ClaimEventImpl struct {
	subject string
	eventAt time.Time
	event   ClaimEventType
}

func NewClaimEvent(subject string, event ClaimEventType) ClaimEvent {
	return &ClaimEventImpl{
		subject: subject,
		eventAt: time.Now(),
		event:   event,
	}
}

func (e ClaimEventImpl) Subject() string {
	return e.subject
}

func (e ClaimEventImpl) EventAt() time.Time {
	return e.eventAt
}

func (e ClaimEventImpl) EventType() ClaimEventType {
	return e.event
}

func ClaimEventChannel() <-chan ClaimEvent {
	return _claimEventCh
}

func PublishClaimEvent(event ClaimEvent) {
	go func() {
		_claimEventCh <- event
	}()
}

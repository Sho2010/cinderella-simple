package slack

import (
	"fmt"

	"github.com/Sho2010/cinderella-simple/domain/event"
	"github.com/slack-go/slack"
)

type ClaimEventHandler struct {
	ch  <-chan event.ClaimEvent
	api *slack.Client
}

func NewClaimEventHandler(ch <-chan event.ClaimEvent, api *slack.Client) *ClaimEventHandler {
	return &ClaimEventHandler{
		ch:  ch,
		api: api,
	}
}

func (h *ClaimEventHandler) Start() {
	go h.handle()
}

func (h *ClaimEventHandler) handle() {
	//必要であればctx使って途中で止められるようにする
	for ev := range h.ch {
		switch ev.EventType() {
		case event.ClaimEventCreated:
			h.created(ev)
		case event.ClaimEventAceepted:
			h.accepted(ev)
		case event.ClaimEventRejected:
			h.rejected(ev)
		case event.ClaimEventExpired:
			h.expired(ev)
		}
	}
}

func (h *ClaimEventHandler) created(ev event.ClaimEvent) {
	// do nothing
}

func (h *ClaimEventHandler) accepted(ev event.ClaimEvent) {
	h.postMessage(ev.Subject(), "claim accepted")
}

func (h *ClaimEventHandler) rejected(ev event.ClaimEvent) {
	h.postMessage(ev.Subject(), "claim rejected")
}

func (h *ClaimEventHandler) expired(ev event.ClaimEvent) {
	h.postMessage(ev.Subject(), "claim expired")

}

func (h *ClaimEventHandler) postMessage(channel string, text string) {
	_, _, err := h.api.PostMessage(channel, slack.MsgOptionText(text, false))
	if err != nil {
		fmt.Printf("failed posting message: %v", err)
	}
}

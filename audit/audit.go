package audit

import (
	"fmt"
	"os"
	"reflect"
	"time"
)

// Global Channel
var auditCh = make(chan AuditEvent)

type AuditEvent interface {
	GetMessage() string
	GetTime() time.Time
	// GetType()
}

// AuditEventHandler can processing received audit event
type AuditEventHandler interface {
	// Start audit event handler loop
	Start(ch <-chan AuditEvent)
}

// AuditBroadcaster is only reciever of AuditEvent channel
// this is broadcast of the event copy to EventHandlers when audit event received
type AuditBroadcaster struct {
	EventHandlers []AuditEventHandler
}

// PublishEvent method is the only way to publish audit event
func PublishEvent(e AuditEvent) {
	auditCh <- e
}

func (b *AuditBroadcaster) Start() {
	fmt.Println("Start audit event broadcasting")

	//debug
	b.testInit()

	b.initializeTeeChannel(auditCh)

}

func (b *AuditBroadcaster) testInit() {
	b.EventHandlers = []AuditEventHandler{
		NewLogHandler(os.Stdout),
		NewLogHandler(os.Stderr),
		NewSlackHandler(),
	}
}

// NOTE: 可変長channelに対しては
// seletc { case } によるchannel処理ができないためreflect.SelectCaseを使う
// See: https://zenn.dev/imamura_sh/articles/select-arbitary-number-of-channels
//      https://github.com/eapache/channels/blob/master/channels.go#L120-L140
//
// Recieve側のChannelの数が決まってるなら「Go言語による並行処理」に書いてある tee実装で良さそう

func (b *AuditBroadcaster) initializeTeeChannel(in <-chan AuditEvent) {

	cases := make([]reflect.SelectCase, len(b.EventHandlers))
	chs := make([]chan AuditEvent, len(b.EventHandlers))

	for i := range cases {
		cases[i].Dir = reflect.SelectSend
		chs[i] = make(chan AuditEvent)
		go b.EventHandlers[i].Start(chs[i])
	}

	go func() {
		for event := range in {
			for i := range cases {
				// select { case ch1 < event } と同義
				cases[i].Chan = reflect.ValueOf(chs[i])
				cases[i].Send = reflect.ValueOf(event)
			}
			for _ = range cases {
				chosen, _, _ := reflect.Select(cases)
				cases[chosen].Chan = reflect.ValueOf(nil)
			}
		}
	}()
}


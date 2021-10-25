package audit

import (
	"fmt"
	"time"
)

// Global Channel
var AuditCh = make(chan AuditEvent)

type AuditEvent interface {
	GetMessage() string
	GetTime() time.Time
	// GetType()
}

// AuditEventHandler can processing received audit event
type AuditEventHandler interface {
	// Start audit event handler loop
	Start(event <-chan AuditEvent)
}

// AuditBroadcaster is only reciever of AuditEvent channel
// this is broadcast of the event copy to EventHandlers when audit event received
type AuditBroadcaster struct {
	EventHandlers []AuditEventHandler
}

// Start is
func (b *AuditBroadcaster) Start(event <-chan AuditEvent) {
	fmt.Println("Start audit event broadcasting")

	for e := range event {
		fmt.Printf("Recieve audit event:%s\n", e.GetMessage())

		//TODO: tee実装で他のevent handler channelにブロードキャストする
	}
}

// 書籍:Go言語による並行処理 tee実装
// func Tee(done <-chan struct{}, in <-chan interface{}) (<-chan interface{}, <-chan interface{}) {
//     out1 := make(chan interface{})
//     out2 := make(chan interface{})
//
//     go func() {
//         defer close(out1)
//         defer close(out2)
//
//         for v := range OrDone(done, in) {
//             var ch1, ch2 = out1, out2
//             for i := 0; i < 2; i++ {
//                 select {
//                 case ch1 <- v:
//                     ch1 = nil
//                 case ch2 <- v:
//                     ch2 = nil
//                 }
//             }
//         }
//     }()
//
//     return out1, out2
// }

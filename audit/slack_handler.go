package audit

import (
	"fmt"
	"log"
	"os"

	"github.com/slack-go/slack"
)

type SlackHandler struct {
	AuditEventHandler
	client *slack.Client
}

func NewSlackHandler() *SlackHandler {
	client := slack.New(
		os.Getenv("SLACK_BOT_TOKEN"),
		slack.OptionAppLevelToken(os.Getenv("SLACK_APP_TOKEN")),
	)
	return &SlackHandler{
		client: client,
	}
}

func (h *SlackHandler) Start(event <-chan AuditEvent) {

	fmt.Println("Start SlackHandler")

	for e := range event {
		log.Println("[SlackHandler]Audit event received and post message to slack")

		//TODO: メッセージの整形
		_, _, err := h.client.PostMessage("#bot-test", slack.MsgOptionText(e.GetMessage(), true))
		if err != nil {
			//TODO: error handle
			panic(err)
		}

	}
}

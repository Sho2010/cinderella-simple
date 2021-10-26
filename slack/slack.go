package slack

import (
	"fmt"
	"log"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type Slack struct {
	Api    *slack.Client
	Socket *socketmode.Client
}

func NewSlack() *Slack {
	api := slack.New(
		os.Getenv("SLACK_BOT_TOKEN"),
		slack.OptionAppLevelToken(os.Getenv("SLACK_APP_TOKEN")),
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)

	socketMode := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "sm: ", log.Lshortfile|log.LstdFlags)),
	)
	return &Slack{
		Api:    api,
		Socket: socketMode,
	}
}

func (s *Slack) Start() {
	go func() {
		for evt := range s.Socket.Events {
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				fmt.Println("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				fmt.Println("Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				fmt.Println("Connected to Slack with Socket Mode.")
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)

					continue
				}

				fmt.Printf("Event received: %+v\n", eventsAPIEvent)

				s.Socket.Ack(*evt.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.AppMentionEvent:
						_, _, err := s.Api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
						if err != nil {
							fmt.Printf("failed posting message: %v", err)
						}
					case *slackevents.MemberJoinedChannelEvent:
						fmt.Printf("user %q joined to channel %q", ev.User, ev.Channel)

					case *slackevents.AppHomeOpenedEvent:
						println("----------------------------------------AppHomeOpened")
						s.appHomeOpenedHandler(ev)
					}

				default:
					s.Socket.Debugf("unsupported Events API event received")
				}
			case socketmode.EventTypeInteractive:
				callback, ok := evt.Data.(slack.InteractionCallback)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)

					continue
				}

				//DEBUG
				// fmt.Printf("Interaction received: %+v\n", callback)

				var payload interface{}

				switch callback.Type {
				case slack.InteractionTypeBlockActions:
					// See https://api.slack.com/apis/connections/socket-implement#button
					s.blockActionsHandler(callback)
					s.Socket.Debugf("button clicked!")
				case slack.InteractionTypeViewSubmission:
					// See https://api.slack.com/apis/connections/socket-implement#modal
					// modalに対してレスポンスを返す時のイベント

				case slack.InteractionTypeShortcut:
				case slack.InteractionTypeDialogSubmission:
				default:

				}

				s.Socket.Ack(*evt.Request, payload)
			case socketmode.EventTypeSlashCommand:
				cmd, ok := evt.Data.(slack.SlashCommand)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)

					continue
				}

				s.Socket.Debugf("Slash command received: %+v", cmd)

				payload := map[string]interface{}{
					"blocks": []slack.Block{
						slack.NewSectionBlock(
							&slack.TextBlockObject{
								Type: slack.MarkdownType,
								Text: "foo",
							},
							nil,
							slack.NewAccessory(
								slack.NewButtonBlockElement(
									"",
									"somevalue",
									&slack.TextBlockObject{
										Type: slack.PlainTextType,
										Text: "bar",
									},
								),
							),
						),
					}}

				s.Socket.Ack(*evt.Request, payload)
			default:
				fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
			}
		}
	}()

	s.Socket.Run()
}

func (s *Slack) appHomeOpenedHandler(e *slackevents.AppHomeOpenedEvent) {
	// NOTE: hashに関して
	// A string that represents view state to protect against possible race conditions.

	// NOTE: CallbackID
	// 表示したViewの識別子 以下のようにする予定
	// cinderella_home_general
	// cinderella_home_admin

	// NOTE: ExternalID
	// こちらは種類ごとではなくて、表示されたModal1つ1つの識別コードのイメージ。Workspace内で一意にしておく必要があります。
	//
	// おすすめは、ユーザー名+タイムスタンプ。
	// これだと、よほどのことがない限り被らず安心です。

	blocks, err := BuildHomeView()
	if err != nil {
		panic(err)
	}

	_, err = s.Api.PublishView(e.User,
		slack.HomeTabViewRequest{
			Type:            slack.VTHomeTab,
			Blocks:          *blocks,
			PrivateMetadata: "",
			CallbackID:      "cinderella_home_general_02",
			ExternalID:      "cinderella_home_general_dakfjda",
		}, "")

	if err != nil {
		panic(err)
	}
}

func (s *Slack) blockActionsHandler(callback slack.InteractionCallback) {
	fmt.Printf("ActionID:%s\n", callback.ActionID)
	fmt.Printf("TriggerID:%s\n", callback.TriggerID)

	for _, v := range callback.ActionCallback.AttachmentActions {
		fmt.Println("---AttachmentActions element")
		fmt.Printf("%#v\n", v.Name)
	}

	for _, v := range callback.ActionCallback.BlockActions {
		fmt.Println("---Block Actions element")
		fmt.Printf("BlockID: %#v\n", v.BlockID)
		fmt.Printf("ActionID: %#v\n", v.ActionID)
		fmt.Printf("Text: %#v\n", v.Text.Text)
		fmt.Printf("Value: %#v\n", v.Value)
		fmt.Println("---")
	}

	s.Show(callback.User.ID, callback.TriggerID)

}

//TODO
func (s *Slack) Show(userID, triggerID string) {

	blocks, err := BuildClaimModalView()
	if err != nil {
		panic(err)
	}
	fmt.Println(blocks)

	modal := slack.ModalViewRequest{
		Type:       slack.VTModal,
		Title:      slack.NewTextBlockObject("plain_text", "権限申請", false, false),
		Blocks:     *blocks,
		Close:      slack.NewTextBlockObject("plain_text", "close", false, false),
		Submit:     slack.NewTextBlockObject("plain_text", "submit", false, false),
		CallbackID: "cinderella_claim",
		ExternalID: "cinderella_home_general_adfdaeda",
	}

	r, err := s.Api.OpenView(triggerID, modal)
	if err != nil {
		log.Printf("とりあえずデバッグの為握りつぶす %v ", err)
	}
	println(r)
}

// NOTE;
// validation errorとかでModalの更新が必要な場合
// s.Api.UpdateView()
// modalから更にモーダルを呼ぶ
// s.Api.PublishView

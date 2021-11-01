package slack

import (
	"fmt"
	"log"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// SlackApp App
type SlackApp struct {
	Api    *slack.Client
	Socket *socketmode.Client
	// Administrators are slack app administrators. they can accept/reject permission claim.
	administrators []string
	ClaimManager
}

const (
	ViewHomeCallbackID        = "cinderella_home"
	ViewClaimShowCallbackID   = "cinderella_claim_show"
	ViewClaimSubmitCallbackID = "cinderella_claim_submit"
	ViewRejectCallbackID      = "cinderella_claim_reject"
)

func NewSlackApp(administrators []string) *SlackApp {
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

	s := SlackApp{
		Api:            api,
		Socket:         socketMode,
		administrators: administrators,
		ClaimManager:   ClaimManager{},
	}

	return &s
}

func (s *SlackApp) Start() {
	go func() {
		for evt := range s.Socket.Events {
			switch evt.Type {
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
				case slack.InteractionTypeViewSubmission:
					// See https://api.slack.com/apis/connections/socket-implement#modal
					s.viewSubmissionHandler(callback)
				case slack.InteractionTypeShortcut:
				case slack.InteractionTypeDialogSubmission:
				// Deprecated
				default:

				}

				s.Socket.Ack(*evt.Request, payload)

			case socketmode.EventTypeConnecting:
				fmt.Println("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				fmt.Println("Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				fmt.Println("Connected to Slack with Socket Mode.")
			case socketmode.EventTypeHello:
				fmt.Println("Hello event received.")
			default:
				fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
			}
		}
	}()

	s.Socket.Run()
}

func (s *SlackApp) appHomeOpenedHandler(e *slackevents.AppHomeOpenedEvent) {
	c := HomeController{
		slack: s.Api,
	}
	c.Show(e)
}

func (s *SlackApp) blockActionsHandler(callback slack.InteractionCallback) {

	s.Socket.Debugf("button clicked!")

	dumpInteractionCallback(callback)

	for _, v := range callback.ActionCallback.BlockActions {
		switch v.ActionID {
		case "create_claim":
			c := ClaimController{
				Slack: s,
			}
			c.Show(callback.User.ID, callback.TriggerID)
		case "create_kubeconfig":
			c := KubeconfigController{
				slack: s.Api,
			}

			claim, err := s.FindClaim(callback.User.ID)
			if err != nil {
				if err := c.CallbackClaimNotFound(callback.User.ID); err != nil {
					panic(err)
				}
				return
			}

			if err := c.SendSlackDM(*claim); err != nil {
				panic(err)
			}

		// case "open_settings":
		// case "claim_details":
		// case "claim_reject":
		default:
			s.Socket.Debugf("unsupported block action: %s", v.ActionID)
		}
	}
}

func (s *SlackApp) viewSubmissionHandler(callback slack.InteractionCallback) {
	// InteractionTypeViewSubmission

	s.Socket.Debugf("interaction event received!")
	dumpInteractionCallback(callback)

	switch callback.View.CallbackID {
	case "cinderella_claim_show":
		s.Socket.Debugf("Close claim modal views")

		c := ClaimController{
			Slack: s,
		}
		claim, err := c.Create(callback)
		if err != nil {
			return
		}
		s.ClaimManager.AddClaim(claim)

	default:
		s.Socket.Debugf("unsupported callbackID: %s", callback.View.CallbackID)
	}
}

func (s *SlackApp) IsAdmin(slackID string) bool {
	for _, v := range s.administrators {
		if v == slackID {
			return true
		}
	}
	return false
}

// NOTE;

// views.update、views.push API メソッドはモーダル内での "block_actions" リクエストを受信したときに使用するものであり、
// "view_submission" 時にモーダルを操作するための API ではありません
//
// validation errorとかでModalの更新が必要な場合
// s.Api.UpdateView()
//
// modalから更にモーダルを呼ぶ
// s.Api.PublishView

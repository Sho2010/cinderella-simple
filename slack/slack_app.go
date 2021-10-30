package slack

import (
	"fmt"
	"log"
	"os"

	"github.com/Sho2010/cinderella-simple/claim"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// SlackApp App
type SlackApp struct {
	Api    *slack.Client
	Socket *socketmode.Client
	Claims []claim.SlackClaim
	// Administrators are slack app administrators. they can accept/reject permission claim.
	administrators []string
}

const (
	ViewHomeCallbackID   = "cinderella_home"
	ViewClaimCallbackID  = "cinderella_claim_modal"
	ViewRejectCallbackID = "cinderella_claim_reject"
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
	}

	return &s
}

func (s *SlackApp) Start() {
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
			claim, err := s.findClaim(callback.User.ID)
			if err != nil {
				if err := c.CallbackClaimNotFound(callback.User.ID); err != nil {
					panic(err)
				}
				fmt.Println(err)
				return
			}
			c.SendSlackDM(claim)
		// case "open_settings":
		// case "claim_details":
		// case "claim_reject":
		default:
			s.Socket.Debugf("unsupported block action: %s", v.ActionID)

		}
	}
}

func (s *SlackApp) findClaim(userId string) (claim.SlackClaim, error) {
	for _, claim := range s.Claims {
		if claim.SlackUser.ID == userId {
			return claim, nil
		}
	}
	return claim.SlackClaim{}, fmt.Errorf("Could not find claim")
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
// validation errorとかでModalの更新が必要な場合
// s.Api.UpdateView()
// modalから更にモーダルを呼ぶ
// s.Api.PublishView

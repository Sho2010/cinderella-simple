package slack

import (
	"log"

	"github.com/slack-go/slack"
)

type ClaimController struct {
	Slack *SlackApp
}

func (c *ClaimController) Show(userID, triggerID string) {
	blocks, err := BuildClaimModalView()
	if err != nil {
		panic(err)
	}

	modal := slack.ModalViewRequest{
		Type:       slack.VTModal,
		Title:      slack.NewTextBlockObject("plain_text", "権限申請", false, false),
		Blocks:     *blocks,
		Close:      slack.NewTextBlockObject("plain_text", "close", false, false),
		Submit:     slack.NewTextBlockObject("plain_text", "submit", false, false),
		CallbackID: ViewClaimCallbackID,
		ExternalID: c.generateExternalID(),
	}

	r, err := c.Slack.Api.OpenView(triggerID, modal)
	if err != nil {
		//TODO error handling
		log.Printf("とりあえずデバッグの為握りつぶす %v ", err)
	}
	println(r)
}

func (c *ClaimController) Create() {

}

func (c *ClaimController) generateExternalID() string {
	return generateExternalID(ViewClaimCallbackID)
}


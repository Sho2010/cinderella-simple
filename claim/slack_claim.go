package claim

import (
	"fmt"

	"github.com/slack-go/slack"
)

type SlackClaim struct {
	Claim `json:"claim"`
	User  slack.User `json:"slack_user"`
}

func (c *SlackClaim) GetLabels() map[string]string {
	//TODO: implement me
	return make(map[string]string)
}

func (c *SlackClaim) GetAnnotations() map[string]string {
	//TODO: implement me
	return map[string]string{
		"": "",
	}
}

func (c *SlackClaim) GetSubject() string {
	return c.User.ID
}

func (c *SlackClaim) GetName() string {
	return c.User.Name
}

func (c *SlackClaim) GetEmail() string {
	return c.User.Profile.Email
}

func (c *SlackClaim) ToBlock() (slack.Block, error) {
	//  "text": "Claimer: *<@U04L97CP5>*\nPeriod: *30min*\nClaim Date: *2021/10/26 12:00:00*\nNamespace: *awesome*\nShort description: デバッグしたい"
	text := fmt.Sprintf("Claimer: *%s*\nPeriod: *%s*\nClaim Date: *%s*\nNamespace: *%s*\nShort description: %s",
		c.User.ID,
		"30min", // TODO: implement me
		c.GetClaimDate().Format("2006/01/02 15:04:05"),
		fmt.Sprintf("%+q", c.GetNamespaces()),
		c.GetDescription(),
	)

	block := slack.NewSectionBlock(
		&slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: text,
		},
		nil,
		slack.NewAccessory(
			slack.NewImageBlockElement(
				"https://api.slack.com/img/blocks/bkb_template_images/creditcard.png", //TODO: いい感じのアイコン
				"claim_image",
			),
		),
	)

	return block, nil
}

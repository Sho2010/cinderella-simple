package claim

import (
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

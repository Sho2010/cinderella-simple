package slack

import (
	"fmt"

	"github.com/Sho2010/cinderella-simple/claim"
	"github.com/Sho2010/cinderella-simple/encrypt"
	"github.com/Sho2010/cinderella-simple/k8s"
	"github.com/sethvargo/go-password/password"
	"github.com/slack-go/slack"
)

type KubeconfigController struct {
	slack *slack.Client
}

func (c *KubeconfigController) CallbackClaimNotFound(channelId string) error {
	_, _, err := c.slack.PostMessage(channelId, slack.MsgOptionText("権限要求が見つかりませんでした。申し訳ございませんがもう一度申請を行い、改善しないようであれば管理者に連絡してください。", false))
	if err != nil {
		return fmt.Errorf("slack PostMessage failed, %w", err)
	}
	return nil
}

func (c *KubeconfigController) Create(claim.Claim) {
}

func (c *KubeconfigController) SendSlackDM(claim claim.SlackClaim) error {

	if claim.EncryptType == encrypt.EncryptTypeZip {
		passwd, err := password.Generate(32, 10, 0, false, false)
		if err != nil {
			panic(err)
		}
		claim.ZipPassword = passwd
	}

	filePath, err := k8s.CreateEncryptedFile(claim.Claim)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	//TODO: わかりやすいメッセージ

	f, err := c.slack.UploadFile(
		slack.FileUploadParameters{
			File: filePath,
			// Reader:   file,
			Filename:       "kubeconfig.zip",
			Channels:       []string{claim.SlackUser.ID},
			Filetype:       string(claim.EncryptType),
			Title:          "kubeconfig",
			InitialComment: fmt.Sprintf("password: `%s`", claim.ZipPassword), //FIXME: Zip password 前提のメッセージを返してしまっている
		})

	if err != nil {
		println(f)
		return fmt.Errorf("Send slack DM fail: %w", err)
	}
	return nil
}

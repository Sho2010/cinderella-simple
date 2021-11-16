package slack

import (
	"fmt"

	"github.com/Sho2010/cinderella-simple/domain/repository"
	"github.com/Sho2010/cinderella-simple/encrypt"
	"github.com/Sho2010/cinderella-simple/k8s"
	"github.com/sethvargo/go-password/password"
	"github.com/slack-go/slack"
)

type KubeconfigController struct {
	slack           *slack.Client
	claimRepository repository.ClaimRepository
}

func NewKubeconfigController(slack *slack.Client, claimRepository repository.ClaimRepository) *KubeconfigController {
	return &KubeconfigController{
		slack:           slack,
		claimRepository: claimRepository,
	}
}

func (c *KubeconfigController) CallbackClaimNotFound(channelId string) error {
	_, _, err := c.slack.PostMessage(channelId, slack.MsgOptionText("権限要求が見つかりませんでした。申し訳ございませんがもう一度申請を行い、改善しないようであれば管理者に連絡してください。", false))
	if err != nil {
		return fmt.Errorf("slack PostMessage failed, %w", err)
	}
	return nil
}

func (c *KubeconfigController) SendSlackDM(subject string) error {
	claim, err := c.claimRepository.FindBySubject("")
	if err != nil {
		return err
	}

	if claim.GetEncryptType() == encrypt.EncryptTypeZip && claim.ZipEncryptOption.ZipPassword == "" {
		passwd, err := password.Generate(32, 10, 0, false, false)
		if err != nil {
			return err
		}
		claim.ZipEncryptOption.ZipPassword = passwd
	}

	//TODO slack packageをk8Sに依存させない　適切なlayer
	filePath, err := k8s.CreateEncryptedFile(*claim)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	//TODO: わかりやすいメッセージ

	f, err := c.slack.UploadFile(
		slack.FileUploadParameters{
			File: filePath,
			// Reader:   file,
			Filename:       "kubeconfig.zip",
			Channels:       []string{claim.GetSubject()},
			Filetype:       string(claim.GetEncryptType()),
			Title:          "kubeconfig",
			InitialComment: fmt.Sprintf("password: `%s`", claim.GetZipPassword()), //FIXME: Zip password 前提のメッセージを返してしまっている
		})

	if err != nil {
		println(f)
		return fmt.Errorf("Send slack DM fail: %w", err)
	}
	return nil
}

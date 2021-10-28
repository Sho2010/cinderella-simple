package slack

import (
	"fmt"
	"os"

	"github.com/Sho2010/cinderella-simple/claim"
	"github.com/Sho2010/cinderella-simple/k8s"
	"github.com/sethvargo/go-password/password"
	"github.com/slack-go/slack"
)

type KubeconfigController struct {
	//TODO private
	Slack *slack.Client
}

func (c *KubeconfigController) Create(claim.Claim) {
}

func (c *KubeconfigController) SendSlackDM(claim claim.SlackClaim) {
	passwd, err := password.Generate(32, 10, 0, false, false)
	if err != nil {
		panic(err)
	}
	claim.ZipPassword = passwd

	tmpFile, _ := os.CreateTemp("", "kubeconfig")
	fmt.Println(tmpFile.Name())
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	k8s.CreateEncryptedFile(tmpFile, claim.Claim)

	//TODO: わかりやすいメッセージ

	f, err := c.Slack.UploadFile(
		slack.FileUploadParameters{
			File: tmpFile.Name(),
			// Reader:   file,
			Filename:       "kubeconfig.zip",
			Channels:       []string{claim.SlackUser.ID},
			Filetype:       string(claim.EncryptType),
			Title:          "kubeconfig",
			InitialComment: fmt.Sprintf("password: `%s`", passwd),
		})

	if err != nil {
		panic(err)
	}

	fmt.Println(f)
}

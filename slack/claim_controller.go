package slack

import (
	"fmt"
	"time"

	"github.com/Sho2010/cinderella-simple/claim"
	"github.com/Sho2010/cinderella-simple/encrypt"
	"github.com/slack-go/slack"
)

type ClaimController struct {
	Slack *SlackApp
}

func (c *ClaimController) Show(userID, triggerID string) error {

	if claim := claim.FindClaim(userID); claim != nil {
		//TODO: 既に申請済みの場合エラーを返す
	}

	blocks, err := BuildClaimModalView()
	if err != nil {
		return fmt.Errorf("BuildClaimModalView error: %w", err)
	}

	modal := slack.ModalViewRequest{
		Type:       slack.VTModal,
		Title:      slack.NewTextBlockObject("plain_text", "権限申請", false, false),
		Blocks:     *blocks,
		Close:      slack.NewTextBlockObject("plain_text", "close", false, false),
		Submit:     slack.NewTextBlockObject("plain_text", "submit", false, false),
		CallbackID: ViewClaimShowCallbackID,
		ExternalID: c.generateExternalID(),
	}
	if _, err := c.Slack.Api.OpenView(triggerID, modal); err != nil {
		return fmt.Errorf("slack API OpenView error: %w", err)
	}
	return nil
}

func (c *ClaimController) Create(callback slack.InteractionCallback) (claim.Claim, error) {
	// Viewからの値の取り出し方
	// callback.View.State.Values[blockID][actionID].Value
	vlaues := callback.View.State.Values

	if claim := claim.FindClaim(callback.User.ID); claim != nil {
		return nil, fmt.Errorf("Claim already exist!")
	}

	encryptType := encrypt.EncryptType(vlaues["radio-encrypt-type"]["encrypt-type"].SelectedOption.Value)

	claim := claim.SlackClaim{
		Claim: &claim.ClaimBase{
			//FIXME: GetSubject()はslack.User から返すのにClaimBase.Validation()がSubjectを参照するせいでエラーになる
			Subject:          callback.User.ID,
			ClaimDate:        time.Now(),
			EncryptType:      encryptType,
			Namespaces:       []string{vlaues["input-namespace"]["namespace"].Value},
			State:            claim.Pending,
			ZipEncryptOption: claim.ZipEncryptOption{},
			GPGEncryptOption: claim.GPGEncryptOption{},
		},
		User: callback.User,
	}

	if err := claim.Validate(); err != nil {
		return nil, err
	}

	return &claim, nil
}

func (c *ClaimController) generateExternalID() string {
	return generateExternalID(ViewClaimShowCallbackID)
}

func (c *ClaimController) gpg() {
	//TODO: Implement gpg
	// if true {
	// 	gh := encrypt.GithubKey{
	// 		User: view.State.Values["input-github-account"]["github-account"].Value,
	// 	}
	//
	// 	gpgOpt = claim.GPGEncryptOption{
	// 		PublicKey: gh.PublicKeyString(),
	// 	}
	// }
}

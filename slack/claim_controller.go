package slack

import (
	"log"
	"time"

	"github.com/Sho2010/cinderella-simple/claim"
	"github.com/Sho2010/cinderella-simple/encrypt"
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
		CallbackID: ViewClaimShowCallbackID,
		ExternalID: c.generateExternalID(),
	}

	r, err := c.Slack.Api.OpenView(triggerID, modal)
	if err != nil {
		//TODO error handling
		log.Printf("とりあえずデバッグの為握りつぶす %v ", err)
	}
	println(r)
}

func (c *ClaimController) Create(callback slack.InteractionCallback) (claim.Claim, error) {

	// callback.View.State.Values[blockID][actionID].Value
	vlaues := callback.View.State.Values


	encryptType := encrypt.EncryptType(vlaues["radio-encrypt-type"]["encrypt-type"].SelectedOption.Value)

	claim := claim.SlackClaim{
		Claim: &claim.ClaimBase{
			ClaimDate:        time.Now(),
			EncryptType:      encryptType,
			Namespaces:       []string{vlaues["input-namespace"]["namespace"].Value},
			State:            claim.Pending,
			ZipEncryptOption: claim.ZipEncryptOption{},
			GPGEncryptOption: claim.GPGEncryptOption{},
		},
		User: callback.User,
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

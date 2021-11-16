package slack

import (
	"errors"
	"fmt"
	"time"

	"github.com/Sho2010/cinderella-simple/domain/model"
	"github.com/Sho2010/cinderella-simple/domain/repository"
	"github.com/Sho2010/cinderella-simple/encrypt"
	"github.com/slack-go/slack"
)

var ErrClaimAlreadyExist = fmt.Errorf("Claim already exist")

type ClaimController struct {
	Slack           *SlackApp
	claimRepository repository.ClaimRepository
}

func NewClaimController(slack *SlackApp, claimRepository repository.ClaimRepository) *ClaimController {
	return &ClaimController{
		Slack:           slack,
		claimRepository: claimRepository,
	}
}

func (c *ClaimController) Show(userID, triggerID string) error {
	if claim, _ := c.claimRepository.FindBySubject(userID); claim != nil {
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
		CallbackID: ViewClaimCallbackID,
		ExternalID: c.generateExternalID(),
	}
	if _, err := c.Slack.Api.OpenView(triggerID, modal); err != nil {
		return fmt.Errorf("slack API OpenView error: %w", err)
	}
	return nil
}

func (c *ClaimController) create(callback slack.InteractionCallback) (*model.Claim, error) {
	// Viewからの値の取り出し方
	// callback.View.State.Values[blockID][actionID].Value
	values := callback.View.State.Values

	// なかったのか、datastoreのエラーなのかわからんのなんとかする
	claim, _ := c.claimRepository.FindBySubject(callback.User.ID)
	if claim != nil {
		return claim, ErrClaimAlreadyExist
	}

	encryptType := encrypt.EncryptType(values["radio-encrypt-type"]["encrypt-type"].SelectedOption.Value)

	//TODO: NewClaim
	claim = &model.Claim{
		Subject:          callback.User.ID,
		Description:      values["input-description"]["description"].Value,
		ClaimAt:          time.Now(),
		EncryptType:      encryptType,
		Namespaces:       []string{values["input-namespace"]["namespace"].Value},
		State:            model.ClaimStatusPending,
		ZipEncryptOption: model.ZipEncryptOption{},
		GPGEncryptOption: model.GPGEncryptOption{},
	}
	if err := claim.Validate(); err != nil {
		return nil, err
	}

	return claim, nil
}

func (c *ClaimController) RegisterClaim(callback slack.InteractionCallback) (*model.Claim, error) {
	claim, err := c.create(callback)
	if err != nil {
		return nil, fmt.Errorf("create error: %w", err)
	}

	if err := repository.DefaultClaimRepository().Add(claim); err != nil {
		return nil, fmt.Errorf("registerClaim error: %w", err)
	}

	return claim, nil
}

func (c *ClaimController) generateExternalID() string {
	return generateExternalID(ViewClaimCallbackID)
}

func (c *ClaimController) IsClaimAlreadyExistError(err error) bool {
	return errors.Is(err, ErrClaimAlreadyExist)
}

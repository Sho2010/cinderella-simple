package slack

import (
	_ "embed"
	"fmt"

	"github.com/Sho2010/cinderella-simple/domain/model"
	"github.com/Sho2010/cinderella-simple/domain/repository"
	"github.com/slack-go/slack"
)

// NOTE PublishView の引数に関して
// - hash
//   A string that represents view state to protect against possible race conditions.
// - CallbackID
//   表示したViewの識別子 以下のようにする予定
//   cinderella_home_general
//   cinderella_home_admin
// - ExternalID
//   こちらは種類ごとではなくて、表示されたModal1つ1つの識別コードのイメージ。Workspace内で一意にしておく必要があります。
//
//   おすすめは、ユーザー名+タイムスタンプ。
//   これだと、よほどのことがない限り被らず安心です。

var (
	//go:embed views/home.json
	homeViewJson []byte

	//go:embed views/home_general.json
	homeGeneralBlocksJson []byte

	//go:embed views/home_admin.json
	homeAdminBlocksJson []byte
)

type HomeController struct {
	slack           *slack.Client
	claimRepository repository.ClaimRepository
}

func NewHomeController(slack *slack.Client, claimRepository repository.ClaimRepository) *HomeController {
	return &HomeController{
		slack:           slack,
		claimRepository: claimRepository,
	}
}

func (c *HomeController) Show(userID string) error {
	myClaim, _ := c.claimRepository.FindBySubject(userID)

	blocks, err := c.buildHomeView(myClaim)
	if err != nil {
		return fmt.Errorf("HomeView build failed: %w", err)
	}

	_, err = c.slack.PublishView(userID,
		slack.HomeTabViewRequest{
			Type:            slack.VTHomeTab,
			Blocks:          *blocks,
			PrivateMetadata: "",
			CallbackID:      ViewHomeCallbackID,
			ExternalID:      c.generateExternalID(),
		}, "")

	if err != nil {
		return fmt.Errorf("Slack PublishView failed: %w", err)
	}
	return nil
}

func (c *HomeController) Update(userID string) error {
	return c.Show(userID)
}

func (c *HomeController) buildHomeView(claim *model.Claim) (*slack.Blocks, error) {
	blocks := slack.Blocks{}

	if err := blocks.UnmarshalJSON(homeViewJson); err != nil {
		return nil, fmt.Errorf("Blocks marshal error %w", err)
	}

	//TODO: 管理者固定になってる
	isAdmin := true

	if isAdmin {
		adminBlocks, err := c.buildAdminView()
		if err != nil {
			return nil, err
		}
		blocks.BlockSet = append(blocks.BlockSet, adminBlocks...)
	} else {
		generalBlocks, err := c.buildGeneralView(claim)
		if err != nil {
			return nil, err
		}
		blocks.BlockSet = append(blocks.BlockSet, generalBlocks...)
		fmt.Printf("%#v\n", blocks)
	}

	fmt.Println(len(blocks.BlockSet))
	return &blocks, nil
}

func (c *HomeController) buildClaims() ([]slack.Block, error) {
	list, err := c.claimRepository.List()
	if err != nil {
		return nil, err
	}

	blocks := []slack.Block{}

	for _, v := range list {
		b := ClaimToBlock(v)

		acceptText := slack.NewTextBlockObject("plain_text", "Accept", false, false)
		rejectText := slack.NewTextBlockObject("plain_text", "Reject", false, false)

		accept := slack.NewButtonBlockElement(ActionAccept, "accept-claim", acceptText)
		reject := slack.NewButtonBlockElement(ActionReject, "reject-claim", rejectText)

		permitBlock := slack.NewActionBlock(BlockPermit, accept, reject)

		blocks = append(blocks, b, permitBlock)
	}
	return blocks, nil
}

func (c *HomeController) buildAdminView() ([]slack.Block, error) {
	headerBlocks := slack.Blocks{}

	if err := headerBlocks.UnmarshalJSON(homeAdminBlocksJson); err != nil {
		return nil, fmt.Errorf("Home admin view header blocks marshal error %w", err)
	}

	claimBlocks, err := c.buildClaims()
	if err != nil {
		return nil, fmt.Errorf("Home admin view claim blocks build error %w", err)
	}

	blocks := headerBlocks.BlockSet
	blocks = append(blocks, claimBlocks...)

	return blocks, nil
}

func (c *HomeController) buildGeneralView(myClaim *model.Claim) ([]slack.Block, error) {
	if myClaim == nil {
		return []slack.Block{slack.NewHeaderBlock(
			slack.NewTextBlockObject(slack.PlainTextType, "現在申請中の権限請求は存在しません。", false, false),
		)}, nil
	}

	var headerText string
	var downloadBlock slack.Block
	switch myClaim.GetState() {
	case model.ClaimStatusAccepted:
		headerText = ":memo: 承認済みの申請が存在します。\nKUBECONFIGをダウンロードし、作業を行ってください。"

		buttonText := slack.NewTextBlockObject("plain_text", "Download KUBECONFIG :arrow_down:", false, false)
		dlButton := slack.NewButtonBlockElement(ActionDownloadKubeconfig, "download-kubeconfig-value", buttonText)
		downloadBlock = slack.NewActionBlock(BlockDownloadKubeconfig, dlButton)

	case model.ClaimStatusPending:
		headerText = "保留中の申請が存在します。"
	case model.ClaimStatusExpired:
		headerText = "期限切れの申請が存在します。"
	case model.ClaimStatusRejected:
		headerText = "申請が拒否されました。"
	}

	header := slack.NewHeaderBlock(
		slack.NewTextBlockObject(slack.PlainTextType, headerText, false, false),
	)

	blocks := append([]slack.Block{}, header, ClaimToBlock(*myClaim))
	if downloadBlock != nil {
		blocks = append(blocks, downloadBlock)
	}

	return blocks, nil
}

func (c *HomeController) generateExternalID() string {
	return generateExternalID(ViewHomeCallbackID)
}

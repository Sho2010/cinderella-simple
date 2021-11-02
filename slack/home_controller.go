package slack

import (
	_ "embed"
	"fmt"

	"github.com/Sho2010/cinderella-simple/claim"
	"github.com/slack-go/slack"
)

var (
	//go:embed views/home.json
	homeViewJson []byte

	//go:embed views/home_general.json
	homeGeneralBlocksJson []byte

	//go:embed views/home_admin.json
	homeAdminBlocksJson []byte
)

type HomeController struct {
	slack *slack.Client
}

func (c *HomeController) Show(userID string) error {
	// NOTE: hashに関して
	// A string that represents view state to protect against possible race conditions.

	// NOTE: CallbackID
	// 表示したViewの識別子 以下のようにする予定
	// cinderella_home_general
	// cinderella_home_admin

	// NOTE: ExternalID
	// こちらは種類ごとではなくて、表示されたModal1つ1つの識別コードのイメージ。Workspace内で一意にしておく必要があります。
	//
	// おすすめは、ユーザー名+タイムスタンプ。
	// これだと、よほどのことがない限り被らず安心です。

	//TODO: 管理者固定になってる
	blocks, err := c.buildHomeView(true)
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

func (c *HomeController) buildHomeView(isAdmin bool) (*slack.Blocks, error) {
	blocks := slack.Blocks{}

	if err := blocks.UnmarshalJSON(homeViewJson); err != nil {
		return nil, fmt.Errorf("Blocks marshal error %w", err)
	}

	if isAdmin {
		adminBlocks, err := c.buildAdminView()
		if err != nil {
			return nil, err
		}
		blocks.BlockSet = append(blocks.BlockSet, adminBlocks...)

	} else {
		generalBlocks, err := c.buildGeneralView()
		if err != nil {
			return nil, err
		}
		blocks.BlockSet = append(blocks.BlockSet, generalBlocks...)
	}

	return &blocks, nil
}

func (c *HomeController) buildClaims() ([]slack.Block, error) {
	list := claim.ListClaims()
	blocks := make([]slack.Block, len(list))

	for i, v := range list {
		if slackClaim, ok := v.(*claim.SlackClaim); ok {
			b, err := slackClaim.ToBlock()
			if err != nil {
				return nil, fmt.Errorf("claim to slack.block fail: %w", err)
			}
			blocks[i] = b
		}
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

func (c *HomeController) buildGeneralView() ([]slack.Block, error) {
	blocks := slack.Blocks{}

	if err := blocks.UnmarshalJSON(homeViewJson); err != nil {
		return nil, fmt.Errorf("Home general blocks marshal error %w", err)
	}

	return blocks.BlockSet, nil
}

func (c *HomeController) generateExternalID() string {
	return generateExternalID(ViewHomeCallbackID)
}

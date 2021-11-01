package slack

import (
	_ "embed"
	"fmt"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
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

// AppHomeOpenedEventに依存関係持ってるの微妙な気もする...
func (c *HomeController) Show(e *slackevents.AppHomeOpenedEvent) error {
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

	blocks, err := c.buildHomeView()
	if err != nil {
		panic(err)
	}

	_, err = c.slack.PublishView(e.User,
		slack.HomeTabViewRequest{
			Type:            slack.VTHomeTab,
			Blocks:          *blocks,
			PrivateMetadata: "",
			CallbackID:      ViewHomeCallbackID,
			ExternalID:      c.generateExternalID(),
		}, "")

	if err != nil {
		panic(err)
	}
	return nil
}

func (c *HomeController) buildHomeView() (*slack.Blocks, error) {
	blocks := slack.Blocks{}

	if err := blocks.UnmarshalJSON(homeViewJson); err != nil {
		return nil, fmt.Errorf("Blocks marshal error %w", err)
	}

	return &blocks, nil
}

func (c *HomeController) generateExternalID() string {
	return generateExternalID(ViewHomeCallbackID)
}

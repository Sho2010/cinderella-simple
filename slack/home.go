package slack

import "github.com/slack-go/slack"

func buildHomeView() []slack.Block {
	return []slack.Block{
		slack.NewSectionBlock(
			&slack.TextBlockObject{
				Type: slack.MarkdownType,
				Text: "foo",
			},
			nil,
			slack.NewAccessory(
				slack.NewButtonBlockElement(
					"",
					"somevalue",
					&slack.TextBlockObject{
						Type: slack.PlainTextType,
						Text: "bar",
					},
				),
			),
		),
	}
}

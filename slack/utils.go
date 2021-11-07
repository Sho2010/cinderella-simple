package slack

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Sho2010/cinderella-simple/claim"
	"github.com/slack-go/slack"
)

func generateExternalID(base string) string {
	// これあってんのかな？
	return fmt.Sprintf("%s.%s", base, strconv.FormatInt(time.Now().UTC().UnixNano(), 10))
}

func dumpInteractionCallback(callback slack.InteractionCallback) {
	fmt.Printf("ActionID:%s\n", callback.ActionID)
	fmt.Printf("TriggerID:%s\n", callback.TriggerID)
	fmt.Printf("View.CallbackID:%s \n", callback.View.CallbackID)

	for _, v := range callback.ActionCallback.AttachmentActions {
		fmt.Println("---AttachmentActions element")
		fmt.Printf("%#v\n", v.Name)
	}

	for _, v := range callback.ActionCallback.BlockActions {
		fmt.Println("---Block Actions element")
		fmt.Printf("BlockID: %#v\n", v.BlockID)
		fmt.Printf("ActionID: %#v\n", v.ActionID)
		fmt.Printf("Text: %#v\n", v.Text.Text)
		fmt.Printf("Value: %#v\n", v.Value)
		fmt.Println("---")
	}
}

func debugPushViewSubmissionResponse() *slack.ViewSubmissionResponse {
	return slack.NewPushViewSubmissionResponse(&slack.ModalViewRequest{
		Type:  slack.VTModal,
		Title: slack.NewTextBlockObject("plain_text", "Test update view submission response", false, false),
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewContextBlock(
					"context_block_id",
					slack.NewTextBlockObject("plain_text", "Context text", false, false),
					slack.NewImageBlockElement("image_url", "alt_text"),
				),
			},
		},
	})
}

func debugErrorsViewSubmissionResponse() *slack.ViewSubmissionResponse {
	// このエラー処理は、block: { "type": "input" }にしか反応しないことに注意する。
	// ! select box, radio とかには利用できない
	return slack.NewErrorsViewSubmissionResponse(
		map[string]string{
			"input-namespace": "test_error",
		},
	)
}

func debugUpdateViewSubmissionResponse() *slack.ViewSubmissionResponse {
	return slack.NewUpdateViewSubmissionResponse(&slack.ModalViewRequest{
		Type:       slack.VTModal,
		Title:      slack.NewTextBlockObject("plain_text", "Test update view submission response", false, false),
		CallbackID: ViewHomeCallbackID,
		Blocks: slack.Blocks{BlockSet: []slack.Block{
			// slack.NewFileBlock("", "external_string", "source_string"),
			// slack.NewTextBlockObject("plain_text", "Test update view submission response", false, false),
			slack.NewSectionBlock(slack.NewTextBlockObject("plain_text", "Test update view submission response", false, false), nil, nil),
		}}})
}

// _, err := s.Api.UpdateView(
// 	slack.ModalViewRequest{
// 		Type:  slack.VTModal,
// 		Title: slack.NewTextBlockObject("plain_text", "Test update view submission response", false, false),
// 		// CallbackID: ViewHomeCallbackID,
// 		CallbackID: "test",
// 		Blocks: slack.Blocks{BlockSet: []slack.Block{
// 			// slack.NewFileBlock("", "external_string", "source_string"),
// 			// slack.NewTextBlockObject("plain_text", "Test update view submission response", false, false),
// 			slack.NewSectionBlock(slack.NewTextBlockObject("plain_text", "Test update view submission response", false, false), nil, nil),
// 		}}}, "dkajfjdajfda", "130321089730192739172392", callback.Container.ViewID)

func ClaimToBlock(c claim.Claim) slack.Block {
	// Block example
	// "text": "Claimer: *<@U04L97CP5>*\nPeriod: *30min*\nClaim Date: *2021/10/26 12:00:00*\nNamespace: *awesome*\nShort description: デバッグしたい"
	text := fmt.Sprintf("Claimer: *@%s*\nPeriod: *%s*\nClaim Date: *%s*\nNamespace: *%s*\nShort description: %s",
		c.GetName(),
		"30min", // TODO: implement me
		c.GetClaimAt().Format("2006/01/02 15:04:05"),
		fmt.Sprintf("%+q", c.GetNamespaces()),
		c.GetDescription(),
	)

	block := slack.NewSectionBlock(
		&slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: text,
		},
		nil,
		slack.NewAccessory(
			slack.NewImageBlockElement(
				"https://api.slack.com/img/blocks/bkb_template_images/creditcard.png", //TODO: いい感じのアイコン
				"claim_image",
			),
		),
	)

	return block
}

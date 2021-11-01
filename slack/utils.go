package slack

import (
	"fmt"
	"strconv"
	"time"

	"github.com/slack-go/slack"
)

func generateExternalID(base string) string {
	// これあってんのかな？
	return fmt.Sprintf("%s.%s", base, strconv.FormatInt(time.Now().UTC().UnixNano(), 10))
}

func dumpInteractionCallback(callback slack.InteractionCallback) {
	fmt.Printf("ActionID:%s\n", callback.ActionID)
	fmt.Printf("TriggerID:%s\n", callback.TriggerID)

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

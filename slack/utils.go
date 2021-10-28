package slack

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/slack-go/slack"
)

func generateExternalID(base string) string {
	// これあってんのかな？
	t := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	return fmt.Sprintf("%s.%s", base, t)
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

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

func Marshal(b *slack.Blocks, file string) {
	data, err := b.MarshalJSON()
	if err != nil {
		log.Printf("marshal err")
	}

	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write(data)
}

func UnMarshalTest() *slack.Blocks {
	file := os.Args[1]

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	blocks := slack.Blocks{}

	if err := blocks.UnmarshalJSON(bytes); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return &blocks
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

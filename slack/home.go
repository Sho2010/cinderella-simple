package slack

import (
	"fmt"
	"io/ioutil"

	"github.com/slack-go/slack"
)

func BuildHomeView() (*slack.Blocks, error) {

	file := "slack/views/home.json"

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("Blocks JSON read fail  %w", err)
	}

	blocks := slack.Blocks{}

	if err := blocks.UnmarshalJSON(bytes); err != nil {
		return nil, fmt.Errorf("Blocks marshal error %w", err)
	}

	return &blocks, nil
}

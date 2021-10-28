package slack

import (
	_ "embed"
	"fmt"

	"github.com/slack-go/slack"
)

//go:embed views/home.json
var homeViewJson []byte

func BuildHomeView() (*slack.Blocks, error) {

	// file := "slack/views/home.json"
	// bytes, err := ioutil.ReadFile(file)
	// if err != nil {
	// 	return nil, fmt.Errorf("Blocks JSON read fail  %w", err)
	// }

	blocks := slack.Blocks{}

	if err := blocks.UnmarshalJSON(homeViewJson); err != nil {
		return nil, fmt.Errorf("Blocks marshal error %w", err)
	}

	return &blocks, nil
}
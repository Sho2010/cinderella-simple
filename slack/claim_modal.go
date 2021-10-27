package slack

import (
	_ "embed"
	"fmt"

	"github.com/slack-go/slack"
)

//go:embed views/claim_modal.json
var claimViewJson []byte

func BuildClaimModalView() (*slack.Blocks, error) {
	blocks := slack.Blocks{}
	if err := blocks.UnmarshalJSON(claimViewJson); err != nil {
		return nil, fmt.Errorf("Blocks marshal error %w", err)
	}

	return &blocks, nil
}

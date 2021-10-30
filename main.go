package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Sho2010/cinderella-simple/audit"
	"github.com/Sho2010/cinderella-simple/k8s"
	"github.com/Sho2010/cinderella-simple/slack"
)

func main() {
	fmt.Println("hello cinderella")

	ctx := context.Background()
	client, _ := k8s.GetDefaultClient()

	cleaner, err := k8s.NewCleaner(client)
	if err != nil {
		panic(err)
	}

	go cleaner.Start(ctx)

	broadcast := audit.AuditBroadcaster{}
	broadcast.Start()

	//tokenが設定されてた場合slack socket mode を起動
	if len(os.Getenv("SLACK_BOT_TOKEN")) != 0 && len(os.Getenv("SLACK_APP_TOKEN")) != 0 {
		s := slack.NewSlackApp().Slack.Administrators)
		s.Start()
	}

	select {} // Block all
}

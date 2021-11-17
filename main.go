package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Sho2010/cinderella-simple/audit"
	"github.com/Sho2010/cinderella-simple/config"
	"github.com/Sho2010/cinderella-simple/domain/event"
	"github.com/Sho2010/cinderella-simple/k8s"
	"github.com/Sho2010/cinderella-simple/slack"
)

func main() {
	fmt.Println("hello cinderella")

	fmt.Println("load config")
	config.LoadConfig()

	client, _ := k8s.GetDefaultClient()

	cleaner, err := k8s.NewCleaner(client, 0)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	go cleaner.Start(ctx)

	broadcast := audit.AuditBroadcaster{}
	broadcast.Start()

	//tokenが設定されてた場合slack socket mode を起動
	if len(os.Getenv("SLACK_BOT_TOKEN")) != 0 && len(os.Getenv("SLACK_APP_TOKEN")) != 0 {
		s := slack.NewSlackApp(config.GetConfig().Slack.Administrators)

		handler := slack.NewClaimEventHandler(event.ClaimEventChannel(), s.Api)
		handler.Start()

		fmt.Printf("U04L97CP5  is admin?: %v \n", s.IsAdmin("U04L97CP5"))
		fmt.Printf("U04KMRW1Y  is admin?: %v \n", s.IsAdmin("U04KMRW1Y"))
		s.Start()
	}

	select {} // Block all
}

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/codingpop/workflow-trigger"
)

var (
	owner       string
	repo        string
	accessToken string
	eventType   string
)

func init() {
	flag.StringVar(&owner, "owner", "", "Repository owner")
	flag.StringVar(&repo, "repo", "", "GitHub repository name")
	flag.StringVar(&accessToken, "access-token", "", "Personal access token")
	flag.StringVar(&eventType, "event-type", "", "A custom webhook event name")
}

func main() {
	flag.Parse()

	if owner == "" {
		log.Fatalln("must provide repository owner")
	}
	if repo == "" {
		log.Fatalln("must provide repository name")
	}
	if accessToken == "" {
		log.Fatalln("must provide personal access token")
	}
	if eventType == "" {
		log.Fatalln("must provide event type")
	}

	p := workflow.Params{
		Owner:       owner,
		Repo:        repo,
		AccessToken: accessToken,
		EventType:   eventType,
	}
	w := workflow.Configure(p)

	err := w.Trigger()

	fmt.Println(err)
}

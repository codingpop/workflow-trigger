package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/codingpop/workflow-trigger"
)

var owner = flag.String("owner", "", "Repository owner")
var repo = flag.String("repo", "", "GitHub repository name")
var accessToken = flag.String("access-token", "", "Personal access token")
var eventType = flag.String("event-type", "", "A custom webhook event name")

func main() {
	flag.Parse()

	if *owner == "" {
		log.Fatalln("must provide repository owner")
	}
	if *repo == "" {
		log.Fatalln("must provide repository name")
	}
	if *accessToken == "" {
		log.Fatalln("must provide personal access token")
	}
	if *eventType == "" {
		log.Fatalln("must provide event type")
	}

	p := workflow.Params{
		Owner:       *owner,
		Repo:        *repo,
		AccessToken: *accessToken,
		EventType:   *eventType,
	}
	w := workflow.Configure(p)

	err := w.Trigger()

	fmt.Println(err)
}

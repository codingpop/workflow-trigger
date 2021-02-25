package workflow_test

import (
	"context"
	"flag"
	"testing"

	"github.com/codingpop/workflow-trigger"
)

var owner = flag.String("owner", "", "Repository owner")
var repo = flag.String("repo", "", "GitHub repository name")
var accessToken = flag.String("access-token", "", "Personal access token")
var eventType = flag.String("event-type", "", "A custom webhook event name")

func TestTrigger(t *testing.T) {
	p := workflow.Params{
		Owner:       *owner,
		Repo:        *repo,
		AccessToken: *accessToken,
		EventType:   *eventType,
	}
	w := workflow.Configure(p)

	got := w.Trigger()

	if got != nil {
		t.Fail()
	}
}

func TestTriggerContext(t *testing.T) {
	p := workflow.Params{
		Owner:       *owner,
		Repo:        *repo,
		AccessToken: *accessToken,
		EventType:   *eventType,
	}
	w := workflow.Configure(p)

	ctx := context.Background()
	got := w.TriggerContext(ctx)

	if got != nil {
		t.Fail()
	}
}

func TestTriggerError(t *testing.T) {
	p := workflow.Params{
		Owner:       *owner,
		Repo:        "unknown-repo",
		AccessToken: *accessToken,
		EventType:   *eventType,
	}
	w := workflow.Configure(p)

	got := w.Trigger()

	if got == nil {
		t.Fail()
	}
}

func TestTriggerContextError(t *testing.T) {
	p := workflow.Params{
		Owner:       *owner,
		Repo:        "unknown-repo",
		AccessToken: *accessToken,
		EventType:   *eventType,
	}
	w := workflow.Configure(p)

	ctx := context.Background()
	got := w.TriggerContext(ctx)

	if got == nil {
		t.Fail()
	}
}

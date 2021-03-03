package workflow_test

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/codingpop/workflow-trigger"
)

var owner = flag.String("owner", "", "Repository owner")
var repo = flag.String("repo", "", "GitHub repository name")
var accessToken = flag.String("access-token", "", "Personal access token")
var eventType = flag.String("event-type", "", "A custom webhook event name")

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func TestWorkflow(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	baseURL := strings.Replace(server.URL, "https://", "", 1)

	t.Run("TestTrigger", func(t *testing.T) {
		p := workflow.Params{
			Owner:       *owner,
			Repo:        *repo,
			AccessToken: *accessToken,
			EventType:   *eventType,
			BaseURL:     baseURL,
		}
		w := workflow.Configure(p)

		got := w.Trigger()

		if got != nil {
			t.Errorf("expecting nil err, got %v", got)
		}
	})

	t.Run("TestTriggerContext", func(t *testing.T) {
		p := workflow.Params{
			Owner:       *owner,
			Repo:        *repo,
			AccessToken: *accessToken,
			EventType:   *eventType,
			BaseURL:     baseURL,
		}

		w := workflow.Configure(p)

		ctx := context.Background()
		got := w.TriggerContext(ctx)

		if got != nil {
			t.Errorf("expecting nil err, got %v", got)
		}
	})
}

func TestWorkflowNotFound(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	baseURL := strings.Replace(server.URL, "https://", "", 1)

	t.Run("TestTriggeNotFound", func(t *testing.T) {
		p := workflow.Params{
			Owner:       *owner,
			Repo:        "unknown-repo",
			AccessToken: *accessToken,
			EventType:   *eventType,
			BaseURL:     baseURL,
		}
		w := workflow.Configure(p)

		got := w.Trigger()

		if got == nil {
			t.Fail()
		}
	})

	t.Run("TestTriggerContextNotFound", func(t *testing.T) {
		p := workflow.Params{
			Owner:       *owner,
			Repo:        "unknown-repo",
			AccessToken: *accessToken,
			EventType:   *eventType,
			BaseURL:     baseURL,
		}
		w := workflow.Configure(p)

		ctx := context.Background()
		got := w.TriggerContext(ctx)

		if got == nil {
			t.Fail()
		}
	})
}

func TestWorkflowUnauthorized(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)

		errResponse := map[string]string{
			"message": "Bad credentials",
		}

		json.NewEncoder(w).Encode(errResponse)
	}))
	defer server.Close()

	baseURL := strings.Replace(server.URL, "https://", "", 1)
	expected := "Bad credentials"

	t.Run("TestTriggerUnauthorized", func(t *testing.T) {
		p := workflow.Params{
			Owner:       *owner,
			Repo:        "unknown-repo",
			AccessToken: *accessToken,
			EventType:   *eventType,
			BaseURL:     baseURL,
		}
		w := workflow.Configure(p)

		got := w.Trigger()

		if got == nil {
			t.Fail()
		}

		if got.Error() != expected {
			t.Errorf("expected an error message: %s, got: %s", expected, got.Error())
		}
	})

	t.Run("TestTriggerContextUnauthorized", func(t *testing.T) {
		p := workflow.Params{
			Owner:       *owner,
			Repo:        "unknown-repo",
			AccessToken: *accessToken,
			EventType:   *eventType,
			BaseURL:     baseURL,
		}
		w := workflow.Configure(p)

		ctx := context.Background()
		got := w.TriggerContext(ctx)

		if got == nil {
			t.Fail()
		}

		if got.Error() != expected {
			t.Errorf("expected an error message: %s, got: %s", expected, got.Error())
		}
	})
}

func TestWorkflowUnproccessableEntity(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)

		errResponse := map[string]string{
			"message": "Invalid request.\n\n\"event_type\" wasn't supplied.",
		}

		json.NewEncoder(w).Encode(errResponse)
	}))
	defer server.Close()

	baseURL := strings.Replace(server.URL, "https://", "", 1)
	expected := "Invalid request.\n\n\"event_type\" wasn't supplied."

	t.Run("TestTriggerUnproccessableEntity", func(t *testing.T) {
		p := workflow.Params{
			Owner:       *owner,
			Repo:        "unknown-repo",
			AccessToken: *accessToken,
			EventType:   *eventType,
			BaseURL:     baseURL,
		}
		w := workflow.Configure(p)

		got := w.Trigger()

		if got == nil {
			t.Fail()
		}

		if got.Error() != expected {
			t.Errorf("expected an error message: %s, got: %s", expected, got.Error())
		}
	})

	t.Run("TestTriggerContextUnproccessableEntity", func(t *testing.T) {
		p := workflow.Params{
			Owner:       *owner,
			Repo:        "unknown-repo",
			AccessToken: *accessToken,
			EventType:   *eventType,
			BaseURL:     baseURL,
		}
		w := workflow.Configure(p)

		ctx := context.Background()
		got := w.TriggerContext(ctx)

		if got == nil {
			t.Fail()
		}

		if got.Error() != expected {
			t.Errorf("expected an error message: %s, got: %s", expected, got.Error())
		}
	})
}

func TestWorkflowServerError(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	baseURL := strings.Replace(server.URL, "https://", "", 1)
	expected := "server error"

	t.Run("TestTriggerServerError", func(t *testing.T) {
		p := workflow.Params{
			Owner:       *owner,
			Repo:        "unknown-repo",
			AccessToken: *accessToken,
			EventType:   *eventType,
			BaseURL:     baseURL,
		}
		w := workflow.Configure(p)

		got := w.Trigger()

		if got == nil {
			t.Fail()
		}

		if got.Error() != expected {
			t.Errorf("expected an error message: %s, got: %s", expected, got.Error())
		}
	})

	t.Run("TestTriggerContextServerError", func(t *testing.T) {
		p := workflow.Params{
			Owner:       *owner,
			Repo:        "unknown-repo",
			AccessToken: *accessToken,
			EventType:   *eventType,
			BaseURL:     baseURL,
		}
		w := workflow.Configure(p)

		ctx := context.Background()
		got := w.TriggerContext(ctx)

		if got == nil {
			t.Fail()
		}

		if got.Error() != expected {
			t.Errorf("expected an error message: %s, got: %s", expected, got.Error())
		}
	})
}

package workflow

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/sync/errgroup"
)

// Workflow defines the trigger functions
type Workflow struct {
	api         string
	accessToken string
	eventType   string
	client      *http.Client
}

// Params collects repository credentials and workflow information
type Params struct {
	Repo        string
	Owner       string
	AccessToken string
	EventType   string
	TimeOut     time.Duration
}

// Configure creates an instance of Workflow
func Configure(p Params) *Workflow {
	c := &http.Client{
		Timeout: p.TimeOut,
	}

	u := url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   fmt.Sprintf("/repos/%s/%s/dispatches", p.Owner, p.Repo),
	}

	return &Workflow{
		api:         u.String(),
		accessToken: p.AccessToken,
		eventType:   p.EventType,
		client:      c,
	}
}

func (w *Workflow) trigger(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		b, err := json.Marshal(map[string]string{
			"event_type": w.eventType,
		})
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}

		body := bytes.NewReader(b)

		req, err := http.NewRequestWithContext(ctx, "POST", w.api, body)
		if err != nil {
			return fmt.Errorf("failed to create new request: %w", err)
		}

		req.Header.Set("Accept", "application/vnd.github.v3+json")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("token %s", w.accessToken))

		resp, err := w.client.Do(req)
		if err != nil {
			return fmt.Errorf("unexpected error triggering workflow: %w", err)
		}

		return handleResponse(resp)
	})

	return eg.Wait()
}

// Trigger triggers a GitHub Action workflow
func (w *Workflow) Trigger() error {
	return w.trigger(context.Background())
}

// TriggerContext accepts context.Context and triggers a GitHub Action workflow
func (w *Workflow) TriggerContext(ctx context.Context) error {
	return w.trigger(ctx)
}

func handleResponse(r *http.Response) (retErr error) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			retErr = err
		}
	}()

	var response struct {
		Message string `json:"message"`
	}

	if r.StatusCode < http.StatusBadRequest {
		return nil
	}

	if r.StatusCode >= http.StatusInternalServerError {
		return errors.New("server error")
	}

	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("unexpected error parsing api response: %w", err)
	}

	return errors.New(response.Message)
}

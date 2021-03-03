package workflow

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// Workflow defines the trigger functions
type Workflow struct {
	api         string
	accessToken string
	eventType   string
}

// Params collects repository credentials and workflow information
type Params struct {
	BaseURL     string
	Repo        string
	Owner       string
	AccessToken string
	EventType   string
	Retries     string
}

// Configure creates an instance of Workflow
func Configure(p Params) *Workflow {
	host := "api.github.com"
	if p.BaseURL != "" {
		host = p.BaseURL
	}

	u := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   fmt.Sprintf("/repos/%s/%s/dispatches", p.Owner, p.Repo),
	}

	return &Workflow{
		api:         u.String(),
		accessToken: p.AccessToken,
		eventType:   p.EventType,
	}
}

func (w *Workflow) trigger(ctx context.Context) (retErr error) {
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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unexpected error triggering workflow: %w", err)
	}

	return handleResponse(resp)
}

// TriggerContext accepts context.Context and triggers a GitHub Action workflow
func (w *Workflow) TriggerContext(ctx context.Context) error {
	return w.trigger(ctx)
}

// Trigger triggers a GitHub Action workflow
func (w *Workflow) Trigger() error {
	return w.TriggerContext(context.Background())
}

func handleResponse(r *http.Response) (retErr error) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			retErr = err
		}
	}()

	if r.StatusCode < http.StatusBadRequest {
		return nil
	}

	if r.StatusCode >= http.StatusInternalServerError {
		return errors.New("server error")
	}

	switch r.StatusCode {
	case http.StatusNoContent, http.StatusOK:
		return nil
	case http.StatusNotFound, http.StatusUnauthorized, http.StatusUnprocessableEntity, http.StatusInternalServerError:
		var body struct {
			Message string `json:"message"`
		}

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			return fmt.Errorf("unexpected error parsing api response: %w", err)
		}

		return errors.New(body.Message)
	}

	return nil
}

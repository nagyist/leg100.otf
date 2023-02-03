package github

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/v41/github"
	"github.com/leg100/otf/cloud"
	"github.com/leg100/otf/events"
)

type Handler struct {
	cloud.HandlerOptions
}

// HandleEvent handles incoming events from github
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	event, err := h.handle(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	// relay event to subscribers
	h.Publish(events.Event{
		Type:    events.EventVCS,
		Payload: event,
	})
}

func (h *Handler) handle(w http.ResponseWriter, r *http.Request) (cloud.VCSEvent, error) {
	payload, err := github.ValidatePayload(r, []byte(h.Secret))
	if err != nil {
		return nil, err
	}

	rawEvent, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		return nil, err
	}

	switch event := rawEvent.(type) {
	case *github.PushEvent:
		// a github.PushEvent includes tag events but otf categorises them as separate
		// event types
		parts := strings.Split(event.GetRef(), "/")
		if len(parts) != 3 || parts[0] != "refs" {
			return nil, fmt.Errorf("malformed ref: %s", event.GetRef())
		}
		switch parts[1] {
		case "tags":
			var action cloud.VCSTagEventAction
			switch {
			case event.GetCreated():
				action = cloud.VCSTagEventCreatedAction
			case event.GetDeleted():
				action = cloud.VCSTagEventDeletedAction
			default:
				return nil, fmt.Errorf("no action specified for tag event")
			}

			return cloud.VCSTagEvent{
				WebhookID:  h.WebhookID,
				Tag:        parts[2],
				Action:     action,
				Identifier: event.GetRepo().GetFullName(),
				CommitSHA:  event.GetAfter(),
			}, nil
		case "heads":
			return cloud.VCSPushEvent{
				WebhookID:  h.WebhookID,
				Branch:     parts[2],
				Identifier: event.GetRepo().GetFullName(),
				CommitSHA:  event.GetAfter(),
			}, nil
		default:
			return nil, fmt.Errorf("malformed ref: %s", event.GetRef())
		}
	case *github.PullRequestEvent:
		var action cloud.VCSPullEventAction
		switch event.GetAction() {
		case "opened":
			action = cloud.VCSPullEventOpened
		case "closed":
			if event.PullRequest.GetMerged() {
				action = cloud.VCSPullEventMerged
			} else {
				action = cloud.VCSPullEventClosed
			}
		case "synchronised":
			action = cloud.VCSPullEventUpdated
		}

		return cloud.VCSPullEvent{
			WebhookID:  h.WebhookID,
			Action:     action,
			Identifier: event.GetRepo().GetFullName(),
			Branch:     event.PullRequest.Head.GetRef(),
			CommitSHA:  event.GetPullRequest().GetHead().GetSHA(),
		}, nil
	default:
		return nil, nil
	}
}

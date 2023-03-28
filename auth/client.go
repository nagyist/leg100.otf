package auth

import (
	"context"
	"fmt"
	"net/url"

	"github.com/leg100/otf"
	"github.com/leg100/otf/http/jsonapi"
)

type Client struct {
	otf.JSONAPIClient
}

// CreateUser creates a user via HTTP/JSONAPI. Options are ignored.
func (c *Client) CreateUser(ctx context.Context, username string, _ ...NewUserOption) (*User, error) {
	req, err := c.NewRequest("POST", "admin/users", &jsonapi.CreateUserOptions{
		Username: otf.String(username),
	})
	if err != nil {
		return nil, err
	}
	user := &jsonapi.User{}
	err = c.Do(ctx, req, user)
	if err != nil {
		return nil, err
	}
	return &User{ID: user.ID, Username: user.Username}, nil
}

// DeleteUser deletes a user via HTTP/JSONAPI.
func (c *Client) DeleteUser(ctx context.Context, username string) error {
	u := fmt.Sprintf("admin/users/%s", url.QueryEscape(username))
	req, err := c.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}
	if err := c.Do(ctx, req, nil); err != nil {
		return err
	}
	return nil
}

// CreateRegistrySession creates a registry session via HTTP/JSONAPI
func (c *Client) CreateRegistrySession(ctx context.Context, opts CreateRegistrySessionOptions) (*RegistrySession, error) {
	req, err := c.NewRequest("POST", "registry/sessions/create", &jsonapi.RegistrySessionCreateOptions{
		Organization: opts.Organization,
	})
	if err != nil {
		return nil, err
	}
	session := &jsonapi.RegistrySession{}
	err = c.Do(ctx, req, session)
	if err != nil {
		return nil, err
	}
	return &RegistrySession{
		Organization: session.OrganizationName,
		Token:        session.Token,
	}, nil
}

func (c *Client) CreateAgentToken(ctx context.Context, options CreateAgentTokenOptions) (*AgentToken, error) {
	req, err := c.NewRequest("POST", "agent/create", &jsonapi.AgentTokenCreateOptions{
		Description:  options.Description,
		Organization: options.Organization,
	})
	if err != nil {
		return nil, err
	}
	at := &jsonapi.AgentToken{}
	err = c.Do(ctx, req, at)
	if err != nil {
		return nil, err
	}
	return &AgentToken{ID: at.ID, Token: *at.Token, Organization: at.Organization}, nil
}

func (c *Client) GetAgentToken(ctx context.Context, token string) (*AgentToken, error) {
	req, err := c.NewRequest("GET", "agent/details", nil)
	if err != nil {
		return nil, err
	}

	at := &jsonapi.AgentToken{}
	err = c.Do(ctx, req, at)
	if err != nil {
		return nil, err
	}

	return &AgentToken{ID: at.ID, Organization: at.Organization}, nil
}

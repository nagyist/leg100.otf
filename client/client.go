package client

import (
	"github.com/leg100/go-tfe"
)

const (
	DefaultAddress = "localhost:8080"
)

type Client interface {
	Organizations() tfe.Organizations
	Workspaces() tfe.Workspaces
	ConfigurationVersions() tfe.ConfigurationVersions
	ConfigurationVersionsExt() ConfigurationVersionsExt
}

type client struct {
	*tfe.Client

	configurationVersionsExt ConfigurationVersionsExt
}

func NewClient(cfg *config) (*client, error) {
	if err := cfg.sanitizeAddress(); err != nil {
		return nil, err
	}

	creds, err := NewCredentialsStore(&SystemDirectories{})
	if err != nil {
		return nil, err
	}

	// If token isn't set then load from DB
	if cfg.Token == "" {
		cfg.Token, err = creds.Load(cfg.Address)
		if err != nil {
			return nil, err
		}
	}

	client := client{}

	client.Client, err = tfe.NewClient(&cfg.Config)
	if err != nil {
		return nil, err
	}

	client.configurationVersionsExt = &configurationVersionsExt{client: &client}

	return &client, nil
}

func (c *client) Organizations() tfe.Organizations {
	return c.Client.Organizations
}

func (c *client) Workspaces() tfe.Workspaces {
	return c.Client.Workspaces
}

func (c *client) ConfigurationVersions() tfe.ConfigurationVersions {
	return c.Client.ConfigurationVersions
}

func (c *client) ConfigurationVersionsExt() ConfigurationVersionsExt {
	return c.configurationVersionsExt
}

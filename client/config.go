package client

import (
	"fmt"
	"net/url"

	"github.com/leg100/go-tfe"
)

type Config interface {
	NewClient() (Client, error)
}

type config struct {
	tfe.Config
}

func (c *config) NewClient() (Client, error) {
	return NewClient(c)
}

// Ensure address is in format https://<host>:<port>
func (c *config) sanitizeAddress() error {
	u, err := url.ParseRequestURI(c.Address)
	if err != nil || u.Host == "" {
		u, er := url.ParseRequestURI("https://" + c.Address)
		if er != nil {
			return fmt.Errorf("could not parse hostname: %w", err)
		}
		c.Address = u.String()
		return nil
	}

	u.Scheme = "https"
	c.Address = u.String()

	return nil
}

package agentservice

import (
	"context"

	"github.com/leg100/otf/internal"
)

type Client struct {
	internal.JSONAPIClient
}

func (c *Client) GetJob(ctx context.Context, agentID string) (<-chan *Job, error) {
}

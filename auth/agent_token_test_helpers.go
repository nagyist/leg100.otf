package auth

import (
	"context"
	"testing"

	"github.com/leg100/otf"
	"github.com/stretchr/testify/require"
)

func NewTestAgentToken(t *testing.T, org string) *otf.AgentToken {
	token, err := otf.NewAgentToken(otf.CreateAgentTokenOptions{
		Organization: org,
		Description:  "lorem ipsum...",
	})
	require.NoError(t, err)
	return token
}

type fakeAgentTokenService struct {
	token *otf.AgentToken

	agentTokenService
}

func (f *fakeAgentTokenService) createAgentToken(context.Context, otf.CreateAgentTokenOptions) (*otf.AgentToken, error) {
	return f.token, nil
}

func (f *fakeAgentTokenService) listAgentTokens(context.Context, string) ([]*otf.AgentToken, error) {
	return []*otf.AgentToken{f.token}, nil
}

func (f *fakeAgentTokenService) deleteAgentToken(context.Context, string) (*otf.AgentToken, error) {
	return f.token, nil
}
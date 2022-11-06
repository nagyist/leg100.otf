package otf

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/v41/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestGithub_GetUser(t *testing.T) {
	ctx := context.Background()
	org := NewTestOrganization(t)
	team := NewTeam("fake-team", org)
	want := NewUser("fake-user", WithOrganizationMemberships(org), WithTeamMemberships(team))
	srv := NewTestGithubServer(t, want)

	u, err := url.Parse(srv.URL)
	require.NoError(t, err)

	cloud := &GithubCloud{
		&GithubConfig{
			cloudConfig{
				hostname:            u.Host,
				skipTLSVerification: true,
			},
		},
	}
	client, err := cloud.NewDirectoryClient(ctx, DirectoryClientOptions{
		OAuthToken: &oauth2.Token{AccessToken: "fake-token"},
	})
	require.NoError(t, err)

	got, err := client.GetUser(ctx)
	require.NoError(t, err)

	assert.Equal(t, want.Username(), got.Username())
	if assert.Equal(t, 1, len(got.Organizations())) {
		assert.Equal(t, org.Name(), got.Organizations()[0].Name())
	}
	if assert.Equal(t, 1, len(got.Teams())) {
		assert.Equal(t, team.Name(), got.Teams()[0].Name())
	}
}

func NewTestGithubServer(t *testing.T, user *User) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/login/oauth/authorize", func(w http.ResponseWriter, r *http.Request) {
		q := url.Values{}
		q.Add("state", r.URL.Query().Get("state"))
		q.Add("code", GenerateRandomString(10))

		referrer, err := url.Parse(r.Referer())
		require.NoError(t, err)

		callback := url.URL{
			Scheme:   referrer.Scheme,
			Host:     referrer.Host,
			Path:     "/oauth/github/callback",
			RawQuery: q.Encode(),
		}

		http.Redirect(w, r, callback.String(), http.StatusFound)
	})
	mux.HandleFunc("/login/oauth/access_token", func(w http.ResponseWriter, r *http.Request) {
		out, err := json.Marshal(&oauth2.Token{AccessToken: "stub_token"})
		require.NoError(t, err)
		w.Header().Add("Content-Type", "application/json")
		w.Write(out)
	})
	mux.HandleFunc("/api/v3/user", func(w http.ResponseWriter, r *http.Request) {
		out, err := json.Marshal(&github.User{Login: String(user.Username())})
		require.NoError(t, err)
		w.Header().Add("Content-Type", "application/json")
		w.Write(out)
	})
	mux.HandleFunc("/api/v3/user/orgs", func(w http.ResponseWriter, r *http.Request) {
		var orgs []*github.Organization
		for _, org := range user.Organizations() {
			orgs = append(orgs, &github.Organization{Login: String(org.Name())})
		}
		out, err := json.Marshal(orgs)
		require.NoError(t, err)
		w.Header().Add("Content-Type", "application/json")
		w.Write(out)
	})
	for _, org := range user.Organizations() {
		mux.HandleFunc("/api/v3/user/memberships/orgs/"+org.Name(), func(w http.ResponseWriter, r *http.Request) {
			out, err := json.Marshal(&github.Membership{
				Role: String("member"),
			})
			require.NoError(t, err)
			w.Header().Add("Content-Type", "application/json")
			w.Write(out)
		})
	}
	mux.HandleFunc("/api/v3/user/teams", func(w http.ResponseWriter, r *http.Request) {
		var teams []*github.Team
		for _, team := range user.Teams() {
			teams = append(teams, &github.Team{
				Name: String(team.Name()),
				Organization: &github.Organization{
					Login: String(team.OrganizationName()),
				},
			})
		}
		out, err := json.Marshal(teams)
		require.NoError(t, err)
		w.Header().Add("Content-Type", "application/json")
		w.Write(out)
	})

	srv := httptest.NewTLSServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

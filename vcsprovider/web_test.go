package vcsprovider

import (
	"context"
	"fmt"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/leg100/otf/http/html"
	"github.com/leg100/otf/inmem"
	"github.com/leg100/otf/organization"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTML_New(t *testing.T) {
	svc := fakeWebServices(t, NewTestVCSProvider(t, nil))

	for _, cloud := range []string{"github", "gitlab"} {
		t.Run(cloud, func(t *testing.T) {
			q := "/?organization_name=acme-corp&cloud=" + cloud
			r := httptest.NewRequest("GET", q, nil)
			w := httptest.NewRecorder()
			svc.new(w, r)
			t.Log(w.Body.String())
			assert.Equal(t, 200, w.Code)
		})
	}
}

func TestCreateVCSProviderHandler(t *testing.T) {
	org := organization.NewTestOrganization(t)
	svc := fakeWebServices(t, NewTestVCSProvider(t, org))

	form := strings.NewReader(url.Values{
		"organization_name": {"acme-corp"},
		"token":             {"secret-token"},
		"name":              {"my-new-vcs-provider"},
		"cloud":             {"fake-cloud"},
	}.Encode())

	r := httptest.NewRequest("POST", "/organization/acme-corp/vcs-providers/create", form)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	svc.create(w, r)

	if assert.Equal(t, 302, w.Code) {
		redirect, err := w.Result().Location()
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("/app/organizations/%s/vcs-providers", org.Name), redirect.Path)
	} else {
		t.Log(w.Body.String())
	}
}

func TestListVCSProvidersHandler(t *testing.T) {
	org := organization.NewTestOrganization(t)
	app := fakeWebServices(t, NewTestVCSProvider(t, org))

	r := httptest.NewRequest("GET", "/?organization_name=acme-corp", nil)
	w := httptest.NewRecorder()
	app.list(w, r)

	assert.Equal(t, 200, w.Code)
}

func TestDeleteVCSProvidersHandler(t *testing.T) {
	org := organization.NewTestOrganization(t)
	app := fakeWebServices(t, NewTestVCSProvider(t, org))

	form := strings.NewReader(url.Values{
		"vcs_provider_id": {"fake-id"},
	}.Encode())

	r := httptest.NewRequest("POST", "/?", form)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	app.delete(w, r)

	assert.Equal(t, 302, w.Code)
}

func fakeWebServices(t *testing.T, provider *VCSProvider) *webHandlers {
	renderer, err := html.NewViewEngine(false)
	require.NoError(t, err)
	return &webHandlers{
		Renderer:     renderer,
		svc:          &fakeApp{provider: provider},
		CloudService: inmem.NewTestCloudService(),
	}
}

type fakeApp struct {
	provider *VCSProvider

	Service
}

func (f *fakeApp) create(ctx context.Context, opts createOptions) (*VCSProvider, error) {
	return f.provider, nil
}

func (f *fakeApp) list(context.Context, string) ([]*VCSProvider, error) {
	return []*VCSProvider{f.provider}, nil
}

func (f *fakeApp) delete(context.Context, string) (*VCSProvider, error) {
	return f.provider, nil
}

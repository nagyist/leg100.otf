package http

import (
	"context"
	"fmt"
	"html/template"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/handlers"

	"net"
	"net/http"
	"time"

	"github.com/allegro/bigcache"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/leg100/jsonapi"
	"github.com/leg100/otf"
	"github.com/leg100/otf/http/assets"
)

const (
	// ShutdownTimeout is the time given for outstanding requests to finish
	// before shutdown.
	ShutdownTimeout = 1 * time.Second

	jsonApplication = "application/json"

	UploadConfigurationVersionRoute WebRoute = "/configuration-versions/%v/upload"
	GetPlanLogsRoute                WebRoute = "plans/%v/logs"
	GetApplyLogsRoute               WebRoute = "applies/%v/logs"

	FlashKey = "flash"
)

var (
	store = sessions.NewCookieStore([]byte("a-random-key"))

	embeddedAssetServer assets.Server
)

// Load embedded templates at startup
func init() {
	server, err := assets.NewEmbeddedServer()
	if err != nil {
		panic("unable to load embedded assets: " + err.Error())
	}

	embeddedAssetServer = server
}

type WebRoute string

// Server provides an HTTP/S server
type Server struct {
	server *http.Server
	ln     net.Listener
	err    chan error

	logr.Logger

	EnableRequestLogging bool

	SSL               bool
	CertFile, KeyFile string

	// Listening Address in the form <ip>:<port>
	Addr string

	OrganizationService         otf.OrganizationService
	WorkspaceService            otf.WorkspaceService
	StateVersionService         otf.StateVersionService
	ConfigurationVersionService otf.ConfigurationVersionService
	EventService                otf.EventService
	RunService                  otf.RunService
	PlanService                 otf.PlanService
	ApplyService                otf.ApplyService
	TokenService                otf.TokenService
	CacheService                *bigcache.BigCache

	assets.Server
}

// NewServer is the constructor for Server
func NewServer() *Server {
	s := &Server{
		server: &http.Server{},
		err:    make(chan error),
		Server: embeddedAssetServer,
	}

	return s
}

// NewRouter constructs an HTTP router
func NewRouter(server *Server) *mux.Router {
	router := mux.NewRouter()

	// Catch panics and return 500s
	router.Use(handlers.RecoveryHandler())

	// Optionally enable HTTP request logging
	if server.EnableRequestLogging {
		router.Use(server.loggingMiddleware)
	}

	router.HandleFunc("/.well-known/terraform.json", server.WellKnown)
	router.HandleFunc("/metrics/cache.json", server.CacheStats)

	router.HandleFunc("/state-versions/{id}/download", server.DownloadStateVersion).Methods("GET")
	router.HandleFunc("/configuration-versions/{id}/upload", server.UploadConfigurationVersion).Methods("PUT")
	router.HandleFunc("/plans/{id}/logs", server.GetPlanLogs).Methods("GET")
	router.HandleFunc("/plans/{id}/logs", server.UploadPlanLogs).Methods("PUT")
	router.HandleFunc("/applies/{id}/logs", server.GetApplyLogs).Methods("GET")
	router.HandleFunc("/applies/{id}/logs", server.UploadApplyLogs).Methods("PUT")
	router.HandleFunc("/runs/{id}/plan", server.UploadPlanFile).Methods("PUT")
	router.HandleFunc("/runs/{id}/plan", server.GetPlanFile).Methods("GET")

	router.HandleFunc("/app/settings/tokens", server.ListTokens).Methods("GET")
	router.HandleFunc("/app/settings/tokens", server.CreateToken).Methods("POST")
	router.HandleFunc("/app/settings/tokens/delete", server.DeleteToken).Methods("POST")
	router.HandleFunc("/healthz", server.Healthz).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(server.GetStaticFS()))).Methods("GET")

	router.HandleFunc("/app/{org}/{workspace}/runs/{id}", server.GetRunLogs).Methods("GET")

	// Websocket connections
	server.registerEventRoutes(router)

	// Filter json-api requests
	sub := router.Headers("Accept", jsonapi.MediaType).Subrouter()

	// Filter api v2 requests
	sub = sub.PathPrefix("/api/v2").Subrouter()

	// Require valid API token for all json-api requests
	sub.Use(newAuthMiddleware(server.TokenService).Middleware)

	sub.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Organization routes
	sub.HandleFunc("/organizations", server.ListOrganizations).Methods("GET")
	sub.HandleFunc("/organizations", server.CreateOrganization).Methods("POST")
	sub.HandleFunc("/organizations/{name}", server.GetOrganization).Methods("GET")
	sub.HandleFunc("/organizations/{name}", server.UpdateOrganization).Methods("PATCH")
	sub.HandleFunc("/organizations/{name}", server.DeleteOrganization).Methods("DELETE")
	sub.HandleFunc("/organizations/{name}/entitlement-set", server.GetEntitlements).Methods("GET")

	// Workspace routes
	sub.HandleFunc("/organizations/{org}/workspaces", server.ListWorkspaces).Methods("GET")
	sub.HandleFunc("/organizations/{org}/workspaces/{name}", server.GetWorkspace).Methods("GET")
	sub.HandleFunc("/organizations/{org}/workspaces", server.CreateWorkspace).Methods("POST")
	sub.HandleFunc("/organizations/{org}/workspaces/{name}", server.UpdateWorkspace).Methods("PATCH")
	sub.HandleFunc("/organizations/{org}/workspaces/{name}", server.DeleteWorkspace).Methods("DELETE")
	sub.HandleFunc("/workspaces/{id}", server.UpdateWorkspaceByID).Methods("PATCH")
	sub.HandleFunc("/workspaces/{id}", server.GetWorkspaceByID).Methods("GET")
	sub.HandleFunc("/workspaces/{id}", server.DeleteWorkspaceByID).Methods("DELETE")
	sub.HandleFunc("/workspaces/{id}/actions/lock", server.LockWorkspace).Methods("POST")
	sub.HandleFunc("/workspaces/{id}/actions/unlock", server.UnlockWorkspace).Methods("POST")

	// StateVersion routes
	sub.HandleFunc("/workspaces/{workspace_id}/state-versions", server.CreateStateVersion).Methods("POST")
	sub.HandleFunc("/workspaces/{workspace_id}/current-state-version", server.CurrentStateVersion).Methods("GET")
	sub.HandleFunc("/state-versions/{id}", server.GetStateVersion).Methods("GET")
	sub.HandleFunc("/state-versions", server.ListStateVersions).Methods("GET")

	// ConfigurationVersion routes
	sub.HandleFunc("/workspaces/{workspace_id}/configuration-versions", server.CreateConfigurationVersion).Methods("POST")
	sub.HandleFunc("/configuration-versions/{id}", server.GetConfigurationVersion).Methods("GET")
	sub.HandleFunc("/workspaces/{workspace_id}/configuration-versions", server.ListConfigurationVersions).Methods("GET")

	// Run routes
	sub.HandleFunc("/runs", server.CreateRun).Methods("POST")
	sub.HandleFunc("/runs/{id}/actions/apply", server.ApplyRun).Methods("POST")
	sub.HandleFunc("/workspaces/{workspace_id}/runs", server.ListRuns).Methods("GET")
	sub.HandleFunc("/runs/{id}", server.GetRun).Methods("GET")
	sub.HandleFunc("/runs/{id}/actions/discard", server.DiscardRun).Methods("POST")
	sub.HandleFunc("/runs/{id}/actions/cancel", server.CancelRun).Methods("POST")
	sub.HandleFunc("/runs/{id}/actions/force-cancel", server.ForceCancelRun).Methods("POST")
	sub.HandleFunc("/runs/{id}/plan/json-output", server.GetJSONPlanByRunID).Methods("GET")

	// Plan routes
	sub.HandleFunc("/plans/{id}", server.GetPlan).Methods("GET")
	sub.HandleFunc("/plans/{id}/json-output", server.GetPlanJSON).Methods("GET")

	// Apply routes
	sub.HandleFunc("/applies/{id}", server.GetApply).Methods("GET")

	return router
}

func (s *Server) SetupRoutes() {
	http.Handle("/", NewRouter(s))
}

// Open validates the server options and begins listening on the bind address.
func (s *Server) Open() (err error) {

	if s.ln, err = net.Listen("tcp", s.Addr); err != nil {
		return err
	}

	// Begin serving requests on the listener. We use Serve() instead of
	// ListenAndServe() because it allows us to check for listen errors (such as
	// trying to use an already open port) synchronously.
	go func() {
		if s.SSL {
			s.err <- s.server.ServeTLS(s.ln, s.CertFile, s.KeyFile)
		} else {
			s.err <- s.server.Serve(s.ln)
		}
	}()

	return nil
}

// Port returns the TCP port for the running server.  This is useful in tests
// where we allocate a random port by using ":0".
func (s *Server) Port() int {
	if s.ln == nil {
		return 0
	}
	return s.ln.Addr().(*net.TCPAddr).Port
}

// Wait blocks until server stops listening or context is cancelled.
func (s *Server) Wait(ctx context.Context) error {
	select {
	case err := <-s.err:
		return err
	case <-ctx.Done():
		return s.server.Close()
	}
}

// Close gracefully shuts down the server.
func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *Server) SetFlashMessage(w http.ResponseWriter, r *http.Request, msg string) error {
	session, _ := store.Get(r, FlashKey)
	session.AddFlash(msg)

	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("unable to save flash message: %w", err)
	}

	return nil
}

func (s *Server) GetFlashMessages(w http.ResponseWriter, r *http.Request) []template.HTML {
	session, _ := store.Get(r, FlashKey)

	flashes := strSliceToHTMLTemplateSlice(interfaceSliceToStringSlice(session.Flashes()))

	// Having read flashes we can now clear them from the session.
	session.Save(r, w)

	return flashes
}

// newLoggingMiddleware returns middleware that logs HTTP requests
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(next, w, r)

		s.Logger.Info("request",
			"duration", fmt.Sprintf("%dms", m.Duration.Milliseconds()),
			"status", m.Code,
			"method", r.Method,
			"path", fmt.Sprintf("%s?%s", r.URL.Path, r.URL.RawQuery))
	})
}

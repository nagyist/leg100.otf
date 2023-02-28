package configversion

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leg100/otf"
	"github.com/leg100/otf/http/decode"
	"github.com/leg100/otf/http/jsonapi"
	"github.com/leg100/surl"
)

type api struct {
	app          app
	otf.Verifier // for verifying upload url

	jsonapiMarshaler

	max int64 // Maximum permitted config upload size in bytes
}

func newAPI(opts apiOptions) *api {
	return &api{
		app:              opts.app,
		max:              opts.max,
		jsonapiMarshaler: jsonapiMarshaler{opts.Signer},
		Verifier:         opts.Signer,
	}
}

type apiOptions struct {
	app
	max int64
	*surl.Signer
}

func (s *api) AddHandlers(r *mux.Router) {
	// ConfigurationVersion routes
	r.HandleFunc("/workspaces/{workspace_id}/configuration-versions", s.CreateConfigurationVersion)
	r.HandleFunc("/configuration-versions/{id}", s.GetConfigurationVersion)
	r.HandleFunc("/workspaces/{workspace_id}/configuration-versions", s.ListConfigurationVersions)
	r.HandleFunc("/configuration-versions/{id}/download", s.DownloadConfigurationVersion)

	signed := r.PathPrefix("/signed/{signature.expiry}").Subrouter()
	signed.Use(otf.VerifySignedURL(s.Verifier))
	signed.HandleFunc("/configuration-versions/{id}/upload", s.UploadConfigurationVersion()).Methods("PUT")
}

func (s *api) CreateConfigurationVersion(w http.ResponseWriter, r *http.Request) {
	workspaceID, err := decode.Param("workspace_id", r)
	if err != nil {
		jsonapi.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	opts := jsonapi.ConfigurationVersionCreateOptions{}
	if err := jsonapi.UnmarshalPayload(r.Body, &opts); err != nil {
		jsonapi.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	cv, err := s.app.create(r.Context(), workspaceID, otf.ConfigurationVersionCreateOptions{
		AutoQueueRuns: opts.AutoQueueRuns,
		Speculative:   opts.Speculative,
	})
	if err != nil {
		jsonapi.Error(w, http.StatusNotFound, err)
		return
	}

	jsonapi.WriteResponse(w, r, s.toMarshalable(cv), jsonapi.WithCode(http.StatusCreated))
}

func (s *api) GetConfigurationVersion(w http.ResponseWriter, r *http.Request) {
	id, err := decode.Param("id", r)
	if err != nil {
		jsonapi.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	cv, err := s.app.get(r.Context(), id)
	if err != nil {
		jsonapi.Error(w, http.StatusNotFound, err)
		return
	}
	jsonapi.WriteResponse(w, r, s.toMarshalable(cv))
}

func (s *api) ListConfigurationVersions(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		WorkspaceID     string `schema:"workspace_id,required"`
		otf.ListOptions        // Pagination
	}
	var params parameters
	if err := decode.All(&params, r); err != nil {
		jsonapi.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	cvl, err := s.app.list(r.Context(), params.WorkspaceID, ConfigurationVersionListOptions{
		ListOptions: params.ListOptions,
	})
	if err != nil {
		jsonapi.Error(w, http.StatusNotFound, err)
		return
	}

	jsonapi.WriteResponse(w, r, s.toMarshableList(cvl))
}

func (s *api) UploadConfigurationVersion() http.HandlerFunc {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := decode.Param("id", r)
		if err != nil {
			jsonapi.Error(w, http.StatusUnprocessableEntity, err)
			return
		}

		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, r.Body); err != nil {
			jsonapi.Error(w, http.StatusUnprocessableEntity, err)
			return
		}
		if err := s.app.upload(r.Context(), id, buf.Bytes()); err != nil {
			jsonapi.Error(w, http.StatusNotFound, err)
			return
		}
	})
	return http.MaxBytesHandler(h, s.max).ServeHTTP
}

func (s *api) DownloadConfigurationVersion(w http.ResponseWriter, r *http.Request) {
	id, err := decode.Param("id", r)
	if err != nil {
		jsonapi.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	resp, err := s.app.download(r.Context(), id)
	if err != nil {
		jsonapi.Error(w, http.StatusNotFound, err)
		return
	}

	w.Write(resp)
}
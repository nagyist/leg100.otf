package http

import (
	"net/http"
	"strconv"

	"github.com/google/jsonapi"
	"github.com/leg100/ots"
)

const (
	DefaultPageNumber = 1
	DefaultPageSize   = 20
	MaxPageSize       = 100
)

// GetObject is a helper for getting an object and marshalling back to JSON-API
func GetObject(w http.ResponseWriter, r *http.Request, getter func() (interface{}, error)) {
	obj, err := getter()
	if err != nil {
		ErrNotFound(w)
		return
	}

	w.Header().Set("Content-type", jsonapi.MediaType)
	if err := MarshalPayload(w, r, obj); err != nil {
		ErrServerFailure(w, err)
	}
}

// ListObjects is a helper for listing objects and marshalling back to JSON-API
func ListObjects(w http.ResponseWriter, r *http.Request, lister func() (interface{}, error)) {
	obj, err := lister()
	if err != nil {
		ErrNotFound(w)
		return
	}

	w.Header().Set("Content-type", jsonapi.MediaType)
	if err := MarshalPayload(w, r, obj); err != nil {
		ErrServerFailure(w, err)
	}
}

// CreateObject is a helper for creating an object, unmarshalling and
// marshalling the request and response from/to JSON-API.
func CreateObject(w http.ResponseWriter, r *http.Request, opts interface{}, creator func(opts interface{}) (interface{}, error)) {
	if err := jsonapi.UnmarshalPayload(r.Body, opts); err != nil {
		ErrUnprocessable(w, err)
		return
	}

	obj, err := creator(opts)
	if err != nil {
		ErrNotFound(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayloadWithoutIncluded(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UpdateObject is a helper for updating an object, unmarshalling and
// marshalling the request and response from/to JSON-API.
func UpdateObject(w http.ResponseWriter, r *http.Request, opts interface{}, updater func(opts interface{}) (interface{}, error)) {
	if err := jsonapi.UnmarshalPayload(r.Body, opts); err != nil {
		ErrUnprocessable(w, err)
		return
	}

	obj, err := updater(opts)
	if err != nil {
		ErrNotFound(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayloadWithoutIncluded(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SanitizeListOptions(o *ots.ListOptions) {
	if o.PageNumber <= 0 {
		o.PageNumber = DefaultPageNumber
	}

	if o.PageSize <= 0 {
		o.PageSize = DefaultPageSize
	} else if o.PageSize > 100 {
		o.PageSize = MaxPageSize
	}
}

type ErrOption func(*jsonapi.ErrorObject)

func WithDetail(detail string) ErrOption {
	return func(err *jsonapi.ErrorObject) {
		err.Detail = detail
	}
}

func ErrNotFound(w http.ResponseWriter, opts ...ErrOption) {
	err := &jsonapi.ErrorObject{
		Status: "404",
		Title:  "not found",
	}

	for _, o := range opts {
		o(err)
	}

	w.WriteHeader(http.StatusNotFound)
	jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{
		err,
	})
}

func WriteError(w http.ResponseWriter, code int) {
	w.Header().Set("Content-type", jsonapi.MediaType)

	w.WriteHeader(code)
	jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{
		{
			Status: strconv.Itoa(code),
		},
	})
}

func ErrUnprocessable(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
		Status: "422",
		Title:  "unable to process payload",
		Detail: err.Error(),
	}})
}

func ErrServerFailure(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
		Status: "500",
		Title:  "internal server failure",
		Detail: err.Error(),
	}})
}

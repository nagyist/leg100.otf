package http

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var (
	// Query schema Encoder, caches structs, and safe for sharing
	Encoder = schema.NewEncoder()
	// Query schema Decoder: caches structs, and safe for sharing.
	Decoder = schema.NewDecoder()
)

// decodeAll collectively decodes route params, query params and form params into
// the obj struct
func decodeAll(r *http.Request, obj interface{}) error {
	if err := decodeForm(r, obj); err != nil {
		return err
	}

	if err := decodeQuery(r, obj); err != nil {
		return err
	}

	if err := decodeRouteVars(r, obj); err != nil {
		return err
	}

	return nil
}

func decodeForm(r *http.Request, obj interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if err := decoder.Decode(obj, r.PostForm); err != nil {
		return err
	}

	return nil
}

// DecodeQuery unmarshals a query string (k1=v1&k2=v2...) into a struct.
func DecodeQuery(opts interface{}, query url.Values) error {
	if err := Decoder.Decode(opts, query); err != nil {
		return fmt.Errorf("unable to decode query string: %w", err)
	}
	return nil
}

func DecodeRoute(obj interface{}, r *http.Request) error {
	// decoder only takes map[string][]string, not map[string]string
	vars := convertStrMapToStrSliceMap(mux.Vars(r))

	if err := Decoder.Decode(obj, vars); err != nil {
		return err
	}

	return nil
}

func convertStrMapToStrSliceMap(m map[string]string) map[string][]string {
	mm := make(map[string][]string, len(m))
	for k, v := range m {
		mm[k] = []string{v}
	}
	return mm
}

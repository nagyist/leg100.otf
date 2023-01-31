// Package state manages terraform state.
package state

import (
	"encoding/base64"
	"encoding/json"

	"github.com/go-logr/logr"
	"github.com/leg100/otf"
)

const (
	DefaultStateVersion = 4
)

var _ otf.StateVersionService = (*service)(nil)

type service struct {
	*app
	*handlers
}

func NewService(opts ServiceOptions) *service {
	app := &app{
		Authorizer: opts.Authorizer,
		Logger:     opts.Logger,
		db:         newPGDB(opts.Database),
		cache:      opts.Cache,
	}
	return &service{
		app:      app,
		handlers: &handlers{app},
	}
}

type ServiceOptions struct {
	otf.Authorizer
	otf.Database
	otf.Cache
	logr.Logger
}

// State is terraform state.
type State struct {
	Version int
	Serial  int64
	Lineage string
	Outputs map[string]StateOutput
}

// StateOutput is a terraform state output.
type StateOutput struct {
	Name      string
	Value     string
	Type      string
	Sensitive bool
}

// StateCreateOptions are options for creating state
type StateCreateOptions struct {
	Version *int
	Serial  *int64
	Lineage *string
}

// unmarshalState unmarshals terraform state from a raw byte slice.
func unmarshalState(data []byte) (*State, error) {
	state := State{}
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

// NewState constructs a new state
func NewState(opts StateCreateOptions, outputs ...StateOutput) *State {
	state := State{
		Version: DefaultStateVersion,
		Serial:  1,
	}
	if opts.Lineage != nil {
		state.Lineage = *opts.Lineage
	}
	if opts.Serial != nil {
		state.Serial = *opts.Serial
	}
	if opts.Version != nil {
		state.Version = *opts.Version
	}
	state.Outputs = make(map[string]StateOutput, len(outputs))
	for _, out := range outputs {
		state.Outputs[out.Name] = out
	}
	return &state
}

// Marshal serializes state as a base64-encoded json string.
func (s *State) Marshal() (string, error) {
	js, err := json.Marshal(s)
	if err != nil {
		return "", nil
	}
	return base64.StdEncoding.EncodeToString(js), nil
}
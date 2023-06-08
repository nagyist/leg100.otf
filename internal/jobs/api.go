package jobs

import (
	"net/http"

	"github.com/gorilla/mux"
	otfhttp "github.com/leg100/otf/internal/http"
)

type api struct {
	Service
}

func (a *api) addHandlers(r *mux.Router) {
	r = otfhttp.APIRouter(r)
	r.HandleFunc("/agent/jobs", a.getAssignedJob).Methods("GET")
}

// long polling...
func (a *api) getAssignedJob(w http.ResponseWriter, r *http.Request) {
}

package v1

import (
	"encoding/json"
	"net/http"

	"github.com/dimfeld/httptreemux"
)

type Responce struct {
	Status string
}

type healthImpl struct{}

func registerHeatlthGroup(router *httptreemux.Group) {
	impl := &healthImpl{}

	router.GET("/health", impl.health)
}

func (impl *healthImpl) health(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	o := Responce{"OK"}

	js, err := json.Marshal(o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

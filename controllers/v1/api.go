package v1

import (
	"net/http"

	"github.com/suarezgary/GolangApi/utils/jsonhttp"
)

// API API
func API(w http.ResponseWriter, r *http.Request) {
	// TODO check service health, send some metrics, etc.
	jsonhttp.JSONSuccess(w, nil, "Server healthy")
}

package v1

import (
	"net/http"

	"github.com/suarezgary/GolangApi/models"
	"github.com/suarezgary/GolangApi/utils/jsonhttp"
	"github.com/suarezgary/GolangApi/utils/queryparams"
)

// GetGophers GetGophers
func GetGophers(w http.ResponseWriter, r *http.Request) {
	limit, offset := queryparams.GetLimitOffsetQueryParametersDefaults(r)

	gophers, err := models.GetGophers(limit, offset)
	if err != nil {
		jsonhttp.JSONNotFoundError(w, "Error fetching gophers", "")
		return
	}

	jsonhttp.JSONSuccess(w, map[string]interface{}{"gophers": gophers}, "Successfully queried gophers")
}

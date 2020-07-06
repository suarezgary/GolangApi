package httpmiddleware

import (
	"net/http"

	"github.com/suarezgary/GolangApi/utils/jsonhttp"
)

// RecoverInternalServerError RecoverInternalServerError
func RecoverInternalServerError(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			if r != nil {
				jsonhttp.JSONInternalError(w, "Internal Server error", "")
			}
		}()
		handler.ServeHTTP(w, r)
	})
}

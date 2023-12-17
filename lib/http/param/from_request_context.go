package param

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// FromRequestContext can be used to decouple the handler from the julienschmidt/httprouter package.
func FromRequestContext(r *http.Request, key string) string {
	return httprouter.ParamsFromContext(r.Context()).ByName(key)
}

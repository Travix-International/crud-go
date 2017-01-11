package crud

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	service "github.com/Travix-International/go-servicefoundation"
)

// ActionNotAvailableHandler should be routed for CRUD actions that do not need to be implemented for a particular resource. It basically prevents
// a 404 and turns it into a 405 instead.
func ActionNotAvailableHandler(w http.ResponseWriter, r *http.Request) {
	ww := service.NewWrappedResponseWriter(w)
	ww.WriteHeader(http.StatusMethodNotAllowed)
}

// Recovery is a sample recovery function for CRUD operations. It logs the error, and writes a 500 to the output, according to the CRUD protocol.
func Recovery(name string, ctx service.AppContext, w http.ResponseWriter, r *http.Request) {
	if rec := recover(); rec != nil {
		ctx.Logger().Error("CrudRecovery", fmt.Sprintf("PANIC recovered for %v: %v\n%v", name, rec, string(debug.Stack())))
		WriteOperationResult(w, r, ErrorResult(errors.New("Panic while handling request")))
	}
}

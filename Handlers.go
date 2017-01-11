package crud

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"strings"

	"strconv"

	servicefoundation "github.com/Travix-International/go-servicefoundation"
	"github.com/gorilla/mux"
)

// RecoverFunc is the recovery function as needed by the CRUD functionality
type RecoverFunc func(name string, ctx servicefoundation.AppContext, w http.ResponseWriter, r *http.Request)

// CreateCrudHandlerGetList is used to request a list of entities
var CreateCrudHandlerGetList = func(ctx servicefoundation.AppContext, svc Service, resourceName string, recoverFunc RecoverFunc) http.HandlerFunc {
	logger := ctx.Logger()
	return func(w http.ResponseWriter, r *http.Request) {
		defer recoverFunc("crudHandler", ctx, w, r)
		logger.Debug("CrudHandlerStart", fmt.Sprintf("CRUD operation %v requested on %v", r.Method, r.URL.Path))

		dsRequest := ExtractDataSetRequestFromURI(r)

		logger.Debug("CreateCrudHandlerGetList", fmt.Sprintf("Interpreted as GetList command. Arguments: %v", dsRequest))
		opResult := svc.GetAll(dsRequest)
		WriteOperationResult(w, r, opResult)
	}
}

// CreateCrudHandlerGetByID is used to get an item by its identifier
var CreateCrudHandlerGetByID = func(ctx servicefoundation.AppContext, svc Service, resourceName string, recoverFunc RecoverFunc) http.HandlerFunc {
	logger := ctx.Logger()
	return func(w http.ResponseWriter, r *http.Request) {
		defer recoverFunc("crudHandler", ctx, w, r)

		vars := mux.Vars(r)
		idVar := vars["id"]

		logger.Debug("CreateCrudHandlerGetByID", fmt.Sprintf("Interpreted as GetById command. ID: %v", idVar))
		opResult := svc.GetByID(idVar)
		WriteOperationResult(w, r, opResult)
	}
}

// CreateCrudHandlerCreateEntity is used to create a new entity
var CreateCrudHandlerCreateEntity = func(ctx servicefoundation.AppContext,
	svc Service, resourceName string, recoverFunc RecoverFunc,
	createFunc func() Entity) http.HandlerFunc {
	logger := ctx.Logger()
	return func(w http.ResponseWriter, r *http.Request) {
		defer recoverFunc("crudHandler", ctx, w, r)

		// Read the entity from the HTTP body
		entity := createFunc()
		err := ReadEntityFromBody(r.Body, entity)
		if err != nil {
			opResult := ValidationFailedResult(errors.New("Failed to parse entity from HTTP body: " + err.Error()))
			WriteOperationResult(w, r, opResult)
			return
		}

		logger.Debug("CreateCrudHandlerCreate", fmt.Sprintf("Interpreted as create command. Entity: %s", entity))
		opResult := svc.Add(entity)
		WriteOperationResult(w, r, opResult)
	}
}

// CreateCrudHandlerDeleteByID is used to delete an item by its identifier
var CreateCrudHandlerDeleteByID = func(ctx servicefoundation.AppContext, svc Service, resourceName string, recoverFunc RecoverFunc) http.HandlerFunc {
	logger := ctx.Logger()
	return func(w http.ResponseWriter, r *http.Request) {
		defer recoverFunc("crudHandler", ctx, w, r)

		vars := mux.Vars(r)
		idVar := vars["id"]

		logger.Debug("CreateCrudHandlerDeleteById", fmt.Sprintf("Interpreted as DeleteByID command. ID: %v", idVar))
		opResult := svc.Delete(idVar)
		WriteOperationResult(w, r, opResult)
	}
}

// CreateCrudHandlerUpdateEntity is used to update a new entity
var CreateCrudHandlerUpdateEntity = func(ctx servicefoundation.AppContext,
	svc Service, resourceName string, recoverFunc RecoverFunc,
	createFunc func() Entity) http.HandlerFunc {
	logger := ctx.Logger()
	return func(w http.ResponseWriter, r *http.Request) {
		defer recoverFunc("crudHandler", ctx, w, r)

		// Read the ID from the path
		vars := mux.Vars(r)
		idVar := vars["id"]

		// Read the entity from the HTTP body
		entity := createFunc()
		err := ReadEntityFromBody(r.Body, entity)
		if err != nil {
			opResult := ValidationFailedResult(errors.New("Failed to parse entity from HTTP body: " + err.Error()))
			WriteOperationResult(w, r, opResult)
			return
		}

		logger.Debug("CreateCrudHandlerUpdateEntity", fmt.Sprintf("Interpreted as update command. Id: %v Entity: %s", idVar, entity))
		opResult := svc.Update(idVar, entity)
		WriteOperationResult(w, r, opResult)
	}
}

// ReadEntityFromBody is a utility function to parse the given entity from the request body
func ReadEntityFromBody(body io.ReadCloser, entity interface{}) error {
	decoder := json.NewDecoder(body)

	defer body.Close()
	if err := decoder.Decode(&entity); err != nil {
		return err
	}

	return nil
}

// WriteOperationResult is a utility function that takes the result of a CRUD operation, and writes
// the corresponding HTTP response, according to the CRUD protocol.
func WriteOperationResult(w http.ResponseWriter, r *http.Request, opResult OperationResult) {

	// Determine HTTP status code, plus the response object
	statusCode := http.StatusInternalServerError
	var responseObject interface{}
	if opResult.Error() != nil {
		responseObject = &ErrorDetails{Message: opResult.Error().Error()}
	}

	switch opResult.State() {
	case Ok:
		statusCode = http.StatusOK
		if opResult.Value() != nil {
			responseObject = opResult.Value()
		}
	case Created:
		statusCode = http.StatusCreated
	case ValidationFailed:
		statusCode = http.StatusBadRequest
	case NotFound:
		statusCode = http.StatusNotFound
		responseObject = nil
	case Conflict:
		statusCode = http.StatusConflict
	case NotSupportedByResource:
		statusCode = http.StatusMethodNotAllowed
	case Error:
		// Nothing to do, accept defaults
	}

	ww := servicefoundation.NewWrappedResponseWriter(w)
	if responseObject == nil {
		ww.WriteHeader(statusCode)
	} else {
		servicefoundation.WriteResponse(ww, r, statusCode, responseObject)
	}
}

// ExtractDataSetRequestFromURI is a helper function to parse URI parameters into a DataSet request. Any
// parameters that are omitted from the URL are substituted with default values.
func ExtractDataSetRequestFromURI(r *http.Request) *DataSetRequest {
	dsReq := &DataSetRequest{
		PageSize:      15,
		PageNumber:    1,
		SortColumn:    "id",
		SortDirection: string(Asc),
		Filters:       nil,
	}

	query := r.URL.Query()
	for k, v := range query {
		switch strings.ToLower(k) {
		case "pagesize":
			intVal, err := strconv.Atoi(v[0])
			if err == nil && intVal > 0 {
				dsReq.PageSize = intVal
			}
		case "pagenumber":
			intVal, err := strconv.Atoi(v[0])
			if err == nil && intVal > 0 {
				dsReq.PageNumber = intVal
			}
		case "sortcolumn":
			dsReq.SortColumn = v[0]
		case "sortdirection":
			if strings.ToLower(v[0]) == "desc" {
				dsReq.SortDirection = string(Desc)
			}
		case "filters":
			// is encoded json, actually
			jsonText := v[0]
			filterKvs := make(map[string]string)
			err := json.Unmarshal([]byte(jsonText), filterKvs)
			if err == nil {
				dsReq.Filters = filterKvs
			}
		}
	}

	return dsReq
}

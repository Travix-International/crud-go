package crud

// SortDirection describes how results are sorted
type SortDirection string

const (
	// Asc sorts ascending
	Asc SortDirection = "Asc"
	// Desc sorts decending
	Desc SortDirection = "Desc"
)

// State is used to describe the high level state of a requested CRUD operation.
type State int

const (
	// Ok means that no issues were found
	Ok State = 1
	// ValidationFailed means that something invalid in the request was encountered
	ValidationFailed State = 2
	// Error is a generic error, such as an unhandled error. Use when none of the other states are applicable.
	Error State = 3
	// NotFound means that the entity with the given identity does not exist (anymore).
	NotFound State = 4
	// Conflict means that the requested operation was understood, but conflicts with the current state of the entity.
	Conflict State = 5
	// NotSupportedByResource means that this operation is not supported - for example, we're trying to delete an item,
	// but the resource doesn't support deletes at all (example: a read-only lookup table).
	NotSupportedByResource State = 6
	// Created is the same as OK, but used to signify that the entity was created
	Created State = 7
)

// PagingInfo is used to describe how results are being pages
type PagingInfo struct {
	// SupportsPaging indicates whether the datasource even supports paging. If false, it means that all the applicable
	// results are rendered on one page, page #1.
	SupportsPaging bool `json:"supportsPaging"`
	// DoesKnowTotalRecords indicates whether the datasource can accurately provide the exact amount of items in the datasource.
	// if false, paging might still be available (see: SupportsPaging).
	DoesKnowTotalRecords bool `json:"doesKnowTotalRecords"`
	// PageSize is the maximum number of items per page. Minimum 1.
	PageSize int `json:"pageSize"`
	// PageNumber is the number of the page. One-based.
	PageNumber int `json:"pageNumber"`
	// TotalRecordsCount indicates exactly the amount of items found. Will be zero, is DoesKnowTotalRecords is false.
	TotalRecordsCount int `json:"totalRecordCount"`
}

// DataSetRequest describes the parameters used to search for a set of results. It describes things like the desired page, sorting, filtering, etc.
type DataSetRequest struct {
	// PageSize is the amount of items that is requested in one page. Note: the datasource will reply with the actual page size. This can happen when paging
	// is for some reason not available, is limited, or when an invalid page size is requested. Minimum of 1 item per page. Ignored if the datasource does not
	// support paging.
	PageSize int `json:"pageSize"`
	// PageNumber is the one-based number of the page being requested.
	PageNumber int `json:"pageNumber"`
	// SortColumn is the name of the column to sort the results by. Required field, defaults to "id".
	SortColumn string `json:"sortColumn"`
	// SortDirection describes how to sort the results. Required field.
	SortDirection string `json:"sortDirection"`
	// Filters are key/value pairs that describe which fields to apply. Each key is a column name, while the values are what to search for in those columns.
	//
	// Note that the exact filter implementation is fully dependent on the datasource, and this protocol does not guarantee how the filtering is applied. E.g. when searching
	// for a string, the implementer can choose whether this is 'start with', 'contains', etc. If the datasource does not support filtering on a given column, then that
	// filter is to be ignored.
	Filters map[string]string `json:"filters"`
}

// DataSet is a set of items
type DataSet struct {
	// Items is the collection of results. Filtering and paging (if supported) have been applied, and the results are sorted.
	Items []interface{} `json:"items"`
	// PagingInfo describes which page of results is being returned.
	PagingInfo PagingInfo `json:"pagingInfo"`
}

// ErrorDetails is used to communicate the details of an error back to the caller
type ErrorDetails struct {
	Message string `json:"message"`
}

// OperationResult is the result of a CRUD operation
type OperationResult interface {
	// State is the high level result of the operation - further details are to be returned in separate fields.
	State() State
	// Error can be nil, if no error occurred
	Error() error
	// Value can be nil in case of errors, or when the operation doesn't return a value.
	Value() interface{}
}

// Service is the interface that needs to be implemented to provide the actual implementation to do CRUD on a resource.
//
// Note that not all the operations need to be implemented
type Service interface {
	// GetAll is used to get a list of entities, optionally applying paging, filtering, sorting.
	GetAll(request *DataSetRequest) OperationResult

	// GetByID returns the entity with the specified ID.
	GetByID(id EntityKey) OperationResult

	// Add will add the given entity. The returned value is the ID of the new entity.
	Add(entity Entity) OperationResult

	// Update will update an existing entity. The returned value is the new state of the entity.
	Update(id EntityKey, entity Entity) OperationResult

	// Delete will delete the entity with the specified ID.
	Delete(id EntityKey) OperationResult
}

// Entity describes what CRUD entities must implement to be used in the standardized CRUD functionality
type Entity interface {
	// Validates an entity to see if it's correct
	Validate() error
	// Format is used as a way to 'fix' things like casing. It'll help make it friendlier for the end user, and/or
	// give a chance to set calculated fields and stuff
	//
	// isNewEntity will be true is this is a completely new entity, as opposed to an existing one
	Format(isNewEntity bool)
}

// EntityKey is the unique key of an entity. Typically an int, string, etc.
type EntityKey interface{}

type crudOperationResult struct {
	OperationResult
	state State
	error error
	value interface{}
}

func (c *crudOperationResult) State() State {
	return c.state
}

func (c *crudOperationResult) Error() error {
	return c.error
}

func (c *crudOperationResult) Value() interface{} {
	return c.value
}

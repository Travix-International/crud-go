package crud

// OkResult constructs an operation result for State 'Ok'.
//
// value can be nil.
func OkResult(value interface{}) OperationResult {
	result := &crudOperationResult{
		state: Ok,
		error: nil,
		value: value,
	}
	return result
}

// CreatedResult constructs an operation result for State 'Created'.
func CreatedResult() OperationResult {
	result := &crudOperationResult{
		state: Created,
		error: nil,
		value: nil,
	}
	return result
}

// ErrorResult constructs an operation result for State 'Error'.
func ErrorResult(err error) OperationResult {
	result := &crudOperationResult{
		state: Error,
		error: err,
		value: nil,
	}
	return result
}

// ValidationFailedResult constructs an operation result for State 'ValidationFailed'.
func ValidationFailedResult(err error) OperationResult {
	result := &crudOperationResult{
		state: ValidationFailed,
		error: err,
		value: nil,
	}
	return result
}

// ConflictResult constructs an operation result for State 'Conflict'.
func ConflictResult(err error) OperationResult {
	result := &crudOperationResult{
		state: Conflict,
		error: err,
		value: nil,
	}
	return result
}

// NotFoundResult constructs an operation result for State 'NotFound'.
func NotFoundResult() OperationResult {
	result := &crudOperationResult{
		state: NotFound,
		error: nil,
		value: nil,
	}
	return result
}

// NotSupportedByResourceResult constructs an operation result for State 'NotSupportedByResource'.
func NotSupportedByResourceResult() OperationResult {
	result := &crudOperationResult{
		state: NotSupportedByResource,
		error: nil,
		value: nil,
	}
	return result
}

// ConstrainPagingRequest clips the request values to reasonable amounts
func ConstrainPagingRequest(r *DataSetRequest, minPageSize, maxPageSize int) {
	if r.PageSize < minPageSize {
		r.PageSize = minPageSize
	}
	if r.PageSize > maxPageSize {
		r.PageSize = maxPageSize
	}
}

// ConstrainSortColumns restricts the sorting of a request to the allowed columns, defaulting to a specific one
func ConstrainSortColumns(r *DataSetRequest, defaultColumn string, allowedColumns ...string) {
	for _, v := range allowedColumns {
		if v == r.SortColumn {
			return
		}
	}
	r.SortColumn = defaultColumn
}

// ConstrainFilterColumns restricts the filtering of a request to the allowed columns, removing invalid ones
func ConstrainFilterColumns(r *DataSetRequest, allowedColumns ...string) {
	if r.Filters == nil || len(r.Filters) == 0 {
		return
	}
	// Simple linear stuff, but these are small collections
	for k := range r.Filters {
		isAllowed := false
		for _, v := range allowedColumns {
			if v == k {
				isAllowed = true
				break
			}
		}
		if !isAllowed || len(r.Filters[k]) == 0 {
			delete(r.Filters, k)
		}
	}
}

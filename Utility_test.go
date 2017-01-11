package crud_test

import (
	"errors"
	"testing"

	crud "github.com/Travix-International/crud-go"
	"github.com/stretchr/testify/assert"
)

func TestOkResultNotNil(t *testing.T) {
	value := "hello"
	result := crud.OkResult(value)
	assert.NotNil(t, result)

	assert.Equal(t, crud.Ok, result.State())
	assert.Nil(t, result.Error())
	assert.Equal(t, value, result.Value())
}

func TestOkResultNil(t *testing.T) {
	result := crud.OkResult(nil)
	assert.NotNil(t, result)

	assert.Equal(t, crud.Ok, result.State())
	assert.Nil(t, result.Error())
	assert.Nil(t, result.Value())
}

func TestCreatedResult(t *testing.T) {
	result := crud.CreatedResult()
	assert.NotNil(t, result)

	assert.Equal(t, crud.Created, result.State())
	assert.Nil(t, result.Error())
	assert.Nil(t, result.Value())
}

func TestErrorResultNotNil(t *testing.T) {
	value := errors.New("Sample error")
	result := crud.ErrorResult(value)
	assert.NotNil(t, result)

	assert.Equal(t, crud.Error, result.State())
	assert.Nil(t, result.Value())
	assert.Equal(t, value, result.Error())
}

func TestValidationFailedResult(t *testing.T) {
	value := errors.New("Sample error")
	result := crud.ValidationFailedResult(value)
	assert.NotNil(t, result)

	assert.Equal(t, crud.ValidationFailed, result.State())
	assert.Nil(t, result.Value())
	assert.Equal(t, value, result.Error())
}

func TestConflictResult(t *testing.T) {
	value := errors.New("Sample error")
	result := crud.ConflictResult(value)
	assert.NotNil(t, result)

	assert.Equal(t, crud.Conflict, result.State())
	assert.Nil(t, result.Value())
	assert.Equal(t, value, result.Error())
}

func TestNotFoundResult(t *testing.T) {
	result := crud.NotFoundResult()
	assert.NotNil(t, result)

	assert.Equal(t, crud.NotFound, result.State())
	assert.Nil(t, result.Value())
	assert.Nil(t, result.Error())
}

func TestNotSupportedByResourceResult(t *testing.T) {
	result := crud.NotSupportedByResourceResult()
	assert.NotNil(t, result)

	assert.Equal(t, crud.NotSupportedByResource, result.State())
	assert.Nil(t, result.Value())
	assert.Nil(t, result.Error())
}

func TestConstrainPagingRequestMin(t *testing.T) {
	ds := &crud.DataSetRequest{
		PageSize: 0,
	}
	crud.ConstrainPagingRequest(ds, 1, 100)
	assert.Equal(t, 1, ds.PageSize)
}

func TestConstrainPagingRequestMax(t *testing.T) {
	ds := &crud.DataSetRequest{
		PageSize: 105,
	}
	crud.ConstrainPagingRequest(ds, 1, 100)
	assert.Equal(t, 100, ds.PageSize)
}

func TestConstrainPagingRequestInBounds(t *testing.T) {
	ds := &crud.DataSetRequest{
		PageSize: 50,
	}
	crud.ConstrainPagingRequest(ds, 1, 100)
	assert.Equal(t, 50, ds.PageSize)
}

func TestConstrainSortColumns_EmptyInput(t *testing.T) {
	ds := &crud.DataSetRequest{
		SortColumn: "",
	}
	crud.ConstrainSortColumns(ds, "col1", "col1", "col2")
	assert.Equal(t, "col1", ds.SortColumn)
}

func TestConstrainSortColumns_IncorrectInput(t *testing.T) {
	ds := &crud.DataSetRequest{
		SortColumn: "not-allowed-col",
	}
	crud.ConstrainSortColumns(ds, "col1", "col1", "col2")
	assert.Equal(t, "col1", ds.SortColumn)
}

func TestConstrainSortColumns_CorrectInput(t *testing.T) {
	ds := &crud.DataSetRequest{
		SortColumn: "col2",
	}
	crud.ConstrainSortColumns(ds, "col1", "col1", "col2")
	assert.Equal(t, "col2", ds.SortColumn)
}

func TestConstrainFilterColumns_NilFiltersInput(t *testing.T) {
	ds := &crud.DataSetRequest{
		Filters: nil,
	}
	crud.ConstrainFilterColumns(ds, "col1", "col2")
	assert.Nil(t, ds.Filters)
}

func TestConstrainFilterColumns_EmptyFiltersInput(t *testing.T) {
	ds := &crud.DataSetRequest{
		Filters: make(map[string]string, 0),
	}
	crud.ConstrainFilterColumns(ds, "col1", "col2")
	assert.NotNil(t, ds.Filters)
	assert.Equal(t, 0, len(ds.Filters))
}

func TestConstrainFilterColumns_DisallowedFilters(t *testing.T) {
	ds := &crud.DataSetRequest{
		Filters: map[string]string{
			"colnotallowed": "somefiltervalue",
		},
	}
	crud.ConstrainFilterColumns(ds, "col1", "col2")
	assert.NotNil(t, ds.Filters)
	assert.Equal(t, 0, len(ds.Filters))
}

func TestConstrainFilterColumns_AllowedFilters(t *testing.T) {
	ds := &crud.DataSetRequest{
		Filters: map[string]string{
			"col1": "somefiltervalue",
			"col2": "somefiltervalue",
		},
	}
	crud.ConstrainFilterColumns(ds, "col1", "col2")
	assert.NotNil(t, ds.Filters)
	assert.Equal(t, 2, len(ds.Filters))
}

func TestConstrainFilterColumns_EmptyFilters(t *testing.T) {
	ds := &crud.DataSetRequest{
		Filters: map[string]string{
			"col1": "",
		},
	}
	crud.ConstrainFilterColumns(ds, "col1", "col2")
	assert.NotNil(t, ds.Filters)
	assert.Equal(t, 0, len(ds.Filters))
}

package crud_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	crud "github.com/Travix-International/crud-go"
	servicefoundation "github.com/Travix-International/go-servicefoundation"
	"github.com/Travix-International/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type crudServiceMock struct {
	crud.Service
	mock.Mock
}

func (m *crudServiceMock) GetAll(request *crud.DataSetRequest) crud.OperationResult {
	m.Called(request)
	return crud.OkResult(&crud.DataSet{})
}

func TestCreateCrudHandlerGetList_WithoutParameters(t *testing.T) {
	loggy, _ := logger.New(make(map[string]string))
	ctx := &servicefoundation.ContextBase{}
	ctx.SetLogger(loggy)
	crudService := new(crudServiceMock)
	resourceName := "myresource"
	recovery := func(name string, ctx servicefoundation.AppContext, w http.ResponseWriter, r *http.Request) {
		// todo: what goes here?
	}

	crudService.On("GetAll", mock.Anything).Once()

	r, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	fn := crud.CreateCrudHandlerGetList(ctx, crudService, resourceName, recovery)

	//act
	fn(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	crudService.AssertExpectations(t)
}

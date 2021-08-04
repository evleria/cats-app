package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/evleria/mongo-crud/internal/repository"
	"github.com/evleria/mongo-crud/internal/repository/entities"
	"github.com/evleria/mongo-crud/internal/service"
)

var (
	mockContext  = mock.Anything
	errSomeError = errors.New("some error")
	bella        = entities.Cat{
		ID:    uuid.New(),
		Name:  "Ms. Bella",
		Color: "brown",
		Age:   4,
		Price: 9.99,
	}
	zorro = entities.Cat{
		ID:    uuid.New(),
		Name:  "Mr. Zorro",
		Color: "black",
		Age:   8,
		Price: 10.99,
	}
	cats = []entities.Cat{bella, zorro}
)

func TestGetAllCats(t *testing.T) {
	// Arrange
	r := new(repository.MockCats)
	r.On("GetAll", mockContext).Return(cats, nil)
	ctx, rec := setup(http.MethodGet, nil)

	// Act
	err := GetAllCats(r)(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, mustEncodeJSON(mapCats(cats)), rec.Body.String())
}

func TestGetAllCatsRepositoryFailed(t *testing.T) {
	// Arrange
	r := new(repository.MockCats)
	r.On("GetAll", mockContext).Return(nil, errSomeError)
	ctx, _ := setup(http.MethodGet, nil)

	// Act
	err := GetAllCats(r)(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, errSomeError.Error()), err)
}

func TestUpdatePrice(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	id := bella.ID
	req := UpdatePriceRequest{Price: 5.99}
	s.On("UpdatePrice", mockContext, id, req.Price).
		Return(nil)
	ctx, rec := setup(http.MethodGet, req)
	ctx.SetParamNames("id")
	ctx.SetParamValues(id.String())

	// Act
	err := UpdatePrice(s)(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "", rec.Body.String())
}

func setup(method string, body interface{}) (echo.Context, *httptest.ResponseRecorder) {
	jsonBody := ""
	if body != nil {
		jsonBody = mustEncodeJSON(body)
	}
	request := httptest.NewRequest(method, "/", strings.NewReader(jsonBody))
	if body != nil {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)
	return c, recorder
}

func mustEncodeJSON(data interface{}) string {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	err := enc.Encode(data)
	if err != nil {
		panic(err)
	}
	return buffer.String()
}

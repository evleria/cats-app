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

func TestGetCat(t *testing.T) {
	// Arrange
	r := new(repository.MockCats)
	id := bella.ID
	r.On("GetOne", mockContext, id).Return(bella, nil)
	ctx, rec := setup(http.MethodGet, nil)
	ctx.SetParamNames("id")
	ctx.SetParamValues(id.String())

	// Act
	err := GetCat(r)(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, mustEncodeJSON(mapCat(bella)), rec.Body.String())
}

func TestGetCatMalformedId(t *testing.T) {
	// Arrange
	r := new(repository.MockCats)
	ctx, _ := setup(http.MethodGet, nil)
	ctx.SetParamNames("id")
	ctx.SetParamValues("malformed-uuid")

	// Act
	err := GetCat(r)(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusBadRequest), err)
}

func TestGetCatNotFound(t *testing.T) {
	// Arrange
	r := new(repository.MockCats)
	id := uuid.New().String()
	r.On("GetOne", mockContext, mock.AnythingOfType("uuid.UUID")).Return(entities.Cat{}, repository.ErrNotFound)
	ctx, _ := setup(http.MethodGet, nil)
	ctx.SetParamNames("id")
	ctx.SetParamValues(id)

	// Act
	err := GetCat(r)(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusNotFound), err)
}

func TestGetCatRepositoryFailed(t *testing.T) {
	// Arrange
	r := new(repository.MockCats)
	id := uuid.New().String()
	r.On("GetOne", mockContext, mock.AnythingOfType("uuid.UUID")).Return(entities.Cat{}, errSomeError)
	ctx, _ := setup(http.MethodGet, nil)
	ctx.SetParamNames("id")
	ctx.SetParamValues(id)

	// Act
	err := GetCat(r)(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, errSomeError.Error()), err)
}

func TestAddNewCat(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	req := AddNewCatRequest{"Mila", "black", 5, 7.99}
	id := uuid.New()
	s.On("CreateNew", mockContext, req.Name, req.Color, req.Age, req.Price).Return(id, nil)
	ctx, rec := setup(http.MethodPost, req)

	// Act
	err := AddNewCat(s)(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, mustEncodeJSON(AddNewCatResponse{id.String()}), rec.Body.String())
}

func TestAddNewCatServiceFailed(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	req := AddNewCatRequest{"Mila", "black", 5, 7.99}
	s.On("CreateNew", mockContext, req.Name, req.Color, req.Age, req.Price).Return(nil, errSomeError)
	ctx, _ := setup(http.MethodPost, req)

	// Act
	err := AddNewCat(s)(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, errSomeError.Error()), err)
}

func TestDeleteCat(t *testing.T) {
	// Arrange
	s := new(repository.MockCats)
	id := bella.ID
	s.On("Delete", mockContext, id).Return(nil)
	ctx, rec := setup(http.MethodDelete, nil)
	ctx.SetParamNames("id")
	ctx.SetParamValues(id.String())

	// Act
	err := DeleteCat(s)(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "", rec.Body.String())
}

func TestDeleteCatFailed(t *testing.T) {
	// Arrange
	r := new(repository.MockCats)
	ctx, _ := setup(http.MethodDelete, nil)
	ctx.SetParamNames("id")
	ctx.SetParamValues("malformed-uuid")

	// Act
	err := DeleteCat(r)(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusBadRequest), err)
}

func TestDeleteCatNotFound(t *testing.T) {
	// Arrange
	r := new(repository.MockCats)
	id := uuid.New().String()
	r.On("Delete", mockContext, mock.AnythingOfType("uuid.UUID")).Return(repository.ErrNotFound)
	ctx, _ := setup(http.MethodDelete, nil)
	ctx.SetParamNames("id")
	ctx.SetParamValues(id)

	// Act
	err := DeleteCat(r)(ctx)

	// Assert
	require.Error(t, err)
	require.Equal(t, echo.NewHTTPError(http.StatusNotFound), err)
}

func TestUpdatePrice(t *testing.T) {
	// Arrange
	s := new(service.MockCats)
	id := bella.ID
	req := UpdatePriceRequest{Price: 5.99}
	s.On("UpdatePrice", mockContext, id, req.Price).Return(nil)
	ctx, rec := setup(http.MethodPut, req)
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

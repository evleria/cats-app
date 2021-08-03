// Package handler encapsulates work with HTTP
package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/evleria/mongo-crud/internal/repository/entities"

	"github.com/labstack/echo/v4"

	"github.com/evleria/mongo-crud/internal/repository"
)

// GetAllCats fetches all entities from cats collection
func GetAllCats(catsRepository repository.Cats) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cats, err := catsRepository.GetAll(ctx.Request().Context())
		if err != nil {
			return err
		}

		response := GetAllCatsResponse{}
		for _, cat := range cats {
			response = append(response, mapCat(cat))
		}
		return ctx.JSON(http.StatusOK, response)
	}
}

// GetCat fetches a single cat from cats collection by ID
func GetCat(catsRepository repository.Cats) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		idParam := ctx.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		cat, err := catsRepository.GetOne(ctx.Request().Context(), id)
		if errors.Is(err, repository.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound)
		} else if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		response := GetCatResponse(mapCat(cat))
		return ctx.JSON(http.StatusOK, response)
	}
}

// AddNewCat creates a new entity in cats collection
func AddNewCat(catsRepository repository.Cats) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		request := new(AddNewCatRequest)
		err := ctx.Bind(&request)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		id, err := catsRepository.Insert(ctx.Request().Context(), request.Name, request.Color, request.Age)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		response := AddNewCatResponse{
			ID: id.String(),
		}
		return ctx.JSON(http.StatusCreated, response)
	}
}

// DeleteCat deletes a single cat from cats collection by ID
func DeleteCat(catRepository repository.Cats) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		idParam := ctx.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		err = catRepository.Delete(ctx.Request().Context(), id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return nil
	}
}

func mapCat(cat entities.Cat) Cat {
	return Cat{
		ID:    cat.ID.String(),
		Name:  cat.Name,
		Color: cat.Color,
		Age:   cat.Age,
	}
}

// AddNewCatRequest represents a request to add new cat
type AddNewCatRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Age   int    `json:"age"`
}

// AddNewCatResponse represents a response to add new cat
type AddNewCatResponse struct {
	ID string `json:"id"`
}

// GetAllCatsResponse represents a response to get all cats
type GetAllCatsResponse []Cat

// GetCatResponse represents a response to get a cat
type GetCatResponse Cat

// Cat represents a cat
type Cat struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Age   int    `json:"age"`
}

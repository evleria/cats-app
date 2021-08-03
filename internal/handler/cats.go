// Package handler encapsulates work with HTTP
package handler

import (
	"net/http"

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
			response.Items = append(response.Items, Cat{
				ID:    cat.ID.String(),
				Name:  cat.Name,
				Color: cat.Color,
				Age:   cat.Age,
			})
		}
		return ctx.JSON(http.StatusOK, cats)
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
type GetAllCatsResponse struct {
	Items []Cat `json:"items"`
}

// Cat represents a cat
type Cat struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Age   int    `json:"age"`
}

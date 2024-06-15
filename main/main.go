package main

import (
	"net/http"

	"github.com/labstack/echo"
)

// StatusResponse is the struct that will be serialized and sent back
type StatusResponse struct {
	Status string `json:"status"`
	User   string `json:"user"`
}

// can be returned to handle them easily
func GetHandler(e echo.Context) error {
	// Create response object
	body := &StatusResponse{
		Status: "Hello world from echo!",
		User:   e.Param("user"),
	}

	return e.JSON(http.StatusOK, body)
}

func main() {
	// Create echo instance
	e := echo.New()

	// Add endpoint route
	e.GET("/", GetHandler)

	// Start echo and handle errors
	e.Logger.Fatal(e.Start(":8002"))
}
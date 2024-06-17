package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

// StatusResponse is the struct that will be serialized and sent back
type StatusResponse struct {
	Status string `json:"status"`
	User   string `json:"user"`
}

// Article contains the data of a single blog post
type Article struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

// GetHandler shows home page
func GetHandler(e echo.Context) error {
	// Create response object
	body := &StatusResponse{
		Status: "<h1>Hello world from Grok-the-Casbah!</h1>",
		User:   e.Param("user"),
	}

	return e.HTML(http.StatusOK, body.Status)
}

// GetArticle shows article page
func GetBlog(e echo.Context) error {
	articles := []Article{{1, "title1", "text1"}, {2, "title2", "text2"}}

	var t *template.Template
	var err error

	if t, err = template.ParseFiles("templates/articles.html"); err != nil {
		log.Println("Error parsing template", err)
		e.Error(err)
		return nil
	}
	if err := t.Execute(e.Response().Writer, articles); err != nil {
		log.Println("Error execute template", err)
		e.Error(err)
	}
	return err
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8002"
	}

	// Create echo instance
	e := echo.New()

	// Add endpoint routes
	e.GET("/", GetHandler)
	e.GET("/blog", GetBlog)

	// Start echo and handle errors
	e.Logger.Fatal(e.Start(":" + port))
}

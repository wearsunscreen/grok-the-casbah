package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"github.com/labstack/echo"
)

var (
	db *sql.DB
)

// StatusResponse is the struct that will be serialized and sent back
type StatusResponse struct {
	Status string `json:"status"`
	User   string `json:"user"`
}

// Article contains the data of a single blog post
type Article struct {
	ID        int           `json:"id"`
	Title     string        `json:"title"`
	Content   template.HTML `json:"content"`
	Timestamp int           `json:"timestamp"`
}

// createRoutes creates the echo context and routes for the application.
// This function is called from main and from test suites.
func createRoutes() *echo.Echo {
	e := echo.New()

	e.GET("/", getHomePage)
	e.GET("/blog", getBlogArticles)
	e.GET("/blog/:id", getBlogArticle)

	/*
		// Custom HTTP error handler
		e.HTTPErrorHandler = func(err error, c echo.Context) {
			if he, ok := err.(*echo.HTTPError); ok {
				if he.Code == http.StatusNotFound {
					// Render your 404 page
					c.Render(http.StatusNotFound, "404.html", nil)
					return
				}
			}
			// Handle other errors or pass them to the default handler
			e.DefaultHTTPErrorHandler(err, c)
		}

		// Catch-all route for undefined paths
		e.GET("/*", func(c echo.Context) error {
			return c.Render(http.StatusNotFound, "404.html", nil)
		})
	*/

	return e
}

// deleteArticle deletes an article from the database
func deleteArticle(e echo.Context, id string) error {
	query, err := db.Prepare("delete from articles where id=?")
	if err != nil {
		return err
	}
	defer query.Close()

	_, err = query.Exec(id)
	return err
}

// getBlogArticle shows article page
func getBlogArticle(e echo.Context) error {
	idString := e.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println("Error converting id to int", err)
		return err
	}

	query, err := db.Prepare("SELECT * FROM article where id = ?")
	if err != nil {
		log.Println("Error preparing query", err)
		return err
	}
	defer query.Close()

	result := query.QueryRow(id)
	article := new(Article)
	if err = result.Scan(&article.ID, &article.Title, &article.Content, &article.Timestamp); err != nil {
		log.Println("Error scanning row", err)
		return err
	}
	log.Println(article.ID, article.Title)

	articles := []Article{*article}
	err = renderArticles(e, articles)
	if err != nil {
		log.Println("Error rendering articles", err)
		return err
	}
	return nil
}

// getBlogArticles shows all articles in a single page
func getBlogArticles(e echo.Context) error {
	rows, err := db.Query("SELECT * FROM article")
	if err != nil {
		log.Println("Could not query database: ", err)
		os.Exit(1)
	}
	defer rows.Close()

	var articles []Article

	for rows.Next() {
		var article Article

		if err := rows.Scan(&article.ID, &article.Title, &article.Content, &article.Timestamp); err != nil {
			log.Println("Error scanning row:", err)
			return err
		}

		articles = append(articles, article)
		log.Println(article.ID, article.Title)
	}

	funcMap := template.FuncMap{
		"formatTime": func(ts int) string {
			t := time.Unix(int64(ts), 0)
			return t.Format("Jan 2, 2006 3:04pm")
		},
	}
	var t *template.Template
	if t, err = template.New("articles.html").Funcs(funcMap).ParseFiles("templates/articles.html"); err != nil {
		log.Println("Error parsing template", err)
		e.Error(err)
		return nil
	}
	if err = t.Execute(e.Response().Writer, articles); err != nil {
		log.Println("Error execute template", err)
		e.Error(err)
	}
	return err
}

// getHomePage shows home page
func getHomePage(e echo.Context) error {
	// Create response object
	body := &StatusResponse{
		Status: "<h1>Hello world from Grok-the-Casbah!</h1>",
		User:   e.Param("user"),
	}

	return e.HTML(http.StatusOK, body.Status)
}

// openDB opens a connection to the database and relies on environment variables
// TURSO_DATABASE and TURSO_AUTH_TOKEN being set. openDB is also called from tests.
func openDB() *sql.DB {
	database := os.Getenv("TURSO_DATABASE")
	token := os.Getenv("TURSO_AUTH_TOKEN")

	if database == "" || token == "" {
		panic("TURSO_DATABASE and TURSO_AUTH_TOKEN environment variables must be set")
	}
	url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", database, token)

	var err error
	db, err = sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}
	return db
}

func renderArticles(e echo.Context, article []Article) error {
	var err error
	funcMap := template.FuncMap{
		"formatTime": func(ts int) string {
			t := time.Unix(int64(ts), 0)
			return t.Format("Jan 2, 2006 3:04pm")
		},
	}
	var t *template.Template
	if t, err = template.New("articles.html").Funcs(funcMap).ParseFiles("templates/articles.html"); err != nil {
		log.Println("Error parsing template", err)
		e.Error(err)
		return err
	}
	if err = t.Execute(e.Response().Writer, article); err != nil {
		log.Println("Error execute template", err)
		e.Error(err)
	}
	return err
}

// updateArticle updates an article in the database
func updateArticle(_ echo.Context, id string, article *Article) error {
	query, err := db.Prepare("update articles set (title, content) = (?,?) where id=?")
	if err != nil {
		return err
	}
	defer query.Close()

	_, err = query.Exec(article.Title, article.Content, id)
	return err
}

func main() {
	db = openDB()
	defer db.Close()

	port := os.Getenv("GTC_PORT")
	if port == "" {
		port = "80"
	}

	// Create echo instance and routes
	echo := createRoutes()

	// Start echo and handle errors
	echo.Logger.Fatal(echo.Start(":" + port))
}

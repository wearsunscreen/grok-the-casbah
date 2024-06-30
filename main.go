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

// getHandler shows home page
func getHandler(e echo.Context) error {
	// Create response object
	body := &StatusResponse{
		Status: "<h1>Hello world from Grok-the-Casbah!</h1>",
		User:   e.Param("user"),
	}

	return e.HTML(http.StatusOK, body.Status)
}

// GetArticle shows article page
func getBlogArticles_static(e echo.Context) error {
	articles := []Article{
		{1, "title1", "<h2>h2</h2> <h3>h3</h3> <p>p</p>", int(time.Now().Unix())},
		{2, "title2", "text2", int(time.Now().Unix())},
	}

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

// getBlogArticle shows article page
func getBlogArticle(e echo.Context) error {
	idString := e.Param("id")
	log.Println("idString", idString)
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
		log.Println("failed to execute query: ", err)
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
func updateArticle(e echo.Context, id string, article *Article) error {
	query, err := db.Prepare("update articles set (title, content) = (?,?) where id=?")
	if err != nil {
		return err
	}
	defer query.Close()

	_, err = query.Exec(article.Title, article.Content, id)
	return err
}

func main() {
	// open database
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
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8002"
	}

	// Create echo instance
	e := echo.New()

	// Add endpoint routes
	e.GET("/", getHandler)
	e.GET("/blog", getBlogArticles)
	e.GET("/blog/:id", getBlogArticle)

	// Start echo and handle errors
	e.Logger.Fatal(e.Start(":" + port))
}

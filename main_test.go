package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

// Setup by opening the database connection
func setupSuite(tb testing.TB) func(tb testing.TB) {
	log.Println("setup up test suite")
	database := openDB()

	// Return a function to teardown the test
	return func(tb testing.TB) {
		log.Println("teardown test suite")
		database.Close()
	}
}

func TestGetBlogArticles(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	table := []struct {
		name       string
		path       string
		returnCode int
		expected   string
	}{
		{"homepage", "/", 200, "Hello world from Grok-the-Casbah!"},
		{"/blog", "/blog", 200, "first"},
		{"/blog/1", "/blog/1", 200, "first"},
		{"/blog/2", "/blog/2", 200, "second"},
		{"/junk", "/junk", 404, ""},
	}

	// Initialize Echo and routes
	e := echo.New()
	e.GET("/blog", getBlogArticles)
	e.GET("/blog/:id", getBlogArticle)
	e.GET("/", getHandler)

	for _, tc := range table {
		// Create a new HTTP request
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)

		// Record the response
		rec := httptest.NewRecorder()

		// Serve the request to the recorder
		e.ServeHTTP(rec, req)

		// Assert the status code is 200
		assert.Equal(t, tc.returnCode, rec.Code)

		// Assert on the body content
		if len(tc.expected) > 0 {
			assert.Contains(t, rec.Body.String(), tc.expected)
		}
	}
}

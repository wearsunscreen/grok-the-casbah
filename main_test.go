package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

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
		expected   []string
	}{
		{"homepage", "/", 200, []string{"Hello world from Grok-the-Casbah!"}},
		{"/blog", "/blog", 200, []string{"first", "second"}},
		{"/blog/1", "/blog/1", 200, []string{"first"}},
		{"/blog/2", "/blog/2", 200, []string{"second"}},
		{"/junk", "/junk", 404, []string{"404", "Page not found"}},
	}

	// Initialize Echo and routes
	e := createRoutes()

	for _, tc := range table {
		log.Printf("Testing %s", tc.name)

		// Create a new HTTP request
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)

		// Record the response
		rec := httptest.NewRecorder()

		// Serve the request to the recorder
		e.ServeHTTP(rec, req)

		// Assert the status code is 200
		assert.Equal(t, tc.returnCode, rec.Code)

		// Assert on the body content
		if tc.expected != nil {
			for _, exp := range tc.expected {
				assert.Contains(t, rec.Body.String(), exp)
			}
		}
	}
}

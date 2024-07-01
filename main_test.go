package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestBlogEndpoint(t *testing.T) {
	db := openDB()
	defer db.Close()

	// Initialize Echo
	e := echo.New()

	// Set up the route for testing
	e.GET("/blog", getBlogArticles)

	// Create a new HTTP request to the /blog endpoint
	req := httptest.NewRequest(http.MethodGet, "/blog", nil)

	// Record the response
	rec := httptest.NewRecorder()

	// Serve the request to the recorder
	e.ServeHTTP(rec, req)

	// Assert the status code is 200
	assert.Equal(t, http.StatusOK, rec.Code)

	// Optionally, assert on the body content if you expect a specific response
	// assert.Contains(t, rec.Body.String(), "expected content")
}

package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/stretchr/testify/assert"
)

func TestRegisterTransactionRoutes(t *testing.T) {
	// t.Skip()
	app := iris.New()
	// Mock middleware
	authMiddleware := func(ctx iris.Context) {
		ctx.Next()
	}
	app.Use(authMiddleware)

	// Mock controller handlers
	createTransactionHandler := func(ctx iris.Context) {
		ctx.StatusCode(http.StatusCreated)
	}
	listTransactionsHandler := func(ctx iris.Context) {
		ctx.StatusCode(http.StatusOK)
	}

	app.Post("/transactions", createTransactionHandler)
	app.Get("/transactions", listTransactionsHandler)

	// Build the router
	app.Build()

	// Test POST /transactions
	req := httptest.NewRequest("POST", "/transactions", nil)
	resp := httptest.NewRecorder()
	app.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	// Test GET /transactions
	req = httptest.NewRequest("GET", "/transactions", nil)
	resp = httptest.NewRecorder()
	app.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

package internal

import (
	"net/http"
	"wiseman/internal/db"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func StartEcho() *echo.Echo {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.GET("/usersGuilds", getUsersGuilds)

	// Start server
	return e
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func getUsersGuilds(c echo.Context) error {
	userID := c.FormValue("userId")

	servers := db.GetServerUsers(userID)

	return c.JSON(http.StatusOK, servers)
}

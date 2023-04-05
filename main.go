package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "hello world!")
	})
	h := &hander{}
	e.POST("/login", h.login)
	e.GET("/private", h.private, IsLoggedIn, isAdmin)
	e.POST("/refresh", h.refresh)
	e.Logger.Fatal(e.Start(":1323"))
}

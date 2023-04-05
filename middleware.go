package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var IsLoggedIn = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: []byte("thisismykey"),
})

func isAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := ctx.Get("user").(*jwt.Token)
		claim := user.Claims.(jwt.MapClaims)
		isAdmin := claim["admin"].(bool)
		if !isAdmin {
			return echo.ErrUnauthorized
		}
		return next(ctx)
	}
}

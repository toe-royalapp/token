package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type hander struct{}

func (h *hander) login(ctx echo.Context) error {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")
	if username == "john" && password == "123456" {
		// create token
		token := jwt.New(jwt.SigningMethodHS256)
		fmt.Println(token)
		// set claims - to use from fontend
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = "John Doe"
		claims["pass"] = "password"
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
		fmt.Println(claims)
		// Generate encoded token and send to client
		t, err := token.SignedString([]byte("thisismykey"))
		if err != nil {
			return err
		}
		// Generate refresh token
		refreshToken := jwt.New(jwt.SigningMethodHS256)
		rtClaims := refreshToken.Claims.(jwt.MapClaims)
		rtClaims["sub"] = 1
		// refresh token time > access token time
		rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
		rt, err := refreshToken.SignedString([]byte("my_ref_key"))
		if err != nil {
			return err
		}
		return ctx.JSON(http.StatusOK, map[string]string{
			"access_token":  t,
			"refresh_token": rt,
		})
	}
	return echo.ErrUnauthorized
}

func (h *hander) private(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claim := user.Claims.(jwt.MapClaims)
	fmt.Println("private claim : ", claim)
	name := claim["name"].(string)
	return ctx.String(http.StatusOK, "Welcome "+name)
}

// This is the api to refresh tokens
// Most of the code is taken from the jwt-go package's sample codes
// https://godoc.org/github.com/dgrijalva/jwt-go#example-Parse--Hmac
func (h *hander) refresh(ctx echo.Context) error {
	type tokenReqBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	tokenReq := tokenReqBody{}
	ctx.Bind(&tokenReq)
	fmt.Println("reftoken: ", tokenReq)

	// Parse takes the token string and a function for looking up the key.
	// The latter is especially useful if you use multiple keys for your application.
	// The standard is to use 'kid' in the head of the token to identify
	// which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenReq.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("my_ref_key"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		if int(claims["sub"].(float64)) == 1 {
			newTokenPair, err := generateTokenPair()
			if err != nil {
				return err
			}
			return ctx.JSON(http.StatusOK, newTokenPair)
		}
		return echo.ErrUnauthorized
	}
	return err
}

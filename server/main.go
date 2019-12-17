package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/skanehira/mysqlstore"
)

func chcekSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err)
		}

		if sess.ID == "" {
			return c.JSON(http.StatusUnauthorized, errors.New("no session"))
		}

		return next(c)
	}
}

func printError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func main() {
	e := echo.New()

	// session middleware
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", "gorilla", "gorilla", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), "session")
	store, err := mysqlstore.NewMySQLStore(dsn, "session", "/", 60*60, []byte("secret"))
	if err != nil {
		printError(err)
	}

	e.Use(session.Middleware(store))
	e.Use(middleware.Logger())

	e.POST("/login", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err)
		}

		sess.Values["name"] = "gorilla"
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, map[string]string{"name": "gorilla"})
	})

	e.GET("/info", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		data := map[string]string{"name": sess.Values["name"].(string)}
		return c.JSON(http.StatusOK, data)
	})

	e.Logger.Fatal(e.Start(":80"))
}

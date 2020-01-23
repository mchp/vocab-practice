package main

import (
	"fmt"
	"os"

	"net/http"
	"vocabpractice/data"
	"vocabpractice/translate"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var version string

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	args := os.Args[1:]
	var db data.Database
	var err error
	if len(args) > 0 && args[0] == "sql" {
		username := os.Getenv("DB_USERNAME")
		password := os.Getenv("DB_PASSWORD")
		host := os.Getenv("DB_HOST")
		db, err = data.InitStructured(host, username, password)
		if err != nil {
			e.Logger.Fatalf("Unable to connect to database: %v", err)
			return
		}
	} else {
		local := len(args) > 0 && args[0] == "local"
		db, err = data.InitDynamoDB(local)
	}

	e.File("/", "public/quiz/index.html")
	e.Static("/static", "public/quiz/static")

	e.File("/add", "public/add/index.html")
	e.Static("/add/static", "public/add/static")

	e.GET("/version", func(c echo.Context) error {
		return c.String(http.StatusOK, version)
	})
	e.GET("/next", func(c echo.Context) error {
		nextWord, err := db.FetchNext()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, nextWord)
	})

	e.GET("/list", func(c echo.Context) error {
		words, err := db.List()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, words)
	})

	e.GET("/lookup", func(c echo.Context) error {
		vocab := c.QueryParam("vocab")
		translations, err := translate.Lookup(vocab)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		alreadyExists, err := db.QueryWord(vocab)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		for _, t := range translations {
			for _, ex := range alreadyExists.Translations {
				if t.Translation == ex.Translation {
					t.Exists = true
					break
				}
			}
		}
		return c.JSON(http.StatusOK, translations)
	})

	e.POST("/input", func(c echo.Context) error {
		v := c.FormValue("vocab")
		t := c.FormValue("translation")
		if v == "" || t == "" {
			return c.String(http.StatusBadRequest, fmt.Sprintf("no necessary parameters in request vocab=%s, translation=%s", v, t))
		}
		err := db.Input(v, t)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	})

	e.POST("/pass", func(c echo.Context) error {
		v := c.FormValue("vocab")
		t := c.FormValue("translation")
		if v == "" || t == "" {
			return c.String(http.StatusBadRequest, fmt.Sprintf("no necessary parameters in request vocab=%s, translation=%s", v, t))
		}
		err := db.Pass(v, t)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	})
	e.Logger.Fatal(e.Start(":80"))
}

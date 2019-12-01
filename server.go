package main

import (
	"fmt"

	"net/http"
	"vocabpractice/data"
	"vocabpractice/translate"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	db, err := data.Init()
	if err != nil {
		e.Logger.Fatalf("Unable to connect to database: %v", err)
		return
	}

	e.File("/", "public/quiz/index.html")
	e.Static("/static", "public/quiz/static")

	e.File("/add", "public/input/index.html")
	e.Static("/add/static", "pubic/add/static")

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

	e.Logger.Fatal(e.Start(":1234"))
}

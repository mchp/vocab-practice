package main

import (
  "fmt"

  "vocabpractice/data"
  "net/http"

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

  e.GET("/", func(c echo.Context) error {
    nextWord, err := db.FetchNext()
    if err != nil {
      return c.String(http.StatusInternalServerError, err.Error())
    }
    return c.String(http.StatusOK, nextWord.String())
  })

  e.GET("/list", func(c echo.Context) error {
    words, err := db.List()
    if err != nil {
      return c.String(http.StatusInternalServerError, err.Error())
    }
    return c.String(http.StatusOK, fmt.Sprintf("%v", words))
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
    return c.String(http.StatusOK, fmt.Sprintf("vocab=%s, translation=%s", v, t))
  })
  e.Logger.Fatal(e.Start(":1234"))
}

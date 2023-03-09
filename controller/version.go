package controller

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

var version string 

type Version struct {
	Version  string `json:"version"`
}
func GetVersion(c echo.Context) error {
	return c.JSON(http.StatusAccepted, Version{version})
}

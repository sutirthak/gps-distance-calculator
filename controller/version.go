package controller

import (
	"os"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)
type Version struct {
	Version  string `json:"version"`
}
func readFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
func GetVersion(c echo.Context) error {
	data, err := readFile("./VERSION")
	if err != nil {
		log.WithFields(log.Fields{
			"message": "failed to read VERSION file",
		}).Error(err)
		return c.String(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusAccepted, Version{Version: data})
}

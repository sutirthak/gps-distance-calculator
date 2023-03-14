package models

import (
	"github.com/labstack/echo/v4"
	cache "github.com/sutirthak/gps-distance-calculator/cache/redis"
)

type Context struct {
	echo.Context
	RedisInstance cache.RedisInstance
}

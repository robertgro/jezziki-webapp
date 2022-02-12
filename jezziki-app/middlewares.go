package main

import (
	"database/sql"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
)

type (
	// Stats defines default server status
	Stats struct {
		//Uptime       time.Time      `json:"uptime"`
		// RequestCount uint64 `json:"requestCount"` // SELECT COUNT(DISTINCT real_ip) FROM visitor;
		//Statuses     map[string]int `json:"statuses"`
		RequestID    string `json:"requestID"`
		NewRequest   bool   `json:"newRequest"`
		mutex        sync.RWMutex
		CheckCounter int       `json:"checkcounter"`
		StartTime    time.Time `json:"starttime"`
		StatErr      error
		UniqueIP     string
	}

	// CustomRateLimiter is a default rate limiter
	CustomRateLimiter struct {
		limiter *rate.Limiter
	}
)

// PostCustomContextMW extends the default context, registered after logging module
func (app *Controller) PostCustomContextMW() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}

// PreCustomContextMW extends the default context, registered before logging module, mutex.Lock required
func (app *Controller) PreCustomContextMW() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {

			if app.stats.UniqueIP == "" || app.stats.UniqueIP != c.RealIP() {
				app.stats.UniqueIP = c.RealIP()
				app.stats.NewRequest = true

				PreLogInfo(c, app)
			}

			app.dbWorker.errorVar = app.dbWorker.PrepareDB(c)

			if app.dbWorker.errorVar != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, app.dbWorker.errorVar.Error())
			}

			app.dbWorker.errorVar = app.dbWorker.PingDB()

			if app.dbWorker.errorVar != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, app.dbWorker.errorVar.Error())
			}

			defer app.dbWorker.pool.Close()

			if app.enforcer.accessGranted, app.dbWorker.errorVar = app.dbWorker.getEditSessionStatus(c, app); app.enforcer.accessGranted && app.dbWorker.errorVar == nil {
				return next(c)
			} else if app.dbWorker.errorVar != nil && app.dbWorker.errorVar != sql.ErrNoRows {
				return echo.NewHTTPError(http.StatusInternalServerError, app.dbWorker.errorVar.Error())
			}
			return next(c)
		}
	}
}

// RateLimitMW is a dos protection
// Credits to
// https://www.alexedwards.net/blog/how-to-rate-limit-http-requests
// https://github.com/hrodic/golang-echo-simple-rate-limit-middleware/blob/master/middleware.go
func (app *Controller) RateLimitMW() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			ts, err := getTimestamp()

			Check(err, false)

			if !app.rateLimiter.limiter.Allow() {
				println(ts.Format(timeFormat) + " - INFO  - SECURITY - RATE LIMITER ACTIVE - IP " + c.RealIP())
				return echo.NewHTTPError(http.StatusTooManyRequests, "DDOS-PROTECTION")
			}
			return next(c)
		}
	}
}

// RegisterMiddlewares registers middlewares
func RegisterMiddlewares(e *echo.Echo, c *Controller) {

	e.Pre(
		mw.Rewrite(c.defaultServerConf.Rewrites),
		mw.NonWWWRedirect(),
	)

	if c.defaultServerConf.SSL {
		e.Pre(mw.HTTPSRedirect())
	}

	e.Use(
		mw.SecureWithConfig(getSecureConfig()),
		mw.CORSWithConfig(getCorsConfig()),
		mw.BodyLimit("2M"),
		c.RateLimitMW(),
		c.PreCustomContextMW(),
		mw.RequestIDWithConfig(getRequestIDConfig(c.customGenerator)),
		mw.Recover(),
		mw.Gzip(),
		mw.LoggerWithConfig(getLoggerConfig(c)),
		c.PostCustomContextMW(),
	)

}

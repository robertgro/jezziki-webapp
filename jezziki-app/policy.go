package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/casbin/casbin/v2"

	"github.com/labstack/echo/v4"
)

// CustomEnforcer enforces custom casbin auth
type CustomEnforcer struct {
	enforcer      *casbin.Enforcer
	dashU         string
	dashP         string
	tokenID       string
	accessTokenA  string
	accessTokenB  string
	accessGranted bool
	tNow          time.Time
	tEnd          time.Time
	tDiff         time.Duration
	realm         string
}

// NewEnforcer constructs a news casbin custom enforcer
func NewEnforcer() (cx *CustomEnforcer) {
	ce, err := casbin.NewEnforcer("auth_model.conf", "auth_policy.csv")
	Check(err, true)
	return &CustomEnforcer{
		enforcer:      ce,
		dashU:         "hFM2kdYAHg9j",
		dashP:         "9VTChm9f8SeT",
		tokenID:       "",
		accessTokenA:  "4SEZudcGeuKX",
		accessTokenB:  "R2HLzgResYYE",
		accessGranted: false,
		tNow:          time.Time{},
		tEnd:          time.Time{},
		tDiff:         0,
		realm:         `Basic realm="localhost"`,
	}
}

// PolicyEnforcerMW middleware enforces the casbin policy rules
func (ce *CustomEnforcer) PolicyEnforcerMW(co *Controller) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			ts, err := getTimestamp()
			if err != nil {
				Check(err, false)
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			println(ts.Format(timeFormat) + " - AUTH  - AUDIT - REALM " + ce.realm + " IP " + c.RealIP() + " URL.PATH " + c.Request().URL.Path)

			if ok, err := ce.EvaluateRequest(c); err == nil && ok {
				return next(c)
			}
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
	}
}

// EvaluateRequest evaluates the csp request
func (ce *CustomEnforcer) EvaluateRequest(c echo.Context) (bool, error) {
	u, pw, _ := c.Request().BasicAuth()
	p := c.Request().URL.Path
	m := c.Request().Method

	ts, err := getTimestamp()
	Check(err, false)

	println(ts.Format(timeFormat) + " - AUTH  - ENFORCE - REALM " + ce.realm + " IP " + c.RealIP() + " USER " + u + " PASS " + pw + " PATH " + p + " METHOD " + m)
	return ce.enforcer.Enforce(u, p, m)
}

// RevokeAccessAfterT revokes the dash operator access after time passed
func (app *Controller) RevokeAccessAfterT(c echo.Context) (err error) {
	time.AfterFunc(10*time.Minute, func() {
		app.RevokeAccess(c, false, true)
	})
	return nil
}

// RevokeAccess revokes the dash operator access
func (app *Controller) RevokeAccess(c echo.Context, active bool, timeBased bool) (err error) {

	if timeBased {
		if err = app.dbWorker.PrepareDB(nil); err != nil {
			Check(err, false)
			return
		}
	} else {
		if err = app.dbWorker.PrepareDB(c); err != nil {
			Check(err, false)
			return
		}
	}

	if err = app.dbWorker.PingDB(); err != nil {
		Check(err, false)
		return
	}

	defer app.dbWorker.pool.Close()

	if timeBased {
		if err = app.dbWorker.disableEditSession(c, app.enforcer.tokenID, active, timeBased); err != nil {
			Check(err, false)
			return
		}
	} else {
		if err = app.dbWorker.disableEditSession(c, app.enforcer.tokenID, active, false); err != nil {
			Check(err, false)
			return
		}
	}

	ts, err := getTimestamp()

	if err != nil {
		Check(err, false)
		return
	}
	app.enforcer.tokenID = ""
	app.enforcer.accessGranted = false

	println(ts.Format(timeFormat) + " - AUTH  - EVENT - REALM " + app.enforcer.realm + " MESSAGE Logout IP " + c.RealIP() + " PATH " + c.Request().URL.Path + " ACTIVE " + strconv.FormatBool(active) + " TIMEBASED " + strconv.FormatBool(timeBased))
	return nil
}

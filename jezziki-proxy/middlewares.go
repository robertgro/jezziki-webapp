package main

import (
	"database/sql"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
)

// PreCustomContextMW extends the default context, registered before logging module, mutex.Lock required
func PreCustomContextMW(pc *proxyController) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			pc.reqStats.mutex.Lock()
			defer pc.reqStats.mutex.Unlock()

			if pc.db.errorVar = pc.PrepareDB(c); pc.db.errorVar != nil {
				return echo.NewHTTPError(http.StatusBadRequest, pc.db.errorVar.Error())
			}

			if pc.db.errorVar = pc.PingDB(); pc.db.errorVar != nil {
				return echo.NewHTTPError(http.StatusBadRequest, pc.db.errorVar.Error())
			}

			defer pc.db.pool.Close()

			if pc.reqStats.Blacklisted {
				c.Redirect(http.StatusTemporaryRedirect, "https://www.google.com")
			}

			if pc.enforcer.access, pc.db.errorVar = pc.getEditSessionStatus(c); pc.enforcer.access && pc.db.errorVar == nil {
				c.Response().Header().Del(echo.HeaderContentSecurityPolicy)
			} else if pc.db.errorVar != nil && pc.db.errorVar != sql.ErrNoRows {
				return echo.NewHTTPError(http.StatusBadRequest, pc.db.errorVar.Error())
			}

			if pc.reqStats.UniqueIP == "" || pc.reqStats.UniqueIP != c.RealIP() {
				pc.PreLogInfo(c)
				pc.reqStats.UniqueIP = c.RealIP()
				pc.reqStats.NewRequest = true
				pc.reqStats.Logged = false
				pc.reqStats.WDone = false
				pc.reqStats.CheckCounter = 7 // force BL query
				if err := pc.getWhoisInfo(c); err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, err.Error())
				}
			}

			// check BL every 7 req (= 1 round trip)
			if pc.reqStats.CheckCounter == 7 {

				pc.reqStats.CheckCounter = 0

				if pc.reqStats.Logged {
					if b, err := pc.queryDBNewUserAgent(c); b && err == nil {
						pc.reqStats.Logged = false
					} else if err != nil {
						return echo.NewHTTPError(http.StatusBadRequest, err.Error())
					}
				}

				if pc.db.errorVar = pc.queryDBVisitorBlacklist(c); pc.db.errorVar != nil {
					return echo.NewHTTPError(http.StatusBadRequest, pc.db.errorVar.Error())
				}

				if !pc.reqStats.WDone {

					if pc.db.errorVar = pc.checkWhoisDone(); pc.db.errorVar != nil {
						return echo.NewHTTPError(http.StatusBadRequest, pc.db.errorVar.Error())
					}
				}

				if pc.reqStats.WDone && !pc.reqStats.ScanDone {

					if pc.db.errorVar = pc.getWhoisResult(); pc.db.errorVar != nil {
						return echo.NewHTTPError(http.StatusBadRequest, pc.db.errorVar.Error())
					}
				}
			}

			pc.reqStats.CheckCounter += 1

			return next(c)
		}
	}
}

// PostCustomContextMW extends the default context, registered after logging module
func (pc *proxyController) PostCustomContextMW() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if !pc.reqStats.Logged {

				if pc.db.errorVar = pc.PrepareDB(c); pc.db.errorVar != nil {
					return echo.NewHTTPError(http.StatusBadRequest, pc.db.errorVar.Error())
				}

				if pc.db.errorVar = pc.PingDB(); pc.db.errorVar != nil {
					return echo.NewHTTPError(http.StatusBadRequest, pc.db.errorVar.Error())
				}

				defer pc.db.pool.Close()

				if pc.db.errorVar = pc.queryDBVisitorLog(pc.reqStats.RequestID, c.RealIP(), c.Request().RemoteAddr, c.Request().Host+c.Request().URL.Port(), c.Request().UserAgent(), c.Request().Referer(), c.Request().Method+" "+c.Request().RequestURI, c.Request().URL.Path, pc.reqStats.UProvider, pc.reqStats.UCountry, pc.reqStats.UCity, pc.reqStats.Blacklisted, c.Request().Proto); pc.db.errorVar != nil {
					return echo.NewHTTPError(http.StatusBadRequest, pc.db.errorVar.Error())
				}
				pc.reqStats.Logged = true

			}

			return next(c)
		}
	}
}

func (pc *proxyController) parseProxyTargets(e *echo.Echo) {
	for _, node := range pc.URLTargets.Nodes {
		prot := ""
		if pc.DefaultProxyConf.SSL {
			prot = "https://"
		} else {
			prot = "http://"
		}
		url, err := url.Parse(prot + node.Host + ":" + node.Port)
		if err != nil {
			e.Logger.Fatal(err)
			url = nil
		}
		pc.pURLs = append(pc.pURLs, url)
	}

	if len(pc.pURLs) > 0 {
		for _, element := range pc.pURLs {
			pc.proxyTargets = append(pc.proxyTargets, &mw.ProxyTarget{URL: element})
		}
	} else {
		e.Logger.Fatal("pURLs is 0")
	}
}

func (pc *proxyController) RegisterMiddlewares(e *echo.Echo) {
	e.Pre(
		mw.NonWWWRedirect(),
	)

	e.Use(
		pc.getCorsConfigMW(),
		pc.getRequestIDConfigMW(),
		mw.BodyLimit("2M"),
		PreCustomContextMW(pc),
		mw.Recover(),
		mw.Gzip(),
		pc.getLoggerConfigMW(),
		pc.PostCustomContextMW(),
		mw.Proxy(mw.NewRoundRobinBalancer(pc.proxyTargets)),
	)
}

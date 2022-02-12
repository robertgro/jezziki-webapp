package main

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tylerb/graceful"
)

func main() {

	/* Load jezziki app */

	app := NewController()

	app.initAppConfig()

	/* Load echo server */

	e := echo.New()

	app.initEchoServer(e)

	RegisterMiddlewares(e, app)

	RegisterRoutes(e, app)

	if app.defaultServerConf.SSL {
		e.Logger.Fatal(graceful.ListenAndServeTLS(e.Server, app.defaultServerConf.Paths.Cert, app.defaultServerConf.Paths.Key, 5*time.Second))
	} else {
		e.Logger.Fatal(graceful.ListenAndServe(e.Server, 5*time.Second))
	}

}

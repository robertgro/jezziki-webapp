package main

import (
	"time"

	"github.com/tylerb/graceful"

	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	pc := NewProxyController()

	pc.initApp(e)

	pc.RegisterMiddlewares(e)

	if pc.DefaultProxyConf.SSL {
		e.Logger.Fatal(graceful.ListenAndServeTLS(e.Server, pc.DefaultProxyConf.Paths.Cert, pc.DefaultProxyConf.Paths.Key, 5*time.Second))
	} else {
		e.Logger.Fatal(graceful.ListenAndServe(e.Server, 5*time.Second))
	}

}

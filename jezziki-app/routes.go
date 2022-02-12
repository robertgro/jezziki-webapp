package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RegisterRoutes registers all routes
func RegisterRoutes(e *echo.Echo, c *Controller) {

	// Main react components and utilities, e.g. fetch
	e.Static("/src", c.defaultServerConf.Paths.Static)

	// Main static assets, images, favicon, robots.txt, err.css/jpg
	e.Static("/", c.defaultServerConf.Paths.Public)

	e.GET("", c.getIndexHandler)

	// Internal dash access group
	a := e.Group(c.defaultServerConf.Paths.Dash)

	a.Use(middleware.BasicAuthWithConfig(getBasicAuthConfig(c)))

	a.Use(c.enforcer.PolicyEnforcerMW(c))

	// Internal dash/dist
	a.Static("/dist", c.defaultServerConf.Paths.Dist)

	// Internal dash/index
	a.GET("", c.getDashIndexHandler)

	// Internal post token check
	a.POST("/token", c.postDashTokenHandler)

	// Internal dash/post updates post item
	a.POST("/post", c.postUpdatePost)

	// Internal dash/nav updates nav items
	a.POST("/nav", c.postUpdateNav)

	// Internal dash/aside updates aside items
	a.POST("/aside", c.postUpdateAside)

	// Internal dash/footer updates footer items
	a.POST("/footer", c.postUpdateFooter)

	// Internal dash/logout deletes access
	a.POST("/logout", c.LogOutHandler)

	// Internal dash/test for item updates
	a.POST("/test", c.postUpdateComponentItem)

	// External rest api/v1
	g := e.Group("api/v1/")

	// External get posts
	g.GET("posts/:id", c.getPostByIDHandler)

	// External get page index components
	g.GET("index", c.getComponentsHandler)

	// External wildcard post not found
	e.POST("/*", c.NotFoundHandler)
	e.PUT("/*", c.NotAllowedHandler)
	e.PATCH("/*", c.NotAllowedHandler)
	e.OPTIONS("/*", c.NotAllowedHandler)
	e.HEAD("/*", c.NotAllowedHandler)
	e.DELETE("/*", c.NotAllowedHandler)
	e.CONNECT("/*", c.NotAllowedHandler)
	e.TRACE("/*", c.NotAllowedHandler)
}

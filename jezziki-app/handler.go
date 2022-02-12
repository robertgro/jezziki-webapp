package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

type (
	// CustomValidator creates a custom validator
	CustomValidator struct {
		validator *validator.Validate
	}

	// CustomPageItems contains page component slices
	CustomPageItems struct {
		NavItems    *[]NavbarItem `json:"nav"`
		AsideItems  *[]AsideItem  `json:"aside"`
		FooterItems *[]FooterItem `json:"footer"`
		PostIndex   *PostItem     `json:"indexpost"`
	}

	// CustomPageItem contains single component items
	CustomPageItem struct {
		NavItem *NavbarItem
		AsItem  *AsideItem
		FoItem  *FooterItem
		PoItem  *PostItem
	}

	// NavbarItem defines a navbar item
	NavbarItem struct {
		ID       int    `json:"id"`
		PID      int    `json:"parentid"`
		Title    string `json:"title"`
		FID      int    `json:"foreignid"`
		External string `json:"ext"`
	}

	// AsideItem defines an aside item
	AsideItem struct {
		ID       int    `json:"id"`
		Title    string `json:"title"`
		FID      int    `json:"foreignid"`
		External string `json:"ext"`
	}

	// FooterItem defines a footer item
	FooterItem struct {
		ID       int    `json:"id"`
		Title    string `json:"title"`
		FID      int    `json:"foreignid"`
		External string `json:"ext"`
	}

	// PostItem defines a post item
	PostItem struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	}

	// NavbarItemsNew from POST request, validation required
	NavbarItemsNew struct {
		NavItems []NavItemNew `json:"navItems" validate:"required"`
	}

	// NavItemNew from POST request, validation required
	NavItemNew struct {
		ID       string `json:"id" validate:"required"`
		Title    string `json:"title" validate:"required"`
		ParentID string `json:"parentid" validate:"required"`
	}

	// AsideItemsNew from POST request, validation required
	AsideItemsNew struct {
		AsideItems []AsideItemNew `json:"asideItems" validate:"required"`
	}

	// AsideItemNew from POST request, validation required
	AsideItemNew struct {
		ID       string `json:"id" validate:"required"`
		Title    string `json:"title" validate:"required"`
		ParentID string `json:"parentid" validate:"required"`
	}

	// FooterItemsNew from POST request, validation required
	FooterItemsNew struct {
		FooterItems []AsideItemNew `json:"footerItems" validate:"required"`
	}

	// FooterItemNew from POST request, validation required
	FooterItemNew struct {
		ID       string `json:"id" validate:"required"`
		Title    string `json:"title" validate:"required"`
		ParentID string `json:"parentid" validate:"required"`
	}

	// PostItemNew from POST request, validation required
	PostItemNew struct {
		ID   string `json:"id" validate:"required"`
		Text string `json:"text" validate:"required"`
		Ext  string `json:"ext" validate:"required"`
	}
)

// Validate needs to be implemented to use the validator
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

// UpdateNavItems sets nav items
func (app *Controller) UpdateNavItems(navItems *NavbarItemsNew) (err error) {
	newNavbarItems := &CustomPageItems{NavItems: &[]NavbarItem{}}
	app.dbWorker.postID = 0

	for _, element := range navItems.NavItems {
		for _, se := range *app.customPageItems.NavItems {
			if element.Title == se.Title {
				app.dbWorker.postID = se.FID
				break
			}
		}

		if postid := app.dbWorker.getExistingPostID(element.Title); postid != 0 {
			app.dbWorker.postID = postid
		}

		*newNavbarItems.NavItems = append(*newNavbarItems.NavItems, NavbarItem{ID: getInt(element.ID), Title: element.Title, PID: getInt(element.ParentID), FID: app.dbWorker.postID, External: "false"})
		app.dbWorker.postID = 0
	}

	u, err := json.Marshal(*app.customPageItems.NavItems)

	if err != nil {
		return err
	}

	app.dbWorker.actionLog = "UPDATE NAV FROM " + string(u)

	if u, err = json.Marshal(newNavbarItems); err != nil {
		return err
	}

	app.dbWorker.actionLog += " TO " + string(u)

	if err = app.dbWorker.queryDBUpdateNav(newNavbarItems); err != nil {
		return err
	}

	return nil
}

// UpdateAsideItems sets nav items
func (app *Controller) UpdateAsideItems(asideItems *AsideItemsNew) (err error) {
	newAsideItems := &CustomPageItems{AsideItems: &[]AsideItem{}}
	app.dbWorker.postID = 0

	for _, element := range asideItems.AsideItems {
		for _, se := range *app.customPageItems.AsideItems {
			if element.Title == se.Title {
				app.dbWorker.postID = se.FID
				break
			}
		}

		if postid := app.dbWorker.getExistingPostID(element.Title); postid != 0 {
			app.dbWorker.postID = postid
		}

		*newAsideItems.AsideItems = append(*newAsideItems.AsideItems, AsideItem{ID: getInt(element.ID), Title: element.Title, FID: app.dbWorker.postID, External: "false"})
		app.dbWorker.postID = 0
	}

	u, err := json.Marshal(*app.customPageItems.AsideItems)

	if err != nil {
		return err
	}

	app.dbWorker.actionLog = "UPDATE ASIDE FROM " + string(u)

	if u, err = json.Marshal(newAsideItems); err != nil {
		return err
	}

	app.dbWorker.actionLog += " TO " + string(u)

	if err = app.dbWorker.queryDBUpdateAside(newAsideItems); err != nil {
		return err
	}

	return nil
}

// UpdateComponentItem reflect the value of the component item and queries the db properly
func (app *Controller) UpdateComponentItem(itemName string) (err error) {
	println("itemname", itemName)
	app.dbWorker.queryDBUpdateComponentItem(itemName)
	return nil
}

// UpdateFooterItems sets nav items
func (app *Controller) UpdateFooterItems(footerItems *FooterItemsNew) (err error) {
	newFooterItems := &CustomPageItems{FooterItems: &[]FooterItem{}}
	app.dbWorker.postID = 0

	for _, element := range footerItems.FooterItems {
		for _, se := range *app.customPageItems.FooterItems {
			if element.Title == se.Title {
				app.dbWorker.postID = se.FID
				break
			}
		}

		if postid := app.dbWorker.getExistingPostID(element.Title); postid != 0 {
			app.dbWorker.postID = postid
		}

		*newFooterItems.FooterItems = append(*newFooterItems.FooterItems, FooterItem{ID: getInt(element.ID), Title: element.Title, FID: app.dbWorker.postID, External: "false"})
		app.dbWorker.postID = 0
	}

	u, err := json.Marshal(*app.customPageItems.FooterItems)

	if err != nil {
		return err
	}

	app.dbWorker.actionLog = "UPDATE FOOTER FROM " + string(u)

	if u, err = json.Marshal(newFooterItems); err != nil {
		return err
	}

	app.dbWorker.actionLog += " TO " + string(u)

	if err = app.dbWorker.queryDBUpdateFooter(newFooterItems); err != nil {
		return err
	}

	return nil
}

// getIndexHandler handles GET /
func (app *Controller) getIndexHandler(c echo.Context) error {

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
		c.Response().Header().Del(echo.HeaderContentSecurityPolicy)
	} else if app.dbWorker.errorVar != nil && app.dbWorker.errorVar != sql.ErrNoRows {
		Check(app.dbWorker.errorVar, false)
	}

	return c.Render(http.StatusOK, "base.gohtml", map[string]interface{}{
		"titlepage":     app.defaultServerConf.Title,
		"favicon":       app.defaultServerConf.Links.Favicon,
		"stylesheets":   app.defaultServerConf.Links.Stylesheets,
		"utils":         app.defaultServerConf.Utils,
		"components":    app.defaultServerConf.Components,
		"accessGranted": app.enforcer.accessGranted,
	})
}

func (app *Controller) getComponents(c echo.Context) error {

	*app.customPageItems.AsideItems = *(new([]AsideItem))
	*app.customPageItems.FooterItems = *(new([]FooterItem))
	*app.customPageItems.NavItems = *(new([]NavbarItem))
	app.customPageItems.PostIndex = &PostItem{
		ID:   "",
		Text: "",
	}

	app.dbWorker.errorVar = app.dbWorker.PrepareDB(c)

	if app.dbWorker.errorVar != nil {
		return app.dbWorker.errorVar
	}

	app.dbWorker.errorVar = app.dbWorker.PingDB()

	if app.dbWorker.errorVar != nil {
		return app.dbWorker.errorVar
	}

	defer app.dbWorker.pool.Close()

	app.dbWorker.sqlQuery = "select content from posts where posts_id=$1;"
	app.dbWorker.sqlValue = "1"

	app.dbWorker.errorVar = app.dbWorker.queryDB()

	if app.dbWorker.errorVar != nil {
		return app.dbWorker.errorVar
	}

	app.customPageItems.PostIndex.ID = app.dbWorker.sqlValue
	app.customPageItems.PostIndex.Text = app.dbWorker.result

	// sqlQuery must match sqlObject interface
	app.dbWorker.sqlQuery = `
		SELECT footer.footer_id, posts.title, footer.posts_id, posts.external
		FROM footer
		INNER JOIN posts ON posts.posts_id=footer.posts_id
		ORDER BY footer.footer_id ASC;`

	app.dbWorker.sqlObject = app.customPageItem.FoItem
	app.dbWorker.sqlObjectList = app.customPageItems.FooterItems

	app.dbWorker.errorVar = app.dbWorker.queryDBResults()

	if app.dbWorker.errorVar != nil {
		return app.dbWorker.errorVar
	}

	// sqlQuery must match sqlObject interface
	app.dbWorker.sqlQuery = `
		SELECT aside.aside_id, posts.title, aside.posts_id, posts.external
		FROM aside
		INNER JOIN posts ON posts.posts_id=aside.posts_id
		ORDER BY aside.aside_id ASC;`

	app.dbWorker.sqlObject = app.customPageItem.AsItem
	app.dbWorker.sqlObjectList = app.customPageItems.AsideItems

	app.dbWorker.errorVar = app.dbWorker.queryDBResults()

	if app.dbWorker.errorVar != nil {
		return app.dbWorker.errorVar
	}

	// sqlQuery must match sqlObject interface
	app.dbWorker.sqlQuery = `
		SELECT navbar.navbar_id, navbar.parent_id, posts.title, navbar.posts_id, posts.external
		FROM navbar
		INNER JOIN posts ON posts.posts_id=navbar.posts_id
		ORDER BY navbar.navbar_id ASC;`
	app.dbWorker.sqlObject = app.customPageItem.NavItem
	app.dbWorker.sqlObjectList = app.customPageItems.NavItems

	app.dbWorker.errorVar = app.dbWorker.queryDBResults()

	if app.dbWorker.errorVar != nil {
		return app.dbWorker.errorVar
	}

	return nil
}

func (app *Controller) getComponentsHandler(c echo.Context) error {

	err := app.getComponents(c)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, app.customPageItems)
}

// getPostByIDHandler handles GET api/v1/posts/:id
func (app *Controller) getPostByIDHandler(c echo.Context) error {

	app.dbWorker.sqlQuery = "select content from posts where posts_id=$1;"
	app.dbWorker.sqlValue = c.Param("id")

	app.dbWorker.errorVar = app.dbWorker.PrepareDB(c)

	if app.dbWorker.errorVar != nil {
		return echo.NewHTTPError(http.StatusBadRequest, app.dbWorker.errorVar)
	}

	app.dbWorker.errorVar = app.dbWorker.PingDB()

	if app.dbWorker.errorVar != nil {
		return echo.NewHTTPError(http.StatusBadRequest, app.dbWorker.errorVar)
	}

	defer app.dbWorker.pool.Close()

	app.dbWorker.errorVar = app.dbWorker.queryDB()

	if app.dbWorker.errorVar != nil {
		return echo.NewHTTPError(http.StatusBadRequest, app.dbWorker.errorVar)
	}

	app.customPageItem.PoItem.ID = app.dbWorker.sqlValue
	app.customPageItem.PoItem.Text = app.dbWorker.result

	return c.JSON(http.StatusOK, app.customPageItem.PoItem)
}

// postUpdateItem handles dash/post dash/aside dash/footer dash/nav
func (app *Controller) postUpdateComponentItem(c echo.Context) (err error) {

	componentName := strings.Replace(c.Request().URL.Path[1:], "5fzt78g4A7fnb882/", "", -1)

	switch componentName {
	case "":
		app.UpdateComponentItem(componentName)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "Component name not recognized")
	}

	return c.String(http.StatusOK, "OK")
}

// postUpdatePost handles dash/post
func (app *Controller) postUpdatePost(c echo.Context) (err error) {

	app.dbWorker.errorVar = app.dbWorker.PrepareDB(c)
	if app.dbWorker.errorVar != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, app.dbWorker.errorVar.Error())
	}
	app.dbWorker.errorVar = app.dbWorker.PingDB()
	if app.dbWorker.errorVar != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, app.dbWorker.errorVar.Error())
	}

	defer app.dbWorker.pool.Close()

	if b, err := app.dbWorker.getEditSessionStatus(c, app); b && err == nil {

		app.dbWorker.errorVar = app.dbWorker.queryCtrlID(app.enforcer.tokenID, app.enforcer.tEnd)
		if app.dbWorker.errorVar != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, app.dbWorker.errorVar.Error())
		}

		postBodyNav := &PostItemNew{}

		if err = c.Bind(postBodyNav); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if err = c.Validate(postBodyNav); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if err = app.dbWorker.updatePost(postBodyNav); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, postBodyNav)
	} else {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}
}

// postUpdateAside handles dash/aside
func (app *Controller) postUpdateAside(c echo.Context) (err error) {

	app.dbWorker.errorVar = app.dbWorker.PrepareDB(c)

	if app.dbWorker.errorVar != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, app.dbWorker.errorVar.Error())
	}

	app.dbWorker.errorVar = app.dbWorker.PingDB()

	if app.dbWorker.errorVar != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, app.dbWorker.errorVar.Error())
	}

	defer app.dbWorker.pool.Close()

	if b, err := app.dbWorker.getEditSessionStatus(c, app); b && err == nil {

		app.dbWorker.errorVar = app.dbWorker.queryCtrlID(app.enforcer.tokenID, app.enforcer.tEnd)

		if app.dbWorker.errorVar != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, app.dbWorker.errorVar.Error())
		}

		postBodyNav := &AsideItemsNew{}

		if err = c.Bind(postBodyNav); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if err = c.Validate(postBodyNav); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if err = app.UpdateAsideItems(postBodyNav); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, postBodyNav)

	} else {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}
}

// postUpdateFooter handles dash/footer
func (app *Controller) postUpdateFooter(c echo.Context) (err error) {

	app.dbWorker.errorVar = app.dbWorker.PrepareDB(c)

	if app.dbWorker.errorVar != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, app.dbWorker.errorVar.Error())
	}

	app.dbWorker.errorVar = app.dbWorker.PingDB()

	if app.dbWorker.errorVar != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, app.dbWorker.errorVar.Error())
	}

	defer app.dbWorker.pool.Close()

	if b, err := app.dbWorker.getEditSessionStatus(c, app); b && err == nil {

		app.dbWorker.errorVar = app.dbWorker.queryCtrlID(app.enforcer.tokenID, app.enforcer.tEnd)

		if app.dbWorker.errorVar != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, app.dbWorker.errorVar.Error())
		}

		postBodyNav := &FooterItemsNew{}

		if err = c.Bind(postBodyNav); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if err = c.Validate(postBodyNav); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if err = app.UpdateFooterItems(postBodyNav); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, postBodyNav)

	} else {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}
}

// postUpdateNav handles dash/nav
func (app *Controller) postUpdateNav(c echo.Context) (err error) {

	app.dbWorker.errorVar = app.dbWorker.PrepareDB(c)

	if app.dbWorker.errorVar != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, app.dbWorker.errorVar.Error())
	}

	app.dbWorker.errorVar = app.dbWorker.PingDB()

	if app.dbWorker.errorVar != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, app.dbWorker.errorVar.Error())
	}

	defer app.dbWorker.pool.Close()

	if b, err := app.dbWorker.getEditSessionStatus(c, app); b && err == nil {

		app.dbWorker.errorVar = app.dbWorker.queryCtrlID(app.enforcer.tokenID, app.enforcer.tEnd)

		if app.dbWorker.errorVar != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, app.dbWorker.errorVar.Error())
		}

		postBodyNav := &NavbarItemsNew{}

		if err = c.Bind(postBodyNav); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if err = c.Validate(postBodyNav); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		app.UpdateNavItems(postBodyNav)

		return c.JSON(http.StatusOK, postBodyNav)

	} else {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}
}

// getDashIndexHandler handles admin access https://192.168.0.8/5fzt78g4A7fnb882
func (app *Controller) getDashIndexHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "dash.gohtml", map[string]interface{}{
		"title": "Welcome to Jezziki's Administrator Gateway",
	})
}

// postDashTokenHandler handles POST token access https://192.168.0.8/5fzt78g4A7fnb882/token
func (app *Controller) postDashTokenHandler(c echo.Context) error {

	err := c.Request().ParseForm()

	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	app.enforcer.tokenID = c.Request().FormValue("token")

	if app.enforcer.tokenID != app.enforcer.accessTokenA && app.enforcer.tokenID != app.enforcer.accessTokenB {
		app.logXToFile("TOKEN MISMATCH TOKENID "+app.enforcer.tokenID, http.StatusForbidden, c, app.serviceConf.policyLogFilePath)
		return echo.NewHTTPError(http.StatusForbidden, "TOKEN MISMATCH TOKENID"+app.enforcer.tokenID)
	}

	if err = app.RevokeAccessAfterT(c); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
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

	app.enforcer.tNow = time.Now()
	app.enforcer.tEnd = app.enforcer.tNow.Add(time.Minute * 10)

	app.enforcer.accessGranted = true

	app.dbWorker.errorVar = app.dbWorker.queryDBCtrlLog(app.enforcer.tNow, app.enforcer.tEnd, app.enforcer.tokenID, c.RealIP())
	if app.dbWorker.errorVar != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, app.dbWorker.errorVar.Error())
	}

	if app.defaultServerConf.SSL {
		return c.Redirect(http.StatusSeeOther, "https://localhost/")
	} else {
		return c.Redirect(http.StatusSeeOther, "http://localhost/")
	}
}

// LogOutHandler handles dash/logout
func (app *Controller) LogOutHandler(c echo.Context) error {
	if err := app.RevokeAccess(c, false, false); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.Response().Header().Del(echo.HeaderAuthorization)

	if app.defaultServerConf.SSL {
		return c.Redirect(http.StatusSeeOther, "https://localhost/")
	} else {
		return c.Redirect(http.StatusSeeOther, "http://localhost/")
	}
}

// UndefinedHandler handles /undefined
func (app *Controller) UndefinedHandler(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotFound, "Source undefined")
}

// NotFoundHandler handles /*
func (app *Controller) NotFoundHandler(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotFound, "Resource not found")
}

// NotAllowedHandler handles /*
func (app *Controller) NotAllowedHandler(c echo.Context) error {
	return echo.NewHTTPError(http.StatusForbidden, "Method is forbidden")
}

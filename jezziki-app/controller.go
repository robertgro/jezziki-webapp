package main

import (
	"database/sql"
	"html/template"
	"io"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

type (

	// Controller controlls the flow
	Controller struct {
		defaultServerConf *DefaultServerConf
		templateRenderer  *TemplateRenderer
		stats             *Stats
		dbWorker          *DBWorker
		customValidator   *CustomValidator
		errHandler        *ErrCodes
		enforcer          *CustomEnforcer
		serviceConf       *ServiceConf
		rateLimiter       *CustomRateLimiter
		customPageItems   *CustomPageItems
		customPageItem    *CustomPageItem
	}

	// TemplateRenderer defines the default template renderer for echo golang
	TemplateRenderer struct {
		templates *template.Template
	}
)

// NewController creates a new app specific controller :)
func NewController() *Controller {
	return &Controller{
		defaultServerConf: &DefaultServerConf{
			Title: "",
			Debug: false,
			Port:  "",
			Host:  "",
			SSL:   false,
			Paths: struct {
				Cert   string "yaml:\"cert\""
				Key    string "yaml:\"key\""
				Static string "yaml:\"static\""
				Public string "yaml:\"public\""
				Dash   string "yaml:\"dash\""
				Dist   string "yaml:\"dist\""
			}{
				Cert:   "",
				Key:    "",
				Static: "",
				Public: "",
				Dash:   "",
				Dist:   "",
			},
			Services: map[string]int{},
			Rewrites: map[string]string{},
			Links: struct {
				Favicon     string   "yaml:\"favicon\""
				Stylesheets []string "yaml:\"stylesheets\""
			}{},
			Utils:      []string{},
			Components: []string{},
		},
		templateRenderer: &TemplateRenderer{
			templates: template.Must(template.ParseGlob("gohtml/*.gohtml")),
		},
		stats: &Stats{
			RequestID:    "",
			NewRequest:   false,
			mutex:        sync.RWMutex{},
			CheckCounter: 0,
			StartTime:    time.Time{},
			StatErr:      nil,
			UniqueIP:     "",
		},
		dbWorker: &DBWorker{
			pool:          &sql.DB{},
			reqCtx:        nil,
			errorVar:      nil,
			postID:        0,
			ctrlID:        0,
			sqlQuery:      "",
			sqlValue:      "",
			sqlValues:     []interface{}{},
			sqlObject:     nil,
			sqlObjectList: nil,
			result:        "",
			stats:         sql.DBStats{},
			actionLog:     "",
		},
		customValidator: &CustomValidator{validator: validator.New()},
		errHandler: &ErrCodes{
			Codes: []struct {
				Code  string "json:\"code\""
				Title string "json:\"title\""
				Msg   string "json:\"msg\""
			}{},
		},
		enforcer:        NewEnforcer(),
		serviceConf:     &ServiceConf{appName: "", errLogFilePath: "logs/default.error.log", policyLogFilePath: "logs/default.policy.log"},
		rateLimiter:     &CustomRateLimiter{limiter: rate.NewLimiter(rate.Limit(25), 11)},
		customPageItems: &CustomPageItems{NavItems: &[]NavbarItem{}, AsideItems: &[]AsideItem{}, FooterItems: &[]FooterItem{}, PostIndex: &PostItem{}},
		customPageItem:  &CustomPageItem{NavItem: &NavbarItem{}, AsItem: &AsideItem{}, FoItem: &FooterItem{}, PoItem: &PostItem{}},
	}
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	switch viewContext := data.(type) {
	case map[string]interface{}:
		token, err := GetRandomToken32()
		if err != nil {
			return err
		}
		viewContext["token"] = token

		// add policy header if not in edit mode
		if viewContext["accessGranted"] != nil && !viewContext["accessGranted"].(bool) {
			// child-src 'none'; was removed due to embedding media content
			c.Response().Header().Add(echo.HeaderContentSecurityPolicy, "script-src 'strict-dynamic' 'nonce-"+token+"'; object-src 'none';")
		}

	default:
		c.Logger().Error(data)
		return echo.NewHTTPError(http.StatusBadRequest, data)

	}
	return t.templates.ExecuteTemplate(w, name, data)
}

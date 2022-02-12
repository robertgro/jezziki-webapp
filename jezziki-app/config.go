package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/ntp"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"

	"gopkg.in/yaml.v2"
)

type (
	// ErrCodes stores the err descriptions and titles from err.json
	ErrCodes struct {
		Codes []struct {
			Code  string `json:"code"`
			Title string `json:"title"`
			Msg   string `json:"msg"`
		} `json:"codes"`
	}

	// ServiceConf stores app instance related information
	ServiceConf struct {
		appName           string
		errLogFilePath    string
		policyLogFilePath string
	}

	// DefaultServerConf uses app-config.yaml for init
	DefaultServerConf struct {
		Title string `yaml:"title"`
		Debug bool   `yaml:"debug"`
		Port  string `yaml:"port"`
		Host  string `yaml:"host"`
		SSL   bool   `yaml:"ssl"`
		Paths struct {
			Cert   string `yaml:"cert"`
			Key    string `yaml:"key"`
			Static string `yaml:"static"`
			Public string `yaml:"public"`
			Dash   string `yaml:"dash"`
			Dist   string `yaml:"dist"`
		} `yaml:"paths"`
		Services map[string]int    `yaml:"services"`
		Rewrites map[string]string `yaml:"rewrites"`
		Links    struct {
			Favicon     string   `yaml:"favicon"`
			Stylesheets []string `yaml:"stylesheets"`
		} `yaml:"links"`
		Utils      []string `yaml:"utils"`
		Components []string `yaml:"components"`
	}
)

const (
	configFilePath = "app-config.yaml"
	errorFilePath  = "err.json"
	ntpServer      = "time.google.com"
	cspServer      = "http://localhost:* https://localhost:* http://127.0.0.1:* https://127.0.0.1:* http://0.0.0.0:* https://0.0.0.0:*"
	cspHeader      = "img-src " + cspServer + " 'self' https: data:; font-src https://fonts.gstatic.com https://cdn.jsdelivr.net; style-src https://fonts.googleapis.com https://cdn.jsdelivr.net " + cspServer + "; frame-src https://www.youtube.com https://www.youtube-nocookie.com/ https://player.twitch.tv/; child-src https://www.youtube.com https://www.youtube-nocookie.com/ https://player.twitch.tv/;"
	timeFormat     = "2006-01-02 15:04:05.000"
)

func (app *Controller) initEchoServer(e *echo.Echo) {
	e.Renderer = app.templateRenderer

	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	e.HTTPErrorHandler = app.CustomHTTPErrorHandler

	e.Validator = app.customValidator

	e.Server.Addr = app.defaultServerConf.Host + ":" + app.defaultServerConf.Port
	e.Server.ReadTimeout = 10 * time.Second
	e.Server.WriteTimeout = 10 * time.Second
}

func (c *DefaultServerConf) getAppConf() *DefaultServerConf {
	f, err := os.Open(configFilePath)
	Check(err, true)
	defer f.Close()
	byteValue, err := ioutil.ReadAll(f)
	Check(err, true)
	err = yaml.Unmarshal(byteValue, c)
	Check(err, true)
	return c
}

// Check checks for errors and panics if they not nil and pmode true else prints error msg
func Check(e error, pmode bool) {
	if e != nil {
		if pmode {
			panic(e)
		}
		fmt.Printf("\nERROR=%s\n", e.Error())
	}
}

// Loads the jezziki-config.yaml and replaces values wtih os args, also populates logo, service, host, port
func (app *Controller) initAppConfig() {

	app.loadErrCodes()

	app.defaultServerConf.getAppConf()

	print("\n- WEB APP NODE WORKER -\n\n")

	app.stats.StartTime, app.stats.StatErr = getTimestamp()

	Check(app.stats.StatErr, false)

	app.defaultServerConf.registerFlags()

	if app.serviceConf.appName = app.defaultServerConf.getServiceName(); app.serviceConf.appName != "" {
		println("\nService:\t\t" + app.serviceConf.appName)
		println("Start time:\t\t" + app.stats.StartTime.Format(timeFormat))
		app.serviceConf.errLogFilePath = "logs/" + app.serviceConf.appName + ".error.log"
		app.serviceConf.policyLogFilePath = "logs/" + app.serviceConf.appName + ".policy.log"
		println("ErrLogFilePath:\t\t" + app.serviceConf.errLogFilePath)
		println("PolicyLogFilePath:\t" + app.serviceConf.policyLogFilePath)
	} else {
		println("Error while getting appName")
	}

	println("Listening:\t\t" + app.defaultServerConf.Host + ":" + app.defaultServerConf.Port)
	print("Logger:\n")
}

func (c *DefaultServerConf) getServiceName() string {
	for k, v := range c.Services {
		if c.Port == strconv.Itoa(v) {
			return k
		}
	}
	return ""
}

func (c *DefaultServerConf) registerFlags() {
	// optional flags
	flag.Bool("h", false, "Help")
	flag.Bool("help", false, "Help Message")

	port := flag.String("port", "8080", "Custom Port")
	host := flag.String("host", "192.168.0.8", "Custom Host")
	cert := flag.String("cert", "/var/certs/cert.pem", "SSL Certification Path")
	key := flag.String("key", "/var/certs/key.pem", "SSL Key Path")

	flag.Parse()

	if flag.NFlag() > 0 {
		flag.Visit(func(f *flag.Flag) {
			switch f.Name {
			case "h", "help":
				print(`                                                         
	Usage: $go run . [OPTION] [VALUE]

	jezziki-server - powered by echo/labstack golang, written by RG
	(default config can altered in jezziki-config.yaml)
	
	Options:
	-h, -help, help           	display this usage message and exit
	-port=<number>  		port between 1-65535 (optional) (defaultPort ` + c.Port + `)
	-host=<hostname> 		public or local IP or DNS host entry (optional) (defaulHost ` + c.Host + `)
	-cert=</path/to/cert.pem>	path to the fullchain.pem file (optional) (defaultCertPath ` + c.Paths.Cert + `)
	-key=</path/to/key.pem>		path to your key.pem file (optional) (defaultKeyPath ` + c.Paths.Key + `)
				`)
				println()
				os.Exit(1)
			case "port":
				c.Port = *port
				println("[INFO] Set port to " + c.Port)
			case "host":
				c.Host = *host
				println("[INFO] Set host to " + c.Host)
			case "cert":
				c.Paths.Cert = *cert
				println("[INFO] Set certpath to " + c.Paths.Cert)
			case "key":
				c.Paths.Key = *key
				println("[INFO] Set keypath to " + c.Paths.Key)
			default:
				panic("Error reading os args")
			}
		})

	}
}

func (app *Controller) loadErrCodes() {
	f, err := os.Open(errorFilePath)
	Check(err, true)
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	Check(err, true)
	err = json.Unmarshal(b, app.errHandler)
	Check(err, true)
}

// ServeErrCtx serves the error page context
func ServeErrCtx(errorCodes *ErrCodes, code int, c echo.Context) {
	for _, item := range errorCodes.Codes {
		if item.Code == strconv.Itoa(code) {
			c.Render(code, "error.gohtml", map[string]interface{}{
				"code":  item.Code,
				"title": item.Title,
				"msg":   item.Msg,
			})
			return
		}
	}
}

// CustomHTTPErrorHandler utilizies custom error pages in public/error
func (app *Controller) CustomHTTPErrorHandler(err error, c echo.Context) {

	code := http.StatusInternalServerError
	he, ok := err.(*echo.HTTPError)

	if !ok {
		c.Logger().Error("HTTP CODE ERROR MSG=" + he.Error())
		code = http.StatusInternalServerError
	} else {
		code = he.Code
	}

	app.logXToFile(he.Message, code, c, app.serviceConf.errLogFilePath)

	ServeErrCtx(app.errHandler, code, c)
}

func getTimestamp() (time.Time, error) {
	startT := time.Time{}
	response, err := ntp.Time(ntpServer)
	Check(err, false)
	if err != nil {
		startT = time.Now()
	} else {
		startT = response
	}
	return startT, err
}

func getInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		println(err.Error())
		return 0
	}
	return i
}

// PreLogInfo prints relevant prelog infos
func PreLogInfo(c echo.Context, app *Controller) {
	ts, err := getTimestamp()
	Check(err, false)
	println(ts.Format(timeFormat) + " - INFO  - SERVICE " + app.serviceConf.appName + " " + app.defaultServerConf.Host + ":" + app.defaultServerConf.Port + " CLIENT " + getUserAgent(c.Request().UserAgent()) + " " + c.RealIP())
}

func getUserAgent(useragent string) string {

	s := strings.Split(useragent, " ")

	if len(s) > 1 {
		switch {
		case strings.Contains(s[len(s)-1], "Edg"):
			return "Edge"
		case strings.Contains(s[len(s)-1], "Firefox"):
			return "Firefox"
		case strings.Contains(s[len(s)-2], "Chrome"):
			return "Chrome"
		default:
			return useragent
		}
	}
	return useragent
}

func getLoggerConfig(c *Controller) mw.LoggerConfig {
	return mw.LoggerConfig{
		Skipper:          mw.DefaultSkipper,
		Format:           "${time_custom} - INFO  - ID ${id} IP ${remote_ip} PROT ${protocol} URI ${host}${uri} METHOD ${method} STATUS ${status} ERROR ${error}" + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.000",
	}
}

func getCorsConfig() mw.CORSConfig {
	return mw.CORSConfig{
		Skipper:      mw.DefaultSkipper,
		AllowOrigins: []string{},
		AllowHeaders: []string{},
		AllowMethods: []string{},
	}
}

func getSecureConfig() mw.SecureConfig {
	return mw.SecureConfig{
		Skipper:               mw.DefaultSkipper,
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: cspHeader,
	}
}

func getRequestIDConfig(customGen func() string) mw.RequestIDConfig {
	return mw.RequestIDConfig{
		Skipper: func(c echo.Context) bool {
			return c.Request().Header.Get(echo.HeaderXRequestID) != ""
		},
		Generator: customGen,
	}
}

func getBasicAuthConfig(co *Controller) mw.BasicAuthConfig {
	return mw.BasicAuthConfig{
		Skipper: func(c echo.Context) bool {
			u, pw, ok := c.Request().BasicAuth()
			if !ok || u != co.enforcer.dashU || pw != co.enforcer.dashP {
				return false
			}
			return true
		},
		Validator: func(u, pw string, c echo.Context) (bool, error) {
			if u == co.enforcer.dashU && pw == co.enforcer.dashP {
				return true, nil
			}
			return false, echo.ErrUnauthorized
		},
		Realm: co.enforcer.realm,
	}
}

// logXToFile logs custom http errors to a file
func (app *Controller) logXToFile(msg interface{}, code int, c echo.Context, filePath string) {

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	ts, err := getTimestamp()

	Check(err, false)

	if m, ok := msg.(string); ok {
		msg = ts.Format(timeFormat) + " - SERVICE - EVENT - ID " + app.stats.RequestID + " REMOTE_IP " + c.Request().RemoteAddr + " REAL_IP " + c.RealIP() + " HOST_URL " + c.Request().Host + c.Request().RequestURI + " METHOD " + c.Request().Method + " CODE " + strconv.Itoa(code) + " MESSAGE " + m + "\n"
	}

	if _, err := f.Write([]byte(msg.(string))); err != nil {
		f.Close() // ignore error; Write error takes precedence
		log.Fatal(err)
		os.Exit(-1)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}

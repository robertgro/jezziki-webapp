package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/ntp"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	"github.com/segmentio/ksuid"
	"gopkg.in/yaml.v2"
)

const (
	configFilePath = "proxy-config.yaml"
	timeFormat     = "2006-01-02 15:04:05.000"
	ntpServer      = "time.google.com"
	dbUsername     = "postgres"
	dbPassword     = "vEKsn2b8RaKEQRZ6"
	dbHost         = "localhost"
	dbPort         = 5432
	dbDBname       = "jezziki"
	dbParamSSLMode = "disable"
)

func (pc *proxyController) initApp(e *echo.Echo) {

	pc.DefaultProxyConf.getProxyConf()

	pc.getTargetNodes()

	pc.parseProxyTargets(e)

	e.Server.Addr = pc.DefaultProxyConf.Host + ":" + pc.DefaultProxyConf.Port
	e.Server.ReadTimeout = 10 * time.Second
	e.Server.WriteTimeout = 10 * time.Second

	e.IPExtractor = echo.ExtractIPDirect()

	pc.startTime, pc.db.errorVar = getTimestamp()

	Check(pc.db.errorVar, false)

	print("\nWEB APP NODE PROXY + LOAD BALANCER\n\n")
	println("Service:\tApp proxy\n")
	println("Start time:\t" + pc.startTime.Format(timeFormat))
	print("Pool:\t\t")
	for _, node := range pc.URLTargets.Nodes {
		print(node.Name + " ")
	}
	println("\nListening:\t" + e.Server.Addr)
	print("Logger:\n")

}

func (pc *proxyController) getCorsConfigMW() echo.MiddlewareFunc {
	return mw.CORSWithConfig(mw.CORSConfig{
		Skipper:      mw.DefaultSkipper,
		AllowOrigins: []string{},
		AllowHeaders: []string{},
		AllowMethods: []string{},
	})
}

func (pc *proxyController) getLoggerConfigMW() echo.MiddlewareFunc {
	return mw.LoggerWithConfig(mw.LoggerConfig{
		Skipper:          mw.DefaultSkipper,
		Format:           "${time_custom} - INFO  - ID ${id} IP ${remote_ip} PROT ${protocol} URI ${host}${uri} METHOD ${method} STATUS ${status} ERROR ${error}" + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.000",
	})
}

func (pc *proxyController) getRequestIDConfigMW() echo.MiddlewareFunc {
	pc.reqStats.mutex.Lock()
	defer pc.reqStats.mutex.Unlock()
	return mw.RequestIDWithConfig(mw.RequestIDConfig{
		Skipper: func(c echo.Context) bool {
			return c.Request().Header.Get(echo.HeaderXRequestID) != ""
		},
		Generator: pc.customGenerator,
	})
}

func (pc *proxyController) getTargetNodes() {
	f, err := os.Open(pc.DefaultProxyConf.Paths.Nodes)
	Check(err, true)
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	Check(err, true)
	err = json.Unmarshal(b, pc.URLTargets)
	Check(err, true)
}

func (c *DefaultProxyConf) getProxyConf() *DefaultProxyConf {
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

// PreLogInfo prints relevant prelog infos
func (pc *proxyController) PreLogInfo(c echo.Context) {
	ts, err := getTimestamp()

	Check(err, false)

	println(ts.Format(timeFormat) + " - INFO  - SERVICE " + pc.DefaultProxyConf.Name + " " + pc.DefaultProxyConf.Host + ":" + pc.DefaultProxyConf.Port + " CLIENT " + getUserAgent(c.Request().UserAgent()) + " " + c.RealIP())
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

// logXToFile logs custom http errors to a file
func (pc *proxyController) logXToFile(msg interface{}, code int, c echo.Context, filePath string) {

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	ts, err := getTimestamp()

	Check(err, false)

	if m, ok := msg.(string); ok {
		msg = ts.Format(timeFormat) + " - SERVICE - EVENT - ID" + pc.reqStats.RequestID + " REMOTE_IP " + c.Request().RemoteAddr + " REAL_IP " + c.RealIP() + " HOST_URL " + c.Request().Host + c.Request().RequestURI + " METHOD " + c.Request().Method + " CODE " + strconv.Itoa(code) + " MESSAGE " + m + "\n"
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

// getWhoisInfo queries whoisinfo
func (pc *proxyController) getWhoisInfo(c echo.Context) (err error) {

	wi := struct {
		cmdE       *exec.Cmd
		shellPath  string
		scriptName string
		argIP      string
		argID      string
	}{
		cmdE:       &exec.Cmd{},
		shellPath:  "",
		scriptName: "",
		argIP:      c.RealIP(),
		argID:      "",
	}

	if pc.reqStats.RequestID == "" {
		pc.reqStats.RequestID = pc.customGenerator()
	}

	wi.argID = pc.reqStats.RequestID

	opsys := runtime.GOOS

	switch opsys {
	case "windows":
		wi.scriptName = "getlocbyip"
		wi.cmdE = exec.Command(wi.scriptName, wi.argIP, wi.argID)
	case "darwin":
		fmt.Printf("%s.\n", opsys)
		return echo.NewHTTPError(http.StatusBadRequest, "DarwinOS")
	case "linux":
		wi.shellPath = "/bin/sh"
		wi.scriptName = "querywhois.sh"
		wi.cmdE = exec.Command(wi.shellPath, wi.scriptName, wi.argIP, wi.argID)
	default:
		fmt.Printf("%s.\n", opsys)
		return echo.NewHTTPError(http.StatusBadRequest, "Unknown os")
	}

	if err = wi.cmdE.Start(); err != nil {
		return err
	}

	return nil
}

// for request id generation
func (pc *proxyController) customGenerator() string {
	pc.reqStats.mutex.Lock()
	defer pc.reqStats.mutex.Unlock()
	if pc.reqStats.NewRequest || pc.reqStats.RequestID == "" {
		id := ksuid.New()
		pc.reqStats.NewRequest = false
		pc.reqStats.RequestID = id.String()
		return pc.reqStats.RequestID
	}
	return pc.reqStats.RequestID
}

// checkWhoisDone checks if bat file created whois stdout file
func (pc *proxyController) checkWhoisDone() (err error) {

	filename := "whois/" + pc.reqStats.RequestID + ".log"

	// file does exist
	if fI, err := os.Stat(filename); err == nil {
		if fI.Size() > 0 {
			b, err := ioutil.ReadFile(filename)

			if err != nil {
				return err
			}

			pc.reqStats.WResult = string(b)

			pc.reqStats.WDone = true
		}
	}

	if pc.reqStats.WDone {

		if err = os.Remove(filename); err != nil {
			return err
		}
	}

	return nil
}

// getWhoisResult processes the whois result
func (pc *proxyController) getWhoisResult() (err error) {
	scn := bufio.NewScanner(strings.NewReader(pc.reqStats.WResult))

	fulladdr := ""

	for scn.Scan() {
		t := scn.Text()
		//println("SCAN-LINE=", t)

		if strings.Contains(t, "descr:") { // unix
			pc.reqStats.UProvider = strings.TrimSpace(strings.Replace(t, "descr:", "", -1))
		} else if strings.Contains(t, "address:") { // unix
			fulladdr += strings.TrimSpace(strings.Replace(t, "address:", "", -1)) + " "
		} else if strings.Contains(t, "country:") { // unix
			pc.reqStats.UCountry = strings.TrimSpace(strings.Replace(t, "country:", "", -1))
		} else if strings.Contains(t, "Domain Name:") { // win
			pc.reqStats.UProvider = strings.TrimSpace(strings.Replace(t, "Domain Name:", "", -1))
		} else if strings.Contains(t, "Registrant Country:") { // win
			pc.reqStats.UCountry = strings.TrimSpace(strings.Replace(t, "Registrant Country:", "", -1))
		} else if strings.Contains(t, "Registrant State/Province:") { // win
			fulladdr = strings.TrimSpace(strings.Replace(t, "Registrant State/Province:", "", -1))
		}
	}

	pc.reqStats.UCity = fulladdr

	if err = scn.Err(); err != nil {
		return err
	}

	pc.reqStats.ScanDone = true

	return nil
}

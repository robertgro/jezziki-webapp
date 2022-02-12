package main

import (
	"context"
	"database/sql"
	"net/url"
	"sync"
	"time"

	mw "github.com/labstack/echo/v4/middleware"
)

type (
	proxyController struct {
		pURLs            []*url.URL
		proxyTargets     []*mw.ProxyTarget
		startTime        time.Time
		db               *DBController
		reqStats         *ReqStats
		mutex            sync.RWMutex
		visitors         *Visitors
		URLTargets       *Targets
		DefaultProxyConf *DefaultProxyConf
		enforcer         *Enforcer
	}

	Enforcer struct {
		tokenID string
		tEnd    time.Time
		tNow    time.Time
		tDiff   time.Duration
		access  bool
	}

	DBController struct {
		pool     *sql.DB
		reqCtx   context.Context
		errorVar error
		result   string
		sqlQuery string
	}

	ReqStats struct {
		RequestID    string
		UniqueIP     string
		NewRequest   bool
		Blacklisted  bool
		Logged       bool
		CheckCounter int
		UProvider    string
		UCountry     string
		UCity        string
		WDone        bool
		WResult      string
		ScanDone     bool
		StartTime    time.Time
		StatErr      error
		mutex        sync.RWMutex
	}

	Visitors struct {
		Users *[]Visitor `json:"users"`
	}

	Visitor struct {
		VID int    `json:"visitorid"`
		RID string `json:"reqid"`
		RIP string `json:"realip"`
		UA  string `json:"useragent"`
		BL  bool   `json:"blacklisted"`
	}

	Targets struct {
		Nodes []struct {
			Name string `json:"name"`
			Host string `json:"host"`
			Port string `json:"port"`
		} `json:"nodes"`
	}

	DefaultProxyConf struct {
		Name  string `yaml:"name"`
		Debug bool   `yaml:"debug"`
		Port  string `yaml:"port"`
		Host  string `yaml:"host"`
		SSL   bool   `yaml:"ssl"`
		Paths struct {
			Cert  string `yaml:"cert"`
			Key   string `yaml:"key"`
			Pol   string `yaml:"pol"`
			Nodes string `yaml:"nodes"`
		} `yaml:"paths"`
	}
)

func NewProxyController() *proxyController {
	return &proxyController{
		pURLs:        []*url.URL{},
		proxyTargets: []*mw.ProxyTarget{},
		startTime:    time.Time{},
		db:           &DBController{pool: &sql.DB{}, reqCtx: nil, errorVar: nil, result: "", sqlQuery: ""},
		reqStats: &ReqStats{
			RequestID:    "",
			UniqueIP:     "",
			NewRequest:   false,
			Blacklisted:  false,
			Logged:       false,
			CheckCounter: 0,
			UProvider:    "",
			UCountry:     "",
			UCity:        "",
			WDone:        false,
			WResult:      "",
			ScanDone:     false,
			StartTime:    time.Time{},
			StatErr:      nil,
			mutex:        sync.RWMutex{},
		},
		mutex:    sync.RWMutex{},
		visitors: &Visitors{Users: &[]Visitor{}},
		URLTargets: &Targets{Nodes: []struct {
			Name string "json:\"name\""
			Host string "json:\"host\""
			Port string "json:\"port\""
		}{}},
		DefaultProxyConf: &DefaultProxyConf{Name: "", Debug: false, Port: "", Host: "", SSL: false, Paths: struct {
			Cert  string "yaml:\"cert\""
			Key   string "yaml:\"key\""
			Pol   string "yaml:\"pol\""
			Nodes string "yaml:\"nodes\""
		}{Cert: "", Key: "", Pol: "", Nodes: ""}},
		enforcer: &Enforcer{
			tokenID: "",
			tEnd:    time.Time{},
			tNow:    time.Time{},
			tDiff:   0,
			access:  false,
		},
	}
}

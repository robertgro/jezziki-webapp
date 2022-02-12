package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

// PrepareDB opens the db connection pool and sets limits appropriately
func (pc *proxyController) PrepareDB(c echo.Context) error {

	pc.db.pool, pc.db.errorVar = sql.Open("postgres", pc.getDSN())

	if pc.db.errorVar != nil {
		return pc.db.errorVar
	}

	pc.db.pool.SetConnMaxLifetime(0)
	pc.db.pool.SetMaxIdleConns(3)
	pc.db.pool.SetMaxOpenConns(3)

	if c == nil {
		pc.db.reqCtx = nil
		return nil
	}

	pc.db.reqCtx = c.Request().Context()

	return nil
}

// PingDB audits the db connection
func (pc *proxyController) PingDB() error {

	if pc.db.reqCtx == nil {
		err := pc.db.pool.Ping()

		if err != nil {
			return err
		}

		return nil
	}

	//Prepares a cancel context for db connection audit
	ctx, stop := context.WithCancel(pc.db.reqCtx)
	defer stop()

	err := pc.PingContext(ctx)

	if err != nil {
		return err
	}
	return nil
}

// PingContext pings the database to verify DSN provided by the user is valid and the
// server accessible. If the ping fails exit the program with an error.
func (pc *proxyController) PingContext(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pc.db.pool.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

// getDSN collects the connection parameters
func (pc *proxyController) getDSN() string {
	return fmt.Sprintf("port=%d host=%s user=%s password=%s dbname=%s sslmode=%s", dbPort, dbHost, dbUsername, dbPassword, dbDBname, dbParamSSLMode)
}

// queryDBVisitorLog saves info about visitors visitor_id, request_id, real_ip, remote_addr, host_port, user_agent, referer, method_requri, req_url
func (pc *proxyController) queryDBVisitorLog(request_id string, real_ip string, remote_addr string, host_post string, user_agent string, referer string, method_requri string, req_url string, provider string, country string, city string, blacklist bool, protocol string) error {

	pc.db.sqlQuery = `INSERT INTO visitor (request_id, real_ip, remote_addr, host_port, user_agent, referer, method_requri, req_url, provider, country, city, blacklist, prot) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13);`

	ctx, cancel := context.WithTimeout(pc.db.reqCtx, 10*time.Second)
	defer cancel()

	if _, err := pc.db.pool.ExecContext(ctx, pc.db.sqlQuery, request_id, real_ip, remote_addr, host_post, user_agent, referer, method_requri, req_url, "", "", "", blacklist, protocol); err != nil {
		return err
	}

	ts, err := getTimestamp()

	Check(err, false)

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %.*v...\n", 100, pc.db.sqlQuery)

	return nil
}

// queryDBNewUserAgent queries for existing user agent visitor id to register a new device
func (pc *proxyController) queryDBNewUserAgent(c echo.Context) (bool, error) {

	pc.db.sqlQuery = `SELECT user_agent FROM visitor WHERE real_ip = $1 AND user_agent = $2;`

	ctx, cancel := context.WithTimeout(pc.db.reqCtx, 10*time.Second)
	defer cancel()

	row := pc.db.pool.QueryRowContext(ctx, pc.db.sqlQuery, c.RealIP(), c.Request().UserAgent())

	ts, err := getTimestamp()

	Check(err, false)

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", pc.db.sqlQuery)

	switch pc.db.errorVar = row.Scan(&pc.db.result); pc.db.errorVar {
	case sql.ErrNoRows:
		return true, nil
	}

	// will be executed from third round trip on if no new User Agent
	// 1) getwhoisInfo -> create File
	// 2) readFile, delFile -> ScanDone, WDone
	// 3) Update Last Seen

	if pc.reqStats.WDone && pc.reqStats.ScanDone {
		pc.db.sqlQuery = `UPDATE visitor SET last_seen = $1, provider = $2, country = $3, city = $4, method_requri = $5, req_url = $6, prot = $7 WHERE real_ip = $8 AND user_agent = $9;`

		ts, err := getTimestamp()

		Check(err, false)

		fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", pc.db.sqlQuery)

		if _, pc.db.errorVar = pc.db.pool.ExecContext(ctx, pc.db.sqlQuery, time.Now(), pc.reqStats.UProvider, pc.reqStats.UCountry, pc.reqStats.UCity, c.Request().Method+" "+c.Request().RequestURI, c.Request().URL.Path, c.Request().Proto, c.RealIP(), c.Request().UserAgent()); pc.db.errorVar != nil {
			return false, pc.db.errorVar
		}
	}

	return false, nil
}

// queryDBVisitorBlacklist looks if the ip is blacklisted to refuse a connection in case
func (pc *proxyController) queryDBVisitorBlacklist(c echo.Context) error {

	v := &Visitor{}

	pc.db.sqlQuery = `SELECT visitor_id, request_id, real_ip, user_agent, blacklist FROM visitor WHERE real_ip = $1;`

	ctx, cancel := context.WithTimeout(pc.db.reqCtx, 10*time.Second)
	defer cancel()

	rows, err := pc.db.pool.QueryContext(ctx, pc.db.sqlQuery, pc.reqStats.UniqueIP)

	if err != nil {
		return err

	}

	ts, err := getTimestamp()

	Check(err, false)

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", pc.db.sqlQuery)

	defer rows.Close()
	for rows.Next() {
		switch pc.db.errorVar = rows.Scan(&v.VID, &v.RID, &v.RIP, &v.UA, &v.BL); pc.db.errorVar {
		case sql.ErrNoRows:
			pc.reqStats.Logged = false
			return nil
		case nil:
			*pc.visitors.Users = append(*pc.visitors.Users, *v)
		default:
			return pc.db.errorVar
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	for _, e := range *pc.visitors.Users {
		if e.BL {
			println(ts.Format(timeFormat)+" - INFO  - SECURITY - CLIENT BLACKLISTED", "VID "+strconv.Itoa(v.VID), "RID "+v.RID, "RIP "+v.RIP, "BL "+strconv.FormatBool(v.BL))
			pc.reqStats.Blacklisted = true
			pc.logXToFile("CLIENT BLACKLISTED VID "+strconv.Itoa(v.VID), http.StatusTeapot, c, pc.DefaultProxyConf.Paths.Pol)
			return nil
		}
	}

	return nil
}

func (pc *proxyController) getEditSessionStatus(c echo.Context) (bool, error) {

	pc.db.sqlQuery = `SELECT date_trunc('SECOND', date_expired), token, active FROM control where ip = $1 AND active = $2;`
	ctx, cancel := context.WithTimeout(pc.db.reqCtx, 10*time.Second)
	defer cancel()

	row := pc.db.pool.QueryRowContext(ctx, pc.db.sqlQuery, c.RealIP(), true)

	ts, err := getTimestamp()

	Check(err, false)

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", pc.db.sqlQuery)

	switch pc.db.errorVar = row.Scan(&pc.enforcer.tEnd, &pc.enforcer.tokenID, &pc.enforcer.access); pc.db.errorVar {
	case sql.ErrNoRows:
		return false, sql.ErrNoRows
	case nil:
		pc.enforcer.tNow, pc.db.errorVar = time.Parse(timeFormat, time.Now().Format(timeFormat))
		if pc.db.errorVar != nil {
			return false, pc.db.errorVar
		}

		pc.enforcer.tEnd, pc.db.errorVar = time.Parse(timeFormat, pc.enforcer.tEnd.Format(timeFormat))
		if pc.db.errorVar != nil {
			return false, pc.db.errorVar
		}

		pc.enforcer.tDiff = pc.enforcer.tEnd.Sub(pc.enforcer.tNow)
		if pc.enforcer.tDiff.Seconds() > 0 {
			return true, nil
		}
		return false, nil
	default:
		return false, pc.db.errorVar
	}
}

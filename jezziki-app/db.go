package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type (
	// DBWorker is responsible for handling the db pool
	DBWorker struct {
		pool          *sql.DB
		reqCtx        context.Context
		errorVar      error
		postID        int
		ctrlID        int
		sqlQuery      string
		sqlValue      string
		sqlValues     []interface{}
		sqlObject     interface{}
		sqlObjectList interface{}
		result        string
		stats         sql.DBStats
		actionLog     string
		byteBuffer    []byte
	}
)

const (
	dbUsername     = "postgres"
	dbPassword     = "vEKsn2b8RaKEQRZ6"
	dbHost         = "localhost"
	dbPort         = 5432
	dbDBname       = "jezziki"
	dbParamSSLMode = "disable"
)

// PrepareDB opens the db connection pool and sets limits appropriately
func (db *DBWorker) PrepareDB(c echo.Context) error {

	db.pool, db.errorVar = sql.Open("postgres", db.getDSN())

	if db.errorVar != nil {
		return db.errorVar
	}

	db.pool.SetConnMaxLifetime(0)
	db.pool.SetMaxIdleConns(3)
	db.pool.SetMaxOpenConns(3)

	if c == nil {
		db.reqCtx = nil
		return nil
	}

	db.reqCtx = c.Request().Context()

	return nil
}

// PingDB audits the db connection
func (db *DBWorker) PingDB() error {

	if db.reqCtx == nil {
		err := db.pool.Ping()

		if err != nil {
			return err
		}

		return nil
	}

	//Prepares a cancel context for db connection audit
	ctx, stop := context.WithCancel(db.reqCtx)
	defer stop()

	if err := db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

// PingContext pings the database to verify DSN provided by the user is valid and the
// server accessible. If the ping fails exit the program with an error.
func (db *DBWorker) PingContext(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.pool.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

// queryDB gets a single result
func (db *DBWorker) queryDB() error {

	ctx, cancel := context.WithTimeout(db.reqCtx, 5*time.Second)
	defer cancel()

	ts, err := getTimestamp()

	Check(err, false)

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v VALUE %v\n", db.sqlQuery, db.sqlValue)

	switch err := db.pool.QueryRowContext(ctx, db.sqlQuery, db.sqlValue).Scan(&db.result); err {
	case sql.ErrNoRows:
		db.result = "No entry available."
		return sql.ErrNoRows
	case nil:
		db.getConnectionInfo()
	default:
		return err
	}
	return nil
}

// queryDBresults queries the db and loads multiple object values into object list via reflection
// https://stackoverflow.com/questions/56525471/how-to-use-rows-scan-of-gos-database-sql
func (db *DBWorker) queryDBResults() error {

	var columnNames []string

	ctx, cancel := context.WithTimeout(db.reqCtx, 10*time.Second)
	defer cancel()

	rows, err := db.pool.QueryContext(ctx, db.sqlQuery)

	if err != nil {
		return err
	}

	formatPattern := regexp.MustCompile(`\s+`)
	logQuery := formatPattern.ReplaceAllString(db.sqlQuery, " ")

	ts, err := getTimestamp()

	Check(err, false)

	// https://stackoverflow.com/questions/2239519/is-there-a-way-to-specify-how-many-characters-of-a-string-to-print-out-using-pri
	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %.*v...\n", 100, logQuery)

	if columnNames, db.errorVar = rows.Columns(); db.errorVar != nil {
		return db.errorVar
	}

	sqlObjPtr := ptr(reflect.ValueOf(db.sqlObject))      // get pointer to single struct object
	sqlObjLPtr := ptr(reflect.ValueOf(db.sqlObjectList)) // get pointer to slice struct object

	sqlObj := sqlObjPtr.Elem().Elem()   // get value of single struct object
	sqlObjL := sqlObjLPtr.Elem().Elem() // get value of slice struct object

	ObjNumF := sqlObj.NumField()         // get the number of Fields of single struct object
	cols := make([]interface{}, ObjNumF) // create columns iface slice with x num of fields

	for i := 0; i < ObjNumF; i++ {
		field := sqlObj.Field(i)
		cols[i] = field.Addr().Interface() // transform field into columns using pointer interface
	}

	if len(columnNames) != len(cols) {
		print(ts.Format(timeFormat) + " - ERROR - MESSAGE Query columns length doesn't match struct column count. In need of query or struct fix.")
		println(db.sqlQuery, len(columnNames), len(cols), sqlObj.NumField())
	}

	defer rows.Close()
	for rows.Next() {
		switch db.errorVar = rows.Scan(cols...); db.errorVar {
		case sql.ErrNoRows:
			db.result = "No entry available."
			return sql.ErrNoRows
		case nil:
			db.insertDBresults(sqlObj, cols, sqlObjL)
		default:
			return db.errorVar
		}
	}

	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	if db.errorVar = rows.Close(); db.errorVar != nil {
		return db.errorVar
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	// source https://golang.org/src/database/sql/example_test.go
	if db.errorVar = rows.Err(); db.errorVar != nil {
		return db.errorVar
	}

	return nil
}

// queryDBUpdateComponentItem deletes the old component item table and creates a new item index with known post_id's
// if title doesn't match (foreignid = 0) it creates a new post
func (db *DBWorker) queryDBUpdateComponentItem(itemName string) error {
	return nil
}

// queryDBUpdateNav deletes the old nav table and recreates the nav index with known post_id's
// if title doesn't match (foreignid = 0) it creates a new post
func (db *DBWorker) queryDBUpdateNav(navitems *CustomPageItems) error {

	db.sqlValues = nil

	db.sqlQuery = `DELETE FROM navbar;`

	ctx, cancel := context.WithTimeout(db.reqCtx, 10*time.Second)
	defer cancel()

	if _, err := db.pool.ExecContext(ctx, db.sqlQuery); err != nil {
		return err
	}

	ts, err := getTimestamp()

	Check(err, false)

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", db.sqlQuery)

	db.sqlQuery = `INSERT INTO navbar (navbar_id, posts_id, parent_id) VALUES `

	// columnlen must match sqlQuery's column count
	db.getInsertQuery(len(*navitems.NavItems), 3)

	for _, e := range *navitems.NavItems {
		if e.FID == 0 {
			createPostQuery := `INSERT INTO posts (title) VALUES ($1);`
			selectIDQuery := `SELECT currval('posts_id_seq');`

			fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", createPostQuery)

			if _, err := db.pool.ExecContext(ctx, createPostQuery, e.Title); err != nil {
				return err
			}

			fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", selectIDQuery)

			if err := db.pool.QueryRowContext(ctx, selectIDQuery).Scan(&e.FID); err != nil {
				return err
			}
		}
		db.sqlValues = append(db.sqlValues, e.ID, e.FID, e.PID)
	}

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", db.sqlQuery)

	if _, err := db.pool.ExecContext(ctx, db.sqlQuery, db.sqlValues...); err != nil {
		return err
	}

	db.sqlQuery = "INSERT INTO admlog (control_id, action) VALUES ($1,$2);"

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", db.sqlQuery)

	if _, err := db.pool.ExecContext(ctx, db.sqlQuery, db.ctrlID, db.actionLog); err != nil {
		return err
	}

	return nil
}

// queryDBUpdateAside deletes the old aside table and recreates the aside index with known post_id's
// if title doesn't match (foreignid = 0) it creates a new post
func (db *DBWorker) queryDBUpdateAside(asideitems *CustomPageItems) error {

	db.sqlValues = nil

	db.sqlQuery = `DELETE FROM aside;`

	ctx, cancel := context.WithTimeout(db.reqCtx, 10*time.Second)
	defer cancel()

	if _, err := db.pool.ExecContext(ctx, db.sqlQuery); err != nil {
		return err
	}

	ts, err := getTimestamp()

	Check(err, false)

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", db.sqlQuery)

	db.sqlQuery = `INSERT INTO aside (aside_id, posts_id) VALUES `

	// columnlen must match sqlQuery's column count
	db.getInsertQuery(len(*asideitems.AsideItems), 2)

	for _, e := range *asideitems.AsideItems {
		if e.FID == 0 {
			createPostQuery := `INSERT INTO posts (title) VALUES ($1);`
			selectIDQuery := `SELECT currval('posts_id_seq');`

			fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", createPostQuery)

			if _, err := db.pool.ExecContext(ctx, createPostQuery, e.Title); err != nil {
				return err
			}

			fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", selectIDQuery)

			if err := db.pool.QueryRowContext(ctx, selectIDQuery).Scan(&e.FID); err != nil {
				return err
			}
		}
		db.sqlValues = append(db.sqlValues, e.ID, e.FID)
	}

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", db.sqlQuery)

	if _, err := db.pool.ExecContext(ctx, db.sqlQuery, db.sqlValues...); err != nil {
		return err
	}

	db.sqlQuery = "INSERT INTO admlog (control_id, action) VALUES ($1,$2);"

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", db.sqlQuery)

	if _, err := db.pool.ExecContext(ctx, db.sqlQuery, db.ctrlID, db.actionLog); err != nil {
		return err
	}

	return nil
}

// queryDBUpdateFooter deletes the old footer table and recreates the footer index with known post_id's
// if title doesn't match (foreignid = 0) it creates a new post
func (db *DBWorker) queryDBUpdateFooter(footeritems *CustomPageItems) error {

	db.sqlValues = nil

	db.sqlQuery = `DELETE FROM footer;`

	ctx, cancel := context.WithTimeout(db.reqCtx, 10*time.Second)
	defer cancel()

	if _, err := db.pool.ExecContext(ctx, db.sqlQuery); err != nil {
		return err
	}

	ts, err := getTimestamp()

	Check(err, false)

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", db.sqlQuery)

	db.sqlQuery = `INSERT INTO footer (footer_id, posts_id) VALUES `

	// columnlen must match sqlQuery's column count
	db.getInsertQuery(len(*footeritems.FooterItems), 2)

	for _, e := range *footeritems.FooterItems {
		if e.FID == 0 {
			createPostQuery := `INSERT INTO posts (title) VALUES ($1);`
			selectIDQuery := `SELECT currval('posts_id_seq');`

			fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", createPostQuery)

			if _, err := db.pool.ExecContext(ctx, createPostQuery, e.Title); err != nil {
				return err
			}

			fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", selectIDQuery)

			if err := db.pool.QueryRowContext(ctx, selectIDQuery).Scan(&e.FID); err != nil {
				return err
			}
		}

		db.sqlValues = append(db.sqlValues, e.ID, e.FID)
	}

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", db.sqlQuery)

	if _, err := db.pool.ExecContext(ctx, db.sqlQuery, db.sqlValues...); err != nil {
		return err
	}

	db.sqlQuery = "INSERT INTO admlog (control_id, action) VALUES ($1,$2);"

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", db.sqlQuery)

	if _, err := db.pool.ExecContext(ctx, db.sqlQuery, db.ctrlID, db.actionLog); err != nil {
		return err
	}

	return nil
}

// queryCtrlID gets the control id and logs the action
func (db *DBWorker) queryCtrlID(tokenID string, dateExpired time.Time) error {

	db.sqlQuery = `SELECT control_id FROM control where token = $1 AND date_trunc('second', date_expired) = $2 AND active = $3;`

	ctx, cancel := context.WithTimeout(db.reqCtx, 10*time.Second)
	defer cancel()

	row := db.pool.QueryRowContext(ctx, db.sqlQuery, tokenID, dateExpired.Format(timeFormat), true)

	ts, err := getTimestamp()

	Check(err, false)

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", db.sqlQuery)

	switch db.errorVar = row.Scan(&db.ctrlID); db.errorVar {
	case sql.ErrNoRows:
		return sql.ErrNoRows
	case nil:
		return nil
	default:
		return db.errorVar
	}
}

// queryDBCtrlLog stores token access info in db
func (db *DBWorker) queryDBCtrlLog(dateCreated time.Time, dateExpired time.Time, accessToken string, ipAddr string) error {

	db.sqlQuery = `INSERT INTO control (date_created, date_expired, token, ip, active) VALUES($1,$2,$3,$4,$5);`

	ctx, cancel := context.WithTimeout(db.reqCtx, 10*time.Second)
	defer cancel()

	if _, err := db.pool.ExecContext(ctx, db.sqlQuery, dateCreated, dateExpired, accessToken, ipAddr, true); err != nil {
		return err
	}

	ts, err := getTimestamp()

	Check(err, false)

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", db.sqlQuery)

	return nil
}

func (db *DBWorker) disableEditSession(c echo.Context, tokenID string, active bool, timeBased bool) error {

	if active {
		db.sqlQuery = `UPDATE control SET active = $1 WHERE ip = $2 AND token = $3;`

		if timeBased {
			if _, err := db.pool.Exec(db.sqlQuery, false, c.RealIP(), tokenID); err != nil {
				return err
			}
		} else {
			if _, err := db.pool.ExecContext(db.reqCtx, db.sqlQuery, false, c.RealIP(), tokenID); err != nil {
				return err
			}
		}

	} else {

		db.sqlQuery = `UPDATE control SET date_expired = $1 WHERE ip = $2 AND token = $3;`

		if timeBased {
			if _, err := db.pool.Exec(db.sqlQuery, time.Now(), c.RealIP(), tokenID); err != nil {
				return err
			}
		} else {
			if _, err := db.pool.ExecContext(db.reqCtx, db.sqlQuery, time.Now(), c.RealIP(), tokenID); err != nil {
				return err
			}
		}
	}
	ts, err := getTimestamp()

	Check(err, false)

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", db.sqlQuery)

	return nil
}

func (db *DBWorker) getEditSessionStatus(c echo.Context, app *Controller) (bool, error) {

	db.sqlQuery = `SELECT date_trunc('SECOND', date_expired), token, active FROM control where ip = $1 AND active = $2;`
	ctx, cancel := context.WithTimeout(db.reqCtx, 10*time.Second)
	defer cancel()

	row := db.pool.QueryRowContext(ctx, db.sqlQuery, c.RealIP(), true)

	ts, err := getTimestamp()

	Check(err, false)

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", db.sqlQuery)

	switch db.errorVar = row.Scan(&app.enforcer.tEnd, &app.enforcer.tokenID, &app.enforcer.accessGranted); db.errorVar {
	case sql.ErrNoRows:
		return false, sql.ErrNoRows
	case nil:
		app.enforcer.tNow, db.errorVar = time.Parse(timeFormat, time.Now().Format(timeFormat))
		if db.errorVar != nil {
			return false, db.errorVar
		}

		app.enforcer.tEnd, db.errorVar = time.Parse(timeFormat, app.enforcer.tEnd.Format(timeFormat))
		if db.errorVar != nil {
			return false, db.errorVar
		}

		app.enforcer.tDiff = app.enforcer.tEnd.Sub(app.enforcer.tNow)

		if app.enforcer.tDiff.Seconds() > 0 {
			return true, nil
		} else {
			if db.errorVar = app.RevokeAccess(c, true, false); db.errorVar != nil {
				return false, db.errorVar
			}
		}
		return false, nil
	default:
		return false, db.errorVar
	}
}

func (db *DBWorker) insertDBresults(rv reflect.Value, cols []interface{}, lv reflect.Value) {
	// interface{} v is a pointer to a type, need to dereference properly
	for i, v := range cols {
		switch v := v.(type) {
		case *int:
			rv.Field(i).SetInt(int64(*v))
		case *string:
			rv.Field(i).SetString(*v)
		case *bool:
			rv.Field(i).SetBool(*v)
		default:
			fmt.Printf("\nType assert not defined. Type is %T\n", v)
		}
	}
	lv.Set(reflect.Append(lv, rv))
}

func (db *DBWorker) updatePost(postItem *PostItemNew) (err error) {

	pItem := struct {
		ID    string
		Title string
		Text  string
		Ext   string
	}{
		ID:    postItem.ID,
		Title: "",
		Text:  "",
		Ext:   "",
	}

	db.sqlQuery = `SELECT title, content, external FROM posts WHERE posts_id = $1;`

	ctx, cancel := context.WithTimeout(db.reqCtx, 10*time.Second)
	defer cancel()

	row := db.pool.QueryRowContext(ctx, db.sqlQuery, postItem.ID)

	ts, err := getTimestamp()

	Check(err, false)

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", db.sqlQuery)

	switch db.errorVar = row.Scan(&pItem.Title, &pItem.Text, &pItem.Ext); db.errorVar {
	case sql.ErrNoRows:
		pItem.Title = ""
		pItem.Text = ""
		pItem.Ext = ""
		return sql.ErrNoRows
	case nil:
		break
	default:
		return db.errorVar
	}

	if db.byteBuffer, err = json.Marshal(pItem); err != nil {
		return err
	}

	db.actionLog = "UPDATE POST FROM " + string(db.byteBuffer)

	pItem.Ext = postItem.Ext
	pItem.Text = postItem.Text

	if db.byteBuffer, err = json.Marshal(pItem); err != nil {
		return err
	}

	db.actionLog += " TO " + string(db.byteBuffer)

	db.sqlQuery = `UPDATE posts SET content = $1, external = $2 WHERE posts_id = $3;`

	if _, err = db.pool.ExecContext(ctx, db.sqlQuery, postItem.Text, postItem.Ext, postItem.ID); err != nil {
		return err
	}

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", db.sqlQuery)

	db.sqlQuery = "INSERT INTO admlog (control_id, action) VALUES ($1,$2);"

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", db.sqlQuery)

	if _, err = db.pool.ExecContext(ctx, db.sqlQuery, db.ctrlID, db.actionLog); err != nil {
		return err
	}

	return nil
}

func (db *DBWorker) getExistingPostID(title string) (postid int) {
	db.sqlQuery = `SELECT posts_id FROM posts WHERE title = $1;`

	ctx, cancel := context.WithTimeout(db.reqCtx, 10*time.Second)
	defer cancel()

	row := db.pool.QueryRowContext(ctx, db.sqlQuery, strings.Trim(title, " "))

	ts, err := getTimestamp()

	Check(err, false)

	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - QUERY - %v\n", db.sqlQuery)

	switch db.errorVar = row.Scan(&postid); db.errorVar {
	case sql.ErrNoRows:
		println(ts.Format(timeFormat) + " - ERROR - MESSAGE NO ROW (POST) '" + title + "' YET")
		postid = 0
	case nil:
		return
	default:
		Check(db.errorVar, false)
		return
	}

	return
}

func (db *DBWorker) getInsertQuery(listlen int, columnlen int) {

	db.sqlQuery += "("

	for i := 1; i <= listlen*columnlen; i++ {
		db.sqlQuery += "$" + strconv.Itoa(i) + ","
		if i%columnlen == 0 {
			db.sqlQuery = db.sqlQuery[0 : len(db.sqlQuery)-1]
			db.sqlQuery += "),("
		}
	}

	db.sqlQuery = db.sqlQuery[0:len(db.sqlQuery)-2] + ";"
}

func (db *DBWorker) getConnectionInfo() {
	db.stats = db.pool.Stats()
	ts, err := getTimestamp()

	Check(err, false)
	fmt.Printf(ts.Format(timeFormat)+" - DEBUG - STATS - CONNECTIONS IDLE %v INUSE %v OPEN %v WAITING %v\n", strconv.Itoa(db.stats.Idle), strconv.Itoa(db.stats.InUse), strconv.Itoa(db.stats.OpenConnections), strconv.Itoa(int(db.stats.WaitCount)))
}

// getDSN collects the connection parameters
func (db *DBWorker) getDSN() string {
	return fmt.Sprintf("port=%d host=%s user=%s password=%s dbname=%s sslmode=%s", dbPort, dbHost, dbUsername, dbPassword, dbDBname, dbParamSSLMode)
}

// credits to https://github.com/a8m/reflect-examples
// ptr wraps the given value with pointer: V => *V, *V => **V, etc.
func ptr(v reflect.Value) reflect.Value {
	pt := reflect.PtrTo(v.Type()) // create a *T type.
	pv := reflect.New(pt.Elem())  // create a reflect.Value of type *T.
	pv.Elem().Set(v)              // sets pv to point to underlying value of v.
	return pv
}

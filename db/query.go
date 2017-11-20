package db

import (
	"database/sql"
	"errors"
	"log"
)

/*
	These functions provide a marginally cleaner syntax for SQL queries inside the made FUSE API calls.
	I've written them so the API calls don't get bogged down with SQL error handling and instead can focus on the
	filesystem aspect of the project.
*/

type query struct {
	Row    *sql.Row
	Rows   *sql.Rows
	result sql.Result
	Err    error
	Sql    string
	Inode  int64
}

/////////////
// Execute //
/////////////

func Execute(sql string, params ...interface{}) (sql.Result, error) { return db.Exec(sql, params) }

func Insert(sql string, params ...interface{}) *query {
	q := new(query)
	q.Sql = sql

	q.result, q.Err = db.Exec(sql, params...)
	if q.Err != nil {
		log.Printf("Error executing SQL: %s (%s)", q.Err, q.Sql)
		return q
	}

	q.Inode, q.Err = q.result.LastInsertId()

	return q
}

////////////////
// Single row //
////////////////

// QueryRec is an embeddable SQL query
func QueryRec(sql string, params ...interface{}) *query {
	q := new(query)
	q.Sql = sql

	q.Row = db.QueryRow(sql, params...)

	if q.Row == nil {
		log.Println("SQL returned nothing:", sql, params)
	}

	return q
}

// ScanRec returns the results for a single row SQL query
func ScanRec(q *query, params ...interface{}) error {
	if q.Row == nil {
		return errors.New("no data returned")
	}

	err := q.Row.Scan(params...)
	if err != nil {
		log.Printf("Error single-row scanning: %s (%s)", err, q.Sql)
		return err
	}

	return nil
}

///////////////
// Multi-row //
///////////////

// QueryRows is for multi-row SQL queries
func QueryRows(sql string, params ...interface{}) *query {
	q := new(query)
	q.Sql = sql

	q.Rows, q.Err = db.Query(sql, params...)
	if q.Err != nil {
		log.Println("SQL returned an error:", q.Err, sql, params)
	}

	return q
}

// ScanRows returns the results of a multi-row SQL query via the callback function
func ScanRows(q *query, callback func(), params ...interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	for q.Rows.Next() {
		err := q.Rows.Scan(params...)
		if err != nil {
			log.Printf("Error multi-row scanning: %s (%s)", err, q.Sql)
			return err
		}

		callback()
	}

	return nil
}

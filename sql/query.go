package sql

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
	row  *sql.Row
	rows *sql.Rows
	err  error
	sql  string
}

////////////////
// Single row //
////////////////

// QueryRec is an embeddable SQL query
func QueryRec(sql string, params ...interface{}) *query {
	q := new(query)
	q.sql = sql

	q.row = db.QueryRow(sql, params...)

	if q.row == nil {
		log.Println("SQL returned nothing:", sql, params)
	}

	return q
}

// ScanRec returns the results for a single row SQL query
func ScanRec(q *query, params ...interface{}) error {
	if q.row == nil {
		return errors.New("no data returned")
	}

	err := q.row.Scan(params...)
	if err != nil {
		log.Printf("Error single-row scanning: %s (%s)", err, q.sql)
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
	q.sql = sql

	q.rows, q.err = db.Query(sql, params...)
	if q.err != nil {
		log.Println("SQL returned an error:", q.err, sql, params)
	}

	return q
}

// ScanRows returns the results of a multi-row SQL query via the callback function
func ScanRows(q *query, callback func(), params ...interface{}) error {
	if q.err != nil {
		return q.err
	}

	for q.rows.Next() {
		err := q.rows.Scan(params...)
		if err != nil {
			log.Printf("Error multi-row scanning: %s (%s)", err, q.sql)
			return err
		}

		callback()
	}

	return nil
}

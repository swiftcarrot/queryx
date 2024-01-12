// Code generated by queryx, DO NOT EDIT.

package queryx

import (
	"database/sql"
	"regexp"
)

type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type Adapter struct {
	db DB
}

func NewAdapter(db DB) *Adapter {
	return &Adapter{db: db}
}

func (a *Adapter) Query(query string, args ...interface{}) *Rows {
	matched1, err := regexp.MatchString(`.* IN (.*?)`, query)
	if err != nil {
		return &Rows{
			rows:    nil,
			adapter: a,
			query:   query,
			args:    args,
			err:     err,
		}
	}
	matched2, err := regexp.MatchString(`.* in (.*?)`, query)
	if err != nil {
		return &Rows{
			rows:    nil,
			adapter: a,
			query:   query,
			args:    args,
			err:     err,
		}
	}
	if matched1 || matched2 {
		query, args, err = In(query, args...)
		if err != nil {
			return &Rows{
				rows:    nil,
				adapter: a,
				query:   query,
				args:    args,
				err:     err,
			}
		}
	}
	query, args = rebind(query, args)
	rows, err := a.db.Query(query, args...)
	if err != nil {
		return &Rows{
			rows:    rows,
			adapter: a,
			query:   query,
			args:    args,
			err:     err,
		}
	}
	return &Rows{
		rows:    rows,
		adapter: a,
		query:   query,
		args:    args,
		err:     nil,
	}
}

type Rows struct {
	rows    *sql.Rows
	adapter *Adapter
	query   string
	args    []interface{}
	err     error
}

func (r *Rows) Scan(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	err := ScanSlice(r.rows, v)
	if err != nil {
		return err
	}
	return err
}

type Row struct {
	rows    *sql.Rows
	adapter *Adapter
	query   string
	args    []interface{}
	err     error
}

func (r *Row) Scan(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	err := ScanOne(r.rows, v)
	if err != nil {
		return err
	}
	return err
}

func (a *Adapter) QueryOne(query string, args ...interface{}) *Row {
	matched1, err := regexp.MatchString(`.* IN (.*?)`, query)
	if err != nil {
		return &Row{
			rows:    nil,
			adapter: a,
			query:   query,
			args:    args,
			err:     err,
		}
	}
	matched2, err := regexp.MatchString(`.* in (.*?)`, query)
	if err != nil {
		return &Row{
			rows:    nil,
			adapter: a,
			query:   query,
			args:    args,
			err:     err,
		}
	}
	if matched1 || matched2 {
		query, args, err = In(query, args...)
		if err != nil {
			return &Row{
				rows:    nil,
				adapter: a,
				query:   query,
				args:    args,
				err:     err,
			}
		}
	}
	query, args = rebind(query, args)
	rows, err := a.db.Query(query, args...)
	if err != nil {
		return &Row{
			rows:    rows,
			adapter: a,
			query:   query,
			args:    args,
			err:     err,
		}
	}

	return &Row{
		rows:    rows,
		adapter: a,
		query:   query,
		args:    args,
		err:     err,
	}
}

// Code generated by queryx, DO NOT EDIT.

package queryx

import (
	"regexp"

	_ "github.com/mattn/go-sqlite3"
)

func (a *Adapter) Exec(query string, args ...interface{}) (int64, error) {
	matched1, err := regexp.MatchString(`.* IN (.*?)`, query)
	if err != nil {
		return 0, err
	}
	matched2, err := regexp.MatchString(`.* in (.*?)`, query)
	if err != nil {
		return 0, err
	}
	if matched1 || matched2 {
		query, args, err = In(query, args...)
		if err != nil {
			return 0, err
		}
	}
	query, args = rebind(query, args)
	result, err := a.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func rebind(query string, args []interface{}) (string, []interface{}) {
	return query, args
}

package queryx

import (
	"database/sql"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

func execUtils(query string, args ...interface{}) (string, []interface{}, error) {
	matched1, err := regexp.MatchString(`.* IN (.*?)`, query)
	if err != nil {
		return "", nil, err
	}
	matched2, err := regexp.MatchString(`.* in (.*?)`, query)
	if err != nil {
		return "", nil, err
	}
	if matched1 || matched2 {
		query, args, err = In(query, args...)
		if err != nil {
			return "", nil, err
		}
	}
	return query, args, err
}

func (a *Adapter) Exec(query string, args ...interface{}) (int64, error) {
	query, args, err := execUtils(query, args...)
	if err != nil {
		return 0, err
	}

	result, err := a.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (a *Adapter) ExecInternal(query string, args ...interface{}) (sql.Result, error) {
	query, args, err := execUtils(query, args...)
	if err != nil {
		return nil, err
	}
	result, err := a.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return result, err
}

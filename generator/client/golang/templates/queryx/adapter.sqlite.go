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
	result, err := a.db.Exec(rebind(DOLLAR, query), args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

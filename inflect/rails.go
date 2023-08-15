package inflect

import (
	"strings"
)

func RailsAdapter(a string) string {
	switch a {
	case "mysql":
		return "mysql2"
	case "sqlite":
		return "sqlite3"
	}
	return ""
}

func RailsURL(a string, u string) string {
	switch a {
	case "mysql":
		return strings.Replace(u, "mysql://", "mysql2://", 1)
	}
	return u

}

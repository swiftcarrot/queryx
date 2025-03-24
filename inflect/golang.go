package inflect

import (
	"fmt"
	"log"
	"strings"
)

// type abbreviation used for a method's receiver variable
// https://github.com/golang/go/wiki/CodeReviewComments#receiver-names
func GoReceiver(typ string) string {
	return strings.ToLower(typ[0:1])
}

// avoid go keyword with syntax error
func goKeywordFix(s string) string {
	switch s {
	case "type":
		return "typ"
	default:
		return s
	}
}

// TODO: use a map
// convert column type to corresponding queryx type
func goModelType(t string, null bool, array bool) string {
	if null {
		switch t {
		case "bigint":
			return "queryx.BigInt"
		case "uuid":
			return "queryx.UUID"
		case "string", "text":
			if array {
				return "queryx.StringArray"
			}
			return "queryx.String"
		case "datetime":
			return "queryx.Datetime"
		case "date":
			return "queryx.Date"
		case "time":
			return "queryx.Time"
		case "integer":
			return "queryx.Integer"
		case "boolean":
			return "queryx.Boolean"
		case "float":
			return "queryx.Float"
		case "json", "jsonb":
			return "queryx.JSON"
		default:
			log.Fatal(fmt.Errorf("unhandled data type %s in goModelType", t))
			return ""
		}
	} else {
		switch t {
		case "bigint":
			return "int64"
		case "uuid":
			return "string"
		case "string", "text":
			return "string"
		case "datetime":
			return "queryx.Datetime"
		case "date":
			return "queryx.Date"
		case "time":
			return "queryx.Time"
		case "integer":
			return "int32"
		case "boolean":
			return "bool"
		case "float":
			return "float"
		case "json", "jsonb":
			return "queryx.JSON"
		default:
			log.Fatal(fmt.Errorf("unhandled data type %s in goModelType", t))
			return ""
		}
	}
}

// convert column type to corresponding queryx type
func goType(t string, null bool, array bool) string {
	switch t {
	case "bigint":
		return "BigInt"
	case "uuid":
		return "UUID"
	case "string", "text":
		if array {
			return "StringArray"
		}
		return "String"
	case "datetime":
		return "Datetime"
	case "date":
		return "Date"
	case "time":
		return "Time"
	case "integer":
		return "Integer"
	case "boolean":
		return "Boolean"
	case "float":
		return "Float"
	case "json", "jsonb":
		return "JSON"
	default:
		log.Fatal(fmt.Errorf("unhandled data type %s in goType", t))
		return ""
	}
}

// convert column type to go type in setter method of change object
func goChangeSetType(t string, null bool, array bool) string {
	switch t {
	case "bigint":
		return "int64"
	case "boolean":
		return "bool"
	case "integer":
		return "int32"
	case "string", "text", "date", "time", "datetime", "uuid":
		if array {
			return "[]string"
		}
		return "string"
	case "float":
		return "float64"
	case "json", "jsonb":
		return "map[string]interface{}"
	default:
		log.Fatal(fmt.Errorf("unhandled data type %s in goChangeSetType", t))
		return ""
	}
}

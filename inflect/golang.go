package inflect

import (
	"log"
)

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
func goModelType(t string, null bool) string {
	if null {
		switch t {
		case "bigint":
			return "queryx.BigInt"
		case "uuid":
			return "queryx.UUID"
		case "string", "text":
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
			log.Fatal("not found", t)
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
			log.Fatal("not found", t)
			return ""
		}
	}
}

// convert column type to corresponding queryx type
func goType(t string) string {
	switch t {
	case "bigint":
		return "BigInt"
	case "uuid":
		return "UUID"
	case "string", "text":
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
		return ""
	}
}

// convert column type to go type in setter method of change object
func goChangeSetType(t string) string {
	switch t {
	case "bigint":
		return "int64"
	case "boolean":
		return "bool"
	case "integer":
		return "int32"
	case "string":
		return "string"
	case "text":
		return "string"
	case "datetime":
		return "string"
	case "date":
		return "string"
	case "time":
		return "string"
	case "uuid":
		return "string"
	case "float":
		return "float64"
	case "json", "jsonb":
		return "map[string]interface{}"
	default:
		log.Fatal("unknown type in goChangeSetType", t)
		return ""
	}
}

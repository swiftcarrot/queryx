package inflect

import (
	"fmt"
	"log"
)

func tsType(t string, null bool, array bool) string {
	switch t {
	case "uuid", "string", "text":
		if array {
			return "string[]"
		}
		return "string"
	case "datetime", "date":
		return "Date"
	case "time":
		return "string"
	case "bigint", "integer", "float":
		return "number"
	case "boolean":
		return "boolean"
	case "json", "jsonb":
		return "{ [key: string]: any }"
	default:
		log.Fatal(fmt.Errorf("unhandled data type %s in tsType", t))
		return ""
	}
}

func tsChangeSetType(t string, null bool, array bool) string {
	switch t {
	case "bigint", "integer", "float":
		return "number"
	case "uuid":
		return "string"
	case "string", "text":
		if array {
			return "string[]"
		}
		return "string"
	case "datetime", "date", "time":
		return "string"
	case "boolean":
		return "boolean"
	case "json", "jsonb":
		return "object"
	default:
		log.Fatal(fmt.Errorf("unhandled data type %s in tsChangeSetType", t))
		return ""
	}
}

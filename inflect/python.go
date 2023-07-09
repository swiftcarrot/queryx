package inflect

import "log"

func pyModelType(typ string, null bool) string {
	switch typ {
	case "bigint", "integer":
		return "int"
	case "uuid":
		return "str"
	case "string", "text":
		return "str"
	case "datetime", "date", "time":
		return "datetime"
	case "boolean":
		return "bool"
	case "float":
		return "float"
	case "json", "jsonb":
		return "Dict"
	default:
		log.Fatal("not found", typ)
		return ""
	}
}

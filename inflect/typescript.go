package inflect

func tsType(t string) string {
	switch t {
	case "bigint":
		return "number"
	case "uuid":
		return "string"
	case "string", "text":
		return "string"
	case "datetime":
		return "Date"
	case "date":
		return "Date"
	case "time":
		return "Date"
	case "integer":
		return "number"
	case "boolean":
		return "boolean"
	case "float":
		return "number"
	case "json", "jsonb":
		return "object"
	default:
		return "" // TODO: raise error
	}
}

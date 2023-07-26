package inflect

func tsType(t string) string {
	switch t {
	case "uuid", "string", "text":
		return "string"
	case "datetime", "date", "time":
		return "Date"
	case "bigint", "integer", "float":
		return "number"
	case "boolean":
		return "boolean"
	case "json", "jsonb":
		return "object"
	default:
		return ""
	}
}

func tsChangeSetType(t string) string {
	switch t {
	case "bigint", "integer", "float":
		return "number"
	case "uuid":
		return "string"
	case "string", "text":
		return "string"
	case "datetime", "date", "time":
		return "string"
	case "boolean":
		return "boolean"
	case "json", "jsonb":
		return "object"
	default:
		return ""
	}
}

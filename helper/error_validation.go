package helper

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func MessageForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "min":
		return "This field must be at least [PARAM] characters"
	case "max":
		return "This field must be at most [PARAM] characters"
	case "alphanum":
		return "This field must be alphanumeric"
	case "containsany":
		return "This field must contain at least one special character, one uppercase letter, one lowercase letter, and one number"
	case "alpha":
		return "This field must be alphabetic"
	case "uppercase":
		return "This field must contain at least one uppercase character"
	case "lowercase":
		return "This field must contain at least one lowercase character"
	case "alphanumunicode":
		return "This field must be alphanumeric and unicode"
	case "eqfield":
		return "This field must be equal to [PARAM]"
	case "len":
		return "This field must be [PARAM] characters"
	case "gte":
		return "This field must be greater than or equal to [PARAM]"
	default:
		return "Invalid field " + tag
	}
}

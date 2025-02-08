package util

func NonEmpty(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

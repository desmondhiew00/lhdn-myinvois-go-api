package constant

import "strings"

func normalize(input string) string {
	s := strings.TrimSpace(input)
	s = strings.ToLower(s)
	return strings.Join(strings.Fields(s), "")
}

func valid(input string, constantMap map[string]string) bool {
	key := normalize(input)
	var codes = make(map[string]struct{})
	for _, code := range constantMap {
		codes[code] = struct{}{}
	}
	_, ok := codes[key]
	return ok
}

func Value(input string, constantMap map[string]string, defaultValue string) string {
	key := normalize(input)
	if code, ok := constantMap[key]; ok {
		return code
	}

	if valid(input, constantMap) {
		return input
	}

	return defaultValue
}

package wrapper

func Wrapper(fields []string, prefix string) string {
	result := ""
	for _, field := range fields {
		result = result + prefix + "." + field + ", "
	}
	return result[:len(result) - 2]
}

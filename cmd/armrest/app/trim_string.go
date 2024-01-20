package app

func trimString(input string, length int, postfix string) string {
	switch {
	case length <= 0:
		return ""
	case len(input) > length:
		return input[:length] + postfix
	default:
		return input
	}
}

package new

// isSpace checks if the given byte represents a whitespace character
func isSpace(c byte) bool {
	return c == ' ' || c == '\t'
}

// isLetter checks if the given byte represents a lower- or uppercase letter
func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// isDigit checks if the given byte represents a digit
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// isIdentStart checks if the given byte represents a character that can start
// and identifier
func isIdentStart(c byte) bool {
	return isLetter(c) || c == '_'
}

// isIdenCont checks if the given byte represents a character than continues
// an identifier
func isIdentCont(c byte) bool {
	return isIdentStart(c) || isDigit(c) || c == '-' || c == '.'
}

func isBracket(c byte) bool {
	return isOpenBracket(c) || isCloseBracket(c)
}

func isOpenBracket(c byte) bool {
	return c == '(' || c == '[' || c == '{' || c == '<'
}

func isCloseBracket(c byte) bool {
	return c == ')' || c == ']' || c == '}' || c == '>'
}

// toString turn the given byte into a 1-long string
func toString(c byte) string {
	return string(rune(c))
}

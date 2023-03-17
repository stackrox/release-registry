// Package validate contains utility functions for input validation.
package validate

import "regexp"

// IsValidString does a regex comparison between a passed pattern and a word.
func IsValidString(pattern string, word string) bool {
	match, _ := regexp.MatchString(pattern, word)

	return match
}

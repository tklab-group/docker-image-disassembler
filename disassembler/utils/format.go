package utils

import "strings"

// CleanArgs trims the whitespace from the given set of strings.
func CleanArgs(args []string) []string {
	var cleanedArgs []string
	for _, str := range args {
		if str != "" {
			cleanedArgs = append(cleanedArgs, strings.TrimSpace(str))
		}
	}
	return cleanedArgs
}

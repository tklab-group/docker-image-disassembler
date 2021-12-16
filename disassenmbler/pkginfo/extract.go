package pkginfo

import "regexp"

func extractNamedGroup(reg *regexp.Regexp, s string) map[string]string {
	match := reg.FindStringSubmatch(s)
	result := map[string]string{}
	for i, name := range reg.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return result
}

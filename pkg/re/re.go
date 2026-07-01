package re

import (
	"regexp"
	"strings"
)

// Precompiled regexps for matching domains and IPv4/IPv6 addresses in text.
var (
	DomainRe = regexp.MustCompile(`([a-zA-Z0-9][-a-zA-Z0-9]{0,62}\.)+([a-zA-Z][-a-zA-Z]{0,62})`)

	IPv4Re = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	IPv6Re = regexp.MustCompile(`fe80:(:[0-9a-fA-F]{1,4}){0,4}(%\w+)?|([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|64:ff9b::(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}|::[fF]{4}:(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}|(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?::(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?`)
)

// MaybeRegexp reports whether s contains characters suggesting it is a regexp pattern.
func MaybeRegexp(s string) bool {
	return strings.ContainsAny(s, "[]{}()?")
}

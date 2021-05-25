package utils

import (
	"fmt"
	"regexp"
)

func Match(pattern string, s string) error {
	reg := regexp.MustCompile(pattern)
	if reg.MatchString(s) {
		return nil
	} else {
		return fmt.Errorf("%s unmatched %s", s, pattern)
	}
}

package utils

import (
	"regexp"

	"github.com/isnlan/coral/pkg/errors"
)

func Match(pattern string, s string) error {
	reg := regexp.MustCompile(pattern)
	if reg.MatchString(s) {
		return nil
	} else {
		return errors.Errorf("%s unmatched %s", s, pattern)
	}
}

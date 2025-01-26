package main_test

import (
	"regexp"
	"testing"
	"time"
)

func TestMain01(t *testing.T) {
	t.Run(
		"RFC3339 time formatting and regexp tests",
		func(t *testing.T) {
			rfc3339FormatString := time.RFC3339
			rfc3339RegexString := `^(?P<yyyy>\d{4})-(?P<mm>\d{2})-(?P<dd>\d{2})T(?P<HH>\d{2}):(?P<MM>\d{2}):(?P<SS>\d{2})Z$`
			rfc3339Regex := regexp.MustCompile(rfc3339RegexString)
			n := time.Now().UTC().Format(rfc3339FormatString)
			t.Logf("time.RFC3339 format: '%s'\n", time.RFC3339)
			t.Logf("time.RFC3339 actual: '%s'\n", n)
			if !rfc3339Regex.MatchString(n) {
				t.Errorf("time is not formatted in expected 'YYYY-mm-ddTHH:MM:SSZ' format: '%v'", n)
			}
		},
	)
}

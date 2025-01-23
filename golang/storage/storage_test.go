package storage_test

import (
	"regexp"
	"testing"
	"time"
)

func TestStorageDatetimeChecks(t *testing.T) {
	t.Run(
		"RFC3339 regexp Test",
		func(t *testing.T) {
			formatString := `2006-01-02T15-04-05.000000000Z0700`
			regexString := `^(?P<yyyy>\d{4})-(?P<mm>\d{2})-(?P<dd>\d{2})T(?P<HH>\d{2})-(?P<MM>\d{2})-(?P<SS>\d{2})\.(?P<nnnnnnnnn>\d{9})Z`
			re := regexp.MustCompile(regexString)

			n := time.Now().UTC().Format(formatString)
			t.Logf("time format: '%s'\n", formatString)
			t.Logf("time actual: '%s'\n", n)

			if !re.MatchString(n) {
				t.Errorf("id is not in 'YYYY-mm-ddTHH:MM:SSZ' format: '%v'", n)
			}
		},
	)
}

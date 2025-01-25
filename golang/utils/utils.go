package utils

import (
	"regexp"
)

var DateTimeFormatString = `2006-01-02T15-04-05.000000000Z0700`
var DateTimeRegexpString = `(?P<yyyy>\d{4})-(?P<mm>\d{2})-(?P<dd>\d{2})T(?P<HH>\d{2})-(?P<MM>\d{2})-(?P<SS>\d{2})\.(?P<nnnnnnnnn>\d{9})Z`

func LASFSElectionIDFromString(s string) string {
	m := regexp.MustCompile("[^a-zA-Z0-9]")
	s = m.ReplaceAllString(s, "_")
	return s
}

func LASFSMemberIDFromString(s string) string {
	m := regexp.MustCompile("[^a-zA-Z0-9]")
	s = m.ReplaceAllString(s, "")
	return s
}

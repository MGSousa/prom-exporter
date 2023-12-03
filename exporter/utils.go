package exporter

import "strings"

func split(s, delim string) []string {
	return strings.Split(s, delim)
}

func join(s []string, delim string) string {
	return strings.Join(s, delim)
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

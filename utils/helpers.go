// utils/helpers.go
package utils

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func UppercaseFirst(s string) string {
	c := cases.Title(language.Und, cases.NoLower)
	return c.String(s)
}

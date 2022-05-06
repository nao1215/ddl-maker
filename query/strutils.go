package query

import "fmt"

// Quote encloses the string with ''.
func Quote(s string) string {
	return fmt.Sprintf("`%s`", s)
}

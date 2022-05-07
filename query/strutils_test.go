package query

import "testing"

func TestQuote(t *testing.T) {
	column := "id"

	want := "`id`"
	got := Quote(column)
	if want != got {
		t.Errorf("mismatch want=%s, got=%s", want, got)
	}
}

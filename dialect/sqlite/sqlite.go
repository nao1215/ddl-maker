package sqlite

import "github.com/nao1215/ddl-maker/query"

const (
	autoIncrement = "AUTOINCREMENT"
)

// SQLite is a model with database engine and character code for SQLite
type SQLite struct{}

// HeaderTemplate return string that is sql header template
func (sqlite SQLite) HeaderTemplate() string {
	return `PRAGMA foreign_keys = false;
`
}

// FooterTemplate return string that is sql footer template
func (sqlite SQLite) FooterTemplate() string {
	return `PRAGMA foreign_keys = true;
`
}

// TableTemplate return string that is sql table template
func (sqlite SQLite) TableTemplate() string {
	return `
DROP TABLE IF EXISTS {{ .Name }};

CREATE TABLE {{ .Name }} (
    {{ range .Columns -}}
        {{ .ToSQL }},
    {{ end -}}
    {{ range .ForeignKeys.Sort  -}}
        {{ .ToSQL }},
    {{ end -}}
);

{{ range .Indexes.Sort -}}
    {{ .ToSQL }},
{{ end -}}

`
}

// ToSQL convert sqlite sql string from typeName and size
func (sqlite SQLite) ToSQL(typeName string, size uint64) (string, error) {
	// TODO:
	return "", nil
}

// Quote return string that encloses with ``.
func (sqlite SQLite) Quote(s string) string {
	return query.Quote(s)
}

// AutoIncrement return string for auto-increment setting
func (sqlite SQLite) AutoIncrement() string {
	return autoIncrement
}

package sqlite

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nao1215/ddl-maker/query"
)

// ErrInvalidType means Invalid type specified when parsing
var ErrInvalidType = errors.New("Specified type is invalid")

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

// TableTemplate return string that is sql table template.
//
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
    {{ .PrimaryKey.ToSQL }}
);

{{ range .Indexes.Sort -}}
    {{ .ToSQL }},
{{ end -}}

`
}

// ToSQL convert sqlite sql string from typeName and size
func (sqlite SQLite) ToSQL(typeName string, size uint64) (string, error) {
	switch typeName {
	case "int8", "*int8":
		return "INTEGER", nil
	case "int16", "*int16":
		return "INTEGER", nil
	case "int32", "*int32", "sql.NullInt32":
		return "INTEGER", nil
	case "int64", "*int64", "sql.NullInt64":
		return "INTEGER", nil
	case "uint8", "*uint8":
		return "INTEGER", nil
	case "uint16", "*uint16":
		return "INTEGER", nil
	case "uint32", "*uint32":
		return "INTEGER", nil
	case "uint64", "*uint64":
		return "INTEGER", nil
	case "float32", "*float32":
		return "REAL", nil
	case "float64", "*float64", "sql.NullFloat64":
		return "REAL", nil
	case "string", "*string", "sql.NullString":
		return "TEXT", nil
	case "[]uint8", "sql.RawBytes":
		return "BLOB", nil
	case "bool", "*bool", "sql.NullBool":
		return "INTEGER", nil
	case "tinytext":
		return "TEXT", nil
	case "text":
		return "TEXT", nil
	case "mediumtext":
		return "TEXT", nil
	case "longtext":
		return "TEXT", nil
	case "tinyblob":
		return "BLOB", nil
	case "blob":
		return "BLOB", nil
	case "mediumblob":
		return "BLOB", nil
	case "longblob":
		return "BLOB", nil
	case "time":
		return "INTEGER", nil
	case "time.Time", "*time.Time":
		return "INTEGER", nil
	case "sql.NullTime":
		return "INTEGER", nil
	case "date":
		return "INTEGER", nil
	case "json.RawMessage", "*json.RawMessage":
		return "JSON", nil // from SQLite3.9
	default:
		return "", fmt.Errorf("%w: %s", ErrInvalidType, typeName)
	}
}

// Quote return string that encloses with ``.
func (sqlite SQLite) Quote(s string) string {
	return query.Quote(s)
}

// AutoIncrement return string for auto-increment setting
func (sqlite SQLite) AutoIncrement() string {
	return autoIncrement
}

// PrimaryKey is a model for determining the primary key
type PrimaryKey struct {
	columns []string
}

// AddPrimaryKey return initialized PrimaryKey struct.
func AddPrimaryKey(columns ...string) PrimaryKey {
	return PrimaryKey{
		columns: columns,
	}
}

// Columns returns the columns that will be the primary keys.
func (pk PrimaryKey) Columns() []string {
	return pk.columns
}

// ToSQL return primary key sql string.
func (pk PrimaryKey) ToSQL() string {
	var columnsStr []string
	for _, c := range pk.columns {
		columnsStr = append(columnsStr, query.Quote(c))
	}
	return fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(columnsStr, ", "))
}

// Index is model representing indexes to speed up DB searches
type Index struct {
	columns []string
	table   string
	name    string
}

// Name return index name
func (i Index) Name() string {
	return i.name
}

// Columns return index columns
func (i Index) Columns() []string {
	return i.columns
}

// ToSQL return index sql string
func (i Index) ToSQL() string {
	var columnsStr []string

	for _, c := range i.columns {
		columnsStr = append(columnsStr, query.Quote(c))
	}

	return fmt.Sprintf("CREATE INDEX %s ON %s (%s);",
		query.Quote(i.name), query.Quote(i.table), strings.Join(columnsStr, ", "))
}

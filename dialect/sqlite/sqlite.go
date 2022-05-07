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

// AddIndex returns a new Index
func AddIndex(idxName, table string, columns ...string) Index {
	return Index{
		name:    idxName,
		table:   table,
		columns: columns,
	}
}

// Name return index name
func (i Index) Name() string {
	return query.Quote(i.name)
}

// Table return table name
func (i Index) Table() string {
	return query.Quote(i.table)
}

// Columns return index columns
func (i Index) Columns() []string {
	var columnsStr []string
	for _, c := range i.columns {
		columnsStr = append(columnsStr, query.Quote(c))
	}
	return columnsStr
}

// ToSQL return index sql string
func (i Index) ToSQL() string {
	return fmt.Sprintf("CREATE INDEX %s ON %s (%s);",
		i.Name(), i.Table(), strings.Join(i.Columns(), " "))
}

// UniqueIndex is model that represents unique constraints
type UniqueIndex struct {
	columns []string
	table   string
	name    string
}

// AddUniqueIndex returns a new UniqueIndex
func AddUniqueIndex(idxName, table string, columns ...string) UniqueIndex {
	return UniqueIndex{
		name:    idxName,
		table:   table,
		columns: columns,
	}
}

// Name return unique index name
func (ui UniqueIndex) Name() string {
	return query.Quote(ui.name)
}

// Table return table name
func (ui UniqueIndex) Table() string {
	return query.Quote(ui.table)
}

// Columns return unique index columns
func (ui UniqueIndex) Columns() []string {
	var columnsStr []string
	for _, c := range ui.columns {
		columnsStr = append(columnsStr, query.Quote(c))
	}
	return columnsStr
}

// ToSQL return unique unique index sql string
func (ui UniqueIndex) ToSQL() string {
	return fmt.Sprintf("CREATE UNIQUE INDEX %s ON %s (%s);",
		ui.Name(), ui.Table(), strings.Join(ui.Columns(), " "))
}

// ForeignKey is a model for setting foreign key constraints
type ForeignKey struct {
	foreignColumns     []string
	referenceTableName string
	referenceColumns   []string
	updateOption       string
	deleteOption       string
}

// ForeignKeyOptionType is string that means foreign key otion
// https://www.sqlite.org/foreignkeys.html
type ForeignKeyOptionType string

// ForeignKeyOptionCascade CASCADE
var ForeignKeyOptionCascade ForeignKeyOptionType = "CASCADE"

// ForeignKeyOptionSetNull SET NULL
var ForeignKeyOptionSetNull ForeignKeyOptionType = "SET NULL"

// ForeignKeyOptionRestrict RESTRICT
var ForeignKeyOptionRestrict ForeignKeyOptionType = "RESTRICT"

// ForeignKeyOptionNoAction NO ACTION
var ForeignKeyOptionNoAction ForeignKeyOptionType = "NO ACTION"

// ForeignKeyOptionSetDefault SET DEFAULT
var ForeignKeyOptionSetDefault ForeignKeyOptionType = "SET DEFAULT"

// String Stringer for ForeignKeyOptionType
func (fkopt ForeignKeyOptionType) String() string {
	return string(fkopt)
}

// ForeignKeyOption is an interface for controlling foreign key constraint options.
type ForeignKeyOption interface {
	Apply(*ForeignKey)
}

type withUpdateForeignKeyOption string

// Apply apply foreign key constraint options for Update.
func (o withUpdateForeignKeyOption) Apply(f *ForeignKey) {
	f.updateOption = string(o)
}

// WithUpdateForeignKeyOption manages the foreign key constraint options for Update.
func WithUpdateForeignKeyOption(option ForeignKeyOptionType) ForeignKeyOption {
	switch option {
	// Specifying RESTRICT (or NO ACTION) is the same as omitting the ON DELETE or ON UPDATE clause.
	case ForeignKeyOptionRestrict, ForeignKeyOptionNoAction:
		return withUpdateForeignKeyOption("")
	}
	return withUpdateForeignKeyOption(option)
}

type withDeleteForeignKeyOption string

// Apply apply foreign key constraint options for Delete.
func (o withDeleteForeignKeyOption) Apply(f *ForeignKey) {
	f.deleteOption = string(o)
}

// WithDeleteForeignKeyOption return query that is the foreign key constraint options for Delete.
func WithDeleteForeignKeyOption(option ForeignKeyOptionType) ForeignKeyOption {
	switch option {
	// Specifying RESTRICT (or NO ACTION) is the same as omitting the ON DELETE or ON UPDATE clause.
	case ForeignKeyOptionRestrict, ForeignKeyOptionNoAction:
		return withDeleteForeignKeyOption("")
	}
	return withDeleteForeignKeyOption(option)
}

// AddForeignKey returns a new ForeignKey
func AddForeignKey(foreignColumns, referenceColumns []string, referenceTableName string, option ...ForeignKeyOption) ForeignKey {
	foreignKey := ForeignKey{
		foreignColumns:     foreignColumns,
		referenceTableName: referenceTableName,
		referenceColumns:   referenceColumns,
	}

	for _, o := range option {
		if o != nil {
			o.Apply(&foreignKey)
		}
	}
	return foreignKey
}

// ForeignColumns return slice of foreign key columns
func (fk ForeignKey) ForeignColumns() []string {
	var foreignColumnsStr []string
	for _, fc := range fk.foreignColumns {
		foreignColumnsStr = append(foreignColumnsStr, query.Quote(fc))
	}
	return foreignColumnsStr
}

// ReferenceTableName return reference table name
func (fk ForeignKey) ReferenceTableName() string {
	return query.Quote(fk.referenceTableName)
}

// ReferenceColumns return slice of return foreign key columns
func (fk ForeignKey) ReferenceColumns() []string {
	var referenceColumnsStr []string
	for _, rc := range fk.referenceColumns {
		referenceColumnsStr = append(referenceColumnsStr, query.Quote(rc))
	}
	return referenceColumnsStr
}

// UpdateOption return foreign key constraint option string for update
func (fk ForeignKey) UpdateOption() string {
	return fk.updateOption
}

// DeleteOption return foreign key constraint option string for delete
func (fk ForeignKey) DeleteOption() string {
	return fk.deleteOption
}

// ToSQL return foreign key sql string
func (fk ForeignKey) ToSQL() string {
	sql := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s (%s)",
		strings.Join(fk.ForeignColumns(), ", "),
		fk.ReferenceTableName(),
		strings.Join(fk.ReferenceColumns(), ", "))
	if fk.DeleteOption() != "" {
		sql = sql + fmt.Sprintf(" ON DELETE %s", fk.DeleteOption())
	}
	if fk.UpdateOption() != "" {
		sql = sql + fmt.Sprintf(" ON UPDATE %s", fk.UpdateOption())
	}
	return sql
}

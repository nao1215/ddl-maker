package mysql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nao1215/ddl-maker/query"
)

const (
	defaultVarcharSize   = 191
	defaultVarbinarySize = 767
	autoIncrement        = "AUTO_INCREMENT"
)

// ErrInvalidType means Invalid type specified when parsing
var ErrInvalidType = errors.New("Specified type is invalid")

// MySQL is a model with database engine and character code for MySQL
type MySQL struct {
	Engine  string
	Charset string
}

// Index XXX
type Index struct {
	columns []string
	name    string
}

// UniqueIndex is model that represents unique constraints
type UniqueIndex struct {
	columns []string
	name    string
}

// FullTextIndex XXX
type FullTextIndex struct {
	columns []string
	name    string
	parser  string
}

// SpatialIndex XXX
type SpatialIndex struct {
	columns []string
	name    string
}

// PrimaryKey is a model for determining the primary key
type PrimaryKey struct {
	columns []string
}

// ForeignKeyOptionType XXX
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

func (fkopt ForeignKeyOptionType) String() string {
	return string(fkopt)
}

// ForeignKey XXX
type ForeignKey struct {
	foreignColumns     []string
	referenceTableName string
	referenceColumns   []string
	updateOption       string
	deleteOption       string
}

// ForeignKeyOption XXX
type ForeignKeyOption interface {
	Apply(*ForeignKey)
}

type withUpdateForeignKeyOption string

func (o withUpdateForeignKeyOption) Apply(f *ForeignKey) {
	f.updateOption = string(o)
}

// WithUpdateForeignKeyOption XXX
func WithUpdateForeignKeyOption(option ForeignKeyOptionType) ForeignKeyOption {
	switch option {
	// Specifying RESTRICT (or NO ACTION) is the same as omitting the ON DELETE or ON UPDATE clause.
	case ForeignKeyOptionRestrict, ForeignKeyOptionNoAction:
		return withUpdateForeignKeyOption("")
	}
	return withUpdateForeignKeyOption(option)
}

type withDeleteForeignKeyOption string

func (o withDeleteForeignKeyOption) Apply(f *ForeignKey) {
	f.deleteOption = string(o)
}

// WithDeleteForeignKeyOption XXX
func WithDeleteForeignKeyOption(option ForeignKeyOptionType) ForeignKeyOption {
	switch option {
	// Specifying RESTRICT (or NO ACTION) is the same as omitting the ON DELETE or ON UPDATE clause.
	case ForeignKeyOptionRestrict, ForeignKeyOptionNoAction:
		return withDeleteForeignKeyOption("")
	}
	return withDeleteForeignKeyOption(option)
}

// HeaderTemplate return string that is sql header template
func (mysql MySQL) HeaderTemplate() string {
	return `SET foreign_key_checks=0;
`
}

// FooterTemplate return string that is sql footer template
func (mysql MySQL) FooterTemplate() string {
	return `SET foreign_key_checks=1;
`
}

// TableTemplate return string that is sql table template
func (mysql MySQL) TableTemplate() string {
	return `
DROP TABLE IF EXISTS {{ .Name }};

CREATE TABLE {{ .Name }} (
    {{ range .Columns -}}
        {{ .ToSQL }},
    {{ end -}}
    {{ range .Indexes.Sort -}}
        {{ .ToSQL }},
    {{ end -}}
    {{ range .ForeignKeys.Sort  -}}
        {{ .ToSQL }},
    {{ end -}}
    {{ .PrimaryKey.ToSQL }}
) ENGINE={{ .Dialect.Engine }} DEFAULT CHARACTER SET {{ .Dialect.Charset }};

`
}

// ToSQL convert mysql sql string from typeName and size
func (mysql MySQL) ToSQL(typeName string, size uint64) (string, error) {
	switch typeName {
	case "int8", "*int8":
		return "TINYINT", nil
	case "int16", "*int16":
		return "SMALLINT", nil
	case "int32", "*int32", "sql.NullInt32": // from Go 1.13
		return "INTEGER", nil
	case "int64", "*int64", "sql.NullInt64":
		return "BIGINT", nil
	case "uint8", "*uint8":
		return "TINYINT unsigned", nil
	case "uint16", "*uint16":
		return "SMALLINT unsigned", nil
	case "uint32", "*uint32":
		return "INTEGER unsigned", nil
	case "uint64", "*uint64":
		return "BIGINT unsigned", nil
	case "float32", "*float32":
		return "FLOAT", nil
	case "float64", "*float64", "sql.NullFloat64":
		return "DOUBLE", nil
	case "string", "*string", "sql.NullString":
		return varchar(size), nil
	case "[]uint8", "sql.RawBytes":
		return varbinary(size), nil
	case "bool", "*bool", "sql.NullBool":
		return "TINYINT(1)", nil
	case "tinytext":
		return "TINYTEXT", nil
	case "text":
		return "TEXT", nil
	case "mediumtext":
		return "MEDIUMTEXT", nil
	case "longtext":
		return "LONGTEXT", nil
	case "tinyblob":
		return "TINYBLOB", nil
	case "blob":
		return "BLOB", nil
	case "mediumblob":
		return "MEDIUMBLOB", nil
	case "longblob":
		return "LONGBLOB", nil
	case "time":
		return "TIME", nil
	case "time.Time", "*time.Time":
		return datetime(size), nil
	case "mysql.NullTime": // https://godoc.org/github.com/go-sql-driver/mysql#NullTime
		return datetime(size), nil
	case "sql.NullTime": // from Go 1.13
		return datetime(size), nil
	case "date":
		return "DATE", nil
	case "json.RawMessage", "*json.RawMessage":
		return "JSON", nil
	case "geometry":
		return "GEOMETRY", nil
	default:
		return "", fmt.Errorf("%w: %s", ErrInvalidType, typeName)
	}
}

// Quote encloses the string with ``.
func (mysql MySQL) Quote(s string) string {
	return query.Quote(s)
}

// AutoIncrement return string for auto-increment setting
func (mysql MySQL) AutoIncrement() string {
	return autoIncrement
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
	for _, c := range i.Columns() {
		columnsStr = append(columnsStr, query.Quote(c))
	}
	return fmt.Sprintf("INDEX %s (%s)", query.Quote(i.Name()), strings.Join(columnsStr, ", "))
}

// Name return unique index name
func (ui UniqueIndex) Name() string {
	return ui.name
}

// Columns return unique index columns
func (ui UniqueIndex) Columns() []string {
	return ui.columns
}

// ToSQL return unique index sql string
func (ui UniqueIndex) ToSQL() string {
	var columnsStr []string
	for _, c := range ui.columns {
		columnsStr = append(columnsStr, query.Quote(c))
	}

	return fmt.Sprintf("UNIQUE %s (%s)", query.Quote(ui.name), strings.Join(columnsStr, ", "))
}

// Name return full text index name
func (fi FullTextIndex) Name() string {
	return fi.name
}

// Columns return full text index columns
func (fi FullTextIndex) Columns() []string {
	return fi.columns
}

// WithParser XXX
func (fi FullTextIndex) WithParser(s string) FullTextIndex {
	fi.parser = s
	return fi
}

// ToSQL return full text index sql string
func (fi FullTextIndex) ToSQL() string {
	var columnsStr []string
	for _, c := range fi.columns {
		columnsStr = append(columnsStr, query.Quote(c))
	}

	sql := fmt.Sprintf("FULLTEXT %s (%s)", query.Quote(fi.name), strings.Join(columnsStr, ", "))
	if fi.parser != "" {
		sql += fmt.Sprintf(" WITH PARSER %s", query.Quote(fi.parser))
	}
	return sql
}

// Name XXX
func (si SpatialIndex) Name() string {
	return si.name
}

// Columns XXX
func (si SpatialIndex) Columns() []string {
	return si.columns
}

// ToSQL return unique index sql string
func (si SpatialIndex) ToSQL() string {
	var columnsStr []string
	for _, c := range si.columns {
		columnsStr = append(columnsStr, query.Quote(c))
	}

	return fmt.Sprintf("SPATIAL KEY %s (%s)", query.Quote(si.name), strings.Join(columnsStr, ", "))
}

// Columns returns the columns that will be the primary keys.
func (pk PrimaryKey) Columns() []string {
	return pk.columns
}

// ToSQL return primary key sql string
func (pk PrimaryKey) ToSQL() string {
	var columnsStr []string
	for _, c := range pk.columns {
		columnsStr = append(columnsStr, query.Quote(c))
	}

	return fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(columnsStr, ", "))
}

// ForeignColumns XXX
func (fk ForeignKey) ForeignColumns() []string {
	return fk.foreignColumns
}

// ReferenceTableName XXX
func (fk ForeignKey) ReferenceTableName() string {
	return fk.referenceTableName
}

// ReferenceColumns XXX
func (fk ForeignKey) ReferenceColumns() []string {
	return fk.referenceColumns
}

// UpdateOption XXX
func (fk ForeignKey) UpdateOption() string {
	return fk.updateOption
}

// DeleteOption XXX
func (fk ForeignKey) DeleteOption() string {
	return fk.deleteOption
}

// ToSQL return foreign key sql string
func (fk ForeignKey) ToSQL() string {
	var foreignColumnsStr, referenceColumnsStr []string
	for _, fc := range fk.foreignColumns {
		foreignColumnsStr = append(foreignColumnsStr, query.Quote(fc))
	}
	for _, rc := range fk.referenceColumns {
		referenceColumnsStr = append(referenceColumnsStr, query.Quote(rc))
	}
	sql := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s (%s)",
		strings.Join(foreignColumnsStr, ", "),
		query.Quote(fk.referenceTableName),
		strings.Join(referenceColumnsStr, ", "))
	if fk.deleteOption != "" {
		sql = sql + fmt.Sprintf(" ON DELETE %s", fk.deleteOption)
	}
	if fk.updateOption != "" {
		sql = sql + fmt.Sprintf(" ON UPDATE %s", fk.updateOption)
	}
	return sql

}

// AddIndex returns a new Index
func AddIndex(idxName string, columns ...string) Index {
	return Index{
		name:    idxName,
		columns: columns,
	}
}

// AddUniqueIndex returns a new UniqueIndex
func AddUniqueIndex(idxName string, columns ...string) UniqueIndex {
	return UniqueIndex{
		name:    idxName,
		columns: columns,
	}
}

// AddFullTextIndex returns a new FullTextIndex
func AddFullTextIndex(idxName string, columns ...string) FullTextIndex {
	return FullTextIndex{
		name:    idxName,
		columns: columns,
	}
}

// AddSpatialIndex returns a new SpatialIndex
func AddSpatialIndex(idxName string, columns ...string) SpatialIndex {
	return SpatialIndex{
		name:    idxName,
		columns: columns,
	}
}

// AddPrimaryKey return initialized PrimaryKey struct.
func AddPrimaryKey(columns ...string) PrimaryKey {
	return PrimaryKey{
		columns: columns,
	}
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

func varchar(size uint64) string {
	if size == 0 {
		return fmt.Sprintf("VARCHAR(%d)", defaultVarcharSize)
	}

	return fmt.Sprintf("VARCHAR(%d)", size)
}

func varbinary(size uint64) string {
	if size == 0 {
		return fmt.Sprintf("VARBINARY(%d)", defaultVarbinarySize)
	}

	return fmt.Sprintf("VARBINARY(%d)", size)
}

func datetime(size uint64) string {
	if size == 0 {
		return "DATETIME"
	}

	return fmt.Sprintf("DATETIME(%d)", size)
}

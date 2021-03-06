package mysql

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestToSQL(t *testing.T) {
	m := MySQL{}

	testcases := []struct {
		typeName string
		size     uint64
		output   string
	}{
		{"bool", 0, "TINYINT(1)"},
		{"*bool", 0, "TINYINT(1)"},
		{"sql.NullBool", 0, "TINYINT(1)"},
		{"int8", 0, "TINYINT"},
		{"*int8", 0, "TINYINT"},
		{"int16", 0, "SMALLINT"},
		{"*int16", 0, "SMALLINT"},
		{"int32", 0, "INTEGER"},
		{"*int32", 0, "INTEGER"},
		{"sql.NullInt32", 0, "INTEGER"},
		{"int64", 0, "BIGINT"},
		{"*int64", 0, "BIGINT"},
		{"sql.NullInt64", 0, "BIGINT"},
		{"uint8", 0, "TINYINT unsigned"},
		{"*uint8", 0, "TINYINT unsigned"},
		{"uint16", 0, "SMALLINT unsigned"},
		{"*uint16", 0, "SMALLINT unsigned"},
		{"uint32", 0, "INTEGER unsigned"},
		{"*uint32", 0, "INTEGER unsigned"},
		{"uint64", 0, "BIGINT unsigned"},
		{"*uint64", 0, "BIGINT unsigned"},
		{"float32", 0, "FLOAT"},
		{"*float32", 0, "FLOAT"},
		{"float64", 0, "DOUBLE"},
		{"*float64", 0, "DOUBLE"},
		{"sql.NullFloat64", 0, "DOUBLE"},
		{"string", 0, fmt.Sprintf("VARCHAR(%d)", defaultVarcharSize)},
		{"*string", 0, fmt.Sprintf("VARCHAR(%d)", defaultVarcharSize)},
		{"sql.NullString", 0, fmt.Sprintf("VARCHAR(%d)", defaultVarcharSize)},
		{"string", 10, "VARCHAR(10)"},
		{"*string", 10, "VARCHAR(10)"},
		{"sql.NullString", 10, "VARCHAR(10)"},
		{"[]uint8", 10, "VARBINARY(10)"},
		{"sql.RawBytes", 10, "VARBINARY(10)"},
		{"tinytext", 0, "TINYTEXT"},
		{"text", 0, "TEXT"},
		{"mediumtext", 0, "MEDIUMTEXT"},
		{"longtext", 0, "LONGTEXT"},
		{"tinyblob", 0, "TINYBLOB"},
		{"blob", 0, "BLOB"},
		{"mediumblob", 0, "MEDIUMBLOB"},
		{"longblob", 0, "LONGBLOB"},
		{"time", 0, "TIME"},
		{"time.Time", 0, "DATETIME"},
		{"time.Time", 6, "DATETIME(6)"},
		{"mysql.NullTime", 0, "DATETIME"}, // https://godoc.org/github.com/go-sql-driver/mysql#NullTime
		{"sql.NullTime", 0, "DATETIME"},   // from Go 1.13
		{"date", 0, "DATE"},
		{"geometry", 0, "GEOMETRY"},
		{"json.RawMessage", 0, "JSON"},
	}

	for _, tc := range testcases {
		got, err := m.ToSQL(tc.typeName, tc.size)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.output {
			t.Fatalf("error %s to sql %s. but result %s", tc.typeName, tc.output, got)
		}
	}
}

func TestToSQL2(t *testing.T) {
	m := MySQL{}

	testcases := []struct {
		typeName string
		size     uint64
		output   error
	}{
		{"noExistType", 0, ErrInvalidType},
	}

	for _, tc := range testcases {
		_, got := m.ToSQL(tc.typeName, tc.size)
		if !errors.As(got, &tc.output) {
			t.Errorf("mismatch want=%v, got=%v", tc.output, got)
		}
	}
}

func TestAuotIncrement(t *testing.T) {
	m := MySQL{}
	if m.AutoIncrement() != autoIncrement {
		t.Fatalf("error auto increament: %s. result:%s", autoIncrement, m.AutoIncrement())
	}
}

func TestAddIndex(t *testing.T) {
	index := AddIndex("player_id_idx", "player_id")
	if index.ToSQL() != "INDEX `player_id_idx` (`player_id`)" {
		t.Fatal("[error] parse player_id_idx. ", index.ToSQL())
	}

	index = AddIndex("player_entry_id_idx", "player_id", "entry_id")
	if index.ToSQL() != "INDEX `player_entry_id_idx` (`player_id`, `entry_id`)" {
		t.Fatal("[error] parse player_entry_id_idx", index.ToSQL())
	}
}

func TestAddUniqIndex(t *testing.T) {
	uniqIndex := AddUniqueIndex("player_id_idx", "player_id")
	if uniqIndex.ToSQL() != "UNIQUE `player_id_idx` (`player_id`)" {
		t.Fatal("[error] parse unique player_id_idx", uniqIndex.ToSQL())
	}

	uniqIndex = AddUniqueIndex("player_entry_id_idx", "player_id", "entry_id")
	if uniqIndex.ToSQL() != "UNIQUE `player_entry_id_idx` (`player_id`, `entry_id`)" {
		t.Fatal("[error] parse unique player_entry_id_idx", uniqIndex.ToSQL())
	}
}

func TestAddFullTextIndex(t *testing.T) {
	fullTextIndex := AddFullTextIndex("full_text_idx", "content")
	if fullTextIndex.ToSQL() != "FULLTEXT `full_text_idx` (`content`)" {
		t.Fatal("[error] parse full_text_idx", fullTextIndex.ToSQL())
	}

	fullTextIndex = AddFullTextIndex("full_text_idx", "content", "title")
	if fullTextIndex.ToSQL() != "FULLTEXT `full_text_idx` (`content`, `title`)" {
		t.Fatal("[error] parse full_text_idx", fullTextIndex.ToSQL())
	}

	fullTextIndex = AddFullTextIndex("full_text_idx", "content").WithParser("ngram")
	if fullTextIndex.ToSQL() != "FULLTEXT `full_text_idx` (`content`) WITH PARSER `ngram`" {
		t.Fatal("[error] parse full_text_idx", fullTextIndex.ToSQL())
	}
}

func TestAddAddSpatialIndex(t *testing.T) {
	spatialIndex := AddSpatialIndex("geometry_idx", "g")
	if spatialIndex.ToSQL() != "SPATIAL KEY `geometry_idx` (`g`)" {
		t.Fatal("[error] parse geometry_idx", spatialIndex.ToSQL())
	}

	spatialIndex = AddSpatialIndex("geometry_idx", "g", "g1")
	if spatialIndex.ToSQL() != "SPATIAL KEY `geometry_idx` (`g`, `g1`)" {
		t.Fatal("[error] parse geometry_idx", spatialIndex.ToSQL())
	}
}

func TestAddPrimaryKey(t *testing.T) {
	pk := AddPrimaryKey("id")
	if pk.ToSQL() != "PRIMARY KEY (`id`)" {
		t.Fatal("[error] parse primary key", pk.ToSQL())
	}

	pk = AddPrimaryKey("id", "created_at")
	if pk.ToSQL() != "PRIMARY KEY (`id`, `created_at`)" {
		t.Fatal("[error] parse primary key", pk.ToSQL())
	}

	pk = AddPrimaryKey("created_at", "id")
	if pk.ToSQL() != "PRIMARY KEY (`created_at`, `id`)" {
		t.Fatal("[error] parse primary key", pk.ToSQL())
	}
}

func TestAddForeignKey(t *testing.T) {
	foreignColumns := []string{"player_id"}
	referenceColumns := []string{"id"}
	fk := AddForeignKey(foreignColumns, referenceColumns, "player")
	if fk.ToSQL() != "FOREIGN KEY (`player_id`) REFERENCES `player` (`id`)" {
		t.Fatal("[error] parse foreign key", fk.ToSQL())
	}

	foreignColumns = []string{"product_category", "product_id"}
	referenceColumns = []string{"category", "id"}
	fk = AddForeignKey(foreignColumns, referenceColumns, "product")
	if fk.ToSQL() != "FOREIGN KEY (`product_category`, `product_id`) REFERENCES `product` (`category`, `id`)" {
		t.Fatal("[error] parse foreign key", fk.ToSQL())
	}

	foreignColumns = []string{"product_category", "product_id"}
	referenceColumns = []string{"category", "id"}
	fk = AddForeignKey(foreignColumns, referenceColumns, "product", WithUpdateForeignKeyOption(ForeignKeyOptionNoAction), WithDeleteForeignKeyOption(ForeignKeyOptionNoAction))
	if fk.ToSQL() != "FOREIGN KEY (`product_category`, `product_id`) REFERENCES `product` (`category`, `id`)" {
		t.Fatal("[error] parse foreign key", fk.ToSQL())
	}

	foreignColumns = []string{"product_category", "product_id"}
	referenceColumns = []string{"category", "id"}
	fk = AddForeignKey(foreignColumns, referenceColumns, "product", WithUpdateForeignKeyOption(ForeignKeyOptionCascade))
	if fk.ToSQL() != "FOREIGN KEY (`product_category`, `product_id`) REFERENCES `product` (`category`, `id`) ON UPDATE CASCADE" {
		t.Fatal("[error] parse foreign key", fk.ToSQL())
	}
}

func Test_varbinary(t *testing.T) {
	type args struct {
		size uint64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "[Normal] return default size",
			args: args{
				size: 0,
			},
			want: "VARBINARY(767)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := varbinary(tt.args.size); got != tt.want {
				t.Errorf("varbinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForeignKey_ToSQL(t *testing.T) {
	type fields struct {
		foreignColumns     []string
		referenceTableName string
		referenceColumns   []string
		updateOption       string
		deleteOption       string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "[Normal] Add delete option",
			fields: fields{
				deleteOption: "cascade",
			},
			want: "FOREIGN KEY () REFERENCES `` () ON DELETE cascade",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fk := ForeignKey{
				foreignColumns:     tt.fields.foreignColumns,
				referenceTableName: tt.fields.referenceTableName,
				referenceColumns:   tt.fields.referenceColumns,
				updateOption:       tt.fields.updateOption,
				deleteOption:       tt.fields.deleteOption,
			}
			if got := fk.ToSQL(); got != tt.want {
				t.Errorf("ForeignKey.ToSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDeleteForeignKeyOption(t *testing.T) {
	type args struct {
		option ForeignKeyOptionType
	}
	tests := []struct {
		name string
		args args
		want ForeignKeyOption
	}{
		{
			name: "[Normal] withDeleteForeignKeyOption set default",
			args: args{
				option: ForeignKeyOptionSetDefault,
			},
			want: WithDeleteForeignKeyOption(ForeignKeyOptionSetDefault),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithDeleteForeignKeyOption(tt.args.option)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Compare value is mismatch (-want +got):%s\n", diff)
			}
		})
	}
}

func TestMySQL_HeaderTemplate(t *testing.T) {
	type fields struct {
		Engine  string
		Charset string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "[Normal] return header template",
			fields: fields{},
			want: `SET foreign_key_checks=0;
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mysql := MySQL{
				Engine:  tt.fields.Engine,
				Charset: tt.fields.Charset,
			}
			got := mysql.HeaderTemplate()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Compare value is mismatch (-want +got):%s\n", diff)
			}
		})
	}
}

func TestMySQL_FooterTemplate(t *testing.T) {
	type fields struct {
		Engine  string
		Charset string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "[Normal] return footer template",
			fields: fields{},
			want: `SET foreign_key_checks=1;
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mysql := MySQL{
				Engine:  tt.fields.Engine,
				Charset: tt.fields.Charset,
			}
			got := mysql.FooterTemplate()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Compare value is mismatch (-want +got):%s\n", diff)
			}
		})
	}
}

func TestMySQL_TableTemplate(t *testing.T) {
	type fields struct {
		Engine  string
		Charset string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "[Normal] return table template",
			fields: fields{},
			want: `
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

`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mysql := MySQL{
				Engine:  tt.fields.Engine,
				Charset: tt.fields.Charset,
			}
			got := mysql.TableTemplate()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Compare value is mismatch (-want +got):%s\n", diff)
			}
		})
	}
}

func TestMySQL_Quote(t *testing.T) {
	type fields struct {
		Engine  string
		Charset string
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "[Normal] return quoted string",
			fields: fields{},
			args:   args{s: "test"},
			want:   "`test`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mysql := MySQL{
				Engine:  tt.fields.Engine,
				Charset: tt.fields.Charset,
			}
			if got := mysql.Quote(tt.args.s); got != tt.want {
				t.Errorf("MySQL.Quote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndex_Name(t *testing.T) {
	type fields struct {
		columns []string
		name    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "[Normal] return index name",
			fields: fields{
				columns: []string{},
				name:    "index name",
			},
			want: "index name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Index{
				columns: tt.fields.columns,
				name:    tt.fields.name,
			}
			if got := i.Name(); got != tt.want {
				t.Errorf("Index.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndex_Columns(t *testing.T) {
	type fields struct {
		columns []string
		name    string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "[Normal] return columns",
			fields: fields{
				columns: []string{"aa", "bb", "cc"},
				name:    "index name",
			},
			want: []string{"aa", "bb", "cc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Index{
				columns: tt.fields.columns,
				name:    tt.fields.name,
			}
			if got := i.Columns(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Index.Columns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniqueIndex_Name(t *testing.T) {
	type fields struct {
		columns []string
		name    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "[Normal] return index name",
			fields: fields{
				columns: []string{"aa", "bb", "cc"},
				name:    "uniq index name",
			},
			want: "uniq index name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := UniqueIndex{
				columns: tt.fields.columns,
				name:    tt.fields.name,
			}
			if got := ui.Name(); got != tt.want {
				t.Errorf("UniqueIndex.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniqueIndex_Columns(t *testing.T) {
	type fields struct {
		columns []string
		name    string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "[Normal] return columns",
			fields: fields{
				columns: []string{"aa", "bb", "cc"},
				name:    "index name",
			},
			want: []string{"aa", "bb", "cc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ui := UniqueIndex{
				columns: tt.fields.columns,
				name:    tt.fields.name,
			}
			if got := ui.Columns(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UniqueIndex.Columns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFullTextIndex_Name(t *testing.T) {
	type fields struct {
		columns []string
		name    string
		parser  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "[Normal] return index name",
			fields: fields{
				columns: []string{"aa", "bb", "cc"},
				name:    "index name",
				parser:  "parser",
			},
			want: "index name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fi := FullTextIndex{
				columns: tt.fields.columns,
				name:    tt.fields.name,
				parser:  tt.fields.parser,
			}
			if got := fi.Name(); got != tt.want {
				t.Errorf("FullTextIndex.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFullTextIndex_Columns(t *testing.T) {
	type fields struct {
		columns []string
		name    string
		parser  string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "[Normal] return columns",
			fields: fields{
				columns: []string{"aa", "bb", "cc"},
				name:    "index name",
				parser:  "parser",
			},
			want: []string{"aa", "bb", "cc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fi := FullTextIndex{
				columns: tt.fields.columns,
				name:    tt.fields.name,
				parser:  tt.fields.parser,
			}
			if got := fi.Columns(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FullTextIndex.Columns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpatialIndex_Name(t *testing.T) {
	type fields struct {
		columns []string
		name    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "[Normal] return index name",
			fields: fields{
				columns: []string{"aa", "bb", "cc"},
				name:    "index name",
			},
			want: "index name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si := SpatialIndex{
				columns: tt.fields.columns,
				name:    tt.fields.name,
			}
			if got := si.Name(); got != tt.want {
				t.Errorf("SpatialIndex.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpatialIndex_Columns(t *testing.T) {
	type fields struct {
		columns []string
		name    string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "[Normal] return columns",
			fields: fields{
				columns: []string{"aa", "bb", "cc"},
				name:    "index name",
			},
			want: []string{"aa", "bb", "cc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si := SpatialIndex{
				columns: tt.fields.columns,
				name:    tt.fields.name,
			}
			if got := si.Columns(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SpatialIndex.Columns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrimaryKey_Columns(t *testing.T) {
	type fields struct {
		columns []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "[Normal] return columns",
			fields: fields{
				columns: []string{"aa", "bb", "cc"},
			},
			want: []string{"aa", "bb", "cc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk := PrimaryKey{
				columns: tt.fields.columns,
			}
			if got := pk.Columns(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrimaryKey.Columns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForeignKey_ForeignColumns(t *testing.T) {
	type fields struct {
		foreignColumns     []string
		referenceTableName string
		referenceColumns   []string
		updateOption       string
		deleteOption       string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "[Normal] return columns",
			fields: fields{
				foreignColumns: []string{"aa", "bb", "cc"},
			},
			want: []string{"aa", "bb", "cc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fk := ForeignKey{
				foreignColumns:     tt.fields.foreignColumns,
				referenceTableName: tt.fields.referenceTableName,
				referenceColumns:   tt.fields.referenceColumns,
				updateOption:       tt.fields.updateOption,
				deleteOption:       tt.fields.deleteOption,
			}
			if got := fk.ForeignColumns(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ForeignKey.ForeignColumns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForeignKey_ReferenceTableName(t *testing.T) {
	type fields struct {
		foreignColumns     []string
		referenceTableName string
		referenceColumns   []string
		updateOption       string
		deleteOption       string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "[Normal] return table name",
			fields: fields{
				referenceTableName: "table",
			},
			want: "table",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fk := ForeignKey{
				foreignColumns:     tt.fields.foreignColumns,
				referenceTableName: tt.fields.referenceTableName,
				referenceColumns:   tt.fields.referenceColumns,
				updateOption:       tt.fields.updateOption,
				deleteOption:       tt.fields.deleteOption,
			}
			if got := fk.ReferenceTableName(); got != tt.want {
				t.Errorf("ForeignKey.ReferenceTableName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForeignKey_ReferenceColumns(t *testing.T) {
	type fields struct {
		foreignColumns     []string
		referenceTableName string
		referenceColumns   []string
		updateOption       string
		deleteOption       string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "[Normal] return referaence columns",
			fields: fields{
				referenceColumns: []string{"aa", "bb", "cc"},
			},
			want: []string{"aa", "bb", "cc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fk := ForeignKey{
				foreignColumns:     tt.fields.foreignColumns,
				referenceTableName: tt.fields.referenceTableName,
				referenceColumns:   tt.fields.referenceColumns,
				updateOption:       tt.fields.updateOption,
				deleteOption:       tt.fields.deleteOption,
			}
			if got := fk.ReferenceColumns(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ForeignKey.ReferenceColumns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForeignKey_UpdateOption(t *testing.T) {
	type fields struct {
		foreignColumns     []string
		referenceTableName string
		referenceColumns   []string
		updateOption       string
		deleteOption       string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "[Normal] return update option",
			fields: fields{
				updateOption: ForeignKeyOptionSetDefault.String(),
			},
			want: ForeignKeyOptionSetDefault.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fk := ForeignKey{
				foreignColumns:     tt.fields.foreignColumns,
				referenceTableName: tt.fields.referenceTableName,
				referenceColumns:   tt.fields.referenceColumns,
				updateOption:       tt.fields.updateOption,
				deleteOption:       tt.fields.deleteOption,
			}
			if got := fk.UpdateOption(); got != tt.want {
				t.Errorf("ForeignKey.UpdateOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForeignKey_DeleteOption(t *testing.T) {
	type fields struct {
		foreignColumns     []string
		referenceTableName string
		referenceColumns   []string
		updateOption       string
		deleteOption       string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "[Normal] return delete option",
			fields: fields{
				deleteOption: ForeignKeyOptionSetDefault.String(),
			},
			want: ForeignKeyOptionSetDefault.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fk := ForeignKey{
				foreignColumns:     tt.fields.foreignColumns,
				referenceTableName: tt.fields.referenceTableName,
				referenceColumns:   tt.fields.referenceColumns,
				updateOption:       tt.fields.updateOption,
				deleteOption:       tt.fields.deleteOption,
			}
			if got := fk.DeleteOption(); got != tt.want {
				t.Errorf("ForeignKey.DeleteOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

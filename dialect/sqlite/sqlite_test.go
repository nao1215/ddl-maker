package sqlite

import (
	"reflect"
	"testing"
)

func TestSQLite_HeaderTemplate(t *testing.T) {
	tests := []struct {
		name   string
		sqlite SQLite
		want   string
	}{
		{
			name:   "[Normal] return header",
			sqlite: SQLite{},
			want: `PRAGMA foreign_keys = false;
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlite := SQLite{}
			if got := sqlite.HeaderTemplate(); got != tt.want {
				t.Errorf("SQLite.HeaderTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLite_FooterTemplate(t *testing.T) {
	tests := []struct {
		name   string
		sqlite SQLite
		want   string
	}{
		{
			name:   "[Normal] return footer",
			sqlite: SQLite{},
			want: `PRAGMA foreign_keys = true;
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlite := SQLite{}
			if got := sqlite.FooterTemplate(); got != tt.want {
				t.Errorf("SQLite.FooterTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLite_TableTemplate(t *testing.T) {
	tests := []struct {
		name   string
		sqlite SQLite
		want   string
	}{
		{
			name:   "[Normal] return table",
			sqlite: SQLite{},
			want: `
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

`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlite := SQLite{}
			if got := sqlite.TableTemplate(); got != tt.want {
				t.Errorf("SQLite.TableTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLite_ToSQL(t *testing.T) {
	type args struct {
		typeName string
		size     uint64
	}
	tests := []struct {
		name    string
		sqlite  SQLite
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "[Normal] int8 to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "int8",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] *int8 to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "*int8",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] int16 to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "int16",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] int32 to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "int32",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] *int32 to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "*int32",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] sql.NullInt32 to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "sql.NullInt32",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] int64 to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "int64",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] *int64 to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "*int64",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] sql.NullInt64 to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "sql.NullInt64",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] uint8 to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "uint8",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] *uint8 to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "*uint8",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] uint16 to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "uint16",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] *uint16 to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "*uint16",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] uint32 to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "uint32",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] *uint32 to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "*uint32",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] uint64 to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "uint64",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] *uint64 to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "*uint64",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] float32 to REAL",
			sqlite: SQLite{},
			args: args{
				typeName: "float32",
			},
			want:    "REAL",
			wantErr: false,
		},
		{
			name:   "[Normal] *float32 to REAL",
			sqlite: SQLite{},
			args: args{
				typeName: "*float32",
			},
			want:    "REAL",
			wantErr: false,
		},
		{
			name:   "[Normal] float64 to REAL",
			sqlite: SQLite{},
			args: args{
				typeName: "float64",
			},
			want:    "REAL",
			wantErr: false,
		},
		{
			name:   "[Normal] *float64 to REAL",
			sqlite: SQLite{},
			args: args{
				typeName: "*float64",
			},
			want:    "REAL",
			wantErr: false,
		},
		{
			name:   "[Normal] sql.NullFloat64 to REAL",
			sqlite: SQLite{},
			args: args{
				typeName: "sql.NullFloat64",
			},
			want:    "REAL",
			wantErr: false,
		},
		{
			name:   "[Normal] string to TEXT",
			sqlite: SQLite{},
			args: args{
				typeName: "string",
			},
			want:    "TEXT",
			wantErr: false,
		},
		{
			name:   "[Normal] *string to TEXT",
			sqlite: SQLite{},
			args: args{
				typeName: "*string",
			},
			want:    "TEXT",
			wantErr: false,
		},
		{
			name:   "[Normal] sql.NullString to TEXT",
			sqlite: SQLite{},
			args: args{
				typeName: "sql.NullString",
			},
			want:    "TEXT",
			wantErr: false,
		},
		{
			name:   "[Normal] []uint8 to BLOB",
			sqlite: SQLite{},
			args: args{
				typeName: "[]uint8",
			},
			want:    "BLOB",
			wantErr: false,
		},
		{
			name:   "[Normal] sql.RawBytes to BLOB",
			sqlite: SQLite{},
			args: args{
				typeName: "sql.RawBytes",
			},
			want:    "BLOB",
			wantErr: false,
		},
		{
			name:   "[Normal] bool to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "bool",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] *bool to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "*bool",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] sql.NullBool to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "sql.NullBool",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] text to TEXT",
			sqlite: SQLite{},
			args: args{
				typeName: "text",
			},
			want:    "TEXT",
			wantErr: false,
		},
		{
			name:   "[Normal] tinytext to TEXT",
			sqlite: SQLite{},
			args: args{
				typeName: "tinytext",
			},
			want:    "TEXT",
			wantErr: false,
		},
		{
			name:   "[Normal] mediumtext to TEXT",
			sqlite: SQLite{},
			args: args{
				typeName: "mediumtext",
			},
			want:    "TEXT",
			wantErr: false,
		},
		{
			name:   "[Normal] longtext to TEXT",
			sqlite: SQLite{},
			args: args{
				typeName: "longtext",
			},
			want:    "TEXT",
			wantErr: false,
		},
		{
			name:   "[Normal] tinyblob to BLOB",
			sqlite: SQLite{},
			args: args{
				typeName: "tinyblob",
			},
			want:    "BLOB",
			wantErr: false,
		},
		{
			name:   "[Normal] blob to BLOB",
			sqlite: SQLite{},
			args: args{
				typeName: "blob",
			},
			want:    "BLOB",
			wantErr: false,
		},
		{
			name:   "[Normal] mediumblob to BLOB",
			sqlite: SQLite{},
			args: args{
				typeName: "mediumblob",
			},
			want:    "BLOB",
			wantErr: false,
		},
		{
			name:   "[Normal] longblob to BLOB",
			sqlite: SQLite{},
			args: args{
				typeName: "longblob",
			},
			want:    "BLOB",
			wantErr: false,
		},
		{
			name:   "[Normal] time to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "time",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] time.Time to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "time.Time",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] *time.Time to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "*time.Time",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] sql.NullTime to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "sql.NullTime",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] date to INTEGER",
			sqlite: SQLite{},
			args: args{
				typeName: "date",
			},
			want:    "INTEGER",
			wantErr: false,
		},
		{
			name:   "[Normal] json.RawMessage to JSON",
			sqlite: SQLite{},
			args: args{
				typeName: "json.RawMessage",
			},
			want:    "JSON",
			wantErr: false,
		},
		{
			name:   "[Normal] *json.RawMessage to JSON",
			sqlite: SQLite{},
			args: args{
				typeName: "*json.RawMessage",
			},
			want:    "JSON",
			wantErr: false,
		},
		{
			name:   "[Error] can not convert geometry",
			sqlite: SQLite{},
			args: args{
				typeName: "geometry",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlite := SQLite{}
			got, err := sqlite.ToSQL(tt.args.typeName, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLite.ToSQL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SQLite.ToSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLite_Quote(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		sqlite SQLite
		args   args
		want   string
	}{
		{
			name:   "[Normal] return string that encloses with ``",
			sqlite: SQLite{},
			args: args{
				s: "test_code",
			},
			want: "`test_code`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlite := SQLite{}
			if got := sqlite.Quote(tt.args.s); got != tt.want {
				t.Errorf("SQLite.Quote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLite_AutoIncrement(t *testing.T) {
	tests := []struct {
		name   string
		sqlite SQLite
		want   string
	}{
		{
			name:   "[Normal] return AUTOINCREMENT",
			sqlite: SQLite{},
			want:   "AUTOINCREMENT",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlite := SQLite{}
			if got := sqlite.AutoIncrement(); got != tt.want {
				t.Errorf("SQLite.AutoIncrement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddPrimaryKey(t *testing.T) {
	type args struct {
		columns []string
	}
	tests := []struct {
		name string
		args args
		want PrimaryKey
	}{
		{
			name: "[Normal] return PrimaryKey struct",
			args: args{
				columns: []string{"aa", "bb", "cc"},
			},
			want: PrimaryKey{
				columns: []string{"aa", "bb", "cc"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddPrimaryKey(tt.args.columns...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddPrimaryKey() = %v, want %v", got, tt.want)
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
			name: "[Normal] return PrimaryKey columns",
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

func TestPrimaryKey_ToSQL(t *testing.T) {
	type fields struct {
		columns []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "[Normal] return PRIMARY KEY query",
			fields: fields{
				columns: []string{"aa", "bb", "cc"},
			},
			want: "PRIMARY KEY (`aa`, `bb`, `cc`)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk := PrimaryKey{
				columns: tt.fields.columns,
			}
			if got := pk.ToSQL(); got != tt.want {
				t.Errorf("PrimaryKey.ToSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndex_ToSQL(t *testing.T) {
	type fields struct {
		columns []string
		table   string
		name    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "[Normal] return INDEX query",
			fields: fields{
				columns: []string{"aa", "bb", "cc"},
				table:   "test_table",
				name:    "test_index",
			},
			want: "CREATE INDEX `test_index` ON `test_table` (`aa`, `bb`, `cc`);",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Index{
				columns: tt.fields.columns,
				table:   tt.fields.table,
				name:    tt.fields.name,
			}
			if got := i.ToSQL(); got != tt.want {
				t.Errorf("Index.ToSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

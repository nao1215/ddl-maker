package sqlite

import (
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
			name:    "[Normal] success to convert sql",
			sqlite:  SQLite{},
			args:    args{},
			want:    "",
			wantErr: false,
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

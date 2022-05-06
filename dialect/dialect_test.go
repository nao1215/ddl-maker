package dialect

import (
	"reflect"
	"testing"

	"github.com/nao1215/ddl-maker/dialect/mysql"
	"github.com/nao1215/ddl-maker/dialect/sqlite"
)

func TestSort(t *testing.T) {
	var indexes Indexes

	idx1 := mysql.AddUniqueIndex("fuga_idx", "fuga")
	indexes = append(indexes, idx1)
	idx2 := mysql.AddIndex("hoge_idx", "hoge")
	indexes = append(indexes, idx2)

	idxes := indexes.Sort()
	if len(idxes) != 2 {
		t.Fatal("error sort Indexes", idxes)
	}
	if idxes[0].ToSQL() != idx2.ToSQL() {
		t.Fatal("error sort index", idxes[0].ToSQL())
	}

}

func TestForeignKeys_Sort(t *testing.T) {
	tests := []struct {
		name        string
		foreignKeys ForeignKeys
		want        ForeignKeys
	}{
		{
			name: "[Normal] can sort foreign keys",
			foreignKeys: ForeignKeys{
				mysql.AddForeignKey(
					[]string{"player_id"},
					[]string{"id"},
					"player",
				),
				mysql.AddForeignKey(
					[]string{"entry_id"},
					[]string{"id"},
					"entry",
				),
				mysql.AddForeignKey(
					[]string{"dummy_id"},
					[]string{"id"},
					"dummy",
				),
				mysql.AddForeignKey(
					[]string{"zynq_id"},
					[]string{"id"},
					"zynq",
				),
			},
			want: ForeignKeys{
				mysql.AddForeignKey(
					[]string{"dummy_id"},
					[]string{"id"},
					"dummy",
				),
				mysql.AddForeignKey(
					[]string{"entry_id"},
					[]string{"id"},
					"entry",
				),
				mysql.AddForeignKey(
					[]string{"player_id"},
					[]string{"id"},
					"player",
				),
				mysql.AddForeignKey(
					[]string{"zynq_id"},
					[]string{"id"},
					"zynq",
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.foreignKeys.Sort(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ForeignKeys.Sort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		driver  string
		engine  string
		charset string
	}
	tests := []struct {
		name    string
		args    args
		want    Dialect
		wantErr bool
	}{
		{
			name: "[Normal] return mysql dialect",
			args: args{
				driver:  "mysql",
				engine:  "",
				charset: "",
			},
			want:    &mysql.MySQL{},
			wantErr: false,
		},
		{
			name: "[Normal] return sqlite dialect",
			args: args{
				driver:  "sqlite",
				engine:  "",
				charset: "",
			},
			want:    &sqlite.SQLite{},
			wantErr: false,
		},
		{
			name: "[Error] no such driver",
			args: args{
				driver:  "unknown",
				engine:  "",
				charset: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.driver, tt.args.engine, tt.args.charset)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

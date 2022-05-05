package dialect

import (
	"reflect"
	"testing"

	"github.com/nao1215/ddl-maker/dialect/mysql"
)

func TestNew(t *testing.T) {
	_, err := New("", "", "")
	if err == nil {
		t.Fatal("error not set driver")
	}

	_, err = New("mysql", "", "")
	if err != nil {
		t.Fatalf("error new dialect:%s error", "mysql")
	}
}

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

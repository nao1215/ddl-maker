package ddlmaker

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/nao1215/ddl-maker/dialect"
	"github.com/nao1215/ddl-maker/dialect/mock"
	"github.com/nao1215/ddl-maker/dialect/mysql"
)

type T1 struct {
	ID          uint64 `ddl:"auto"`
	Name        string
	Description sql.NullString `ddl:"null,text"`
	CreatedAt   time.Time
	Binary      []byte
	Ignore      string `ddl:"-"`
}

type T2 struct {
	ID     uint64 `ddl:"auto"`
	Ignore string `ddl:"-"`
	Name   *string
}

func (t1 T1) Table() string {
	return "test_one"
}

func (t1 T1) Indexes() dialect.Indexes {
	return dialect.Indexes{
		mysql.AddUniqueIndex("token_idx", "token"),
	}
}

func (t1 T1) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id", "created_at")
}

func (t1 T1) ForeignKeys() dialect.ForeignKeys {
	return dialect.ForeignKeys{
		mysql.AddForeignKey([]string{"player_id"}, []string{"id"}, "player",
			mysql.WithUpdateForeignKeyOption(mysql.ForeignKeyOptionNoAction),
			mysql.WithDeleteForeignKeyOption(mysql.ForeignKeyOptionNoAction)),
	}
}

func TestParseField(t *testing.T) {
	t1 := T1{}
	idColumn := column{
		name:     "id",
		tag:      "auto",
		typeName: "uint64",
		dialect:  mysql.MySQL{},
	}
	nameColumn := column{
		name:     "name",
		typeName: "string",
		dialect:  mysql.MySQL{},
	}
	descColumn := column{
		name:     "description",
		typeName: "sql.NullString",
		tag:      "null,text",
		dialect:  mysql.MySQL{},
	}
	createdAtColumn := column{
		name:     "created_at",
		typeName: "time.Time",
		dialect:  mysql.MySQL{},
	}
	binaryColumn := column{
		name:     "binary",
		typeName: "[]uint8",
		dialect:  mysql.MySQL{},
	}
	columns := []dialect.Column{idColumn, nameColumn, descColumn, createdAtColumn, binaryColumn}

	rt := reflect.TypeOf(t1)

	if rt.NumField() == 0 {
		t.Fatal("T1 field is 0")
	}

	for i := 0; i < rt.NumField(); i++ {
		column, err := parseField(rt.Field(i), mysql.MySQL{})
		if err != nil {
			if err == ErrIgnoreField {
				continue
			}
			t.Fatal("error parse field", err.Error())
		}

		if !reflect.DeepEqual(columns[i], column) {
			t.Fatalf("parsed %s: %v is different \n %s: %v", column.Name(), column, columns[i].Name(), columns[i])
		}
	}
}

func TestParseTable(t *testing.T) {
	t1 := T1{}
	d := mysql.MySQL{}

	var columns []dialect.Column
	table := parseTable(t1, columns, d)
	if table.Name() != d.Quote(t1.Table()) {
		t.Fatal("error parse table name", table.Name())
	}

	if len(table.Indexes()) != len(t1.Indexes()) {
		t.Fatal("error parse index ", len(table.Indexes()))
	}

	if table.PrimaryKey().ToSQL() != "PRIMARY KEY (`id`, `created_at`)" {
		t.Fatal("error parse pk: ", table.PrimaryKey().ToSQL())
	}

	if len(table.ForeignKeys()) != len(t1.ForeignKeys()) {
		t.Fatal("error parse fk: ", len(table.ForeignKeys()))
	}
}

func TestDDLMaker_parse(t *testing.T) {
	t.Run("[Normal] can parse ignore field and pointer field", func(t *testing.T) {
		dm := DDLMaker{}
		dm.Dialect = &mock.SQLMock{
			Engine:  "dummy",
			Charset: "dummy",
		}

		if err := dm.AddStruct(&T2{}); err != nil {
			t.Fatal(err)
		}

		got := dm.parse()
		if got != nil {
			t.Error(got)
		}
	})
}

package ddlmaker

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/ddl-maker/dialect"
	"github.com/nao1215/ddl-maker/dialect/mock"
	"github.com/nao1215/ddl-maker/dialect/mysql"
	"github.com/nao1215/ddl-maker/dialect/sqlite"
)

type TestOne struct {
	ID        uint64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t1 TestOne) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id")
}

type TestTwo struct {
	ID        uint64
	TestOneID uint64
	Comment   sql.NullString `ddl:"null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type unknown int
type TestThree struct {
	ID unknown
}

func (t2 *TestTwo) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id", "created_at")
}

func TestNew(t *testing.T) {
	conf := Config{}
	_, err := New(conf)
	if err == nil {
		t.Fatal("Not set driver name", err)
	}

	conf = Config{
		DB: DBConfig{Driver: "dummy"},
	}
	_, err = New(conf)
	if err == nil {
		t.Fatal("Set unsupport driver name", err)
	}

	conf = Config{
		DB: DBConfig{Driver: "mysql"},
	}
	_, err = New(conf)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddStruct(t *testing.T) {
	dm, err := New(Config{
		DB: DBConfig{Driver: "mysql"},
	})
	if err != nil {
		t.Fatal(err)
	}
	err = dm.AddStruct(nil)
	if err == nil {
		t.Fatal("nil is not support")
	}

	err = dm.AddStruct(TestOne{}, TestTwo{})
	if err != nil {
		t.Fatal(err)
	}
	if len(dm.Structs) != 2 {
		t.Fatal("[error] add stuct")
	}

	err = dm.AddStruct(TestOne{})
	if err != nil {
		t.Fatal("[error] add duplicate struct")
	}
}

func TestAddStruct2(t *testing.T) {
	t.Run("[Error] add same struct at twice", func(t *testing.T) {
		dm, err := New(Config{
			OutFilePath: "",
			DB: DBConfig{
				Driver:  "mysql",
				Engine:  "InnoDB",
				Charset: "utf8mb4",
			},
		})
		if err != nil {
			t.Fatal("error new maker", err)
		}

		dm.Dialect = &mock.SQLMock{
			Engine:  "dummy",
			Charset: "dummy",
		}

		got := dm.AddStruct(&TestOne{}, &TestOne{})
		if got == nil {
			t.Fatal("add struct error did not occure")
		}
		want := "github.com/nao1215/ddl-maker.TestOne is already added"
		if want != got.Error() {
			t.Errorf("mismatch want:%s, got:%s", want, got.Error())
		}
	})
}

func TestGenerate(t *testing.T) {
	m := mysql.MySQL{}
	generatedDDL := fmt.Sprintf(`%s
DROP TABLE IF EXISTS %s;

CREATE TABLE %s (
    %s BIGINT unsigned NOT NULL,
    %s VARCHAR(191) NOT NULL,
    %s DATETIME NOT NULL,
    %s DATETIME NOT NULL,
    PRIMARY KEY (%s)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;

%s`, m.HeaderTemplate(), m.Quote("test_one"), m.Quote("test_one"), m.Quote("id"), m.Quote("name"), m.Quote("created_at"), m.Quote("updated_at"), m.Quote("id"), m.FooterTemplate())

	generatedDDL2 := fmt.Sprintf(`%s
DROP TABLE IF EXISTS %s;

CREATE TABLE %s (
    %s BIGINT unsigned NOT NULL,
    %s BIGINT unsigned NOT NULL,
    %s VARCHAR(191) NULL,
    %s DATETIME NOT NULL,
    %s DATETIME NOT NULL,
    PRIMARY KEY (%s, %s)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;

%s`, m.HeaderTemplate(), m.Quote("test_two"), m.Quote("test_two"), m.Quote("id"), m.Quote("test_one_id"), m.Quote("comment"), m.Quote("created_at"), m.Quote("updated_at"), m.Quote("id"), m.Quote("created_at"), m.FooterTemplate())

	dm, err := New(Config{
		DB: DBConfig{
			Driver:  "mysql",
			Engine:  "InnoDB",
			Charset: "utf8mb4",
		},
	})
	if err != nil {
		t.Fatal("error new maker", err)
	}

	if err = dm.AddStruct(&TestOne{}); err != nil {
		t.Fatal("error add struct", err)
	}
	if err = dm.parse(); err != nil {
		t.Fatal(err)
	}

	var ddl1 bytes.Buffer
	if err = dm.generate(&ddl1); err != nil {
		t.Fatal("error generate ddl", err)
	}

	if ddl1.String() != generatedDDL {
		t.Log(ddl1.String())
		t.Fatalf("generatedDDL: %s \n checkDDLL: %s \n", ddl1.String(), generatedDDL)
	}

	dm2, err := New(Config{
		DB: DBConfig{
			Driver:  "mysql",
			Engine:  "InnoDB",
			Charset: "utf8mb4",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err = dm2.AddStruct(&TestTwo{}); err != nil {
		t.Fatal("error add pointer struct", err)
	}
	if err = dm2.parse(); err != nil {
		t.Fatal(err)
	}

	var ddl2 bytes.Buffer
	if err = dm2.generate(&ddl2); err != nil {
		t.Fatal("error generate ddl", err)
	}

	if ddl2.String() != generatedDDL2 {
		t.Fatalf("generatedDDL: %s \n checkDDLL: %s \n", ddl2.String(), generatedDDL2)
	}
}

func TestGenerate2(t *testing.T) {
	t.Run("[Normal] generate ddl file", func(t *testing.T) {
		dm, err := New(Config{
			OutFilePath: "./testdata/mysql/test.sql",
			DB: DBConfig{
				Driver:  "mysql",
				Engine:  "InnoDB",
				Charset: "utf8mb4",
			},
		})
		if err != nil {
			t.Fatal("error new maker", err)
		}
		defer os.Remove("./testdata/mysql/test.sql")

		if err = dm.AddStruct(&TestOne{}); err != nil {
			t.Fatal("error add struct", err)
		}
		if err = dm.parse(); err != nil {
			t.Fatal(err)
		}

		if err = dm.Generate(); err != nil {
			t.Fatal(err)
		}

		got, err := os.ReadFile("./testdata/mysql/test.sql")
		if err != nil {
			t.Fatal(err)
		}

		want, err := os.ReadFile("./testdata/mysql/golden.sql")
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(string(want), string(got)); diff != "" {
			t.Errorf("Compare value is mismatch (-want +got):%s\n", diff)
		}
	})

	t.Run("[Error] open file error", func(t *testing.T) {
		dm, err := New(Config{
			OutFilePath: "",
			DB: DBConfig{
				Driver:  "mysql",
				Engine:  "InnoDB",
				Charset: "utf8mb4",
			},
		})
		if err != nil {
			t.Fatal("error new maker", err)
		}

		got := dm.Generate()
		want := "error create ddl file: open : no such file or directory"
		if got == nil {
			t.Fatal("open file error did not occure")
		}
		if want != got.Error() {
			t.Errorf("mismatch want:%s, got:%s", want, got.Error())
		}
	})

	t.Run("[Error] generate ddl error", func(t *testing.T) {
		dm, err := New(Config{
			OutFilePath: "./testdata/mysql/test.sql",
			DB: DBConfig{
				Driver:  "mysql",
				Engine:  "InnoDB",
				Charset: "utf8mb4",
			},
		})
		if err != nil {
			t.Fatal("error new maker", err)
		}
		defer os.Remove("./testdata/mysql/test.sql")

		dm.Dialect = &mock.SQLMock{
			Engine:  "dummy",
			Charset: "dummy",
			MockHeaderTemplate: func() string {
				return "{{"
			},
		}

		got := dm.Generate()
		want := "error generate: error parse header template: template: header:1: unclosed action"
		if got == nil {
			t.Fatal("open file error did not occure")
		}
		if want != got.Error() {
			t.Errorf("mismatch want:%s, got:%s", want, got.Error())
		}
	})

	t.Run("[Error] template execute error", func(t *testing.T) {
		dm, err := New(Config{
			OutFilePath: "./testdata/mysql/test.sql",
			DB: DBConfig{
				Driver:  "mysql",
				Engine:  "InnoDB",
				Charset: "utf8mb4",
			},
		})
		if err != nil {
			t.Fatal("error new maker", err)
		}
		defer os.Remove("./testdata/mysql/test.sql")

		err = dm.AddStruct(&TestThree{})
		if err != nil {
			t.Fatal("error add struct", err)
		}

		got := dm.Generate()
		want := mysql.ErrInvalidType
		if got == nil {
			t.Fatal("template execute error did not occure")
		}
		if !errors.As(got, &want) {
			t.Errorf("mismatch want:%v, got:%v", want, got)
		}
	})
}

func TestDDLMaker_generate(t *testing.T) {
	t.Run("[Error] parse header tamplate error", func(t *testing.T) {
		dm := DDLMaker{}
		dm.Dialect = &mock.SQLMock{
			Engine:  "dummy",
			Charset: "dummy",
			MockHeaderTemplate: func() string {
				return "{{"
			},
		}

		var ddl bytes.Buffer
		got := dm.generate(&ddl)
		want := "error parse header template: template: header:1: unclosed action"
		if got == nil {
			t.Fatal("parse error did not occure")
		}
		if want != got.Error() {
			t.Errorf("mismatch want:%s, got:%s", want, got.Error())
		}
	})

	t.Run("[Error] parse footer tamplate error", func(t *testing.T) {
		dm := DDLMaker{}
		dm.Dialect = &mock.SQLMock{
			Engine:  "dummy",
			Charset: "dummy",
			MockHeaderTemplate: func() string {
				return ""
			},
			MockFooterTemplate: func() string {
				return "{{"
			},
		}

		var ddl bytes.Buffer
		got := dm.generate(&ddl)
		want := "error parse footer template: template: footer:1: unclosed action"
		if got == nil {
			t.Fatal("parse error did not occure")
		}
		if want != got.Error() {
			t.Errorf("mismatch want:%s, got:%s", want, got.Error())
		}
	})

	t.Run("[Error] parse table template error", func(t *testing.T) {
		dm := DDLMaker{}
		dm.Dialect = &mock.SQLMock{
			Engine:  "dummy",
			Charset: "dummy",
			MockHeaderTemplate: func() string {
				return ""
			},
			MockFooterTemplate: func() string {
				return ""
			},
			MockTableTemplate: func() string {
				return "{{"
			},
		}

		var ddl bytes.Buffer
		got := dm.generate(&ddl)
		want := "error parse ddl template: template: ddl:1: unclosed action"
		if got == nil {
			t.Fatal("parse error did not occure")
		}
		if want != got.Error() {
			t.Errorf("mismatch want:%s, got:%s", want, got.Error())
		}
	})
}

type User struct {
	ID                  uint64
	Name                string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Token               string `ddl:"-"`
	DailyNotificationAt string `ddl:"type=time"`
}

func (u *User) Table() string {
	return "player"
}

func (u *User) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id")
}

type Entry struct {
	ID        int32
	PlayerID  int32
	Title     string  `ddl:"size=100"` // not used tag
	Public    bool    `ddl:"default=0"`
	Content   *string `ddl:"type=text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (e Entry) PrimaryKey() dialect.PrimaryKey {
	return sqlite.AddPrimaryKey("id")
}

func (e Entry) Indexes() dialect.Indexes {
	return dialect.Indexes{
		sqlite.AddUniqueIndex("created_at_uniq_idx", "entry", "created_at"),
		sqlite.AddIndex("title_idx", "entry", "title"),
		sqlite.AddIndex("created_at_idx", "entry", "created_at"),
	}
}

func (e Entry) ForeignKeys() dialect.ForeignKeys {
	return dialect.ForeignKeys{
		sqlite.AddForeignKey(
			[]string{"player_id"},
			[]string{"id"},
			"player",
			sqlite.WithDeleteForeignKeyOption(sqlite.ForeignKeyOptionCascade),
		),
	}
}

func TestDDLMaker_GenerateForSQLite(t *testing.T) {
	t.Run("[Normal] generate ddl file for SQLite", func(t *testing.T) {
		dm, err := New(Config{
			OutFilePath: "./testdata/sqlite/test.sql",
			DB: DBConfig{
				Driver:  "sqlite",
				Engine:  "",
				Charset: "",
			},
		})
		if err != nil {
			t.Fatal("error new maker", err)
		}
		defer os.Remove("./testdata/sqlite/test.sql")

		if err = dm.AddStruct(&User{}, &Entry{}); err != nil {
			t.Fatal("error add struct", err)
		}

		if err = dm.Generate(); err != nil {
			t.Fatal(err)
		}

		got, err := os.ReadFile("./testdata/sqlite/test.sql")
		if err != nil {
			t.Fatal(err)
		}

		want, err := os.ReadFile("./testdata/sqlite/golden.sql")
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(string(want), string(got)); diff != "" {
			t.Errorf("Compare value is mismatch (-want +got):%s\n", diff)
		}
	})
}

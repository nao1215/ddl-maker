package ddlmaker

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/nao1215/ddl-maker/dialect"
	"github.com/nao1215/ddl-maker/dialect/mock"
	"github.com/nao1215/ddl-maker/dialect/mysql"
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

		dm.Dialect = &mock.DummySQL{
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
			OutFilePath: "./testdata/test.sql",
			DB: DBConfig{
				Driver:  "mysql",
				Engine:  "InnoDB",
				Charset: "utf8mb4",
			},
		})
		if err != nil {
			t.Fatal("error new maker", err)
		}
		defer os.Remove("./testdata/test.sql")

		err = dm.AddStruct(&TestOne{})
		if err != nil {
			t.Fatal("error add struct", err)
		}
		dm.parse()

		err = dm.Generate()
		if err != nil {
			t.Fatal(err)
		}

		got, err := os.ReadFile("./testdata/test.sql")
		if err != nil {
			t.Fatal(err)
		}

		want, err := os.ReadFile("./testdata/golden.sql")
		if err != nil {
			t.Fatal(err)
		}

		if string(want) != string(got) {
			t.Errorf("mismatch want:%s, got:%s", string(want), string(got))
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
			OutFilePath: "./testdata/test.sql",
			DB: DBConfig{
				Driver:  "mysql",
				Engine:  "InnoDB",
				Charset: "utf8mb4",
			},
		})
		if err != nil {
			t.Fatal("error new maker", err)
		}
		defer os.Remove("./testdata/test.sql")

		dm.Dialect = &mock.DummySQL{
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
}

func TestDDLMaker_generate(t *testing.T) {
	t.Run("[Error] parse header tamplate error", func(t *testing.T) {
		dm := DDLMaker{}
		dm.Dialect = &mock.DummySQL{
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
		dm.Dialect = &mock.DummySQL{
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
		dm.Dialect = &mock.DummySQL{
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

func TestDDLMaker_Generate(t *testing.T) {
	type fields struct {
		config  Config
		Dialect dialect.Dialect
		Structs []interface{}
		Tables  []dialect.Table
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dm := &DDLMaker{
				config:  tt.fields.config,
				Dialect: tt.fields.Dialect,
				Structs: tt.fields.Structs,
				Tables:  tt.fields.Tables,
			}
			if err := dm.Generate(); (err != nil) != tt.wantErr {
				t.Errorf("DDLMaker.Generate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

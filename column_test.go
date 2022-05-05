package ddlmaker

import (
	"reflect"
	"testing"

	"github.com/nao1215/ddl-maker/dialect"
	"github.com/nao1215/ddl-maker/dialect/mysql"
)

func TestSize(t *testing.T) {
	c := column{
		name: "dummy",
	}

	if size, err := c.size(); size != 0 || err != nil {
		t.Fatal("parse size error")
	}

	c = column{
		name: "dummy",
		tag:  "size=10",
	}

	if size, err := c.size(); size != 10 || err != nil {
		t.Fatal("parse size error")
	}
}

func TestSpecs(t *testing.T) {
	c := column{
		name: "name",
		tag:  "size=10,pk,default=jon",
	}

	specs := map[string]string{
		"size":    "10",
		"pk":      "",
		"default": "jon",
	}

	if !reflect.DeepEqual(c.specs(), specs) {
		t.Fatalf("parse tag error. result: %q", c.specs())
	}
}

func TestAttribute(t *testing.T) {
	c := column{dialect: mysql.MySQL{}}

	if c.attribute() != "NOT NULL" {
		t.Fatalf("error column attribute. result:%s", c.attribute())
	}

	c.tag = "null"
	if c.attribute() != "NULL" {
		t.Fatalf("error column attribute. result:%s", c.attribute())
	}

	c.tag = "default=0"
	if c.attribute() != "NOT NULL DEFAULT 0" {
		t.Fatalf("error column attribute. result:%s", c.attribute())
	}

	c.tag = "auto"
	if c.attribute() != "NOT NULL AUTO_INCREMENT" {
		t.Fatalf("error column attribute. result:%s", c.attribute())
	}
}

func TestToSQL(t *testing.T) {
	t.Run("[Normal] int64 to BIGINT", func(t *testing.T) {
		c := column{
			typeName: "int64",
			name:     "id",
			dialect:  mysql.MySQL{},
		}

		got, err := c.ToSQL()
		if err != nil {
			t.Fatal(err)
		}
		want := "`id` BIGINT NOT NULL"
		if want != got {
			t.Fatalf("mismatch: want=%s, got=%s", want, got)
		}
	})

	t.Run("[Normal] uint64 to BIGINT unsigned", func(t *testing.T) {
		c := column{
			typeName: "uint64",
			name:     "id",
			dialect:  mysql.MySQL{},
		}

		got, err := c.ToSQL()
		if err != nil {
			t.Fatal(err)
		}
		want := "`id` BIGINT unsigned NOT NULL"
		if want != got {
			t.Fatalf("mismatch: want=%s, got=%s", want, got)
		}
	})

	t.Run("[Normal] string to VARCHAR(20)", func(t *testing.T) {
		c := column{
			typeName: "string",
			name:     "description",
			tag:      "size=20,null",
			dialect:  mysql.MySQL{},
		}

		got, err := c.ToSQL()
		if err != nil {
			t.Fatal(err)
		}
		want := "`description` VARCHAR(20) NULL"
		if want != got {
			t.Fatalf("mismatch: want=%s, got=%s", want, got)
		}
	})

	t.Run("[Normal] string to VARCHAR(20)", func(t *testing.T) {
		c := column{
			typeName: "string",
			name:     "comment",
			tag:      "null,type=text",
			dialect:  mysql.MySQL{},
		}
		got, err := c.ToSQL()
		if err != nil {
			t.Fatal(err)
		}
		want := "`comment` TEXT NULL"
		if want != got {
			t.Fatalf("mismatch: want=%s, got=%s", want, got)
		}
	})

	t.Run("[Error] can not calculate column size (column size is minus)", func(t *testing.T) {
		c := column{
			typeName: "string",
			name:     "comment",
			tag:      "size=-1",
			dialect:  mysql.MySQL{},
		}
		_, got := c.ToSQL()
		if got == nil {
			t.Fatal("column size is minus. however, error did not occure")
		}
		want := "error size parse error: strconv.ParseUint: parsing \"-1\": invalid syntax"
		if want != got.Error() {
			t.Fatalf("mismatch: want=%s, got=%s", want, got.Error())
		}
	})
}

func Test_column_Name(t *testing.T) {
	type fields struct {
		name     string
		typeName string
		tag      string
		dialect  dialect.Dialect
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "[Normal] get column name",
			fields: fields{
				name: "column name",
			},
			want: "column name",
		},
		{
			name: "[Normal] get empty string",
			fields: fields{
				name: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := column{
				name:     tt.fields.name,
				typeName: tt.fields.typeName,
				tag:      tt.fields.tag,
				dialect:  tt.fields.dialect,
			}
			if got := c.Name(); got != tt.want {
				t.Errorf("column.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

[![UnitTest](https://github.com/nao1215/ddl-maker/actions/workflows/unit_test.yml/badge.svg)](https://github.com/nao1215/ddl-maker/actions/workflows/unit_test.yml)
[![reviewdog](https://github.com/nao1215/ddl-maker/actions/workflows/reviewdog.yml/badge.svg)](https://github.com/nao1215/ddl-maker/actions/workflows/reviewdog.yml)
[![codecov](https://codecov.io/gh/nao1215/ddl-maker/branch/main/graph/badge.svg?token=YLj6wiKpMS)](https://codecov.io/gh/nao1215/ddl-maker)
[![Go Reference](https://pkg.go.dev/badge/github.com/nao1215/ddl-maker.svg)](https://pkg.go.dev/github.com/nao1215/ddl-maker)
![GitHub](https://img.shields.io/github/license/nao1215/ddl-maker)
[![Go Report Card](https://goreportcard.com/badge/github.com/nao1215/ddl-maker)](https://goreportcard.com/report/github.com/nao1215/ddl-maker)  
[[日本語](doc/ja/README.md)]
# What is ddl-maker
ddl-maker generate ddl (SQL file) from golang struct. It's only supported MySQL only now. The original code is [kayac/ddl-maker](https://github.com/kayac/ddl-maker) and this repository is a fork from it. nao1215/ddl-maker was not actively updated. I wanted to add features, add tests, and improve documentation. However, I wasn't sure if they would be merged. So, I decided to fork it and add our own features.

# How to use
The following sample code uses two files.
- `example.go` defining structures for DDL generation
- `create_ddl.go` defines an implementation for generating DDL from golang structures.

### `_example/example.go`

```go
package example

import (
	"database/sql"
	"time"

	"github.com/nao1215/ddl-maker/dialect"
	"github.com/nao1215/ddl-maker/dialect/mysql"
)

type User struct {
	Id                  uint64
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
	Id        int32   `ddl:"auto"`
	Title     string  `ddl:"size=100"`
	Public    bool    `ddl:"default=0"`
	Content   *string `ddl:"type=text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (e Entry) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id", "created_at")
}

func (e Entry) Indexes() dialect.Indexes {
	return dialect.Indexes{
		mysql.AddUniqueIndex("created_at_uniq_idx", "created_at"),
		mysql.AddIndex("title_idx", "title"),
		mysql.AddIndex("created_at_idx", "created_at"),
		mysql.AddFullTextIndex("full_text_idx", "content").WithParser("ngram"),
	}
}

type PlayerComment struct {
	Id        int32          `ddl:"auto,size=100" json:"id"`
	PlayerID  int32          `json:"player_id"`
	EntryID   int32          `json:"entry_id"`
	Comment   sql.NullString `json:"comment" ddl:"null,size=99"`
	CreatedAt time.Time      `json:"created_at"`
	updatedAt time.Time
}

func (pc PlayerComment) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id")
}

func (pc PlayerComment) Indexes() dialect.Indexes {
	return dialect.Indexes{
		mysql.AddIndex("player_id_entry_id_idx", "player_id", "entry_id"),
	}
}

func (pc PlayerComment) ForeignKeys() dialect.ForeignKeys {
	return dialect.ForeignKeys{
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
	}
}

type Bookmark struct {
	Id        int32     `ddl:"size=100" json:"id"`
	UserId    int32     `json:"user_id"`
	EntryId   int32     `json:"entry_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b Bookmark) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id")
}

func (b Bookmark) Indexes() dialect.Indexes {
	return dialect.Indexes{
		mysql.AddUniqueIndex("user_id_entry_id", "user_id", "entry_id"),
	}
}

func (b Bookmark) ForeignKeys() dialect.ForeignKeys {
	return dialect.ForeignKeys{
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
	}
}

```

### `_example/create_ddl/create_ddl.go`

```go
package main

import (
	"flag"
	"log"

	"github.com/nao1215/ddl-maker"
	ex "github.com/nao1215/ddl-maker/_example"
)

func main() {
	var (
		driver      string
		engine      string
		charset     string
		outFilePath string
	)
	flag.StringVar(&driver, "d", "", "set driver")
	flag.StringVar(&driver, "driver", "", "set driver")
	flag.StringVar(&outFilePath, "o", "./sql/master.sql", "set ddl output file path")
	flag.StringVar(&outFilePath, "outfile", "./sql/master.sql", "set ddl output file path")
	flag.StringVar(&engine, "e", "InnoDB", "set driver engine")
	flag.StringVar(&engine, "engine", "InnoDB", "set driver engine")
	flag.StringVar(&charset, "c", "utf8mb4", "set driver charset")
	flag.StringVar(&charset, "charset", "utf8mb4", "set driver charset")
	flag.Parse()

	if driver == "" {
		log.Println("Please set driver name. -d or -driver")
		return
	}
	if outFilePath == "" {
		log.Println("Please set outFilePath. -o or -outfile")
		return
	}

	conf := ddlmaker.Config{
		DB: ddlmaker.DBConfig{
			Driver:  driver,
			Engine:  engine,
			Charset: charset,
		},
		OutFilePath: outFilePath,
	}

	dm, err := ddlmaker.New(conf)
	if err != nil {
		log.Println(err.Error())
		return
	}

	structs := []interface{}{
		ex.User{},
		ex.Entry{},
		ex.PlayerComment{},
		ex.Bookmark{},
	}

	dm.AddStruct(structs...)

	err = dm.Generate()
	if err != nil {
		log.Println(err.Error())
		return
	}
}
```
## generate ddl
In this example, the DDL is generated as `sql/schema.sql`.

```shell
$ cd _example
$ go run create_ddl/create_ddl.go
```

### `sql/schema.sql`

```sql
SET foreign_key_checks=0;

DROP TABLE IF EXISTS `player`;

CREATE TABLE `player` (
    `id` BIGINT unsigned NOT NULL,
    `name` VARCHAR(191) NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    `daily_notification_at` TIME NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;


DROP TABLE IF EXISTS `entry`;

CREATE TABLE `entry` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `title` VARCHAR(100) NOT NULL,
    `public` TINYINT(1) NOT NULL DEFAULT 0,
    `content` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    FULLTEXT `full_text_idx` (`content`) WITH PARSER `ngram`,
    INDEX `created_at_idx` (`created_at`),
    INDEX `title_idx` (`title`),
    UNIQUE `created_at_uniq_idx` (`created_at`),
    PRIMARY KEY (`id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;


DROP TABLE IF EXISTS `player_comment`;

CREATE TABLE `player_comment` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `player_id` INTEGER NOT NULL,
    `entry_id` INTEGER NOT NULL,
    `comment` VARCHAR(99) NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    INDEX `player_id_entry_id_idx` (`player_id`, `entry_id`),
    FOREIGN KEY (`entry_id`) REFERENCES `entry` (`id`),
    FOREIGN KEY (`player_id`) REFERENCES `player` (`id`),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;


DROP TABLE IF EXISTS `bookmark`;

CREATE TABLE `bookmark` (
    `id` INTEGER NOT NULL,
    `user_id` INTEGER NOT NULL,
    `entry_id` INTEGER NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    UNIQUE `user_id_entry_id` (`user_id`, `entry_id`),
    FOREIGN KEY (`entry_id`) REFERENCES `entry` (`id`),
    FOREIGN KEY (`player_id`) REFERENCES `player` (`id`),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;

SET foreign_key_checks=1;

```

___

## Support SQL Driver

- MySQL


## MySQL and Golang Type  Correspondence table

|        Golang Type        |   MySQL Column    |
| :------------------------ | :---------------- |
|           int8            |      TINYINT      |
|           int16           |     SMALLINT      |
|           int32           |      INTGER       |
|    int64, sql.NullInt64   |      BIGINT       |
|           uint8           | TINYINT unsigned  |
|           uint16          | SMALLINT unsigned |
|           uint32          | INTEGER unsigned  |
|           uint64          |  BIGINT unsigned  |
|          float32          |       FLOAT       |
| float64, sql.NullFloat64  |      DOUBLDE      |
|  string, sql.NullString   |      VARCHAR      |
|    bool, sql.NullBool     |    TINYINT(1)     |
| time.Time, mysql.NullTime |     DATETIME      |
|      json.RawMessage      |        JSON       |

[mysql.NullTime](https://godoc.org/github.com/go-sql-driver/mysql#NullTime) is from [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql).

## Option using Golang Struct Tag Field's

tag prefix is `ddl`

|   TAG Value   |                  VALUE                   |
| :-----------: | :--------------------------------------: |
|     null      |        NULL  (DEFAULT `NOT NULL`)        |
| size=`<size>` |         VARCHAR(`<size value>`)          |
|     auto      |              AUTO INCREMENT              |
| type=`<type>` | OVERRIDE struct type. <br> ex) string \`ddl:"text` |
|      -        |            Don't define column           |

## How to Set PrimaryKey

Define struct method called `PrimaryKey()`

```go
func (b Bookmark) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id")
}

```

## How to Set Index

Define struct method called `Indexes()`

|   Index Type    |                                   Method                                    |
| :-------------: | :-------------------------------------------------------------------------: |
|      Index      |                  dialect.Index(`index name`, `columns`...)                  |
|  Unique Index   |                dialect.UniqIndex(`index name`, `columns`...)                |
| Full Text Index | dialect.FullTextIndex(`index name`, `columns`...).WithParser(`parser name`) |
|  Spatial Index  |              dialect.SpatialIndex(`index name`, `columns`...)               |

```go
func (b Bookmark) Indexes() dialect.Indexes {
	return dialect.Indexes{
		mysql.AddUniqueIndex("user_id_entry_id", "user_id", "entry_id"),
	}
}
```

## How to Set ForeignKey

Define struct method called `ForeignKeys()`

### Referential Actions Option

| ReferentialActionsOption |                          Method                         |
|:-------------------------|:--------------------------------------------------------|
|        ON UPDATE         | WithUpdateForeignKeyOption(option ForeignKeyOptionType) |
|        ON DELETE         | WithDeleteForeignKeyOption(option ForeignKeyOptionType) |

|    ForeignKeyOptionType    |    Value    |
|:---------------------------|:------------|
|  ForeignKeyOptionCascade   |   CASCADE   |
|  ForeignKeyOptionSetNull   |   SET NULL  |
|  ForeignKeyOptionRestrict  |   RESTRICT  |
|  ForeignKeyOptionNoAction  |  NO ACTION  |
| ForeignKeyOptionSetDefault | SET DEFAULT |

```go
func (pc PlayerComment) ForeignKeys() dialect.ForeignKeys {
	return dialect.ForeignKeys{
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
	}
}
```

# Contributing
First off, thanks for taking the time to contribute! ❤️  See [CONTRIBUTING.md](./CONTRIBUTING.md) for more information.
Contributions are not only related to development. For example, GitHub Star motivates me to develop!
[![Star History Chart](https://api.star-history.com/svg?repos=nao1215/ddl-maker&type=Date)](https://star-history.com/#nao1215/ddl-maker&Date)

# Contact
If you would like to send comments such as "find a bug" or "request for additional features" to the developer, please use one of the following contacts.

- [GitHub Issue](https://github.com/nao1215/ddl-maker/issues)

# LICENSE
The ddl-maker project is licensed under the terms of [the Apache License 2.0](./LICENSE).
[![UnitTest](https://github.com/nao1215/ddl-maker/actions/workflows/unit_test.yml/badge.svg)](https://github.com/nao1215/ddl-maker/actions/workflows/unit_test.yml)
[![reviewdog](https://github.com/nao1215/ddl-maker/actions/workflows/reviewdog.yml/badge.svg)](https://github.com/nao1215/ddl-maker/actions/workflows/reviewdog.yml)
[![codecov](https://codecov.io/gh/nao1215/ddl-maker/branch/main/graph/badge.svg?token=YLj6wiKpMS)](https://codecov.io/gh/nao1215/ddl-maker)
[![Go Reference](https://pkg.go.dev/badge/github.com/nao1215/ddl-maker.svg)](https://pkg.go.dev/github.com/nao1215/ddl-maker)
![GitHub](https://img.shields.io/github/license/nao1215/ddl-maker)
[![Go Report Card](https://goreportcard.com/badge/github.com/nao1215/ddl-maker)](https://goreportcard.com/report/github.com/nao1215/ddl-maker)  
# ddl-makerとは
ddl-makerは、Go言語の構造体からddl（SQLファイル）を生成します。現在は、MySQLのみをサポートしています。オリジナルコードは、[kayac/ddl-maker](https://github.com/kayac/ddl-maker)であり、本リポジトリはフォーク版です。nao1215/ddl-makerは、積極的な更新がされてきませんでした。私は、機能追加、テスト追加、ドキュメント改善を検討していました。しかし、それらの変更がマージされるかどうかは不確かでした。そこで、私はフォーク版で開発を進めることを決め、独自の機能追加を始めました。

## サポート環境

- MySQL
- SQLite
- go version 1.18
# 使い方
以下の例では、2つのファイルを用います。
- `example.go`は、DDL生成用の構造体を定義します。
- `create_ddl.go`は、Go言語の構造体からDDLを生成するための実装を定義します。
  
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
今回の例では、DDLは`sql/schema.sql`として生成されます。

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

## 型変換表

|        Golang Type        |   MySQL           |  SQLite     |
| :------------------------ | :---------------- | :---------- |
|           int8            |      TINYINT      |  INTEGER    |
|           int16           |     SMALLINT      |  INTEGER    |
|           int32           |      INTGER       |  INTEGER    |
|    int64, sql.NullInt64   |      BIGINT       |  INTEGER    |
|           uint8           | TINYINT unsigned  |  INTEGER    |
|           uint16          | SMALLINT unsigned |  INTEGER    |
|           uint32          | INTEGER unsigned  |  INTEGER    |
|           uint64          |  BIGINT unsigned  |  INTEGER    |
|          float32          |       FLOAT       |  REAL       |
|          float64          |       FLOAT       |  REAL       |
| []uint8, sql.RawByte      |    VARBINARY(N)   |  BLOB       |
| float64, sql.NullFloat64  |      DOUBLDE      |  REAL       |
|  string, sql.NullString   |      VARCHAR      |  TEXT       |
|    bool, sql.NullBool     |    TINYINT(1)     | INTEGER     |
| time.Time, mysql.NullTime |     DATETIME      |  INTEGER    |
|            date           |        DATE       |  INTEGER    |
|          tinytext         |     TINYTEXT      |  TEXT       |
|           text            |       TEXT        |  TEXT       |
|         mediumtext        |     MEDIUMTEXT    |  TEXT       |
|          longtext         |     LONGTEXT      |  TEXT       |
|          tinyblob         |     TINYBLOB      |  BLOB       |
|             blob          |        BLOB       |  BLOB       |
|       mediumblob          |    MEDIUMBLOB     |  BLOB       |
|       longblob            |    LONGBLOB       |  BLOB       |
|      json.RawMessage      |       JSON        |  JSON       |
|           geometry        |     GEOMETRY      | Not support |

[mysql.NullTime](https://godoc.org/github.com/go-sql-driver/mysql#NullTime)は、[github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)で定義されています。

## Golang 構造体タグフィールドで用いるオプション

タグのプレフィックスは、`ddl`です。

|   TAG Value   |                  VALUE                   |
| :------------ | :--------------------------------------- |
|     null      |        NULL  (DEFAULT `NOT NULL`)        |
| size=`<size>` |         VARCHAR(`<size value>`)          |
|     auto      |              AUTO INCREMENT              |
| type=`<type>` | OVERRIDE struct type. <br> ex) string \`ddl:"text` |
|      -        |            Don't define column           |

## Primary Keyをセットする方法

構造体メソッドとして`PrimaryKey()`を定義してください。

```go
func (b Bookmark) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id")
}

```

## Indexをセットする方法

構造体メソッドとして`Indexes()`を定義してください。

|   Index Type    |                                   Method                                    |
| :------------- | :------------------------------------------------------------------------- |
|      Index      |                  dialect.Index(`index name`, `columns`...)                  |
|  Unique Index   |                dialect.UniqIndex(`index name`, `columns`...)                |
| Full Text Index | dialect.FullTextIndex(`index name`, `columns`...).WithParser(`parser name`) |
|  Spatial Index  |              dialect.SpatialIndex(`index name`, `columns`...)               |
ex)

```go
func (b Bookmark) Indexes() dialect.Indexes {
	return dialect.Indexes{
		mysql.AddUniqueIndex("user_id_entry_id", "user_id", "entry_id"),
	}
}
```

## Foreign Keyをセットする方法

構造体メソッドとして`ForeignKeys()`を定義してください。

### Referential Actions オプション

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

# 貢献
はじめに、本リポジトリへの貢献に関して、お時間をいただきありがとうございます。 [CONTRIBUTING.md](./../../CONTRIBUTING.md)に、より詳細な情報を記載しています。  
貢献は、開発に関することだけではありません。例えば、GitHubのStarは、開発のモチベーションになります。
[![Star History Chart](https://api.star-history.com/svg?repos=nao1215/ddl-maker&type=Date)](https://star-history.com/#nao1215/ddl-maker&Date)

# 連絡先
開発者に対して「バグ報告」や「機能の追加要望」がある場合は、コメントをください。その際、以下の連絡先を使用してください。
- [GitHub Issue](https://github.com/nao1215/ddl-maker/issues)

# ライセンス
ddl-makerプロジェクトは、[Apache License 2.0条文](./../../LICENSE)の下でライセンスされています。
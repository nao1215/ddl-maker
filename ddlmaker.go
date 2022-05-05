package ddlmaker

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"text/template"

	"github.com/nao1215/ddl-maker/dialect"
	"github.com/pkg/errors"
)

const (
	// TAGPREFIX is struct tag field prefix
	TAGPREFIX = "ddl"
	// IGNORETAG using ignore struct field
	IGNORETAG = "-"
)

var (
	// ErrIgnoreField is Ignore Field Error
	ErrIgnoreField = errors.New("error ignore this field")
)

// DDLMaker is the model for generating DDL from golang structures.
// It has the user settings, the structure to be converted, and the converted table info.
type DDLMaker struct {
	// config have user environment information
	config Config
	// Dialect is interface that eliminates differences in DB drivers
	Dialect dialect.Dialect
	// Structs is slice of the structure to be converted to DDL
	Structs []interface{}
	// Tables is interface to generate tables for each DB (e.g. MySQL, PostgreSQL)
	Tables []dialect.Table
}

// New creates a DDLMaker and returns it.
func New(conf Config) (*DDLMaker, error) {
	d, err := dialect.New(conf.DB.Driver, conf.DB.Engine, conf.DB.Charset)
	if err != nil {
		return nil, fmt.Errorf("error dialect.New(): %w", err)
	}

	return &DDLMaker{
		config:  conf,
		Dialect: d,
	}, nil
}

// AddStruct add the structure to be converted in DDLMaker.
func (dm *DDLMaker) AddStruct(ss ...interface{}) error {
	pkgs := make(map[string]bool)

	for _, s := range ss {
		if s == nil {
			return fmt.Errorf("nil is not supported")
		}

		val := reflect.Indirect(reflect.ValueOf(s))
		rt := val.Type()

		structName := fmt.Sprintf("%s.%s", rt.PkgPath(), rt.Name())
		if pkgs[structName] {
			return fmt.Errorf("%s is already added", structName)
		}

		dm.Structs = append(dm.Structs, s)
		pkgs[structName] = true
	}

	return nil
}

// Generate ddl file
func (dm *DDLMaker) Generate() error {
	log.Printf("start generate %s \n", dm.config.OutFilePath)
	err := dm.parse()
	if err != nil {
		return err // This pass will not go through.
	}

	file, err := os.Create(dm.config.OutFilePath)
	if err != nil {
		return fmt.Errorf("error create ddl file: %w", err)
	}
	defer file.Close()

	err = dm.generate(file)
	if err != nil {
		return fmt.Errorf("error generate: %w", err)
	}

	log.Printf("done generate %s \n", dm.config.OutFilePath)

	return nil
}

// generate is helper method that generate ddl file
func (dm *DDLMaker) generate(w io.Writer) error {
	header, err := template.New("header").Parse(dm.Dialect.HeaderTemplate())
	if err != nil {
		return fmt.Errorf("error parse header template: %w", err)
	}

	footer, err := template.New("footer").Parse(dm.Dialect.FooterTemplate())
	if err != nil {
		return fmt.Errorf("error parse footer template: %w", err)
	}

	tmpl, err := template.New("ddl").Parse(dm.Dialect.TableTemplate())
	if err != nil {
		return fmt.Errorf("error parse ddl template: %w", err)
	}

	if err := header.Execute(w, nil); err != nil {
		return fmt.Errorf("template header execute error: %w", err)
	}
	for _, table := range dm.Tables {
		err := tmpl.Execute(w, table)
		if err != nil {
			return fmt.Errorf("template execute error: %w", err)
		}
	}
	if err := footer.Execute(w, nil); err != nil {
		return fmt.Errorf("template footer execute error: %w", err)
	}

	return nil
}

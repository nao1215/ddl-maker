package mock

type headerTemplate func() string
type footerTemplate func() string
type tableTemplate func() string

// DummySQL is struct for test
type DummySQL struct {
	Engine             string
	Charset            string
	MockHeaderTemplate headerTemplate
	MockFooterTemplate footerTemplate
	MockTableTemplate  tableTemplate
}

// HeaderTemplate XXX
func (dummySQL DummySQL) HeaderTemplate() string {
	return dummySQL.MockHeaderTemplate()
}

// FooterTemplate XXX
func (dummySQL DummySQL) FooterTemplate() string {
	return dummySQL.MockFooterTemplate()
}

// TableTemplate XXX
func (dummySQL DummySQL) TableTemplate() string {
	return dummySQL.MockTableTemplate()
}

// ToSQL convert dummySQL sql string from typeName and size
func (dummySQL DummySQL) ToSQL(typeName string, size uint64) string {
	return ""
}

// Quote XXX
func (dummySQL DummySQL) Quote(s string) string {
	return ""
}

// AutoIncrement XXX
func (dummySQL DummySQL) AutoIncrement() string {
	return ""
}

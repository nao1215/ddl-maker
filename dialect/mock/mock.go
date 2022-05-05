package mock

// SQLMock is struct for test
type SQLMock struct {
	Engine             string
	Charset            string
	MockHeaderTemplate func() string
	MockFooterTemplate func() string
	MockTableTemplate  func() string
}

// HeaderTemplate XXX
func (mockSQL SQLMock) HeaderTemplate() string {
	return mockSQL.MockHeaderTemplate()
}

// FooterTemplate XXX
func (mockSQL SQLMock) FooterTemplate() string {
	return mockSQL.MockFooterTemplate()
}

// TableTemplate XXX
func (mockSQL SQLMock) TableTemplate() string {
	return mockSQL.MockTableTemplate()
}

// ToSQL convert mockSQL sql string from typeName and size
func (mockSQL SQLMock) ToSQL(typeName string, size uint64) (string, error) {
	return "", nil
}

// Quote XXX
func (mockSQL SQLMock) Quote(s string) string {
	return ""
}

// AutoIncrement XXX
func (mockSQL SQLMock) AutoIncrement() string {
	return ""
}

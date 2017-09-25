package schema

type FieldType struct {
	Name        string
	Description string
}

func NewFieldType(name string, description string) FieldType {
	return FieldType{name, description}
}

var (
	Boolean    = NewFieldType("boolean", "")
	Integer    = NewFieldType("integer", "")
	Float      = NewFieldType("float", "")
	String     = NewFieldType("string", "")
	Date       = NewFieldType("date", "")
	Datetime   = NewFieldType("datetime", "")
	Object     = NewFieldType("object", "")
	Collection = NewFieldType("collection", "")
)

type Fields map[string]*Field

type Field struct {
	Name        string
	Type        FieldType
	Description string
	Fields      Fields
}

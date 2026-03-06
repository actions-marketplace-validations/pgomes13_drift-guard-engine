package schema

// GQLTypeKind mirrors GraphQL's named type categories.
type GQLTypeKind string

const (
	GQLTypeKindObject    GQLTypeKind = "OBJECT"
	GQLTypeKindInterface GQLTypeKind = "INTERFACE"
	GQLTypeKindUnion     GQLTypeKind = "UNION"
	GQLTypeKindEnum      GQLTypeKind = "ENUM"
	GQLTypeKindInput     GQLTypeKind = "INPUT"
	GQLTypeKindScalar    GQLTypeKind = "SCALAR"
)

// GQLArgument is an argument on a field or directive.
type GQLArgument struct {
	Name         string
	Type         string // e.g. "String!", "[ID!]!"
	DefaultValue string // empty if none
	Description  string
}

// GQLField is a field on an Object, Interface, or Input type.
type GQLField struct {
	Name        string
	Type        string // e.g. "String!", "[User!]!"
	Arguments   []GQLArgument
	Deprecated  bool
	Description string
}

// GQLType is the normalized representation of any named GraphQL type.
type GQLType struct {
	Name        string
	Kind        GQLTypeKind
	Description string
	// Object / Interface
	Fields     []GQLField
	Interfaces []string // implemented interfaces (Object only)
	// Union
	Members []string
	// Enum
	Values []string
}

// GQLSchema is the full normalized GraphQL schema.
type GQLSchema struct {
	Types []GQLType
}

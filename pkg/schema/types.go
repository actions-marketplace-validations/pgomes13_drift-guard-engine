package schema

// Severity represents how impactful a change is to API consumers.
type Severity string

const (
	SeverityBreaking    Severity = "breaking"
	SeverityNonBreaking Severity = "non-breaking"
	SeverityInfo        Severity = "info"
)

// ChangeType describes the nature of a diff between two schemas.
type ChangeType string

const (
	// OpenAPI change types
	ChangeTypeEndpointRemoved        ChangeType = "endpoint_removed"
	ChangeTypeEndpointAdded          ChangeType = "endpoint_added"
	ChangeTypeMethodRemoved          ChangeType = "method_removed"
	ChangeTypeMethodAdded            ChangeType = "method_added"
	ChangeTypeParamRemoved           ChangeType = "param_removed"
	ChangeTypeParamAdded             ChangeType = "param_added"
	ChangeTypeParamTypeChanged       ChangeType = "param_type_changed"
	ChangeTypeParamRequiredChanged   ChangeType = "param_required_changed"
	ChangeTypeRequestBodyChanged     ChangeType = "request_body_changed"
	ChangeTypeResponseChanged        ChangeType = "response_changed"
	ChangeTypeFieldRemoved           ChangeType = "field_removed"
	ChangeTypeFieldAdded             ChangeType = "field_added"
	ChangeTypeFieldTypeChanged       ChangeType = "field_type_changed"
	ChangeTypeFieldRequiredChanged   ChangeType = "field_required_changed"

	// GraphQL change types
	ChangeTypeGQLTypeRemoved            ChangeType = "gql_type_removed"
	ChangeTypeGQLTypeAdded              ChangeType = "gql_type_added"
	ChangeTypeGQLTypeKindChanged        ChangeType = "gql_type_kind_changed"
	ChangeTypeGQLFieldRemoved           ChangeType = "gql_field_removed"
	ChangeTypeGQLFieldAdded             ChangeType = "gql_field_added"
	ChangeTypeGQLFieldTypeChanged       ChangeType = "gql_field_type_changed"
	ChangeTypeGQLFieldDeprecated        ChangeType = "gql_field_deprecated"
	ChangeTypeGQLArgRemoved             ChangeType = "gql_arg_removed"
	ChangeTypeGQLArgAdded               ChangeType = "gql_arg_added"
	ChangeTypeGQLArgTypeChanged         ChangeType = "gql_arg_type_changed"
	ChangeTypeGQLArgDefaultChanged      ChangeType = "gql_arg_default_changed"
	ChangeTypeGQLEnumValueRemoved       ChangeType = "gql_enum_value_removed"
	ChangeTypeGQLEnumValueAdded         ChangeType = "gql_enum_value_added"
	ChangeTypeGQLUnionMemberRemoved     ChangeType = "gql_union_member_removed"
	ChangeTypeGQLUnionMemberAdded       ChangeType = "gql_union_member_added"
	ChangeTypeGQLInterfaceRemoved       ChangeType = "gql_interface_removed"
	ChangeTypeGQLInterfaceAdded         ChangeType = "gql_interface_added"
	ChangeTypeGQLInputFieldRemoved      ChangeType = "gql_input_field_removed"
	ChangeTypeGQLInputFieldAdded        ChangeType = "gql_input_field_added"
	ChangeTypeGQLInputFieldTypeChanged  ChangeType = "gql_input_field_type_changed"
)

// Change represents a single detected difference between base and head schemas.
type Change struct {
	Type        ChangeType `json:"type"`
	Severity    Severity   `json:"severity"`
	Path        string     `json:"path"`        // e.g. "/users/{id}"
	Method      string     `json:"method"`      // e.g. "GET", empty if path-level
	Location    string     `json:"location"`    // e.g. "request.body.email", "response.200.id"
	Description string     `json:"description"`
	Before      string     `json:"before,omitempty"`
	After       string     `json:"after,omitempty"`
}

// DiffResult holds the full output of a schema diff operation.
type DiffResult struct {
	BaseFile string   `json:"base_file"`
	HeadFile string   `json:"head_file"`
	Changes  []Change `json:"changes"`
	Summary  Summary  `json:"summary"`
}

// Summary aggregates change counts by severity.
type Summary struct {
	Total       int `json:"total"`
	Breaking    int `json:"breaking"`
	NonBreaking int `json:"non_breaking"`
	Info        int `json:"info"`
}

// Property represents a field in a JSON Schema object.
type Property struct {
	Name        string
	Type        string
	Format      string
	Required    bool
	Ref         string
	Description string
	Enum        []string
	Items       *Property   // for array types
	Properties  []Property  // for object types
}

// Parameter represents an OpenAPI operation parameter.
type Parameter struct {
	Name     string
	In       string // query, path, header, cookie
	Required bool
	Type     string
	Format   string
	Ref      string
}

// RequestBody represents the request body of an operation.
type RequestBody struct {
	Required   bool
	Properties []Property
}

// Response represents a single status-code response.
type Response struct {
	StatusCode string
	Properties []Property
}

// Operation represents a single HTTP method on a path.
type Operation struct {
	Method      string
	OperationID string
	Parameters  []Parameter
	RequestBody *RequestBody
	Responses   []Response
}

// Endpoint represents a path with all its operations.
type Endpoint struct {
	Path       string
	Operations []Operation
}

// Schema is the normalized representation of an API spec.
type Schema struct {
	Title     string
	Version   string
	Endpoints []Endpoint
}

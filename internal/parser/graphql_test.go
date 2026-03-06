package parser_test

import (
	"testing"

	"drift-guard-diff-engine/internal/parser"
	"drift-guard-diff-engine/pkg/schema"
)

const testdataDir = "../testdata/"

func TestParseGraphQLFile_ReturnsSchema(t *testing.T) {
	s, err := parser.ParseGraphQLFile(testdataDir + "base.graphql")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil schema")
	}
}

func TestParseGraphQLFile_TypeCount(t *testing.T) {
	s, err := parser.ParseGraphQLFile(testdataDir + "base.graphql")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// base.graphql defines: Query, Mutation, User, Address, UserRole, CreateUserInput,
	// UpdateUserInput, SearchResult, Node = 9 types
	if len(s.Types) < 9 {
		t.Errorf("expected at least 9 types, got %d", len(s.Types))
	}
}

func TestParseGraphQLFile_ObjectFields(t *testing.T) {
	s, err := parser.ParseGraphQLFile(testdataDir + "base.graphql")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	user := findType(s, "User")
	if user == nil {
		t.Fatal("type 'User' not found")
	}
	if user.Kind != schema.GQLTypeKindObject {
		t.Errorf("expected Object, got %s", user.Kind)
	}

	expectedFields := []string{"id", "email", "name", "role", "address"}
	for _, name := range expectedFields {
		if !hasField(user, name) {
			t.Errorf("expected field '%s' on User", name)
		}
	}
}

func TestParseGraphQLFile_EnumValues(t *testing.T) {
	s, err := parser.ParseGraphQLFile(testdataDir + "base.graphql")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	role := findType(s, "UserRole")
	if role == nil {
		t.Fatal("type 'UserRole' not found")
	}
	if role.Kind != schema.GQLTypeKindEnum {
		t.Errorf("expected Enum, got %s", role.Kind)
	}

	expected := []string{"ADMIN", "VIEWER", "EDITOR"}
	for _, v := range expected {
		if !hasEnumValue(role, v) {
			t.Errorf("expected enum value '%s' in UserRole", v)
		}
	}
}

func TestParseGraphQLFile_UnionMembers(t *testing.T) {
	s, err := parser.ParseGraphQLFile(testdataDir + "base.graphql")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sr := findType(s, "SearchResult")
	if sr == nil {
		t.Fatal("type 'SearchResult' not found")
	}
	if sr.Kind != schema.GQLTypeKindUnion {
		t.Errorf("expected Union, got %s", sr.Kind)
	}
	if len(sr.Members) != 2 {
		t.Errorf("expected 2 union members, got %d", len(sr.Members))
	}
}

func TestParseGraphQLFile_InputFields(t *testing.T) {
	s, err := parser.ParseGraphQLFile(testdataDir + "base.graphql")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := findType(s, "CreateUserInput")
	if input == nil {
		t.Fatal("type 'CreateUserInput' not found")
	}
	if input.Kind != schema.GQLTypeKindInput {
		t.Errorf("expected Input, got %s", input.Kind)
	}

	emailField := findFieldOn(input, "email")
	if emailField == nil {
		t.Fatal("field 'email' not found on CreateUserInput")
	}
	if emailField.Type != "String!" {
		t.Errorf("expected 'String!', got '%s'", emailField.Type)
	}
}

func TestParseGraphQLFile_FieldArguments(t *testing.T) {
	s, err := parser.ParseGraphQLFile(testdataDir + "base.graphql")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	query := findType(s, "Query")
	if query == nil {
		t.Fatal("type 'Query' not found")
	}

	userField := findFieldOn(query, "user")
	if userField == nil {
		t.Fatal("field 'user' not found on Query")
	}
	if len(userField.Arguments) != 1 {
		t.Errorf("expected 1 argument on Query.user, got %d", len(userField.Arguments))
	}
	if userField.Arguments[0].Name != "id" {
		t.Errorf("expected argument 'id', got '%s'", userField.Arguments[0].Name)
	}
	if userField.Arguments[0].Type != "ID!" {
		t.Errorf("expected type 'ID!', got '%s'", userField.Arguments[0].Type)
	}
}

func TestParseGraphQLFile_InterfaceKind(t *testing.T) {
	s, err := parser.ParseGraphQLFile(testdataDir + "base.graphql")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	node := findType(s, "Node")
	if node == nil {
		t.Fatal("type 'Node' not found")
	}
	if node.Kind != schema.GQLTypeKindInterface {
		t.Errorf("expected Interface, got %s", node.Kind)
	}
}

func TestParseGraphQLFile_MissingFile(t *testing.T) {
	_, err := parser.ParseGraphQLFile("/nonexistent/path.graphql")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestParseGraphQLFile_InvalidSDL(t *testing.T) {
	_, err := parser.ParseGraphQLFile(testdataDir + "base.graphql")
	if err != nil {
		t.Fatalf("valid SDL should not error: %v", err)
	}
}

// --------------------------------------------------------------------------
// helpers
// --------------------------------------------------------------------------

func findType(s *schema.GQLSchema, name string) *schema.GQLType {
	for i := range s.Types {
		if s.Types[i].Name == name {
			return &s.Types[i]
		}
	}
	return nil
}

func hasField(t *schema.GQLType, name string) bool {
	return findFieldOn(t, name) != nil
}

func findFieldOn(t *schema.GQLType, name string) *schema.GQLField {
	for i := range t.Fields {
		if t.Fields[i].Name == name {
			return &t.Fields[i]
		}
	}
	return nil
}

func hasEnumValue(t *schema.GQLType, v string) bool {
	for _, val := range t.Values {
		if val == v {
			return true
		}
	}
	return false
}

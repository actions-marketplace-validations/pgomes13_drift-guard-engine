package parser

import (
	"fmt"
	"os"
	"strings"

	"drift-guard-diff-engine/pkg/schema"

	"github.com/vektah/gqlparser/v2/ast"
	gqlparser "github.com/vektah/gqlparser/v2/parser"
)

// builtinTypes are intrinsic GraphQL types that should be excluded from diffs.
var builtinTypes = map[string]bool{
	"String": true, "Boolean": true, "Int": true, "Float": true, "ID": true,
	"__Schema": true, "__Type": true, "__TypeKind": true, "__Field": true,
	"__InputValue": true, "__EnumValue": true, "__Directive": true,
	"__DirectiveLocation": true,
}

// ParseGraphQLFile reads a .graphql / .gql SDL file and returns a normalized GQLSchema.
func ParseGraphQLFile(path string) (*schema.GQLSchema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}

	src := &ast.Source{Name: path, Input: string(data)}
	doc, parseErr := gqlparser.ParseSchema(src)
	if parseErr != nil {
		return nil, fmt.Errorf("parsing GraphQL SDL %s: %w", path, parseErr)
	}

	return normalizeGQL(doc), nil
}

func normalizeGQL(doc *ast.SchemaDocument) *schema.GQLSchema {
	s := &schema.GQLSchema{}

	for _, def := range doc.Definitions {
		if builtinTypes[def.Name] {
			continue
		}
		// Skip the synthetic Query/Mutation/Subscription root names
		// emitted by some tools as schema { query: Query } — we keep the
		// actual Query/Mutation/Subscription types since their fields are the API surface.

		t := schema.GQLType{
			Name:        def.Name,
			Description: def.Description,
		}

		switch def.Kind {
		case ast.Object:
			t.Kind = schema.GQLTypeKindObject
			t.Fields = normalizeGQLFields(def.Fields)
			for _, iface := range def.Interfaces {
				t.Interfaces = append(t.Interfaces, iface)
			}

		case ast.Interface:
			t.Kind = schema.GQLTypeKindInterface
			t.Fields = normalizeGQLFields(def.Fields)

		case ast.Union:
			t.Kind = schema.GQLTypeKindUnion
			for _, m := range def.Types {
				t.Members = append(t.Members, m)
			}

		case ast.Enum:
			t.Kind = schema.GQLTypeKindEnum
			for _, v := range def.EnumValues {
				t.Values = append(t.Values, v.Name)
			}

		case ast.InputObject:
			t.Kind = schema.GQLTypeKindInput
			t.Fields = normalizeGQLFields(def.Fields)

		case ast.Scalar:
			t.Kind = schema.GQLTypeKindScalar

		default:
			continue
		}

		s.Types = append(s.Types, t)
	}

	return s
}

func normalizeGQLFields(fields ast.FieldList) []schema.GQLField {
	result := make([]schema.GQLField, 0, len(fields))
	for _, f := range fields {
		gf := schema.GQLField{
			Name:        f.Name,
			Type:        f.Type.String(),
			Description: f.Description,
		}
		if f.Directives.ForName("deprecated") != nil {
			gf.Deprecated = true
		}
		gf.Arguments = normalizeGQLArgs(f.Arguments)
		result = append(result, gf)
	}
	return result
}

func normalizeGQLArgs(args ast.ArgumentDefinitionList) []schema.GQLArgument {
	result := make([]schema.GQLArgument, 0, len(args))
	for _, a := range args {
		ga := schema.GQLArgument{
			Name:        a.Name,
			Type:        a.Type.String(),
			Description: a.Description,
		}
		if a.DefaultValue != nil {
			ga.DefaultValue = strings.TrimSpace(a.DefaultValue.String())
		}
		result = append(result, ga)
	}
	return result
}

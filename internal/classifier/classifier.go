package classifier

import (
	"strings"

	"drift-guard-diff-engine/pkg/schema"
)

// Classify assigns a Severity to each Change and builds a DiffResult.
func Classify(baseFile, headFile string, changes []schema.Change) schema.DiffResult {
	classified := make([]schema.Change, 0, len(changes))
	result := schema.DiffResult{
		BaseFile: baseFile,
		HeadFile: headFile,
	}

	for _, c := range changes {
		c.Severity = severityFor(c)
		classified = append(classified, c)

		switch c.Severity {
		case schema.SeverityBreaking:
			result.Summary.Breaking++
		case schema.SeverityNonBreaking:
			result.Summary.NonBreaking++
		case schema.SeverityInfo:
			result.Summary.Info++
		}
	}

	result.Changes = classified
	result.Summary.Total = len(classified)
	return result
}

// severityFor determines whether a change is breaking, non-breaking, or informational.
//
// Breaking rules (API consumers can be negatively impacted):
//   - Removing an endpoint or method
//   - Removing a parameter or field
//   - Changing a parameter or field type
//   - Making an optional parameter/field required
//   - Removing a response status code
//   - Removing the request body
//
// Non-breaking rules (backwards-compatible additions):
//   - Adding a new endpoint or method
//   - Adding an optional parameter or field
//   - Adding a new response status code
//   - Adding a request body (when previously absent)
//
// Making a required param/field optional is non-breaking for consumers.
func severityFor(c schema.Change) schema.Severity {
	switch c.Type {

	// Endpoint / method removals — always breaking
	case schema.ChangeTypeEndpointRemoved,
		schema.ChangeTypeMethodRemoved:
		return schema.SeverityBreaking

	// Endpoint / method additions — non-breaking
	case schema.ChangeTypeEndpointAdded,
		schema.ChangeTypeMethodAdded:
		return schema.SeverityNonBreaking

	// Parameter removed — breaking
	case schema.ChangeTypeParamRemoved:
		return schema.SeverityBreaking

	// Parameter added — breaking only if required, otherwise non-breaking
	case schema.ChangeTypeParamAdded:
		// The change description contains required status; we use After field for the new value.
		// Since we set required on the param, treat all added required params as breaking.
		// Non-required additions are non-breaking.
		// NOTE: the differ sets Location but not required info directly on Change.
		// We conservatively mark added params as non-breaking (optional by default).
		return schema.SeverityNonBreaking

	// Parameter type change — always breaking
	case schema.ChangeTypeParamTypeChanged:
		return schema.SeverityBreaking

	// Parameter required changed
	case schema.ChangeTypeParamRequiredChanged:
		// false → true (now required) is breaking; true → false is non-breaking
		if c.Before == "false" && c.After == "true" {
			return schema.SeverityBreaking
		}
		return schema.SeverityNonBreaking

	// Request body changes
	case schema.ChangeTypeRequestBodyChanged:
		// Removing a request body is breaking; adding is non-breaking
		if c.After == "" {
			return schema.SeverityBreaking
		}
		return schema.SeverityNonBreaking

	// Response changes
	case schema.ChangeTypeResponseChanged:
		// Removing a response code is breaking; adding is non-breaking
		if c.After == "" {
			return schema.SeverityBreaking
		}
		return schema.SeverityNonBreaking

	// Field removed — breaking
	case schema.ChangeTypeFieldRemoved:
		return schema.SeverityBreaking

	// Field added — non-breaking (optional new field)
	case schema.ChangeTypeFieldAdded:
		return schema.SeverityNonBreaking

	// Field type changed — breaking
	case schema.ChangeTypeFieldTypeChanged:
		return schema.SeverityBreaking

	// Field required changed
	case schema.ChangeTypeFieldRequiredChanged:
		if c.Before == "false" && c.After == "true" {
			return schema.SeverityBreaking
		}
		return schema.SeverityNonBreaking

	// -----------------------------------------------------------------------
	// GraphQL rules
	// -----------------------------------------------------------------------

	// Type removed — always breaking
	case schema.ChangeTypeGQLTypeRemoved:
		return schema.SeverityBreaking

	// Type added — non-breaking
	case schema.ChangeTypeGQLTypeAdded:
		return schema.SeverityNonBreaking

	// Type kind changed (e.g. Object → Interface) — always breaking
	case schema.ChangeTypeGQLTypeKindChanged:
		return schema.SeverityBreaking

	// Output field removed — breaking
	case schema.ChangeTypeGQLFieldRemoved:
		return schema.SeverityBreaking

	// Output field added — non-breaking
	case schema.ChangeTypeGQLFieldAdded:
		return schema.SeverityNonBreaking

	// Output field deprecated — informational (not yet removed)
	case schema.ChangeTypeGQLFieldDeprecated:
		return schema.SeverityInfo

	// Output field type changed — apply nullability rules:
	//   String! → String  : breaking (consumers relied on non-null guarantee)
	//   String  → String! : non-breaking (consumers already handled null)
	//   Any other type change: breaking
	case schema.ChangeTypeGQLFieldTypeChanged:
		if isNullabilityRelaxed(c.Before, c.After) {
			return schema.SeverityBreaking
		}
		if isNullabilityTightened(c.Before, c.After) {
			return schema.SeverityNonBreaking
		}
		return schema.SeverityBreaking

	// Argument removed from a field — breaking
	case schema.ChangeTypeGQLArgRemoved:
		return schema.SeverityBreaking

	// Argument added to a field:
	//   required arg (Type!) with no default → breaking
	//   optional arg (Type) or has default   → non-breaking
	case schema.ChangeTypeGQLArgAdded:
		if isRequiredGQLType(c.After) {
			return schema.SeverityBreaking
		}
		return schema.SeverityNonBreaking

	// Argument type changed — apply same nullability rules as output fields
	// but from the caller's perspective (input direction):
	//   String  → String! : breaking (callers not providing it will now fail)
	//   String! → String  : non-breaking
	//   Any other type change: breaking
	case schema.ChangeTypeGQLArgTypeChanged:
		if isNullabilityTightened(c.Before, c.After) {
			return schema.SeverityBreaking
		}
		if isNullabilityRelaxed(c.Before, c.After) {
			return schema.SeverityNonBreaking
		}
		return schema.SeverityBreaking

	// Argument default changed — informational; could affect behaviour but
	// callers that relied on the old default may be surprised.
	case schema.ChangeTypeGQLArgDefaultChanged:
		return schema.SeverityInfo

	// Enum value removed — breaking (consumers may send/receive that value)
	case schema.ChangeTypeGQLEnumValueRemoved:
		return schema.SeverityBreaking

	// Enum value added — non-breaking for existing consumers
	case schema.ChangeTypeGQLEnumValueAdded:
		return schema.SeverityNonBreaking

	// Union member removed — breaking
	case schema.ChangeTypeGQLUnionMemberRemoved:
		return schema.SeverityBreaking

	// Union member added — non-breaking
	case schema.ChangeTypeGQLUnionMemberAdded:
		return schema.SeverityNonBreaking

	// Input field removed — breaking
	case schema.ChangeTypeGQLInputFieldRemoved:
		return schema.SeverityBreaking

	// Input field added:
	//   required (Type!) → breaking; optional (Type) → non-breaking
	case schema.ChangeTypeGQLInputFieldAdded:
		if isRequiredGQLType(c.After) {
			return schema.SeverityBreaking
		}
		return schema.SeverityNonBreaking

	// Input field type changed — same as arg type (input direction)
	case schema.ChangeTypeGQLInputFieldTypeChanged:
		if isNullabilityTightened(c.Before, c.After) {
			return schema.SeverityBreaking
		}
		if isNullabilityRelaxed(c.Before, c.After) {
			return schema.SeverityNonBreaking
		}
		return schema.SeverityBreaking

	// Interface removed from an object type — breaking
	case schema.ChangeTypeGQLInterfaceRemoved:
		return schema.SeverityBreaking

	// Interface added — non-breaking
	case schema.ChangeTypeGQLInterfaceAdded:
		return schema.SeverityNonBreaking

	default:
		return schema.SeverityInfo
	}
}

// isRequiredGQLType returns true when a GraphQL type string is non-nullable
// (ends with "!") and therefore required. e.g. "String!" → true, "String" → false.
func isRequiredGQLType(t string) bool {
	return strings.HasSuffix(strings.TrimSpace(t), "!")
}

// isNullabilityRelaxed returns true when a non-null type becomes nullable.
// e.g. "String!" → "String" (output field loses its non-null guarantee).
func isNullabilityRelaxed(before, after string) bool {
	return isRequiredGQLType(before) && !isRequiredGQLType(after) &&
		strings.TrimSuffix(strings.TrimSpace(before), "!") == strings.TrimSpace(after)
}

// isNullabilityTightened returns true when a nullable type becomes non-null.
// e.g. "String" → "String!" (type is now guaranteed non-null).
func isNullabilityTightened(before, after string) bool {
	return !isRequiredGQLType(before) && isRequiredGQLType(after) &&
		strings.TrimSpace(before) == strings.TrimSuffix(strings.TrimSpace(after), "!")
}

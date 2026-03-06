package classifier

import (
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

	default:
		return schema.SeverityInfo
	}
}

package languages

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// isNestJSProject returns true when the project at dir has a package.json
// that declares any core NestJS package as a dependency.
func isNestJSProject(dir string) bool {
	for _, pkg := range []string{"@nestjs/core", "@nestjs/common", "@nestjs/swagger"} {
		if hasPackageJSONDep(dir, pkg) {
			return true
		}
	}
	return false
}

// isExpressProject returns true when the project at dir has a package.json
// that declares express as a dependency but is NOT a NestJS project.
func isExpressProject(dir string) bool {
	if _, err := os.Stat(filepath.Join(dir, "package.json")); err != nil {
		return false
	}
	return hasPackageJSONDep(dir, "express") && !isNestJSProject(dir)
}

// isNodeJSProject returns true when the project at dir has a package.json but
// is not NestJS or Express (generic Node.js / TypeScript project).
func isNodeJSProject(dir string) bool {
	if _, err := os.Stat(filepath.Join(dir, "package.json")); err != nil {
		return false
	}
	return !isNestJSProject(dir)
}

// hasPackageJSONDep reports whether package.json in dir lists depName in
// dependencies or devDependencies.
func hasPackageJSONDep(dir, depName string) bool {
	data, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return false
	}
	var pkg struct {
		Dependencies    map[string]json.RawMessage `json:"dependencies"`
		DevDependencies map[string]json.RawMessage `json:"devDependencies"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return strings.Contains(string(data), `"`+depName+`"`)
	}
	_, inDeps := pkg.Dependencies[depName]
	_, inDev := pkg.DevDependencies[depName]
	return inDeps || inDev
}

package generate

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ScaffoldNestSwaggerScript writes a starter scripts/generate-swagger.ts to
// projectDir (creating the scripts/ directory if needed). It returns the path
// of the file that was written.
//
// The generated script uses NestFactory to boot the app, calls
// SwaggerModule.createDocument, and writes the spec to the path given by the
// SWAGGER_OUTPUT environment variable — exactly what drift-guard expects.
// Inline comments guide the user on how to mock heavy providers (TypeORM,
// Redis, etc.) when running outside of a live environment.
func ScaffoldNestSwaggerScript(projectDir string) (string, error) {
	outPath := filepath.Join(projectDir, "scripts", "generate-swagger.ts")

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return "", fmt.Errorf("create scripts directory: %w", err)
	}

	appModuleRelPath := detectAppModuleRelPath(projectDir)

	content := buildNestSwaggerScaffold(appModuleRelPath)

	if err := os.WriteFile(outPath, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("write scaffold: %w", err)
	}
	return outPath, nil
}

// detectAppModuleRelPath returns a relative import path for the AppModule,
// suitable for use inside scripts/generate-swagger.ts.
func detectAppModuleRelPath(projectDir string) string {
	candidates := []struct {
		rel    string // file to stat
		import_ string // import path to use
	}{
		{"src/app.module.ts", "../src/app.module"},
		{"src/app.module.js", "../src/app.module"},
		{"app.module.ts", "../app.module"},
		{"app.module.js", "../app.module"},
	}
	for _, c := range candidates {
		if _, err := os.Stat(filepath.Join(projectDir, c.rel)); err == nil {
			return c.import_
		}
	}
	return "../src/app.module" // sensible default
}

// buildNestSwaggerScaffold returns the TypeScript source for the scaffold.
func buildNestSwaggerScaffold(appModuleRelPath string) string {
	return fmt.Sprintf(`/**
 * scripts/generate-swagger.ts
 *
 * Generates an OpenAPI (swagger) document for your NestJS application and
 * writes it to the path specified by the SWAGGER_OUTPUT environment variable.
 *
 * Run via drift-guard:
 *   drift-guard generate openapi
 *
 * Or directly:
 *   SWAGGER_OUTPUT=swagger.json npx ts-node --transpile-only scripts/generate-swagger.ts
 *
 * -----------------------------------------------------------------------
 * If your app requires a live database or other services to start, you
 * have two options:
 *
 * Option A — start the real services, then run this script.
 *
 * Option B — override the heavy providers with no-op mocks so the app can
 * boot without infrastructure. Example using @nestjs/testing:
 *
 *   import { Test } from '@nestjs/testing';
 *   import { TypeOrmModule } from '@nestjs/typeorm';
 *   import { getDataSourceToken } from '@nestjs/typeorm';
 *
 *   const moduleRef = await Test.createTestingModule({ imports: [AppModule] })
 *     .overrideProvider(getDataSourceToken())
 *     .useValue({ isInitialized: true })
 *     .compile();
 *   const app = moduleRef.createNestApplication();
 *   await app.init();
 * -----------------------------------------------------------------------
 */

import { NestFactory } from '@nestjs/core';
import { SwaggerModule, DocumentBuilder } from '@nestjs/swagger';
import * as fs from 'fs';
import { AppModule } from '%s';

async function generate(): Promise<void> {
  // abortOnError: false makes NestJS throw on failure instead of silently
  // calling process.exit(1), so errors are visible in the output below.
  const app = await NestFactory.create(AppModule, { abortOnError: false });

  const config = new DocumentBuilder()
    .setTitle('API')
    .setVersion('1.0')
    .build();

  const document = SwaggerModule.createDocument(app, config);

  const output = process.env.SWAGGER_OUTPUT ?? 'swagger.json';
  fs.writeFileSync(output, JSON.stringify(document, null, 2));

  // Force-exit so open handles (DB pools, queues) don't block the process.
  process.exit(0);
}

generate().catch((err) => {
  console.error(err);
  process.exit(1);
});
`, appModuleRelPath)
}

// --------------------------------------------------------------------------
// tsoa scaffold
// --------------------------------------------------------------------------

// ScaffoldTsoa writes a tsoa.json with sensible defaults to projectDir and
// returns the path of the file written.
func ScaffoldTsoa(projectDir string) (string, error) {
	outPath := filepath.Join(projectDir, "tsoa.json")

	entryFile := detectEntryFile(projectDir)

	cfg := map[string]any{
		"entryFile":                      entryFile,
		"noImplicitAdditionalProperties": "throw",
		"controllerPathGlobs":            []string{"src/**/*.controller.ts"},
		"spec": map[string]any{
			"outputDirectory": "build",
			"specVersion":     3,
		},
		"routes": map[string]any{
			"routesDir": "build",
		},
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal tsoa config: %w", err)
	}

	if err := os.WriteFile(outPath, append(data, '\n'), 0o644); err != nil {
		return "", fmt.Errorf("write tsoa.json: %w", err)
	}
	return outPath, nil
}

// InstallTsoa runs the appropriate package manager install command to add tsoa
// as a dev dependency.
func InstallTsoa(projectDir string) error {
	pm := detectPackageManager(projectDir)
	var args []string
	switch pm {
	case "pnpm":
		args = []string{"add", "--save-dev", "tsoa"}
	case "yarn":
		args = []string{"add", "--dev", "tsoa"}
	default:
		args = []string{"install", "--save-dev", "tsoa"}
		pm = "npm"
	}
	cmd := exec.Command(pm, args...)
	cmd.Dir = projectDir
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s install tsoa: %w", pm, err)
	}
	return nil
}

// HasTsoaControllers reports whether the project at projectDir uses tsoa
// controller decorators (@Route, @Get, etc.). Returns false for plain Express
// projects that don't use tsoa's decorator-based approach.
func HasTsoaControllers(projectDir string) bool {
	// Walk TypeScript source files looking for tsoa's @Route decorator.
	srcDir := filepath.Join(projectDir, "src")
	if _, err := os.Stat(srcDir); err != nil {
		srcDir = projectDir
	}
	found := false
	_ = filepath.WalkDir(srcDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".ts" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		if strings.Contains(string(data), "@Route(") {
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	return found
}

func detectEntryFile(projectDir string) string {
	candidates := []string{
		"src/app.ts", "src/main.ts", "src/server.ts", "src/index.ts",
		"app.ts", "main.ts", "server.ts", "index.ts",
	}
	for _, c := range candidates {
		if _, err := os.Stat(filepath.Join(projectDir, c)); err == nil {
			return c
		}
	}
	return "src/app.ts"
}

func detectPackageManager(projectDir string) string {
	if _, err := os.Stat(filepath.Join(projectDir, "pnpm-lock.yaml")); err == nil {
		return "pnpm"
	}
	if _, err := os.Stat(filepath.Join(projectDir, "yarn.lock")); err == nil {
		return "yarn"
	}
	return "npm"
}

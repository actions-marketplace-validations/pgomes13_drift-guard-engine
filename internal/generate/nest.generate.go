package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Nest generates a swagger.json in outputDir for the NestJS project rooted at
// projectDir using @nestjs/swagger.
//
// Strategy (in order):
//  1. tsoa   — if tsoa.json is present, run `npx tsoa spec` and copy the result.
//  2. Script — look for an existing scripts/generate-swagger.ts (or .js).
//  3. Scaffold — write a temporary NestJS bootstrap script and run it.
func Nest(projectDir, outputDir string) error {
	// 1. tsoa
	if _, err := os.Stat(filepath.Join(projectDir, "tsoa.json")); err == nil {
		return tsoaSpec(projectDir, outputDir)
	}

	outputPath := filepath.Join(outputDir, "swagger.json")

	// 2. Existing generation script.
	candidates := []string{
		"scripts/generate-swagger.ts",
		"scripts/generate-swagger.js",
		"src/generate-swagger.ts",
		"generate-swagger.ts",
	}
	for _, rel := range candidates {
		full := filepath.Join(projectDir, rel)
		if _, err := os.Stat(full); err == nil {
			return runScript(projectDir, full, outputPath)
		}
	}

	// 3. Scaffold a temporary NestJS bootstrap script.
	scriptPath := filepath.Join(projectDir, ".drift-guard-swagger-gen.ts")
	defer os.Remove(scriptPath) //nolint:errcheck

	if err := os.WriteFile(scriptPath, []byte(nestJSScript(projectDir)), 0o600); err != nil {
		return fmt.Errorf("scaffold NestJS swagger script: %w", err)
	}
	return runScript(projectDir, scriptPath, outputPath)
}

func nestJSScript(projectDir string) string {
	appModule := "./src/app.module"
	for _, c := range []string{"src/app.module.ts", "src/app.module.js"} {
		if _, err := os.Stat(filepath.Join(projectDir, c)); err == nil {
			appModule = "./" + strings.TrimSuffix(filepath.ToSlash(c), filepath.Ext(c))
			break
		}
	}

	return fmt.Sprintf(`import { NestFactory } from '@nestjs/core';
import { SwaggerModule, DocumentBuilder } from '@nestjs/swagger';
import { writeFileSync } from 'fs';
import { AppModule } from '%s';

async function generate() {
  const app = await NestFactory.create(AppModule, { logger: false });
  const config = new DocumentBuilder()
    .setTitle('API')
    .setVersion('1.0')
    .build();
  const document = SwaggerModule.createDocument(app, config);
  const outputPath = process.env.SWAGGER_OUTPUT || 'swagger.json';
  writeFileSync(outputPath, JSON.stringify(document, null, 2));
  await app.close();
}

generate().catch(err => { console.error(err); process.exit(1); });
`, appModule)
}

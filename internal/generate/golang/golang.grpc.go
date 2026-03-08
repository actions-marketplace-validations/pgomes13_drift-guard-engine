package golang

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GoGRPC finds the primary .proto file for the Go project and copies it to
// outputDir/schema.proto.
func GoGRPC(projectDir, outputDir string) error {
	src := FindProtoFile(projectDir)
	if src == "" {
		return fmt.Errorf(
			"no .proto file found in %s\n\n"+
				"Ensure your proto file is in one of:\n"+
				"  proto/, protos/, src/proto/, or the project root.",
			projectDir,
		)
	}
	return copySchema(src, filepath.Join(outputDir, "schema.proto"))
}

// FindProtoFile returns the absolute path of the first .proto file found in
// dir, checking common locations before falling back to a directory walk.
func FindProtoFile(dir string) string {
	for _, sub := range []string{"proto", "protos", "src/proto", "src/protos", "grpc", "."} {
		subDir := filepath.Join(dir, filepath.FromSlash(sub))
		entries, err := os.ReadDir(subDir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".proto") {
				return filepath.Join(subDir, e.Name())
			}
		}
	}
	return ""
}

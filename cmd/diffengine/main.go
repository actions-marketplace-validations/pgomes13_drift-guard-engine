package main

import (
	"flag"
	"fmt"
	"os"

	"drift-guard-diff-engine/internal/classifier"
	"drift-guard-diff-engine/internal/differ"
	"drift-guard-diff-engine/internal/parser"
	"drift-guard-diff-engine/internal/reporter"
)

func main() {
	var (
		baseFile   = flag.String("base", "", "Path to the base OpenAPI schema (e.g. main branch)")
		headFile   = flag.String("head", "", "Path to the head OpenAPI schema (e.g. PR branch)")
		format     = flag.String("format", "text", "Output format: text, json, github")
		failOnBreak = flag.Bool("fail-on-breaking", false, "Exit with code 1 if breaking changes are detected")
	)
	flag.Parse()

	if *baseFile == "" || *headFile == "" {
		fmt.Fprintln(os.Stderr, "Error: --base and --head are required")
		flag.Usage()
		os.Exit(2)
	}

	baseSchema, err := parser.ParseFile(*baseFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing base schema: %v\n", err)
		os.Exit(2)
	}

	headSchema, err := parser.ParseFile(*headFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing head schema: %v\n", err)
		os.Exit(2)
	}

	changes := differ.Diff(baseSchema, headSchema)
	result := classifier.Classify(*baseFile, *headFile, changes)

	if err := reporter.Write(os.Stdout, result, reporter.Format(*format)); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing report: %v\n", err)
		os.Exit(2)
	}

	if *failOnBreak && reporter.HasBreakingChanges(result) {
		os.Exit(1)
	}
}

package impact

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Report writes the impact report to w in the requested format (text, json, markdown, github).
func Report(w io.Writer, hits []Hit, format string) error {
	if len(hits) == 0 {
		fmt.Fprintln(w, "No references found.")
		return nil
	}
	switch format {
	case "json":
		return reportJSON(w, hits)
	case "markdown":
		return reportMarkdown(w, hits)
	case "github":
		return reportGitHub(w, hits)
	default:
		return reportText(w, hits)
	}
}

// groupByChange buckets hits by their (changeType, changePath) pair.
func groupByChange(hits []Hit) map[string][]Hit {
	m := make(map[string][]Hit)
	for _, h := range hits {
		key := h.ChangeType + "\x00" + h.ChangePath
		m[key] = append(m[key], h)
	}
	return m
}

func sortedKeys(m map[string][]Hit) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func reportText(w io.Writer, hits []Hit) error {
	groups := groupByChange(hits)
	for _, key := range sortedKeys(groups) {
		parts := strings.SplitN(key, "\x00", 2)
		fmt.Fprintf(w, "Breaking change: %s\n", parts[1])
		for _, h := range groups[key] {
			fmt.Fprintf(w, "  %s:%d\t%s\n", h.File, h.LineNum, h.Line)
		}
		fmt.Fprintln(w)
	}
	return nil
}

func reportMarkdown(w io.Writer, hits []Hit) error {
	groups := groupByChange(hits)
	keys := sortedKeys(groups)

	// Summary line
	totalFiles := countDistinctFiles(hits)
	fmt.Fprintf(w, "> **%d** reference(s) to breaking changes across **%d** file(s)\n\n", len(hits), totalFiles)

	for _, key := range keys {
		group := groups[key]
		parts := strings.SplitN(key, "\x00", 2)
		changePath := parts[1]

		// Collapsible section per breaking change
		fmt.Fprintf(w, "<details>\n<summary>🔴 <strong>%s</strong> — %d reference(s)</summary>\n\n", changePath, len(group))
		fmt.Fprintln(w, "| File | Line | Code |")
		fmt.Fprintln(w, "|------|------|------|")
		for _, h := range group {
			fmt.Fprintf(w, "| `%s` | %d | `%s` |\n", h.File, h.LineNum, escapeMarkdown(strings.TrimSpace(h.Line)))
		}
		fmt.Fprintf(w, "\n</details>\n\n")
	}
	return nil
}

func countDistinctFiles(hits []Hit) int {
	seen := make(map[string]struct{}, len(hits))
	for _, h := range hits {
		seen[h.File] = struct{}{}
	}
	return len(seen)
}

// reportGitHub emits GitHub Actions workflow commands so each hit appears as an
// inline annotation in the PR diff view:
//
//	::error file=<path>,line=<n>,title=<change>::<message>
func reportGitHub(w io.Writer, hits []Hit) error {
	for _, h := range hits {
		title := escapeProperty(fmt.Sprintf("Breaking API change: %s", h.ChangePath))
		msg := escapeData(strings.TrimSpace(h.Line))
		fmt.Fprintf(w, "::error file=%s,line=%d,title=%s::%s\n",
			h.File, h.LineNum, title, msg)
	}
	return nil
}

func reportJSON(w io.Writer, hits []Hit) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(hits)
}

func escapeMarkdown(s string) string {
	return strings.NewReplacer("|", "\\|", "`", "'").Replace(s)
}

// escapeProperty escapes special characters in GitHub Actions workflow command
// property values (title, file, etc.).
func escapeProperty(s string) string {
	return strings.NewReplacer("%", "%25", "\r", "%0D", "\n", "%0A", ":", "%3A", ",", "%2C").Replace(s)
}

// escapeData escapes special characters in GitHub Actions workflow command data
// (the message after the final ::).
func escapeData(s string) string {
	return strings.NewReplacer("%", "%25", "\r", "%0D", "\n", "%0A").Replace(s)
}

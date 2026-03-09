package impact

import (
	"strings"
	"testing"
)

var testHits = []Hit{
	{File: "services/client.go", LineNum: 42, Line: "\tclient.Delete(\"/users/\" + id)", ChangeType: "endpoint_removed", ChangePath: "DELETE /users/{id}"},
	{File: "apps/routes.go", LineNum: 17, Line: "\tr.DELETE(\"/users/:id\", handler)", ChangeType: "endpoint_removed", ChangePath: "DELETE /users/{id}"},
}

func TestReportGitHub_EmitsAnnotations(t *testing.T) {
	var b strings.Builder
	if err := Report(&b, testHits, "github"); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	for _, want := range []string{
		"::error file=services/client.go,line=42,",
		"::error file=apps/routes.go,line=17,",
		"title=Breaking API change%3A DELETE /users/{id}",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\ngot:\n%s", want, out)
		}
	}
}

func TestReportGitHub_NoHits(t *testing.T) {
	var b strings.Builder
	if err := Report(&b, nil, "github"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(b.String(), "No references found") {
		t.Errorf("expected no-hit message, got: %s", b.String())
	}
}

func TestReportGitHub_EscapesSpecialChars(t *testing.T) {
	hits := []Hit{
		{File: "src/api.go", LineNum: 1, Line: "100% done", ChangeType: "field_removed", ChangePath: "POST /x > body > email"},
	}
	var b strings.Builder
	if err := Report(&b, hits, "github"); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	if strings.Contains(out, "%\r") || strings.Contains(out, "%\n") {
		t.Errorf("unescaped percent in output: %s", out)
	}
	if !strings.Contains(out, "100%25 done") {
		t.Errorf("percent in message not escaped: %s", out)
	}
}

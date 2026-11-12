// Package conftc1 verifies that a run over the distribution case set reports
// the four groups, attributes every case to a contract entry, and records one
// failing capability case without changing any other result (CONF-AC-1,
// CONF-TC-1).
package conftc1

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/tsic404/kos-conformance"
	"github.com/tsic404/kos-conformance/acceptance"
)

// failingCase is the capability case the run is made to fail.
const failingCase = "volume.bind"

// observed is the difference the failing implementation reports.
const observed = "the claim stays pending: the backend volume was not bound"

// coveredContractEntries is what each group must cover (CONF-C-1).
var coveredContractEntries = map[conformance.Group][]string{
	conformance.GroupCapability: {"KOSD-C-09", "KOSD-C-10", "KOSD-C-11"},
	conformance.GroupProvider:   {"KOSD-C-47", "KOSD-C-48"},
	conformance.GroupExtension:  {"KOSD-C-17", "KOSD-C-18", "KOSD-C-19", "KOSD-C-20"},
	conformance.GroupApp:        {"KOSD-C-73", "KOSD-C-75"},
}

// TestFailureIsListedAndLeavesOtherResultsUnchanged runs the case set against a
// provider whose volume binding fails.
func TestFailureIsListedAndLeavesOtherResultsUnchanged(t *testing.T) {
	distribution := acceptance.Standard()
	runner := acceptance.NewRunner(distribution, acceptance.Registry{
		"fixture.storage": acceptance.TargetFunc(func(_ context.Context, c conformance.AttributedCase) error {
			if c.ID == failingCase {
				return errors.New(observed)
			}
			return nil
		}),
	})

	report := runner.Run(t.Context(), "fixture.storage", distribution.Digest())

	if len(report.Results) != len(distribution.Cases()) {
		t.Fatalf("report carries %d cases, the distribution ships %d", len(report.Results), len(distribution.Cases()))
	}

	seen := make(map[conformance.Group]map[string]bool)
	for _, result := range report.Results {
		if !result.Case.Group.IsKnown() {
			t.Fatalf("case %q carries group %q outside the four groups", result.Case.ID, result.Case.Group)
		}
		if len(result.Case.Objectives) == 0 || result.Case.CriterionVersion == "" {
			t.Fatalf("case %q is reported without contract attribution", result.Case.ID)
		}
		if seen[result.Case.Group] == nil {
			seen[result.Case.Group] = make(map[string]bool, len(coveredContractEntries[result.Case.Group]))
		}
		for _, objective := range result.Case.Objectives {
			seen[result.Case.Group][objective] = true
		}
	}
	for group, objectives := range coveredContractEntries {
		if len(seen[group]) == 0 {
			t.Fatalf("group %q carries no case", group)
		}
		for _, objective := range objectives {
			if !seen[group][objective] {
				t.Fatalf("group %q has no case for %s", group, objective)
			}
		}
	}

	failed := report.Failed()
	if len(failed) != 1 || failed[0].Case.ID != failingCase {
		t.Fatalf("failed list is %v, want the single case %q", caseIDs(failed), failingCase)
	}
	if !slices.Contains(failed[0].Case.Objectives, "KOSD-C-09") {
		t.Fatalf("failed case names objectives %v, want KOSD-C-09", failed[0].Case.Objectives)
	}
	if !strings.Contains(failed[0].Detail, observed) {
		t.Fatalf("failed case records %q, want the observed difference", failed[0].Detail)
	}

	for _, result := range report.Results {
		if result.Case.ID == failingCase {
			continue
		}
		if result.Status != acceptance.StatusPassed {
			t.Fatalf("case %q is %s, want passed: one failure must not change other results", result.Case.ID, result.Status)
		}
	}

	if report.Verdict != acceptance.VerdictFailed {
		t.Fatalf("verdict is %q, want failed: a comparable run with a failing case is not passed", report.Verdict)
	}
}

// caseIDs lists the cases of one result list for failure messages.
func caseIDs(results []acceptance.Result) []string {
	ids := make([]string, 0, len(results))
	for _, result := range results {
		ids = append(ids, result.Case.ID)
	}
	return ids
}

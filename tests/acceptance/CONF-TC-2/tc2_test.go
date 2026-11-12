// Package conftc2 verifies that a case whose prerequisite the provider lacks
// is recorded as not run, kept apart from failures, and never counted as a
// pass (CONF-AC-1, CONF-TC-2).
package conftc2

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/tsic404/kos-conformance"
	"github.com/tsic404/kos-conformance/acceptance"
)

// blockedCase is the case whose prerequisite the provider lacks.
const blockedCase = "endpoint.tls"

// prerequisite is what the provider is missing.
const prerequisite = "no certificate issuer is installed"

// TestMissingPrerequisiteIsNotRun runs the case set against a provider that
// cannot serve TLS because it has no issuer.
func TestMissingPrerequisiteIsNotRun(t *testing.T) {
	distribution := acceptance.Standard()
	runner := acceptance.NewRunner(distribution, acceptance.Registry{
		"fixture.network": acceptance.TargetFunc(func(_ context.Context, c conformance.AttributedCase) error {
			if c.ID == blockedCase {
				return fmt.Errorf("%w: %s", acceptance.ErrPrerequisite, prerequisite)
			}
			return nil
		}),
	})

	report := runner.Run(t.Context(), "fixture.network", distribution.Digest())

	notRun := report.NotRun()
	if len(notRun) != 1 || notRun[0].Case.ID != blockedCase {
		t.Fatalf("not run list is %v, want the single case %q", caseIDs(notRun), blockedCase)
	}
	if !strings.Contains(notRun[0].Detail, prerequisite) {
		t.Fatalf("not run case records %q, want the unmet prerequisite", notRun[0].Detail)
	}

	// 未执行与未通过分列：缺前置不产生失败项。
	if failed := report.Failed(); len(failed) != 0 {
		t.Fatalf("failed list is %v, want empty: a missing prerequisite is not a failure", caseIDs(failed))
	}

	// 未执行不出现在通过项。
	for _, result := range report.Passed() {
		if result.Case.ID == blockedCase {
			t.Fatalf("case %q is counted as passed although it did not run", blockedCase)
		}
	}
	if len(report.Passed()) != len(report.Results)-1 {
		t.Fatalf("passed list carries %d of %d cases, want every case but the not run one", len(report.Passed()), len(report.Results))
	}

	// 其余用例结果不受影响。
	for _, result := range report.Results {
		if result.Case.ID == blockedCase {
			continue
		}
		if result.Status != acceptance.StatusPassed {
			t.Fatalf("case %q is %s, want passed: a not run case must not change other results", result.Case.ID, result.Status)
		}
	}

	if report.Verdict == acceptance.VerdictPassed {
		t.Fatalf("verdict is %q, want a run with a not run case not passed", report.Verdict)
	}
	if report.Verdict != acceptance.VerdictFailed {
		t.Fatalf("verdict is %q, want failed: a comparable run that is not fully passing", report.Verdict)
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

// Package conftc4 verifies that a run whose digest does not match the
// distribution's record is recorded as incomparable, with both digests
// readable and the verdict not passed (CONF-AC-2, CONF-TC-4).
package conftc4

import (
	"context"
	"crypto/sha256"
	"strings"
	"testing"

	"github.com/tsic404/kos-conformance"
	"github.com/tsic404/kos-conformance/acceptance"
)

// passing answers every case successfully, so the verdict cannot come from a
// failed case.
func passing(context.Context, conformance.AttributedCase) error { return nil }

// TestMismatchedDigestIsIncomparable runs a fully passing capability against a
// record from a different case set.
func TestMismatchedDigestIsIncomparable(t *testing.T) {
	distribution := acceptance.Standard()
	runner := acceptance.NewRunner(distribution, acceptance.Registry{
		"fixture.storage": acceptance.TargetFunc(passing),
	})
	recorded := "sha256:" + strings.Repeat("0", sha256.Size*2)

	report := runner.Run(t.Context(), "fixture.storage", recorded)

	if report.Verdict != acceptance.VerdictIncomparable {
		t.Fatalf("verdict is %q, want incomparable: the digest differs from the record", report.Verdict)
	}
	if report.DistributionDigest != recorded {
		t.Fatalf("report carries record %q, want %q", report.DistributionDigest, recorded)
	}
	if report.RunDigest != distribution.Digest() {
		t.Fatalf("report carries run digest %q, want %q", report.RunDigest, distribution.Digest())
	}
	if failed, notRun := report.Failed(), report.NotRun(); len(failed) != 0 || len(notRun) != 0 {
		t.Fatalf("run has %d failed and %d not run cases, want every case passed", len(failed), len(notRun))
	}

	// 结论不可比来自摘要不一致：同一 Provider 用一致摘要运行时记为通过。
	matched := runner.Run(t.Context(), "fixture.storage", distribution.Digest())
	if matched.Verdict != acceptance.VerdictPassed {
		t.Fatalf("verdict with a matching record is %q, want passed", matched.Verdict)
	}
}

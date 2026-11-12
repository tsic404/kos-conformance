// Package conftc3 verifies that the case-set digest stays the same across runs
// of the same distribution and changes when a case's criterion or the case
// collection changes (CONF-AC-2, CONF-TC-3).
package conftc3

import (
	"context"
	"crypto/sha256"
	"strings"
	"testing"

	"github.com/tsic404/kos-conformance"
	"github.com/tsic404/kos-conformance/acceptance"
)

// passing answers every case successfully, so results never explain a digest.
func passing(context.Context, conformance.AttributedCase) error { return nil }

// TestDigestIsStableAndTracksTheCaseSet runs the same distribution twice, then
// changes one case's criterion and runs once more.
func TestDigestIsStableAndTracksTheCaseSet(t *testing.T) {
	distribution := acceptance.Standard()
	runner := acceptance.NewRunner(distribution, acceptance.Registry{
		"fixture.storage": acceptance.TargetFunc(passing),
	})

	first := runner.Run(t.Context(), "fixture.storage", distribution.Digest())
	second := runner.Run(t.Context(), "fixture.storage", distribution.Digest())
	if first.RunDigest != second.RunDigest {
		t.Fatalf("two runs of the same distribution report %q and %q, want one digest", first.RunDigest, second.RunDigest)
	}
	if want := len("sha256:") + sha256.Size*2; !strings.HasPrefix(first.RunDigest, "sha256:") || len(first.RunDigest) != want {
		t.Fatalf("run digest is %q, want sha256:<hex>", first.RunDigest)
	}

	bumped := rebuilt(t, func(set *conformance.CaseSet) bool {
		if set.ID != "kos/conformance/storage-volume-v1" {
			return false
		}
		changed := false
		for i := range set.Cases {
			if set.Cases[i].ID == "volume.bind" {
				set.Cases[i].CriterionVersion = "v2"
				changed = true
			}
		}
		return changed
	})
	third := acceptance.NewRunner(bumped, acceptance.Registry{
		"fixture.storage": acceptance.TargetFunc(passing),
	}).Run(t.Context(), "fixture.storage", bumped.Digest())
	if third.RunDigest == first.RunDigest || third.RunDigest == second.RunDigest {
		t.Fatalf("digest %q survives a criterion change, want one different from the earlier runs", third.RunDigest)
	}

	reduced := rebuilt(t, func(set *conformance.CaseSet) bool {
		if set.ID != "kos/conformance/storage-volume-v1" {
			return false
		}
		before := len(set.Cases)
		set.Cases = set.Cases[:before-1]
		return len(set.Cases) != before
	})
	if reduced.Digest() == distribution.Digest() {
		t.Fatalf("digest %q survives a case leaving the collection", reduced.Digest())
	}
}

// rebuilt copies the shipped distribution, applies mutate to every set, and
// returns a distribution over the result. A mutation that applied nowhere fails
// the test: such a test would assert nothing.
func rebuilt(t *testing.T, mutate func(*conformance.CaseSet) bool) *conformance.Distribution {
	t.Helper()
	sets := acceptance.Standard().Sets()
	applied := 0
	for i := range sets {
		if mutate(&sets[i]) {
			applied++
		}
	}
	if applied == 0 {
		t.Fatalf("the mutation reached no case set")
	}
	distribution, err := conformance.NewDistribution(sets...)
	if err != nil {
		t.Fatalf("rebuild the distribution: %v", err)
	}
	return distribution
}

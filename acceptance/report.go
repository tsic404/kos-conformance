// Package acceptance is the distribution's run entry and result summary: it
// executes a case set against one provider and reports every case as passed,
// failed or not run (CONF-C-2), together with the digest that decides whether
// the run is comparable (CONF-C-3).
package acceptance

import "github.com/tsic404/kos-conformance"

// Status is one case's outcome in a run.
type Status string

// The three outcomes a report carries; they are exclusive.
const (
	StatusPassed Status = "passed"
	StatusFailed Status = "failed"
	StatusNotRun Status = "not-run"
)

// Verdict is the conclusion of one run. Incomparable is not a failure: it says
// the run cannot be compared with the distribution's record (CONF-C-3).
type Verdict string

// The three conclusions a run can reach.
const (
	VerdictPassed       Verdict = "passed"
	VerdictFailed       Verdict = "failed"
	VerdictIncomparable Verdict = "incomparable"
)

// Result is one case's row: the attributed case, its outcome, and the observed
// difference when it did not pass.
type Result struct {
	Case   conformance.AttributedCase
	Status Status
	// Detail names what was observed: the difference for a failed case, the
	// unmet prerequisite for a not run case. It is empty for a pass.
	Detail string
}

// Report is the output of one run (CONF-C-2, CONF-C-3).
type Report struct {
	// Provider is the implementation the run exercised.
	Provider string
	// DistributionDigest is the record the run was compared against.
	DistributionDigest string
	// RunDigest is the digest of the case set this run executed.
	RunDigest string
	// Verdict is the conclusion: only a matching digest and a fully passing
	// run is passed.
	Verdict Verdict
	// Results holds one row per case, in report order.
	Results []Result
}

// Passed returns the cases that passed, in report order.
func (r Report) Passed() []Result {
	return r.byStatus(StatusPassed)
}

// Failed returns the cases that failed, in report order.
func (r Report) Failed() []Result {
	return r.byStatus(StatusFailed)
}

// NotRun returns the cases that could not run, in report order. They are kept
// apart from failures and never counted as passes (CONF-C-2).
func (r Report) NotRun() []Result {
	return r.byStatus(StatusNotRun)
}

// byStatus selects the rows carrying one outcome.
func (r Report) byStatus(status Status) []Result {
	selected := make([]Result, 0, len(r.Results))
	for _, result := range r.Results {
		if result.Status == status {
			selected = append(selected, result)
		}
	}
	return selected
}

// verdict computes the conclusion. A digest that differs from the record makes
// the run incomparable, never passed or failed; otherwise only a fully passing
// run is passed (CONF-C-2, CONF-C-3).
func (r Report) verdict() Verdict {
	if r.RunDigest != r.DistributionDigest {
		return VerdictIncomparable
	}
	for _, result := range r.Results {
		if result.Status != StatusPassed {
			return VerdictFailed
		}
	}
	return VerdictPassed
}

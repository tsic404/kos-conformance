package acceptance

import (
	"context"
	"errors"

	"github.com/tsic404/kos-conformance"
)

// ErrPrerequisite marks a case that could not run because the provider lacks
// what it needs. The runner records it as not run, never as a pass (CONF-C-2).
var ErrPrerequisite = errors.New("prerequisite unmet")

// Target executes one case against the provider under test.
type Target interface {
	Run(ctx context.Context, c conformance.AttributedCase) error
}

// TargetFunc adapts a function to a Target.
type TargetFunc func(ctx context.Context, c conformance.AttributedCase) error

// Run implements Target.
func (f TargetFunc) Run(ctx context.Context, c conformance.AttributedCase) error {
	return f(ctx, c)
}

// Registry maps provider names to their implementations. A provider without an
// entry cannot be run; its cases are not run, never passed (CONF-C-2).
type Registry map[string]Target

// Runner executes one distribution case set against the providers a
// distribution registers.
type Runner struct {
	distribution *conformance.Distribution
	targets      Registry
}

// NewRunner builds a runner over one distribution and one provider registry.
func NewRunner(distribution *conformance.Distribution, targets Registry) *Runner {
	return &Runner{distribution: distribution, targets: targets}
}

// Run executes every case against one provider and records each result. A
// failed case never stops the run; a provider without an implementation leaves
// every case not run, so an absent implementation can never read as a pass
// (CONF-C-2).
func (r *Runner) Run(ctx context.Context, provider, distributionDigest string) Report {
	report := Report{
		Provider:           provider,
		DistributionDigest: distributionDigest,
		RunDigest:          r.distribution.Digest(),
	}
	target, isRegistered := r.targets[provider]
	cases := r.distribution.Cases()
	report.Results = make([]Result, 0, len(cases))
	for _, c := range cases {
		report.Results = append(report.Results, r.record(ctx, provider, target, isRegistered, c))
	}
	report.Verdict = report.verdict()
	return report
}

// record runs one case and classifies its outcome.
func (r *Runner) record(ctx context.Context, provider string, target Target, isRegistered bool, c conformance.AttributedCase) Result {
	if !isRegistered {
		return Result{
			Case:   c,
			Status: StatusNotRun,
			Detail: "no implementation is registered for provider " + provider,
		}
	}
	err := target.Run(ctx, c)
	switch {
	case err == nil:
		return Result{Case: c, Status: StatusPassed}
	case errors.Is(err, ErrPrerequisite):
		return Result{Case: c, Status: StatusNotRun, Detail: err.Error()}
	default:
		return Result{Case: c, Status: StatusFailed, Detail: err.Error()}
	}
}

package acceptance

import (
	"github.com/tsic404/kos-conformance"
	"github.com/tsic404/kos-conformance/app"
	"github.com/tsic404/kos-conformance/capability"
	"github.com/tsic404/kos-conformance/extension"
	"github.com/tsic404/kos-conformance/provider"
)

// Standard returns the distribution case set this release ships: the four
// groups' sets assembled into one comparable collection (CONF-C-1).
func Standard() *conformance.Distribution {
	sets := make([]conformance.CaseSet, 0, 16)
	sets = append(sets, capability.Sets()...)
	sets = append(sets, provider.Sets()...)
	sets = append(sets, extension.Sets()...)
	sets = append(sets, app.Sets()...)
	distribution, err := conformance.NewDistribution(sets...)
	if err != nil {
		panic(err) // the shipped sets are static data; an invalid one is a broken build
	}
	return distribution
}

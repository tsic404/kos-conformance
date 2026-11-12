package conformance_test

import (
	"strings"
	"testing"

	"github.com/tsic404/kos-conformance"
)

// fixtureSet is a valid set every case of the table starts from.
func fixtureSet() conformance.CaseSet {
	return conformance.CaseSet{
		ID:      "kos/conformance/fixture-v1",
		Group:   conformance.GroupCapability,
		Version: "v1",
		Cases: []conformance.Case{{
			ID:               "fixture.check",
			Description:      "the fixture holds",
			Objectives:       []string{"KOSD-C-09"},
			CriterionVersion: "v1",
		}},
	}
}

// TestNewDistributionRefusesUnreportableSets covers the two report promises a
// broken set would violate: every case carries a group and a contract entry,
// and no identifier is ambiguous (CONF-C-1).
func TestNewDistributionRefusesUnreportableSets(t *testing.T) {
	mutated := func(mutate func(*conformance.CaseSet)) []conformance.CaseSet {
		set := fixtureSet()
		mutate(&set)
		return []conformance.CaseSet{set}
	}

	tests := []struct {
		name string
		sets []conformance.CaseSet
		want string
	}{
		{
			name: "no case set",
			sets: nil,
			want: "ships no conformance case set",
		},
		{
			name: "case without a contract entry",
			sets: mutated(func(set *conformance.CaseSet) { set.Cases[0].Objectives = nil }),
			want: "verifies no contract entry",
		},
		{
			name: "case without a criterion version",
			sets: mutated(func(set *conformance.CaseSet) { set.Cases[0].CriterionVersion = "" }),
			want: "names no criterion version",
		},
		{
			name: "group outside the four",
			sets: mutated(func(set *conformance.CaseSet) { set.Group = conformance.Group("kernel") }),
			want: "outside the four groups",
		},
		{
			name: "case declared twice",
			sets: mutated(func(set *conformance.CaseSet) { set.Cases = append(set.Cases, set.Cases[0]) }),
			want: "declares case",
		},
		{
			name: "case set declared twice",
			sets: []conformance.CaseSet{fixtureSet(), fixtureSet()},
			want: "declared twice",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := conformance.NewDistribution(test.sets...); err == nil {
				t.Fatalf("NewDistribution accepted a set a report could not carry")
			} else if !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error is %q, want it to mention %q", err, test.want)
			}
		})
	}
}

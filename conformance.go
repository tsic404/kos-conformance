// Package conformance holds the distribution case set: the four groups of
// cases, the contract entries each case verifies, and the digest that makes a
// run comparable (CONF-C-1, CONF-C-3).
package conformance

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

// Group is one of the four groups a distribution case set organizes cases
// into. The set is closed: a report carries no ungrouped case (CONF-C-1).
type Group string

// The four groups, in report order.
const (
	GroupCapability Group = "capability"
	GroupProvider   Group = "provider"
	GroupExtension  Group = "extension"
	GroupApp        Group = "app"
)

// Groups returns the four groups in report order.
func Groups() [4]Group {
	return [...]Group{GroupCapability, GroupProvider, GroupExtension, GroupApp}
}

// IsKnown reports whether g is one of the four groups.
func (g Group) IsKnown() bool {
	switch g {
	case GroupCapability, GroupProvider, GroupExtension, GroupApp:
		return true
	}
	return false
}

// Case is one check of the promise a case set judges. Objectives are the
// contract entries the case verifies; a case without one would be an
// unattributed result, which the distribution refuses to ship (CONF-C-1).
type Case struct {
	ID string
	// Description states the promise the case judges.
	Description string
	// Objectives are the contract entries the case verifies, e.g. KOSD-C-09.
	Objectives []string
	// CriterionVersion names the revision of the judgement this case
	// embodies. Changing what the case accepts requires a new version, which
	// changes the distribution digest (CONF-C-3).
	CriterionVersion string
}

// CaseSet is one set of cases: the checks a distribution ships for one
// capability version, or for a cross-capability rule when Capability is empty.
type CaseSet struct {
	// ID names the set; a report and a run both cite it.
	ID string
	// Group is the group the set and every case in it belongs to.
	Group Group
	// Capability is the capability the set judges; empty for a set that
	// judges a rule spanning capabilities.
	Capability string
	// Version is the case-set version the cases were written against.
	Version string
	// Cases are the checks, in execution order.
	Cases []Case
}

// AttributedCase is one case with the set that declares it. The runner and the
// report use this shape so the group and the contract attribution never
// detach from a case.
type AttributedCase struct {
	Group      Group
	SetID      string
	Capability string
	Version    string
	Case
}

// Distribution is the case-set collection a distribution ships (CONF-C-1). Its
// digest is the record a run is compared against (CONF-C-3).
type Distribution struct {
	sets   []CaseSet
	digest string
}

// NewDistribution validates the sets and computes the distribution digest. An
// empty collection, an unattributed case, an unknown group or a duplicate
// identifier is refused: such a set could be neither reported nor compared.
func NewDistribution(sets ...CaseSet) (*Distribution, error) {
	if len(sets) == 0 {
		return nil, fmt.Errorf("a distribution ships no conformance case set")
	}
	copied := make([]CaseSet, 0, len(sets))
	seen := make(map[string]bool, len(sets))
	for _, set := range sets {
		if err := validateSet(set); err != nil {
			return nil, err
		}
		if seen[set.ID] {
			return nil, fmt.Errorf("conformance case set %q is declared twice", set.ID)
		}
		seen[set.ID] = true
		copied = append(copied, copySet(set))
	}
	digest, err := digestOf(copied)
	if err != nil {
		return nil, err
	}
	return &Distribution{sets: copied, digest: digest}, nil
}

// Sets returns the case sets in declared order.
func (d *Distribution) Sets() []CaseSet {
	if d == nil {
		return nil
	}
	sets := make([]CaseSet, 0, len(d.sets))
	for _, set := range d.sets {
		sets = append(sets, copySet(set))
	}
	return sets
}

// Cases returns every case with its group and set identity: groups in report
// order, sets in declared order, cases in execution order.
func (d *Distribution) Cases() []AttributedCase {
	if d == nil {
		return nil
	}
	total := 0
	for _, set := range d.sets {
		total += len(set.Cases)
	}
	cases := make([]AttributedCase, 0, total)
	for _, group := range Groups() {
		for _, set := range d.sets {
			if set.Group != group {
				continue
			}
			for _, entry := range set.Cases {
				cases = append(cases, AttributedCase{
					Group:      set.Group,
					SetID:      set.ID,
					Capability: set.Capability,
					Version:    set.Version,
					Case:       entry,
				})
			}
		}
	}
	return cases
}

// Digest returns the case-set digest as sha256:<hex> (CONF-C-3).
func (d *Distribution) Digest() string {
	if d == nil {
		return ""
	}
	return d.digest
}

// validateSet rejects a set a run could not report: it needs an identifier, a
// known group, a version and at least one fully attributed case.
func validateSet(set CaseSet) error {
	if strings.TrimSpace(set.ID) == "" {
		return fmt.Errorf("a conformance case set names no identifier")
	}
	if !set.Group.IsKnown() {
		return fmt.Errorf("case set %q names group %q outside the four groups", set.ID, set.Group)
	}
	if strings.TrimSpace(set.Version) == "" {
		return fmt.Errorf("case set %q names no version", set.ID)
	}
	if len(set.Cases) == 0 {
		return fmt.Errorf("case set %q has no case", set.ID)
	}
	seen := make(map[string]bool, len(set.Cases))
	for _, entry := range set.Cases {
		if strings.TrimSpace(entry.ID) == "" {
			return fmt.Errorf("case set %q has a case without an identifier", set.ID)
		}
		if seen[entry.ID] {
			return fmt.Errorf("case set %q declares case %q twice", set.ID, entry.ID)
		}
		seen[entry.ID] = true
		if strings.TrimSpace(entry.Description) == "" {
			return fmt.Errorf("case %q states no promise", entry.ID)
		}
		if len(entry.Objectives) == 0 {
			return fmt.Errorf("case %q verifies no contract entry", entry.ID)
		}
		for _, objective := range entry.Objectives {
			if strings.TrimSpace(objective) == "" {
				return fmt.Errorf("case %q names an empty contract entry", entry.ID)
			}
		}
		if strings.TrimSpace(entry.CriterionVersion) == "" {
			return fmt.Errorf("case %q names no criterion version", entry.ID)
		}
	}
	return nil
}

// copySet deep-copies a set so a caller cannot mutate a shipped distribution.
func copySet(set CaseSet) CaseSet {
	copied := CaseSet{
		ID:         set.ID,
		Group:      set.Group,
		Capability: set.Capability,
		Version:    set.Version,
		Cases:      make([]Case, 0, len(set.Cases)),
	}
	for _, entry := range set.Cases {
		copied.Cases = append(copied.Cases, Case{
			ID:               entry.ID,
			Description:      entry.Description,
			Objectives:       slices.Clone(entry.Objectives),
			CriterionVersion: entry.CriterionVersion,
		})
	}
	return copied
}

// canonicalSet is what the digest covers of one set: its identity and every
// case's identity, objectives and criterion version. Descriptions are excluded
// because wording is not a criterion (CONF-C-3).
type canonicalSet struct {
	ID         string          `json:"id"`
	Group      Group           `json:"group"`
	Capability string          `json:"capability,omitempty"`
	Version    string          `json:"version"`
	Cases      []canonicalCase `json:"cases"`
}

// canonicalCase is what the digest covers of one case.
type canonicalCase struct {
	ID               string   `json:"id"`
	Objectives       []string `json:"objectives"`
	CriterionVersion string   `json:"criterionVersion"`
}

// digestOf renders the distribution digest: the sets sorted by identifier, so
// declaration order does not change what a distribution ships, with cases in
// execution order.
func digestOf(sets []CaseSet) (string, error) {
	canonical := make([]canonicalSet, 0, len(sets))
	for _, set := range sets {
		entry := canonicalSet{
			ID:         set.ID,
			Group:      set.Group,
			Capability: set.Capability,
			Version:    set.Version,
			Cases:      make([]canonicalCase, 0, len(set.Cases)),
		}
		for _, one := range set.Cases {
			entry.Cases = append(entry.Cases, canonicalCase{
				ID:               one.ID,
				Objectives:       one.Objectives,
				CriterionVersion: one.CriterionVersion,
			})
		}
		canonical = append(canonical, entry)
	}
	slices.SortFunc(canonical, func(a, b canonicalSet) int {
		return strings.Compare(a.ID, b.ID)
	})
	encoded, err := json.Marshal(canonical)
	if err != nil {
		return "", fmt.Errorf("encode the distribution case set: %w", err)
	}
	sum := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

// Package capability ships the capability group of the distribution case set:
// the baseline capabilities' promises, and the rules that make a verified
// capability trustworthy and a demand satisfiable (KOSD-C-09, KOSD-C-10,
// KOSD-C-11).
package capability

import "github.com/tsic404/kos-conformance"

// Sets returns the capability group's case sets.
func Sets() []conformance.CaseSet {
	return []conformance.CaseSet{
		{
			ID:         "kos/conformance/storage-volume-v1",
			Group:      conformance.GroupCapability,
			Capability: "storage.volume",
			Version:    "v1",
			Cases: []conformance.Case{
				{ID: "volume.provision", Description: "a claim is provisioned as a backend volume", Objectives: []string{"KOSD-C-09"}, CriterionVersion: "v1"},
				{ID: "volume.bind", Description: "the claim binds to the provisioned volume", Objectives: []string{"KOSD-C-09"}, CriterionVersion: "v1"},
				{ID: "volume.expand", Description: "the volume expands online without rebinding", Objectives: []string{"KOSD-C-09"}, CriterionVersion: "v1"},
				{ID: "volume.delete", Description: "deleting the claim releases the backend volume", Objectives: []string{"KOSD-C-09"}, CriterionVersion: "v1"},
			},
		},
		{
			ID:         "kos/conformance/storage-snapshot-v1",
			Group:      conformance.GroupCapability,
			Capability: "storage.snapshot",
			Version:    "v1",
			Cases: []conformance.Case{
				{ID: "snapshot.create", Description: "a snapshot captures the volume contents", Objectives: []string{"KOSD-C-09"}, CriterionVersion: "v1"},
				{ID: "snapshot.restore", Description: "a snapshot restores into a new volume", Objectives: []string{"KOSD-C-09"}, CriterionVersion: "v1"},
			},
		},
		{
			ID:         "kos/conformance/network-http-endpoint-v1",
			Group:      conformance.GroupCapability,
			Capability: "network.http-endpoint",
			Version:    "v1",
			Cases: []conformance.Case{
				{ID: "endpoint.host", Description: "a host name reaches the target application", Objectives: []string{"KOSD-C-09"}, CriterionVersion: "v1"},
				{ID: "endpoint.path", Description: "a path prefix reaches the target application", Objectives: []string{"KOSD-C-09"}, CriterionVersion: "v1"},
				{ID: "endpoint.tls", Description: "the served endpoint presents the issued certificate", Objectives: []string{"KOSD-C-09"}, CriterionVersion: "v1"},
			},
		},
		{
			ID:      "kos/conformance/capability-scope-v1",
			Group:   conformance.GroupCapability,
			Version: "v1",
			Cases: []conformance.Case{
				{ID: "scope.baseline", Description: "a baseline capability registered without a scope takes the baseline scope", Objectives: []string{"KOSD-C-09"}, CriterionVersion: "v1"},
				{ID: "scope.reject", Description: "a registration whose scope is outside the two values is refused", Objectives: []string{"KOSD-C-09"}, CriterionVersion: "v1"},
				{ID: "scope.extension-required", Description: "an extension capability registered without a scope is refused", Objectives: []string{"KOSD-C-09"}, CriterionVersion: "v1"},
			},
		},
		{
			ID:      "kos/conformance/capability-trust-v1",
			Group:   conformance.GroupCapability,
			Version: "v1",
			Cases: []conformance.Case{
				{ID: "trust.self-report", Description: "a provider's own result never verifies a capability", Objectives: []string{"KOSD-C-10"}, CriterionVersion: "v1"},
				{ID: "trust.partial", Description: "a partially failed run degrades the capability and lists the failed cases", Objectives: []string{"KOSD-C-10"}, CriterionVersion: "v1"},
				{ID: "trust.digest", Description: "a run whose digest differs from the distribution verifies nothing", Objectives: []string{"KOSD-C-10"}, CriterionVersion: "v1"},
			},
		},
		{
			ID:      "kos/conformance/capability-demand-v1",
			Group:   conformance.GroupCapability,
			Version: "v1",
			Cases: []conformance.Case{
				{ID: "demand.missing", Description: "a required capability with no compatible registration leaves the application pending", Objectives: []string{"KOSD-C-11"}, CriterionVersion: "v1"},
				{ID: "demand.guarantee", Description: "an insufficient guarantee degrades the application and lists the unmet items", Objectives: []string{"KOSD-C-11"}, CriterionVersion: "v1"},
				{ID: "demand.optional", Description: "a missing optional demand only degrades the application", Objectives: []string{"KOSD-C-11"}, CriterionVersion: "v1"},
				{ID: "demand.scope", Description: "a cluster demand is not satisfied by a node-scope registration", Objectives: []string{"KOSD-C-11"}, CriterionVersion: "v1"},
			},
		},
	}
}

// Package provider ships the provider group of the distribution case set: the
// registration and status rules, and the departure rules that protect the
// objects depending on a provider (KOSD-C-47, KOSD-C-48).
package provider

import "github.com/tsic404/kos-conformance"

// Sets returns the provider group's case sets.
func Sets() []conformance.CaseSet {
	return []conformance.CaseSet{
		{
			ID:      "kos/conformance/provider-registration-v1",
			Group:   conformance.GroupProvider,
			Version: "v1",
			Cases: []conformance.Case{
				{ID: "provider.registered", Description: "a provider registers its capability, storage class and version through the component path", Objectives: []string{"KOSD-C-47"}, CriterionVersion: "v1"},
				{ID: "provider.empty-storage-class", Description: "a provider without storage capability registers with an empty storage class", Objectives: []string{"KOSD-C-47"}, CriterionVersion: "v1"},
				{ID: "provider.result-separation", Description: "the self-reported result and the verified result are separately readable", Objectives: []string{"KOSD-C-47"}, CriterionVersion: "v1"},
				{ID: "provider.deprovision-keeps-ready", Description: "a provider whose deletion is blocked stays ready", Objectives: []string{"KOSD-C-47"}, CriterionVersion: "v1"},
			},
		},
		{
			ID:      "kos/conformance/provider-deprovision-v1",
			Group:   conformance.GroupProvider,
			Version: "v1",
			Cases: []conformance.Case{
				{ID: "provider.deprovision-blocked", Description: "deleting a referenced provider keeps it terminating", Objectives: []string{"KOSD-C-48"}, CriterionVersion: "v1"},
				{ID: "provider.dependents-listed", Description: "every dependent of a blocked provider is listed", Objectives: []string{"KOSD-C-48"}, CriterionVersion: "v1"},
				{ID: "provider.withdraw-blocked", Description: "withdrawing a referenced capability blocks the upgrade and keeps the old set effective", Objectives: []string{"KOSD-C-48"}, CriterionVersion: "v1"},
			},
		},
	}
}

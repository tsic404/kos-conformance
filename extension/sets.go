// Package extension ships the extension group of the distribution case set:
// the package, declaration, object-class and controller rules an extension
// must satisfy to be installed and served (KOSD-C-17 to KOSD-C-20).
package extension

import "github.com/tsic404/kos-conformance"

// Sets returns the extension group's case sets.
func Sets() []conformance.CaseSet {
	return []conformance.CaseSet{
		{
			ID:      "kos/conformance/extension-package-v1",
			Group:   conformance.GroupExtension,
			Version: "v1",
			Cases: []conformance.Case{
				{ID: "package.digest", Description: "a package whose digest does not match is refused", Objectives: []string{"KOSD-C-17"}, CriterionVersion: "v1"},
				{ID: "package.contributions", Description: "the contributions the declaration lists are exactly the visible ones", Objectives: []string{"KOSD-C-17"}, CriterionVersion: "v1"},
				{ID: "package.single-version", Description: "one extension identity has at most one installed version", Objectives: []string{"KOSD-C-17"}, CriterionVersion: "v1"},
				{ID: "package.system-app", Description: "a system app installs through the same package format and path", Objectives: []string{"KOSD-C-17"}, CriterionVersion: "v1"},
			},
		},
		{
			ID:      "kos/conformance/extension-declaration-v1",
			Group:   conformance.GroupExtension,
			Version: "v1",
			Cases: []conformance.Case{
				{ID: "declaration.required", Description: "a declaration missing a required item is refused", Objectives: []string{"KOSD-C-18"}, CriterionVersion: "v1"},
				{ID: "declaration.undeclared-invisible", Description: "a contribution the declaration omits never appears", Objectives: []string{"KOSD-C-18"}, CriterionVersion: "v1"},
			},
		},
		{
			ID:      "kos/conformance/extension-object-class-v1",
			Group:   conformance.GroupExtension,
			Version: "v1",
			Cases: []conformance.Case{
				{ID: "objectclass.served-when-installed", Description: "an object class is served only while its component is installed", Objectives: []string{"KOSD-C-19"}, CriterionVersion: "v1"},
				{ID: "objectclass.reference-blocks-delete", Description: "deleting an object an extension object references is blocked and lists the referrers", Objectives: []string{"KOSD-C-19"}, CriterionVersion: "v1"},
			},
		},
		{
			ID:      "kos/conformance/extension-controller-v1",
			Group:   conformance.GroupExtension,
			Version: "v1",
			Cases: []conformance.Case{
				{ID: "controller.write-scope", Description: "a write outside the granted scope is refused naming the subject and object", Objectives: []string{"KOSD-C-20"}, CriterionVersion: "v1"},
				{ID: "controller.unknown-isolation", Description: "stopping one controller turns only its objects unknown", Objectives: []string{"KOSD-C-20"}, CriterionVersion: "v1"},
			},
		},
	}
}

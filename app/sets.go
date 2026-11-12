// Package app ships the app group of the distribution case set: the window
// application manifest and the session rules that isolate an application from
// the console (KOSD-C-73, KOSD-C-75).
package app

import "github.com/tsic404/kos-conformance"

// Sets returns the app group's case sets.
func Sets() []conformance.CaseSet {
	return []conformance.CaseSet{
		{
			ID:      "kos/conformance/app-manifest-v1",
			Group:   conformance.GroupApp,
			Version: "v1",
			Cases: []conformance.Case{
				{ID: "manifest.required-fields", Description: "a manifest missing a required field is invalid", Objectives: []string{"KOSD-C-73"}, CriterionVersion: "v1"},
				{ID: "manifest.origin", Description: "an entry that shares the console origin is refused", Objectives: []string{"KOSD-C-73"}, CriterionVersion: "v1"},
				{ID: "manifest.asset-digest", Description: "a UI asset whose digest or signature does not match is not loaded", Objectives: []string{"KOSD-C-73"}, CriterionVersion: "v1"},
			},
		},
		{
			ID:      "kos/conformance/app-session-v1",
			Group:   conformance.GroupApp,
			Version: "v1",
			Cases: []conformance.Case{
				{ID: "session.token-binding", Description: "the window token binds account, application, window, permissions and expiry", Objectives: []string{"KOSD-C-75"}, CriterionVersion: "v1"},
				{ID: "session.origin-check", Description: "a message whose origin does not match is ignored", Objectives: []string{"KOSD-C-75"}, CriterionVersion: "v1"},
				{ID: "session.token-expiry", Description: "a request after the token expires is refused", Objectives: []string{"KOSD-C-75"}, CriterionVersion: "v1"},
				{ID: "session.permission-revoke", Description: "revoking a permission disables the capability immediately", Objectives: []string{"KOSD-C-75"}, CriterionVersion: "v1"},
				{ID: "session.top-navigation", Description: "top-level navigation from the application frame is blocked", Objectives: []string{"KOSD-C-75"}, CriterionVersion: "v1"},
			},
		},
	}
}

package kosdtc155

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/tsic404/kos-sdk/manifest"
	"github.com/tsic404/kos-sdk/packaging"
)

// host is the observation point of these cases. It plays the kernel's component
// surface: it holds the artifacts a registry serves, accepts an install only
// after the package passes the checks an install path runs, projects the
// contributions the installed declarations name, and withdraws one component on
// uninstall (KOSD-C-17, KOSD-C-29).
type host struct {
	artifacts  map[string][]byte
	components map[string]manifest.ExtensionManifest
	digests    map[string]packaging.Digest
	leftovers  []string
}

func newHost() *host {
	return &host{
		artifacts:  map[string][]byte{},
		components: map[string]manifest.ExtensionManifest{},
		digests:    map[string]packaging.Digest{},
	}
}

// publish places one artifact where an install request can reference it.
func (h *host) publish(reference string, archive []byte) {
	h.artifacts[reference] = archive
}

// install runs one install request: the artifact the reference names is read,
// verified against the digest the request names, and its declaration becomes the
// component's declared surface. A refused package installs nothing.
func (h *host) install(reference string, expected packaging.Digest) (manifest.ExtensionManifest, error) {
	archive, ok := h.artifacts[reference]
	if !ok {
		return manifest.ExtensionManifest{}, fmt.Errorf("artifact %q is not published", reference)
	}
	pkg, err := packaging.Verify(archive, expected)
	if err != nil {
		return manifest.ExtensionManifest{}, err
	}
	h.components[pkg.Manifest.ID] = pkg.Manifest
	h.digests[pkg.Manifest.ID] = pkg.Digest
	return pkg.Manifest, nil
}

// uninstall withdraws one component and records what its retention declaration
// keeps, which is the leftover list an uninstall reports (KOSD-C-29).
func (h *host) uninstall(identifier string) {
	declared, ok := h.components[identifier]
	if !ok {
		return
	}
	for class, action := range declared.UninstallRetention.Contributions {
		if action == manifest.Retain {
			h.leftovers = append(h.leftovers, "retained/"+class)
		}
	}
	for _, external := range declared.UninstallRetention.ExternalResources {
		h.leftovers = append(h.leftovers, "external/"+external)
	}
	sort.Strings(h.leftovers)
	delete(h.components, identifier)
	delete(h.digests, identifier)
}

func (h *host) component(identifier string) (manifest.ExtensionManifest, bool) {
	declared, ok := h.components[identifier]
	return declared, ok
}

// digestOf reports the package digest the installed component names, which is
// the digest its install request carried (KOSD-C-17).
func (h *host) digestOf(identifier string) packaging.Digest {
	return h.digests[identifier]
}

// installed names the installed components, so a case states which ones survive.
func (h *host) installed() []string {
	identifiers := make([]string, 0, len(h.components))
	for identifier := range h.components {
		identifiers = append(identifiers, identifier)
	}
	sort.Strings(identifiers)
	return identifiers
}

// surface lists what the installed components make visible, one "class/name"
// entry per declared contribution, in a stable order.
func (h *host) surface() []string {
	entries := make([]string, 0, len(h.components))
	for _, declared := range h.components {
		entries = append(entries, surfaceOf(declared)...)
	}
	sort.Strings(entries)
	return entries
}

// kernelSurface is what the installed components add to the kernel itself: the
// served kinds, the capability baseline and the status writers.
type kernelSurface struct {
	Kinds        []string
	Capabilities []string
	Writers      []string
}

func (h *host) kernel() kernelSurface {
	surface := kernelSurface{}
	for _, declared := range h.components {
		for _, class := range declared.Contributions.ObjectClasses {
			surface.Kinds = append(surface.Kinds, class.Kind)
		}
		for _, capability := range declared.Contributions.Capabilities {
			name, _ := capability["name"].(string)
			surface.Capabilities = append(surface.Capabilities, name)
		}
		for kind, status := range declared.Status.ObjectClasses {
			surface.Writers = append(surface.Writers, kind+"/"+status.Writer)
		}
	}
	sort.Strings(surface.Kinds)
	sort.Strings(surface.Capabilities)
	sort.Strings(surface.Writers)
	return surface
}

// surfaceOf lists what one declaration makes visible. A class the declaration
// does not list contributes nothing, so comparing this list with the expected
// surface checks the declaration against what the kernel serves (KOSD-C-17).
func surfaceOf(declared manifest.ExtensionManifest) []string {
	encoded, err := json.Marshal(declared.Contributions)
	if err != nil {
		return nil
	}
	record := make(map[string]any)
	if err := json.Unmarshal(encoded, &record); err != nil {
		return nil
	}
	entries := make([]string, 0, len(manifest.ContributionClasses))
	for _, class := range manifest.ContributionClasses {
		switch items := record[class].(type) {
		case []any:
			for _, item := range items {
				entries = append(entries, class+"/"+nameOf(item))
			}
		case map[string]any:
			if len(items) > 0 {
				entries = append(entries, class)
			}
		}
	}
	sort.Strings(entries)
	return entries
}

// nameOf reads the identifying field a contribution's payload carries; the
// classes name their items by id, name or kind (KOSD-C-18).
func nameOf(item any) string {
	fields, ok := item.(map[string]any)
	if !ok {
		return ""
	}
	for _, key := range []string{"id", "name", "kind"} {
		if name, ok := fields[key].(string); ok && name != "" {
			return name
		}
	}
	return ""
}

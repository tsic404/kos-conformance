// Package kosdtc155 verifies this product's slice of KOSD-TC-155: the artifact
// installs as one system component, serves exactly the surface its declaration
// names and nothing else, and uninstalls without residue while the other
// components and the kernel objects stay as they were (KOSD-AC-72, KOSD-C-17,
// KOSD-C-29).
package kosdtc155

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"slices"
	"testing"

	"github.com/tsic404/kos-sdk/manifest"
	"github.com/tsic404/kos-sdk/packaging"

	"github.com/tsic404/kos-conformance/release"
)

const (
	conformanceReference = "registry.kos.example/kos-conformance-0.1.0.tgz"
	companionReference   = "registry.kos.example/notes-1.0.0.tgz"
	companionID          = "kos.example.notes"
	companionPage        = "pages/notes"
)

// TestInstallServesTheDeclaredSurface pins the install half of KOSD-TC-155 for
// this product: the artifact passes the checks an install path runs, the
// component exists afterwards, and the kernel serves its declarations, which
// name no contribution at all.
func TestInstallServesTheDeclaredSurface(t *testing.T) {
	host := newHost()
	before := host.kernel()
	packaged := publishConformance(t, host)

	declared, err := host.install(conformanceReference, packaged.Digest)
	if err != nil {
		t.Fatalf("install %s: %v", conformanceReference, err)
	}
	if declared.ID != release.Identifier || declared.Version != release.Version {
		t.Fatalf("component = %s@%s, want %s@%s",
			declared.ID, declared.Version, release.Identifier, release.Version)
	}
	if _, found := host.component(release.Identifier); !found {
		t.Fatalf("component %s is not installed; installed = %v", release.Identifier, host.installed())
	}
	if recorded := host.digestOf(release.Identifier); recorded != packaged.Digest {
		t.Fatalf("component names digest %s, want the requested %s", recorded, packaged.Digest)
	}
	// The product's API surface is the runner operators execute: a page, a
	// window or any other contribution would appear here.
	if surface := host.surface(); len(surface) != 0 {
		t.Fatalf("visible surface = %v, want the empty surface its declaration names", surface)
	}
	assertKernelUnchanged(t, before, host.kernel())
}

// TestPackageCarriesTheDeclaredFieldSet pins the declaration the artifact
// distributes: exactly the manifest field set the SDK defines, with the empty
// blocks present rather than defaulted away, and no asset outside the
// declaration.
func TestPackageCarriesTheDeclaredFieldSet(t *testing.T) {
	_, archive, err := release.Package()
	if err != nil {
		t.Fatalf("pack %s: %v", release.Identifier, err)
	}
	entries, declared := packageContent(t, archive)
	if want := []string{"digest", "manifest.json"}; !slices.Equal(entries, want) {
		t.Fatalf("package entries = %v, want %v", entries, want)
	}
	fields := make([]string, 0, len(declared))
	for field := range declared {
		fields = append(fields, field)
	}
	slices.Sort(fields)
	want := slices.Clone(manifest.ExtensionManifestFields)
	slices.Sort(want)
	if !slices.Equal(fields, want) {
		t.Fatalf("manifest fields = %v, want %v", fields, want)
	}
	if id, _ := declared["id"].(string); id != release.Identifier {
		t.Fatalf("declared id = %q, want %q", id, release.Identifier)
	}
	contributions, _ := declared["contributions"].(map[string]any)
	if len(contributions) != 0 {
		t.Fatalf("declared contributions = %v, want none", contributions)
	}
}

// TestUninstallLeavesNoResidue pins the uninstall half of KOSD-TC-155: the
// component and its registrations disappear, nothing is retained, and the other
// component and the kernel stay as they were (KOSD-C-29).
func TestUninstallLeavesNoResidue(t *testing.T) {
	host := newHost()
	before := host.kernel()
	conformance := publishConformance(t, host)
	companion := publishCompanion(t, host)

	// Both artifacts install through the same request shape; nothing about the
	// author changes the channel (KOSD-C-17).
	if _, err := host.install(conformanceReference, conformance.Digest); err != nil {
		t.Fatalf("install %s: %v", release.Identifier, err)
	}
	if _, err := host.install(companionReference, companion.Digest); err != nil {
		t.Fatalf("install %s: %v", companionID, err)
	}
	if surface := host.surface(); !slices.Equal(surface, []string{companionPage}) {
		t.Fatalf("visible surface before uninstall = %v, want %v", surface, []string{companionPage})
	}

	host.uninstall(release.Identifier)

	if _, found := host.component(release.Identifier); found {
		t.Fatalf("component %s survived its uninstall", release.Identifier)
	}
	if installed := host.installed(); !slices.Equal(installed, []string{companionID}) {
		t.Fatalf("installed components after uninstall = %v, want %v", installed, []string{companionID})
	}
	if surface := host.surface(); !slices.Equal(surface, []string{companionPage}) {
		t.Fatalf("visible surface after uninstall = %v, want the other component's %v", surface, []string{companionPage})
	}
	if len(host.leftovers) != 0 {
		t.Fatalf("leftover list = %v, want none", host.leftovers)
	}
	assertKernelUnchanged(t, before, host.kernel())
}

// TestInstallRefusesAForeignDigest pins the install request: an artifact that
// does not match the digest the request names is refused, and nothing is
// installed (KOSD-C-17).
func TestInstallRefusesAForeignDigest(t *testing.T) {
	host := newHost()
	publishConformance(t, host)
	other := publishCompanion(t, host)

	if _, err := host.install(conformanceReference, other.Digest); !errors.Is(err, packaging.ErrDigestMismatch) {
		t.Fatalf("install with a foreign digest = %v, want %v", err, packaging.ErrDigestMismatch)
	}
	if installed := host.installed(); len(installed) != 0 {
		t.Fatalf("installed components = %v after a refused request, want none", installed)
	}
}

// publishConformance packs this product and publishes the artifact, returning
// the package whose digest a component would name.
func publishConformance(t *testing.T, host *host) packaging.Package {
	t.Helper()
	packaged, archive, err := release.Package()
	if err != nil {
		t.Fatalf("pack %s: %v", release.Identifier, err)
	}
	host.publish(conformanceReference, archive)
	return packaged
}

// publishCompanion publishes an unrelated extension whose only contribution is
// one page, so the uninstall case has another component to leave untouched.
func publishCompanion(t *testing.T, host *host) packaging.Package {
	t.Helper()
	packaged, archive, err := packaging.Pack(map[string]any{
		"id":                 companionID,
		"version":            "1.0.0",
		"frameworkVersion":   ">=1.1 <2.0",
		"contributions":      map[string]any{"pages": []any{map[string]any{"id": "notes"}}},
		"dependencies":       map[string]any{},
		"permissions":        map[string]any{},
		"config":             map[string]any{},
		"status":             map[string]any{},
		"uninstallRetention": map[string]any{},
	}, nil)
	if err != nil {
		t.Fatalf("pack %s: %v", companionID, err)
	}
	host.publish(companionReference, archive)
	return packaged
}

func assertKernelUnchanged(t *testing.T, before, after kernelSurface) {
	t.Helper()
	if !slices.Equal(before.Kinds, after.Kinds) {
		t.Fatalf("kernel kinds = %v, want %v", after.Kinds, before.Kinds)
	}
	if !slices.Equal(before.Capabilities, after.Capabilities) {
		t.Fatalf("capability baseline = %v, want %v", after.Capabilities, before.Capabilities)
	}
	if !slices.Equal(before.Writers, after.Writers) {
		t.Fatalf("writer table = %v, want %v", after.Writers, before.Writers)
	}
}

// packageContent reads the entry names a package carries and the declaration
// inside it, so a case asserts both the structure and the manifest it delivers.
func packageContent(t *testing.T, archive []byte) ([]string, map[string]any) {
	t.Helper()
	compressed, err := gzip.NewReader(bytes.NewReader(archive))
	if err != nil {
		t.Fatalf("read package archive: %v", err)
	}
	defer compressed.Close()
	names := make([]string, 0, 2)
	var declared map[string]any
	reader := tar.NewReader(compressed)
	for {
		header, err := reader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("read package archive: %v", err)
		}
		names = append(names, header.Name)
		if header.Name != packaging.ManifestEntry {
			continue
		}
		if err := json.NewDecoder(reader).Decode(&declared); err != nil {
			t.Fatalf("decode %s: %v", packaging.ManifestEntry, err)
		}
	}
	slices.Sort(names)
	if declared == nil {
		t.Fatalf("package carries no %s", packaging.ManifestEntry)
	}
	return names, declared
}

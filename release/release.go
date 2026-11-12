// Package release declares this product as an extension artifact: the
// declaration of KOSD-C-18 and the package KOSD-C-17 distributes. It supplies
// only the declaration; the SDK's one packing path produces the artifact.
package release

import (
	"encoding/json"

	"github.com/tsic404/kos-sdk/manifest"
	"github.com/tsic404/kos-sdk/packaging"
)

const (
	// Identifier is the extension identity this product installs under; the
	// kernel keeps at most one installed version per identity (KOSD-C-17).
	Identifier = "kos-conformance"
	// Version is this artifact's version, independent of the framework release.
	Version = "0.1.0"
	// FrameworkVersion is the framework API range the artifact targets; the
	// kernel serves major.minor framework releases (KOSD-C-16).
	FrameworkVersion = ">=1.1 <2.0"
)

// declaration is the declaration surface of KOSD-C-18. It contributes
// nothing: the only API surface is the test runner operators execute, so
// installing the component adds no page, window, object class, controller,
// capability or setting to the kernel.
func declaration() manifest.ExtensionManifest {
	return manifest.ExtensionManifest{
		ID:               Identifier,
		Version:          Version,
		FrameworkVersion: FrameworkVersion,
		Dependencies:     manifest.Dependencies{Framework: FrameworkVersion},
		// The empty blocks are declared rather than defaulted: the framework
		// requires each field and this product supplies no content for it
		// (SDK-C-3, KOSD-C-18).
		Contributions:      manifest.Contributions{},
		Permissions:        manifest.ExtensionPermissions{},
		Config:             map[string]any{},
		Status:             manifest.StatusDeclaration{},
		UninstallRetention: manifest.UninstallRetention{},
	}
}

// assets are the UI assets the package carries. The product has no UI, so the
// digest covers the declaration alone.
func assets() []packaging.Asset { return nil }

// Package builds the artifact through the SDK packing path every author uses.
// The declaration is marshalled from the framework's own type, so the package
// carries exactly the manifest field set the SDK defines.
func Package() (packaging.Package, []byte, error) {
	encoded, err := json.Marshal(declaration())
	if err != nil {
		return packaging.Package{}, nil, err
	}
	var declared map[string]any
	if err := json.Unmarshal(encoded, &declared); err != nil {
		return packaging.Package{}, nil, err
	}
	return packaging.Pack(declared, assets())
}

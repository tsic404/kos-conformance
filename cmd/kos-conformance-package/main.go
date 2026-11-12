// Command kos-conformance-package writes this product's extension artifact
// through the SDK packing path. It is author tooling: distributing the artifact
// and installing it as a system component use the standard channel, not this
// command.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tsic404/kos-conformance/release"
)

func main() {
	out := flag.String("out", "dist", "directory the artifact is written to")
	flag.Parse()

	pkg, archive, err := release.Package()
	if err != nil {
		fail("pack %s: %v", release.Identifier, err)
	}
	if err := os.MkdirAll(*out, 0o755); err != nil {
		fail("create output directory %s: %v", *out, err)
	}
	name := filepath.Join(*out, release.Identifier+"-"+release.Version+".tgz")
	if err := os.WriteFile(name, archive, 0o644); err != nil {
		fail("write artifact %s: %v", name, err)
	}
	fmt.Printf("artifact: %s\ndigest: %s\n", name, pkg.Digest)
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

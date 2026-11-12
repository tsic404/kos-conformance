// Command kos-conformance is the distribution's run entry: it prints the case
// set and its digest, and runs the case set against one provider under the
// operator's identity. It has no product page and writes no capability state.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/tsic404/kos-conformance/acceptance"
)

// registry holds the provider implementations this binary links. A distribution
// that ships implementations registers them here; this repository ships the
// suite alone, so a run reports every case as not run rather than passing it
// (CONF-C-2).
var registry = acceptance.Registry{}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "cases":
		err = cases(os.Args[2:])
	case "digest":
		err = digest()
	case "run":
		err = run(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "kos-conformance:", err)
		os.Exit(1)
	}
}

// usage states the three commands and the run flags.
func usage() {
	fmt.Fprint(os.Stderr, `usage: kos-conformance <command> [flags]

commands:
  cases    print the distribution case set and its digest
  digest   print the case-set digest a run is compared with
  run      run the case set against one provider

run flags:
  -provider string              provider under test (required)
  -distribution-digest string   the distribution's recorded digest (required)
  -json                         print the report as JSON
`)
}

// cases prints the shipped case set with its digest.
func cases(_ []string) error {
	return acceptance.Cases(os.Stdout, acceptance.Standard())
}

// digest prints the case-set digest alone, which a distribution records.
func digest() error {
	fmt.Println(acceptance.Standard().Digest())
	return nil
}

// run executes the case set against one provider and writes the report.
func run(args []string) error {
	flags := flag.NewFlagSet("run", flag.ContinueOnError)
	provider := flags.String("provider", "", "provider under test")
	distributionDigest := flags.String("distribution-digest", "", "the distribution's recorded digest")
	asJSON := flags.Bool("json", false, "print the report as JSON")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *provider == "" {
		return fmt.Errorf("run needs -provider")
	}
	if *distributionDigest == "" {
		return fmt.Errorf("run needs -distribution-digest: a run without the record cannot be comparable")
	}
	report := acceptance.NewRunner(acceptance.Standard(), registry).Run(context.Background(), *provider, *distributionDigest)
	var err error
	if *asJSON {
		err = acceptance.JSON(os.Stdout, report)
	} else {
		err = acceptance.Text(os.Stdout, report)
	}
	if err != nil {
		return err
	}
	if report.Verdict != acceptance.VerdictPassed {
		os.Exit(1) // the report is written; a caller scripting the entry sees the conclusion
	}
	return nil
}

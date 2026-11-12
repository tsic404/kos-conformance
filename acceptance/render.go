package acceptance

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/tsic404/kos-conformance"
)

// Cases writes the case set a distribution ships: its digest, then every case
// with its group, set, objectives and criterion version.
func Cases(w io.Writer, distribution *conformance.Distribution) error {
	cases := distribution.Cases()
	setWidth, caseWidth := 0, 0
	for _, c := range cases {
		setWidth = max(setWidth, len(c.SetID))
		caseWidth = max(caseWidth, len(c.ID))
	}
	var out strings.Builder
	fmt.Fprintf(&out, "digest: %s\n", distribution.Digest())
	for _, group := range conformance.Groups() {
		fmt.Fprintf(&out, "\n%s\n", group)
		for _, c := range cases {
			if c.Group != group {
				continue
			}
			fmt.Fprintf(&out, "  %-*s  %-*s  %s@%s\n", setWidth, c.SetID, caseWidth, c.ID, strings.Join(c.Objectives, ","), c.CriterionVersion)
		}
	}
	_, err := io.WriteString(w, out.String())
	return err
}

// Text writes one report: both digests, the verdict, every case row, and the
// failed and not run lists kept apart (CONF-C-2).
func Text(w io.Writer, report Report) error {
	var out strings.Builder
	fmt.Fprintf(&out, "provider:            %s\n", report.Provider)
	fmt.Fprintf(&out, "run digest:          %s\n", report.RunDigest)
	fmt.Fprintf(&out, "distribution digest: %s\n", report.DistributionDigest)
	fmt.Fprintf(&out, "verdict:             %s\n", report.Verdict)
	for _, group := range conformance.Groups() {
		fmt.Fprintf(&out, "\n%s\n", group)
		for _, result := range report.Results {
			if result.Case.Group != group {
				continue
			}
			fmt.Fprintf(&out, "  %-7s  %s %s  %s@%s", result.Status, result.Case.SetID, result.Case.ID, strings.Join(result.Case.Objectives, ","), result.Case.CriterionVersion)
			if result.Detail != "" {
				fmt.Fprintf(&out, "  observed: %s", result.Detail)
			}
			out.WriteString("\n")
		}
	}
	writeList(&out, "failed", report.Failed())
	writeList(&out, "not run", report.NotRun())
	_, err := io.WriteString(w, out.String())
	return err
}

// writeList renders one of the two lists a report keeps apart.
func writeList(out *strings.Builder, title string, results []Result) {
	fmt.Fprintf(out, "\n%s (%d)\n", title, len(results))
	for _, result := range results {
		fmt.Fprintf(out, "  [%s] %s %s (%s@%s): %s\n", result.Case.Group, result.Case.SetID, result.Case.ID, strings.Join(result.Case.Objectives, ","), result.Case.CriterionVersion, result.Detail)
	}
}

// jsonReport is the machine-readable shape of one report: the same rows as
// Text, with the failed and not run lists written out separately.
type jsonReport struct {
	Provider           string       `json:"provider"`
	DistributionDigest string       `json:"distributionDigest"`
	RunDigest          string       `json:"runDigest"`
	Verdict            Verdict      `json:"verdict"`
	Results            []jsonResult `json:"results"`
	Failed             []jsonResult `json:"failed"`
	NotRun             []jsonResult `json:"notRun"`
}

// jsonResult is one case row in the machine-readable report.
type jsonResult struct {
	Group            conformance.Group `json:"group"`
	Set              string            `json:"set"`
	Capability       string            `json:"capability,omitempty"`
	Version          string            `json:"version"`
	Case             string            `json:"case"`
	Description      string            `json:"description"`
	Objectives       []string          `json:"objectives"`
	CriterionVersion string            `json:"criterionVersion"`
	Status           Status            `json:"status"`
	Detail           string            `json:"detail,omitempty"`
}

// JSON writes one report as JSON.
func JSON(w io.Writer, report Report) error {
	view := jsonReport{
		Provider:           report.Provider,
		DistributionDigest: report.DistributionDigest,
		RunDigest:          report.RunDigest,
		Verdict:            report.Verdict,
		Results:            jsonResults(report.Results),
		Failed:             jsonResults(report.Failed()),
		NotRun:             jsonResults(report.NotRun()),
	}
	encoded, err := json.MarshalIndent(view, "", "  ")
	if err != nil {
		return fmt.Errorf("encode the report: %w", err)
	}
	encoded = append(encoded, '\n')
	_, err = w.Write(encoded)
	return err
}

// jsonResults converts rows for encoding.
func jsonResults(results []Result) []jsonResult {
	rows := make([]jsonResult, 0, len(results))
	for _, result := range results {
		rows = append(rows, jsonResult{
			Group:            result.Case.Group,
			Set:              result.Case.SetID,
			Capability:       result.Case.Capability,
			Version:          result.Case.Version,
			Case:             result.Case.ID,
			Description:      result.Case.Description,
			Objectives:       result.Case.Objectives,
			CriterionVersion: result.Case.CriterionVersion,
			Status:           result.Status,
			Detail:           result.Detail,
		})
	}
	return rows
}

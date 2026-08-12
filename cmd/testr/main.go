package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/eidos-agi/testr/internal/testr"
)

func printOut(v any, asJSON bool) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if asJSON {
		_ = enc.Encode(v)
		return
	}
	b, _ := json.MarshalIndent(v, "", "  ")
	var m map[string]any
	if json.Unmarshal(b, &m) == nil {
		for k, val := range m {
			fmt.Printf("%s: %v\n", k, val)
		}
		return
	}
	fmt.Println(string(b))
}

func usage() {
	fmt.Fprintf(os.Stderr, `testr %s — AI testing config + proof memory (Go)

testr does NOT run tests.
It stores how-this-product-is-proven config so AI agents can test repeatedly.

2027 kickstart contract:
  • Keyword "testr" means: load THIS product's committed model and obey it.
  • Committed .testr/product-test-model.json ALWAYS wins over auto-detect.
  • Detection is greenfield-only. Never treat detect output as product policy.
  • Prefer path-relevant tests. Do not invent full-suite ceremony for every change.

Usage:
  testr model [--project PATH] [--description TEXT] [--write] [--force] [--json]
  testr attempt --goal TEXT [--project PATH]
                [--status planned|running|passed|failed|blocked|skipped]
                [--notes TEXT] [--proof TEXT ...] [--json]
  testr frontier [--project PATH] [--json]

  model          Print the model agents should use (committed file if present).
  model --write  Create model only if missing. Refuses to clobber without --force.
  model --force  Required with --write to replace an existing committed model.

Config file:  .testr/product-test-model.json  (committed; not gitignored)  ← SOURCE OF TRUTH
Ledger:       .testr/test-attempts/              (committed)
Sibling:      shipr (.shipr/) — proofs align with test_commands

Agent workflow:
  1. testr model --project . --json     # reads committed model; does NOT re-ceremony
  2. YOU run path-relevant test_commands only
  3. testr attempt --goal "…" --status passed --proof "…" --json
`, testr.Version)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	if os.Args[1] == "--version" || os.Args[1] == "version" {
		fmt.Println("testr", testr.Version)
		return
	}
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		usage()
		return
	}
	cmd := os.Args[1]
	args := os.Args[2:]
	project, desc, goal, status, notes := ".", "", "", "planned", ""
	asJSON, write, force := false, false, false
	var proofs []string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--json":
			asJSON = true
		case "--write":
			write = true
		case "--force":
			force = true
		case "--project":
			i++
			if i < len(args) {
				project = args[i]
			}
		case "--description":
			i++
			if i < len(args) {
				desc = args[i]
			}
		case "--goal":
			i++
			if i < len(args) {
				goal = args[i]
			}
		case "--status":
			i++
			if i < len(args) {
				status = args[i]
			}
		case "--notes":
			i++
			if i < len(args) {
				notes = args[i]
			}
		case "--proof":
			i++
			if i < len(args) {
				proofs = append(proofs, args[i])
			}
		case "-h", "--help":
			usage()
			return
		}
	}
	switch cmd {
	case "model":
		var model testr.TestModel
		var source string
		if force && write {
			model = testr.DetectTestModel(project, desc)
			source = "detected"
			model["model_source"] = "detected"
		} else {
			model, source = testr.ResolveTestModel(project, desc)
		}
		if write {
			path, err := testr.WriteTestModelForced(project, model, force)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			model["written_to"] = path
			model["model_source"] = source
			if force {
				model["model_source"] = "detected"
				model["write_mode"] = "forced_replace"
			} else {
				model["write_mode"] = "created"
			}
		}
		printOut(model, asJSON)
	case "attempt":
		if goal == "" {
			fmt.Fprintln(os.Stderr, "testr attempt requires --goal")
			os.Exit(1)
		}
		path, attempt, err := testr.RecordAttempt(project, goal, status, notes, proofs)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		attempt["written_to"] = path
		printOut(attempt, asJSON)
	case "frontier":
		printOut(testr.Frontier(project), asJSON)
	default:
		usage()
		os.Exit(2)
	}
}

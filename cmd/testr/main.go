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

Usage:
  testr model [--project PATH] [--description TEXT] [--write] [--json]
  testr attempt --goal TEXT [--project PATH]
                [--status planned|running|passed|failed|blocked|skipped]
                [--notes TEXT] [--proof TEXT ...] [--json]
  testr frontier [--project PATH] [--json]

Config file:  .testr/product-test-model.json
Ledger:       .testr/test-attempts/
Sibling:      shipr (.shipr/) — uses test_commands as proof_commands when shipping
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
	asJSON, write := false, false
	var proofs []string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--json":
			asJSON = true
		case "--write":
			write = true
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
		model := testr.DetectTestModel(project, desc)
		if write {
			path, err := testr.WriteTestModel(project, model)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			model["written_to"] = path
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

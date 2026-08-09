"""Testr CLI — per-product test model and run memory."""

from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path

from . import __version__
from .core import detect_test_model, frontier, record_attempt, write_test_model


def _print(payload: dict, as_json: bool) -> None:
    if as_json:
        print(json.dumps(payload, indent=2, sort_keys=True))
        return
    for key, value in payload.items():
        print(f"{key}: {value}")


def main(argv: list[str] | None = None) -> None:
    parser = argparse.ArgumentParser(
        prog="testr",
        description="Persistent Eidos testing operator — learns how each product is proven.",
    )
    parser.add_argument("--version", action="version", version=f"testr {__version__}")
    sub = parser.add_subparsers(dest="cmd", required=True)

    m = sub.add_parser("model", help="Detect or refresh a product test model")
    m.add_argument("--project", type=Path, default=Path("."))
    m.add_argument("--description", default="")
    m.add_argument("--write", action="store_true")
    m.add_argument("--json", action="store_true")

    a = sub.add_parser("attempt", help="Record a test attempt")
    a.add_argument("--project", type=Path, default=Path("."))
    a.add_argument("--goal", required=True)
    a.add_argument(
        "--status",
        choices=["planned", "running", "passed", "failed", "blocked", "skipped"],
        default="planned",
    )
    a.add_argument("--notes", default="")
    a.add_argument("--proof", action="append", default=[])
    a.add_argument("--json", action="store_true")

    f = sub.add_parser("frontier", help="Show current test frontier")
    f.add_argument("--project", type=Path, default=Path("."))
    f.add_argument("--json", action="store_true")

    args = parser.parse_args(argv)
    if args.cmd == "model":
        model = detect_test_model(args.project, args.description or "")
        if args.write:
            path = write_test_model(args.project, model)
            model["written_to"] = str(path)
        _print(model, args.json)
    elif args.cmd == "attempt":
        path, attempt = record_attempt(
            args.project,
            goal=args.goal,
            status=args.status,
            notes=args.notes,
            proofs=args.proof or [],
        )
        attempt["written_to"] = str(path)
        _print(attempt, args.json)
    elif args.cmd == "frontier":
        _print(frontier(args.project), args.json)
    else:
        raise SystemExit(2)


if __name__ == "__main__":
    main()

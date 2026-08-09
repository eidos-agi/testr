"""Testr core — product test models and attempt ledger."""

from __future__ import annotations

import datetime as dt
import json
import re
import subprocess
from pathlib import Path
from typing import Any

MODEL_PATH = Path(".testr/product-test-model.json")
ATTEMPTS_DIR = Path(".testr/test-attempts")


def ensure_testr_ignored(root: Path) -> None:
    gi = root / ".gitignore"
    line = ".testr/"
    if gi.exists():
        text = gi.read_text()
        if line not in text.splitlines() and ".testr" not in text:
            with gi.open("a") as f:
                if text and not text.endswith("\n"):
                    f.write("\n")
                f.write(line + "\n")
    else:
        gi.write_text(line + "\n")


def _exists(root: Path, *parts: str) -> bool:
    return (root.joinpath(*parts)).exists()


def detect_test_model(project: Path, description: str = "") -> dict[str, Any]:
    root = project.resolve()
    product = root.name
    suites: list[str] = []
    commands: list[str] = []
    evidence: list[str] = []

    if _exists(root, "pyproject.toml") or _exists(root, "pytest.ini") or _exists(root, "tests"):
        suites.append("pytest")
        commands.append("python -m pytest -q")
    if _exists(root, "package.json"):
        suites.append("node")
        commands.append("npm test")
    if list(root.glob("**/go.mod")) or _exists(root, "go.mod"):
        suites.append("go")
        commands.append("go test ./...")
    if _exists(root, "docs", "emf") or _exists(root, "docs", "emf", "index.md"):
        suites.append("emf")
        commands.append(
            "cd ~/repos-eidos-agi/emf && PYTHONPATH=. python3 -m emf.validate "
            f"{root}/docs/emf/"
        )
    if _exists(root, "Makefile"):
        text = (root / "Makefile").read_text(errors="ignore")
        if re.search(r"^test:", text, re.M):
            suites.append("make-test")
            commands.append("make test")

    if not commands:
        commands.append("define product-specific test command before claiming green")

    return {
        "schema_version": 1,
        "product_id": product,
        "project_root": str(root),
        "description": description,
        "test_suites": sorted(set(suites)) or ["unknown"],
        "test_commands": list(dict.fromkeys(commands)),
        "evidence_paths": evidence or ["test output / junit / coverage when configured"],
        "methods_source": "test-forge (retired) → testr",
        "related_operators": ["shipr"],
        "learning_questions": [
            "What failed that the ship path assumed green?",
            "Which test should become a shipr proof_command?",
            "What is flaky vs broken?",
        ],
        "memory_paths": {
            "model": str(MODEL_PATH),
            "attempts_dir": str(ATTEMPTS_DIR),
        },
        "updated_at": dt.datetime.now(dt.timezone.utc).isoformat(),
    }


def write_test_model(project: Path, model: dict[str, Any]) -> Path:
    root = project.resolve()
    ensure_testr_ignored(root)
    path = root / MODEL_PATH
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(model, indent=2, sort_keys=True) + "\n")
    return path


def record_attempt(
    project: Path,
    *,
    goal: str,
    status: str,
    notes: str = "",
    proofs: list[str] | None = None,
) -> tuple[Path, dict[str, Any]]:
    root = project.resolve()
    ensure_testr_ignored(root)
    attempts = root / ATTEMPTS_DIR
    attempts.mkdir(parents=True, exist_ok=True)
    model_path = root / MODEL_PATH
    model = json.loads(model_path.read_text()) if model_path.exists() else detect_test_model(root)
    stamp = dt.datetime.now(dt.timezone.utc).strftime("%Y%m%dT%H%M%SZ")
    slug = re.sub(r"[^a-z0-9]+", "-", goal.lower()).strip("-")[:60] or "test"
    path = attempts / f"{stamp}-{slug}.json"
    attempt = {
        "schema_version": 1,
        "product_id": model.get("product_id", root.name),
        "goal": goal,
        "status": status,
        "notes": notes,
        "proofs": proofs or [],
        "test_commands_snapshot": model.get("test_commands", []),
        "created_at": dt.datetime.now(dt.timezone.utc).isoformat(),
    }
    path.write_text(json.dumps(attempt, indent=2, sort_keys=True) + "\n")
    return path, attempt


def frontier(project: Path) -> dict[str, Any]:
    root = project.resolve()
    model_path = root / MODEL_PATH
    model = json.loads(model_path.read_text()) if model_path.exists() else detect_test_model(root)
    attempts_dir = root / ATTEMPTS_DIR
    latest = None
    if attempts_dir.exists():
        files = sorted(attempts_dir.glob("*.json"), reverse=True)
        if files:
            latest = json.loads(files[0].read_text())
    return {
        "product_id": model.get("product_id"),
        "model_path": str(model_path),
        "test_commands": model.get("test_commands", []),
        "latest_status": (latest or {}).get("status"),
        "latest_attempt": str(attempts_dir / latest["goal"]) if latest else None,
        "latest": latest,
        "status": "model_ready" if model_path.exists() else "model_missing",
        "related": {"shipr": "use test_commands as shipr proof_commands when shipping"},
    }

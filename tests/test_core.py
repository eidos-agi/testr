from pathlib import Path

from testr.core import detect_test_model, write_test_model, record_attempt, frontier


def test_detect_and_write(tmp_path: Path):
    (tmp_path / "tests").mkdir()
    (tmp_path / "pyproject.toml").write_text("[project]\nname='x'\n")
    model = detect_test_model(tmp_path, "demo")
    assert "pytest" in model["test_suites"]
    path = write_test_model(tmp_path, model)
    assert path.exists()
    assert ".testr/" not in (tmp_path / ".gitignore").read_text() if (tmp_path / ".gitignore").exists() else True
    p, attempt = record_attempt(tmp_path, goal="unit", status="passed", proofs=["pytest -q"])
    assert p.exists()
    assert attempt["status"] == "passed"
    front = frontier(tmp_path)
    assert front["latest_status"] == "passed"

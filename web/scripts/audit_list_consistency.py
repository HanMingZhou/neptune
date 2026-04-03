from __future__ import annotations

from pathlib import Path


ROOT = Path(__file__).resolve().parents[1] / "src" / "view"

MANUAL_CHIP_TOKENS = (
    "bg-emerald-500/10",
    "bg-blue-500/10",
    "bg-red-500/10",
    "bg-amber-500/10",
    "rounded-full text-xs font-bold",
)

VALID_SHELL_TOKENS = (
    "TableCard",
    "ManagementListShell",
    "console-detail-card",
    "detail-table-shell",
    "console-list-footer",
)


def classify(path: Path, text: str) -> list[str]:
    issues: list[str] = []

    is_table = "console-table" in text
    has_actions = "list-row-button" in text
    has_action_wrap = "list-row-actions" in text
    has_action_header = "console-actions-header" in text
    has_action_cell = "console-actions-cell" in text
    uses_table_shell = any(token in text for token in VALID_SHELL_TOKENS)
    has_pagination = "ListPaginationBar" in text
    has_manual_chip = any(token in text for token in MANUAL_CHIP_TOKENS)

    if is_table and has_actions and not has_action_wrap:
        issues.append("action-buttons-outside-list-row-actions")
    if is_table and has_actions and not has_action_header:
        issues.append("missing-console-actions-header")
    if is_table and has_actions and not has_action_cell:
        issues.append("missing-console-actions-cell")
    if is_table and not uses_table_shell:
        issues.append("table-not-using-shared-card-shell")
    if is_table and not has_pagination and "totalRecords" in text:
        issues.append("table-has-total-but-no-shared-pagination")
    if is_table and has_manual_chip:
        issues.append("manual-status-chip-markup")

    return issues


def main() -> int:
    findings: list[tuple[Path, list[str]]] = []
    for path in sorted(ROOT.rglob("*.vue")):
        try:
            text = path.read_text(encoding="utf-8")
        except UnicodeDecodeError:
            continue
        issues = classify(path, text)
        if issues:
            findings.append((path, issues))

    if not findings:
        print("No list consistency issues found.")
        return 0

    print("List consistency audit")
    print("=" * 80)
    for path, issues in findings:
        print(path)
        for issue in issues:
            print(f"  - {issue}")
    print("=" * 80)
    print(f"Total files flagged: {len(findings)}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

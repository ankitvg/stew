#!/usr/bin/env python3
"""One-off importer for old StewReads docs ledgers.

This is intentionally not a public Stew command. It handles the historical
formats used in StewReads repos before the Stew CLI existed, plus old
monolithic `.stew/<ledger>.md` files.
"""

from __future__ import annotations

import argparse
import os
import re
import subprocess
import sys
from collections import Counter
from dataclasses import dataclass
from datetime import datetime, timezone
from pathlib import Path
from zoneinfo import ZoneInfo


LEDGERS = ("iterations", "decisions")
UTC_RE = r"20\d\d-\d\d-\d\dT\d\d:\d\d:\d\dZ"
DATE_RE = r"20\d\d-\d\d-\d\d"
CANONICAL_HEADING = re.compile(rf"^## ({UTC_RE}) — (.+)$", re.M)
LEGACY_HEADINGS = [
    re.compile(rf"^(?P<marker>###) \[(?P<stamp>{UTC_RE})\] - (?P<summary>.+)$"),
    re.compile(rf"^(?P<marker>###) \[(?P<stamp>{DATE_RE})\] - (?P<summary>.+)$"),
    re.compile(rf"^(?P<marker>###) (?P<stamp>{UTC_RE}) - (?P<summary>.+)$"),
    re.compile(rf"^(?P<marker>###) (?P<stamp>{DATE_RE}) - (?P<summary>.+)$"),
    re.compile(rf"^(?P<marker>##) (?P<stamp>{UTC_RE})(?: [—-] (?P<summary>.+))?$"),
    re.compile(rf"^(?P<marker>##) (?P<stamp>{DATE_RE} \d\d:\d\d (?:EDT|EST))(?: [—-] (?P<summary>.+))?$"),
    re.compile(rf"^(?P<marker>##) (?P<stamp>{DATE_RE})(?: [—-] (?P<summary>.+))?$"),
]
MANAGED_START = "<!-- BEGIN STEW (managed) -->"
MANAGED_END = "<!-- END STEW (managed) -->"
MANAGED_BLOCK = """<!-- BEGIN STEW (managed) -->
## Stew

This repo uses stew to maintain durable project memory in append-only markdown ledger entries.

Run `stew help` to discover available commands.
Run `stew full-spec` before working to load the workflow and ledger contract.
<!-- END STEW (managed) -->
"""


@dataclass
class Entry:
    ledger: str
    timestamp: str
    summary: str
    content: str
    source_path: Path
    original_heading: str
    exact_timestamp: bool


@dataclass
class ImportPlan:
    repo: Path
    source_kind: str
    entries: list[Entry]
    legacy_paths: list[Path]


def main() -> int:
    parser = argparse.ArgumentParser(description="Import legacy StewReads ledger docs into atomic .stew entries.")
    parser.add_argument("repos", nargs="+", type=Path, help="Repository roots to inspect/import")
    parser.add_argument("--stew-bin", type=Path, default=Path("dist/stew"), help="Path to the built stew binary")
    mode = parser.add_mutually_exclusive_group()
    mode.add_argument("--dry-run", action="store_true", help="Inspect only; this is the default")
    mode.add_argument("--apply", action="store_true", help="Write .stew atomic entries")
    parser.add_argument("--remove-legacy", action="store_true", help="Remove imported legacy ledger files after write")
    args = parser.parse_args()

    apply = bool(args.apply)
    stew_bin = args.stew_bin.resolve()
    if apply and not stew_bin.exists():
        raise SystemExit(f"missing stew binary: {stew_bin}")

    plans = [build_plan(repo.resolve()) for repo in args.repos]
    for plan in plans:
        print_plan(plan)

    if not apply:
        print("\nDry run only. Re-run with --apply to write.")
        return 0

    for plan in plans:
        apply_plan(plan, stew_bin, remove_legacy=args.remove_legacy)
    return 0


def build_plan(repo: Path) -> ImportPlan:
    if not (repo / ".git").exists():
        raise SystemExit(f"not a git checkout: {repo}")

    monolithic = [repo / ".stew" / f"{ledger}.md" for ledger in LEDGERS if (repo / ".stew" / f"{ledger}.md").exists()]
    docs = [repo / "docs" / f"{ledger}.md" for ledger in LEDGERS if (repo / "docs" / f"{ledger}.md").exists()]

    if monolithic:
        entries = []
        for path in monolithic:
            entries.extend(parse_canonical(path.stem, path))
        return ImportPlan(repo=repo, source_kind="stew-monolithic", entries=entries, legacy_paths=monolithic)

    if docs:
        entries = []
        for path in docs:
            entries.extend(parse_legacy_docs(path.stem, path, repo))
        return ImportPlan(repo=repo, source_kind="pre-cli-docs", entries=entries, legacy_paths=docs)

    raise SystemExit(f"no legacy ledger files found in {repo}")


def parse_canonical(ledger: str, path: Path) -> list[Entry]:
    text = path.read_text()
    matches = list(CANONICAL_HEADING.finditer(text))
    if not matches:
        if text.strip() and not has_only_stew_preamble(text):
            raise SystemExit(f"{path} has no canonical entries and non-preamble content")
        return []

    entries = []
    for i, match in enumerate(matches):
        start = match.start()
        end = matches[i + 1].start() if i + 1 < len(matches) else len(text)
        content = text[start:end].strip() + "\n"
        entries.append(
            Entry(
                ledger=ledger,
                timestamp=match.group(1),
                summary=match.group(2).strip(),
                content=content,
                source_path=path,
                original_heading=match.group(0),
                exact_timestamp=True,
            )
        )
    return entries


def parse_legacy_docs(ledger: str, path: Path, repo: Path) -> list[Entry]:
    text = path.read_text()
    heading_matches = []
    for line_match in re.finditer(r"^(?:##|###) .+$", text, flags=re.M):
        parsed = parse_legacy_heading(line_match.group(0), ledger)
        if parsed is not None:
            heading_matches.append((line_match, parsed))

    if not heading_matches:
        if text.strip():
            raise SystemExit(f"{path} has content but no recognized legacy headings")
        return []

    preamble = text[: heading_matches[0][0].start()].strip()
    if preamble:
        raise SystemExit(f"{path} has preamble content before first entry; refusing to drop it")

    entries = []
    for i, (match, parsed) in enumerate(heading_matches):
        start = match.start()
        end = heading_matches[i + 1][0].start() if i + 1 < len(heading_matches) else len(text)
        legacy_content = text[start:end].strip() + "\n"
        timestamp, summary, exact = parsed
        rendered = render_imported_entry(timestamp, summary, path.relative_to(repo), legacy_content)
        entries.append(
            Entry(
                ledger=ledger,
                timestamp=timestamp,
                summary=summary,
                content=rendered,
                source_path=path,
                original_heading=match.group(0),
                exact_timestamp=exact,
            )
        )
    return entries


def parse_legacy_heading(line: str, ledger: str) -> tuple[str, str, bool] | None:
    for pattern in LEGACY_HEADINGS:
        match = pattern.match(line)
        if not match:
            continue
        stamp = match.group("stamp")
        summary = (match.groupdict().get("summary") or f"Legacy {ledger} entry").strip()
        timestamp, exact = normalize_timestamp(stamp)
        return timestamp, summary, exact
    return None


def normalize_timestamp(stamp: str) -> tuple[str, bool]:
    if re.fullmatch(UTC_RE, stamp):
        return stamp, True

    date_time = re.fullmatch(rf"({DATE_RE}) (\d\d):(\d\d) (EDT|EST)", stamp)
    if date_time:
        date_part, hour, minute, _tz_label = date_time.groups()
        year, month, day = [int(part) for part in date_part.split("-")]
        local = datetime(year, month, day, int(hour), int(minute), tzinfo=ZoneInfo("America/New_York"))
        return local.astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"), False

    if re.fullmatch(DATE_RE, stamp):
        return f"{stamp}T00:00:00Z", False

    raise ValueError(f"unsupported timestamp: {stamp}")


def render_imported_entry(timestamp: str, summary: str, source_path: Path, legacy_content: str) -> str:
    return (
        f"## {timestamp} — {summary}\n\n"
        f"**Prompt:** Legacy import from {source_path.as_posix()}\n\n"
        f"{legacy_content.rstrip()}\n\n"
        "---\n"
    )


def apply_plan(plan: ImportPlan, stew_bin: Path, remove_legacy: bool) -> None:
    run([str(stew_bin), "init", "--path", str(plan.repo), "--no-agents-md", "--quiet"])
    update_agents(plan.repo)

    by_ledger: dict[str, list[Entry]] = {ledger: [] for ledger in LEDGERS}
    for entry in plan.entries:
        by_ledger.setdefault(entry.ledger, []).append(entry)

    for ledger, entries in by_ledger.items():
        entry_dir = plan.repo / ".stew" / "ledgers" / ledger
        existing = list(entry_dir.glob("*.md")) if entry_dir.exists() else []
        if existing:
            raise SystemExit(f"{entry_dir} already has {len(existing)} markdown entries; refusing duplicate import")
        entry_dir.mkdir(parents=True, exist_ok=True)
        for entry in entries:
            write_entry(entry_dir, entry)

    verify_import(plan)

    if remove_legacy:
        for path in plan.legacy_paths:
            if path.exists():
                path.unlink()


def update_agents(repo: Path) -> None:
    path = repo / "AGENTS.md"
    if not path.exists():
        path.write_text(MANAGED_BLOCK)
        return

    existing = path.read_text()
    if MANAGED_START in existing and MANAGED_END in existing:
        start = existing.index(MANAGED_START)
        end = existing.index(MANAGED_END, start) + len(MANAGED_END)
        next_text = existing[:start] + MANAGED_BLOCK.rstrip() + existing[end:]
        path.write_text(next_text.rstrip() + "\n")
        return

    if "docs/iterations.md" in existing and "docs/decisions.md" in existing:
        path.write_text(MANAGED_BLOCK)
        return

    path.write_text(existing.rstrip() + "\n\n" + MANAGED_BLOCK)


def write_entry(entry_dir: Path, entry: Entry) -> None:
    suffix = 1
    while True:
        name = filename(entry.timestamp, entry.summary, suffix)
        path = entry_dir / name
        try:
            fd = os.open(path, os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o644)
        except FileExistsError:
            suffix += 1
            continue
        with os.fdopen(fd, "w") as handle:
            handle.write(entry.content.rstrip() + "\n")
        return


def verify_import(plan: ImportPlan) -> None:
    expected = Counter((entry.ledger, entry.timestamp, entry.summary) for entry in plan.entries)
    actual = Counter()
    for ledger in LEDGERS:
        entry_dir = plan.repo / ".stew" / "ledgers" / ledger
        for path in sorted(entry_dir.glob("*.md")):
            text = path.read_text()
            match = CANONICAL_HEADING.search(text)
            if not match:
                raise SystemExit(f"imported file has no canonical heading: {path}")
            actual[(ledger, match.group(1), match.group(2).strip())] += 1
    if sum(actual.values()) != len(plan.entries):
        raise SystemExit(
            f"verification failed for {plan.repo}: wrote {sum(actual.values())} entries, expected {len(plan.entries)}"
        )
    missing = expected - actual
    extra = actual - expected
    if missing or extra:
        raise SystemExit(
            f"verification failed for {plan.repo}: missing {sum(missing.values())}, extra {sum(extra.values())}"
        )


def filename(timestamp: str, summary: str, suffix: int) -> str:
    compact = datetime.strptime(timestamp, "%Y-%m-%dT%H:%M:%SZ").strftime("%Y-%m-%dT%H%M%SZ")
    base = f"{compact}-{slug(summary)}"
    if suffix > 1:
        base += f"-{suffix}"
    return base + ".md"


def slug(value: str) -> str:
    chars = []
    previous_hyphen = False
    for char in value.lower():
        if "a" <= char <= "z" or "0" <= char <= "9":
            chars.append(char)
            previous_hyphen = False
        elif not previous_hyphen:
            chars.append("-")
            previous_hyphen = True
    result = "".join(chars).strip("-")
    if not result:
        return "entry"
    result = result[:80].rstrip("-")
    return result or "entry"


def has_only_stew_preamble(text: str) -> bool:
    for raw in text.splitlines():
        line = raw.strip()
        if not line or line.startswith("# ") or line == "<!-- Managed by stew -->":
            continue
        return False
    return True


def print_plan(plan: ImportPlan) -> None:
    total = len(plan.entries)
    exact = sum(1 for entry in plan.entries if entry.exact_timestamp)
    print(f"{plan.repo}")
    print(f"  source: {plan.source_kind}")
    print(f"  entries: {total} ({exact} exact UTC timestamps, {total - exact} normalized timestamps)")
    for ledger in LEDGERS:
        count = sum(1 for entry in plan.entries if entry.ledger == ledger)
        print(f"  {ledger}: {count}")
    for path in plan.legacy_paths:
        print(f"  legacy: {path.relative_to(plan.repo)}")


def run(command: list[str]) -> None:
    completed = subprocess.run(command, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    if completed.returncode != 0:
        raise SystemExit(
            f"command failed ({completed.returncode}): {' '.join(command)}\n"
            f"stdout:\n{completed.stdout}\n"
            f"stderr:\n{completed.stderr}"
        )


if __name__ == "__main__":
    sys.exit(main())

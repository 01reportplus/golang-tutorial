"""Validate documentation navigation, local Markdown links, and snippet inputs."""
from pathlib import Path
import re
import sys
from urllib.parse import unquote, urlsplit

import yaml

ROOT = Path(__file__).resolve().parents[1]
DOCS = ROOT / "docs"


def leaves(value):
    if isinstance(value, str):
        yield value
    elif isinstance(value, list):
        for item in value:
            yield from leaves(item)
    elif isinstance(value, dict):
        for item in value.values():
            yield from leaves(item)


def prose_only(text):
    text = re.sub(r"^```[^\n]*\n.*?^```\s*$", "", text, flags=re.M | re.S)
    return re.sub(r"`[^`\n]*`", "", text)


def main():
    config = yaml.safe_load((ROOT / "mkdocs.yml").read_text(encoding="utf-8"))
    nav = list(leaves(config["nav"]))
    pages = sorted(DOCS.rglob("*.md"))
    errors = []
    for entry in nav:
        if not (DOCS / entry).is_file():
            errors.append(f"missing navigation target: {entry}")
    if len(nav) != len(set(nav)):
        errors.append("duplicate navigation targets")
    for page in pages:
        rel = page.relative_to(DOCS).as_posix()
        text = page.read_text(encoding="utf-8")
        prose = prose_only(text)
        if rel not in nav:
            errors.append(f"page absent from navigation: {rel}")
        if len(re.findall(r"^# ", prose, re.M)) != 1:
            errors.append(f"expected exactly one page title: {rel}")
        if len(re.findall(r"^```", text, re.M)) % 2:
            errors.append(f"unbalanced fenced code: {rel}")
        for target in re.findall(r"\[[^\]\n]*\]\(([^)\s]+)\)", prose):
            parsed = urlsplit(target)
            if parsed.scheme or parsed.netloc or not parsed.path:
                continue
            path = (page.parent / unquote(parsed.path)).resolve()
            if not path.is_file():
                errors.append(f"broken local Markdown link: {rel} -> {target}")
        for snippet in re.findall(r'^--8<-- "([^"]+)"', text, re.M):
            path = (ROOT / snippet).resolve()
            if not path.is_relative_to(ROOT) or not path.is_file():
                errors.append(f"missing or invalid snippet: {rel} -> {snippet}")
    if errors:
        print("\n".join(errors), file=sys.stderr)
        return 1
    chars = sum(len(re.findall(r"[\u4e00-\u9fff]", p.read_text(encoding="utf-8"))) for p in pages)
    print(f"OK: {len(pages)} pages, all navigable; {chars:,} Chinese characters; links and snippet inputs valid")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

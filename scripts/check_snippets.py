"""Compile standalone Go code fences. Partial fragments are intentionally excluded."""
from pathlib import Path
import re
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[1]


def main():
    checked = 0
    cache = ROOT / ".cache"
    cache.mkdir(exist_ok=True)
    with tempfile.TemporaryDirectory(prefix="snippets-", dir=cache) as directory:
        for page in sorted((ROOT / "docs").rglob("*.md")):
            text = page.read_text(encoding="utf-8")
            for number, match in enumerate(re.finditer(r"^```go[^\n]*\n(.*?)^```", text, re.M | re.S), 1):
                code = match.group(1)
                if not re.match(r"\s*package main\b", code):
                    continue
                checked += 1
                file = Path(directory) / f"snippet_{checked}.go"
                file.write_text(code, encoding="utf-8")
                result = subprocess.run(
                    ["go", "test", str(file)], cwd=ROOT, capture_output=True,
                    text=True, encoding="utf-8", errors="replace", timeout=120,
                )
                if result.returncode:
                    print(f"FAIL: {page.relative_to(ROOT)} fence {number}")
                    print(result.stdout + result.stderr)
                    return 1
    print(f"OK: compiled {checked} standalone Go fences; partial fragments and source includes are tested separately")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

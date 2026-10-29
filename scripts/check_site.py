"""Check generated HTML links, anchors and the Chinese search index offline."""
from html.parser import HTMLParser
import json
from pathlib import Path
from urllib.parse import unquote, urlsplit

import yaml

ROOT = Path(__file__).resolve().parents[1]
SITE = ROOT / "site"


class Page(HTMLParser):
    def __init__(self, text):
        super().__init__(convert_charrefs=True)
        self.ids = set()
        self.links = []
        self.feed(text)

    def handle_starttag(self, tag, attributes):
        attrs = dict(attributes)
        if attrs.get("id"):
            self.ids.add(attrs["id"])
        for name in ("href", "src"):
            if attrs.get(name):
                self.links.append(attrs[name])


def main():
    config = yaml.safe_load((ROOT / "mkdocs.yml").read_text(encoding="utf-8"))
    base = urlsplit(config.get("site_url", "/")).path
    pages = {p.resolve(): Page(p.read_text(encoding="utf-8")) for p in SITE.rglob("*.html")}
    errors = []
    if not pages:
        errors.append("no generated pages; run mkdocs build first")
    for path, page in pages.items():
        for link in page.links:
            url = urlsplit(link)
            if url.scheme or url.netloc:
                continue
            raw = unquote(url.path)
            if not raw:
                target = path
            elif raw.startswith("/"):
                # MkDocs emits deployment-root URLs in 404.html.
                if base != "/" and raw.startswith(base):
                    raw = raw[len(base):]
                target = (SITE / raw.lstrip("/")).resolve()
            else:
                target = (path.parent / raw).resolve()
            if target.is_dir():
                target = target / "index.html"
            if not target.is_file():
                errors.append(f"{path.relative_to(SITE)} -> missing {link}")
                continue
            if url.fragment and target in pages and unquote(url.fragment) not in pages[target].ids:
                errors.append(f"{path.relative_to(SITE)} -> missing anchor {link}")
    index = SITE / "search/search_index.json"
    if not index.is_file():
        errors.append("missing search index")
    else:
        search = json.loads(index.read_text(encoding="utf-8"))
        if not search.get("docs") or not any("切片" in d.get("text", "") for d in search["docs"]):
            errors.append("search index has no Chinese chapter content")
        if "zh" not in search.get("config", {}).get("lang", []):
            errors.append("search index does not declare Chinese support")
    if errors:
        print("\n".join(errors[:60]))
        print(f"FAIL: {len(errors)} generated site issues")
        return 1
    print(f"OK: {len(pages)} HTML pages; internal resources and anchors valid; Chinese search index populated")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

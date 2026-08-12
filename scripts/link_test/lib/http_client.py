"""Shared HTTP client for new_shop link tests (x-token)."""

from __future__ import annotations

import json
import os
from typing import Any, Optional

import requests


class ApiClient:
    def __init__(self, base: str) -> None:
        self.base = base.rstrip("/")
        self.token = ""
        self.session = requests.Session()

    def set_token(self, token: str) -> None:
        self.token = token or ""

    def request(self, method: str, path: str, **kwargs: Any) -> requests.Response:
        headers = dict(kwargs.pop("headers", {}) or {})
        if self.token:
            headers.setdefault("x-token", self.token)
        url = path if path.startswith("http") else self.base + path
        return self.session.request(method, url, headers=headers, timeout=15, **kwargs)


def load_token_from_manifest(path: Optional[str] = None) -> str:
    env_token = os.environ.get("TOKEN", "").strip()
    if env_token:
        return env_token
    manifest = path or os.environ.get("SEED_MANIFEST", "scripts/seed/seed_manifest.json")
    if not os.path.isfile(manifest):
        return ""
    with open(manifest, "r", encoding="utf-8") as f:
        data = json.load(f)
    user = data.get("primary_user") or {}
    return str(user.get("token") or "")

#!/usr/bin/env python3
"""fulfillment v0.3.0 confirmOrder REST link tests."""

from __future__ import annotations

import os
import sys

sys.path.insert(0, os.path.join(os.path.dirname(__file__), ".."))

from lib.http_client import ApiClient, load_token_from_manifest  # noqa: E402


def _run(name: str, fn) -> bool:
    try:
        fn()
        print(f"[PASS] {name}")
        return True
    except Exception as exc:  # noqa: BLE001
        print(f"[FAIL] {name}: {exc}")
        return False


def test_no_token(client: ApiClient) -> None:
    bare = ApiClient(client.base)
    r = bare.request("POST", "/order/confirmOrder", json={"ID": 1})
    assert r.status_code == 200, r.text
    data = r.json()
    assert data.get("code") == 401, data


def test_id_required(client: ApiClient) -> None:
    r = client.request("POST", "/order/confirmOrder", json={"ID": 0})
    assert r.status_code == 200, r.text
    data = r.json()
    assert data.get("code") == 7, data
    assert "订单ID" in str(data.get("msg", "")), data


def test_missing_order(client: ApiClient) -> None:
    r = client.request("POST", "/order/confirmOrder", json={"ID": 2147483646})
    assert r.status_code == 200, r.text
    data = r.json()
    assert data.get("code") == 7, data
    msg = str(data.get("msg", ""))
    assert any(s in msg for s in ("不存在", "无权", "不允许", "权限不足")), data


def test_happy_path(client: ApiClient, order_id: int) -> None:
    r = client.request("POST", "/order/confirmOrder", json={"ID": order_id})
    assert r.status_code == 200, r.text
    data = r.json()
    assert data.get("code") == 0, data
    assert data.get("msg") == "确认收货成功", data
    r2 = client.request("POST", "/order/confirmOrder", json={"ID": order_id})
    assert r2.status_code == 200, r2.text
    data2 = r2.json()
    assert data2.get("code") == 7, data2


def main() -> int:
    base = os.environ.get("API_BASE", "http://127.0.0.1:48888")
    client = ApiClient(base)
    token = load_token_from_manifest()
    failed = 0
    total = 0

    total += 1
    if not _run("no token => 401", lambda: test_no_token(client)):
        failed += 1

    if not token:
        print("[FAIL] missing TOKEN / seed_manifest; skip authenticated cases")
        print(f"Summary: {total - failed}/{total} passed")
        return 1

    client.set_token(token)
    cases = [
        ("ID=0 => 订单ID不能为空", lambda: test_id_required(client)),
        ("missing order => business error", lambda: test_missing_order(client)),
    ]
    order_id = os.environ.get("ORDER_ID", "").strip()
    if order_id:
        oid = int(order_id)
        cases.append(
            (f"ORDER_ID={oid} confirm then reject duplicate", lambda: test_happy_path(client, oid))
        )

    for name, fn in cases:
        total += 1
        if not _run(name, fn):
            failed += 1

    print(f"Summary: {total - failed}/{total} passed")
    return 0 if failed == 0 else 1


if __name__ == "__main__":
    raise SystemExit(main())

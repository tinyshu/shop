# 链路接口测试

依赖：Python 3.9+、`requests`。服务默认 `http://127.0.0.1:48888`，鉴权头 `x-token`。

```bash
pip install -r scripts/link_test/requirements.txt
set API_BASE=http://127.0.0.1:48888
set TOKEN=your_jwt
python scripts/link_test/fulfillment/test_confirm_order_v030.py
```

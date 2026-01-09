import os
import time
import random
import itertools
import json
from dataclasses import dataclass
from pathlib import Path
from typing import Any, Dict, Iterator, Optional, Tuple

import pytest
import requests


DEFAULT_BASE_URL = "http://localhost:8080/api"


def pytest_configure(config: pytest.Config) -> None:
    # 让 pytest 的报告输出目录稳定存在（例如 pytest.ini 里默认写 junitxml 到 test/out/）
    out_dir = Path(__file__).resolve().parents[1] / "out"
    out_dir.mkdir(parents=True, exist_ok=True)


def pytest_runtest_makereport(item: pytest.Item, call: pytest.CallInfo) -> None:
    # 失败时把最近几次 API 请求/响应附在报告里（命令行输出 & junit/html 报告都能看到）
    if call.when != "call" or call.excinfo is None:
        return
    funcargs = getattr(item, "funcargs", {}) or {}
    sess = funcargs.get("session")
    trace = getattr(sess, "_vvechat_api_trace", None)
    if not trace:
        return
    try:
        content = json.dumps(trace[-5:], ensure_ascii=False, indent=2)
    except Exception:
        content = str(trace[-5:])
    item.add_report_section(call.when, "vvechat-api-trace(last-5)", content)


@dataclass(frozen=True)
class ApiConfig:
    base_url: str
    timeout: float


class ApiResponseError(RuntimeError):
    def __init__(
        self,
        *,
        http_status: int,
        code: Optional[int],
        message: str,
        payload: Any,
        method: Optional[str] = None,
        url: Optional[str] = None,
    ):
        prefix = f"{method} {url} -> " if method and url else ""
        super().__init__(f"{prefix}HTTP {http_status} code={code} message={message}")
        self.http_status = http_status
        self.code = code
        self.message = message
        self.payload = payload
        self.method = method
        self.url = url


@pytest.fixture(scope="session")
def api_config() -> ApiConfig:
    base_url = os.getenv("VVECHAT_BASE_URL", DEFAULT_BASE_URL).rstrip("/")
    timeout = float(os.getenv("VVECHAT_TIMEOUT", "8"))
    return ApiConfig(base_url=base_url, timeout=timeout)


@pytest.fixture()
def session() -> Iterator[requests.Session]:
    s = requests.Session()
    setattr(s, "_vvechat_api_trace", [])
    try:
        yield s
    finally:
        s.close()


def _redact_headers(headers: Dict[str, str]) -> Dict[str, str]:
    redacted = dict(headers)
    if "Authorization" in redacted:
        redacted["Authorization"] = "Bearer <redacted>"
    return redacted


def _api_json(
    session: requests.Session,
    cfg: ApiConfig,
    method: str,
    path: str,
    *,
    token: Optional[str] = None,
    json_body: Optional[Dict[str, Any]] = None,
) -> Tuple[int, Dict[str, Any]]:
    url = f"{cfg.base_url}{path}"
    headers: Dict[str, str] = {"Content-Type": "application/json"}
    if token:
        headers["Authorization"] = f"Bearer {token}"

    try:
        res = session.request(method=method, url=url, headers=headers, json=json_body, timeout=cfg.timeout)
    except requests.RequestException as e:
        pytest.fail(f"无法连接后端：{e}. 请先启动 Go 服务，并确认 VVECHAT_BASE_URL={cfg.base_url}", pytrace=False)

    try:
        payload = res.json()
    except Exception as e:
        body = res.text
        trace = getattr(session, "_vvechat_api_trace", None)
        if isinstance(trace, list):
            trace.append(
                {
                    "request": {
                        "method": method,
                        "url": url,
                        "headers": _redact_headers(headers),
                        "json": json_body,
                        "timeout": cfg.timeout,
                    },
                    "response": {
                        "http_status": res.status_code,
                        "body": body[:2000],
                    },
                }
            )
        raise RuntimeError(f"响应不是 JSON: {e}; body={body[:300]}")

    if not isinstance(payload, dict) or "code" not in payload:
        trace = getattr(session, "_vvechat_api_trace", None)
        if isinstance(trace, list):
            trace.append(
                {
                    "request": {
                        "method": method,
                        "url": url,
                        "headers": _redact_headers(headers),
                        "json": json_body,
                        "timeout": cfg.timeout,
                    },
                    "response": {
                        "http_status": res.status_code,
                        "payload": payload,
                    },
                }
            )
        raise RuntimeError(f"响应结构异常: http={res.status_code} payload={payload}")

    trace = getattr(session, "_vvechat_api_trace", None)
    if isinstance(trace, list):
        trace.append(
            {
                "request": {
                    "method": method,
                    "url": url,
                    "headers": _redact_headers(headers),
                    "json": json_body,
                    "timeout": cfg.timeout,
                },
                "response": {
                    "http_status": res.status_code,
                    "payload": payload,
                },
            }
        )

    return res.status_code, payload


def api_ok(
    session: requests.Session,
    cfg: ApiConfig,
    method: str,
    path: str,
    *,
    token: Optional[str] = None,
    json_body: Optional[Dict[str, Any]] = None,
    expected_code: int,
) -> Any:
    http_status, payload = _api_json(session, cfg, method, path, token=token, json_body=json_body)
    if http_status != 200:
        raise ApiResponseError(
            http_status=http_status,
            code=payload.get("code"),
            message=str(payload.get("message", "")),
            payload=payload,
            method=method,
            url=f"{cfg.base_url}{path}",
        )

    code = payload.get("code")
    if code != expected_code:
        raise ApiResponseError(
            http_status=http_status,
            code=code,
            message=str(payload.get("message", "")),
            payload=payload,
            method=method,
            url=f"{cfg.base_url}{path}",
        )

    return payload.get("data")


def api_fail(
    session: requests.Session,
    cfg: ApiConfig,
    method: str,
    path: str,
    *,
    token: Optional[str] = None,
    json_body: Optional[Dict[str, Any]] = None,
    expected_http: int,
    expected_code: Optional[int] = None,
) -> Dict[str, Any]:
    http_status, payload = _api_json(session, cfg, method, path, token=token, json_body=json_body)
    assert http_status == expected_http, f"期望 HTTP {expected_http}, 实际 {http_status}, payload={payload}"
    if expected_code is not None:
        assert payload.get("code") == expected_code, f"期望 code={expected_code}, 实际 {payload.get('code')}, payload={payload}"
    return payload


def gen_unique_phone() -> str:
    # 11 位“手机号”字符串：尽量保证在同一进程内绝对不重复。
    # 规则：后 10 位 = (纳秒时间戳 + 进程扰动 + 自增计数器) mod 1e10
    # 说明：我们不追求“像真实手机号”，只需要满足后端格式校验且不撞库。
    if not hasattr(gen_unique_phone, "_counter"):
        gen_unique_phone._counter = itertools.count()  # type: ignore[attr-defined]
        # 进程扰动：避免并发/多进程场景下的同一时间戳碰撞
        gen_unique_phone._pid_jitter = (os.getpid() % 10000) * 100000  # type: ignore[attr-defined]

    counter = next(gen_unique_phone._counter)  # type: ignore[attr-defined]
    now_ns = time.time_ns()
    suffix = (now_ns + gen_unique_phone._pid_jitter + counter) % (10**10)  # type: ignore[attr-defined]
    return "1" + f"{suffix:010d}"


@dataclass
class TestUser:
    name: str
    phone_number: str
    password: str
    uid: Optional[str] = None
    access_token: Optional[str] = None
    refresh_token: Optional[str] = None


def ensure_user_registered_and_logged_in(
    session: requests.Session,
    cfg: ApiConfig,
    *,
    name_prefix: str = "pytest",
    password: str = "password123",
    max_attempts: int = 5,
) -> TestUser:
    last_err: Optional[Exception] = None
    for attempt in range(1, max_attempts + 1):
        phone = gen_unique_phone()
        user = TestUser(name=f"{name_prefix}_{phone[-6:]}", phone_number=phone, password=password)

        http_status, payload = _api_json(
            session,
            cfg,
            "POST",
            "/register",
            json_body={"name": user.name, "password": user.password, "phone_number": user.phone_number},
        )

        if http_status == 200 and payload.get("code") == 201:
            # ok
            pass
        elif http_status == 400 and "手机号已存在" in str(payload.get("message", "")):
            # 极小概率撞库，重试
            continue
        else:
            last_err = ApiResponseError(
                http_status=http_status,
                code=payload.get("code"),
                message=str(payload.get("message", "")),
                payload=payload,
            )
            break

        # login
        data = api_ok(
            session,
            cfg,
            "POST",
            "/login/phone_number",
            json_body={"phone_number": user.phone_number, "password": user.password},
            expected_code=200,
        )

        assert isinstance(data, dict)
        user_info = data.get("user_info")
        token_class = data.get("token_class")
        assert isinstance(user_info, dict) and isinstance(token_class, dict)
        user.uid = str(user_info.get("uid"))
        user.access_token = str(token_class.get("token"))
        user.refresh_token = str(token_class.get("refresh_token"))
        assert user.uid and user.access_token and user.refresh_token
        return user

    raise RuntimeError(f"创建测试用户失败（尝试 {max_attempts} 次）: {last_err}")

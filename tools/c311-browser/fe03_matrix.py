#!/usr/bin/env python3
"""Repeatable FE-03 request submission, draft, and staff-assist checks."""

from __future__ import annotations

import json
import os
import sys
import tempfile
from datetime import datetime, timezone
from pathlib import Path
from urllib.parse import urlparse

from playwright.sync_api import Browser, Page, sync_playwright


VIEWPORTS = ((1440, 900), (768, 900), (390, 844))
ARTIFACT_DIR_ENV = "C311_ARTIFACT_DIR"
COMPOSE_URL = os.environ.get("C311_COMPOSE_URL", "http://127.0.0.1:18082").rstrip("/")
ADMIN_URL = os.environ.get("C311_ADMIN_URL", "http://127.0.0.1:18083").rstrip("/")
MAIN = "[data-c311-main]"
SUBMIT_PATH = "/c311/submit"
STAFF_SUBMIT_PATH = "/c311/staff/submit"
SUMMARY_FIELD = "#c311-summary"
DESCRIPTION_FIELD = "#c311-description"
CONSENT_FIELD = "#c311-consent"
SAVE_DRAFT_ACTION = '[data-c311-action="save-draft"]'
ERROR_SUMMARY = "[data-c311-error-summary]"
SUBMISSION_RESULT = "[data-c311-submission-result]"
ATTACHMENT_RECOVERY = "[data-c311-attachment-recovery]"
SUBMIT_ACTION = '[data-c311-action="submit-request"]'
STAFF_SUBMIT_ACTION = '[data-c311-action="submit-staff-request"]'
KNOWN_DEV_RESOURCE_FAILURES = {"/code-snippets.js", "/custom.css"}
KNOWN_DEV_WEBSOCKET_URL = "wss://api.cortezaproject.your-domain.tld/websocket"
KNOWN_DEV_LOCALE_FAILURES = {
    "/system/locale/en-US/corteza-webapp-compose",
    "/system/locale/en-US/corteza-webapp-admin",
    "/system/locale/en-US+en/corteza-webapp-compose",
    "/system/locale/en-US+en/corteza-webapp-admin",
    "/system/locale/en/corteza-webapp-compose",
    "/system/locale/en/corteza-webapp-admin",
    "/system/locale/en+en-US/corteza-webapp-compose",
    "/system/locale/en+en-US/corteza-webapp-admin",
}
# These are emitted by the shared development shell in the local test services.
# HTTP responses remain checked separately against the exact resource allow-list.
KNOWN_DEV_CONSOLE_ERRORS = {
    "Failed to load resource: the server responded with a status of 500 (Internal Server Error)",
    "WebSocket connection to 'wss://api.cortezaproject.your-domain.tld/websocket' failed: Error in connection establishment: net::ERR_CONNECTION_CLOSED",
}
GENERIC_CONNECTION_ERROR = "Failed to load resource: net::ERR_CONNECTION_CLOSED"


def check(condition: bool, message: str) -> None:
    if not condition:
        raise AssertionError(message)


def artifact_directory() -> Path:
    configured = os.environ.get(ARTIFACT_DIR_ENV)
    directory = Path(configured).expanduser() if configured else Path(tempfile.mkdtemp(prefix="c311-fe03-gate-"))
    if directory.is_symlink():
        raise ValueError(f"{ARTIFACT_DIR_ENV} must not point to a symbolic link")
    directory.mkdir(parents=True, exist_ok=True, mode=0o700)
    directory = directory.resolve()
    if directory == Path(directory.anchor):
        raise ValueError(f"{ARTIFACT_DIR_ENV} must point to a child directory")
    directory.chmod(0o700)
    return directory


def bootstrap(page: Page, scenario: str = "success", role: str = "public_visitor", session: str = "current") -> None:
    page.add_init_script(
        "window.C311Mode = 'mock'; "
        f"window.C311MockRole = {json.dumps(role)}; "
        f"window.C311MockScenario = {json.dumps(scenario)}; "
        f"window.C311MockSession = {json.dumps(session)};"
    )


def open_page(page: Page, base_url: str, path: str, scenario: str = "success", role: str = "public_visitor", session: str = "current") -> None:
    bootstrap(page, scenario, role, session)
    page.goto(f"{base_url}{path}", wait_until="domcontentloaded")
    page.locator(MAIN).first.wait_for(state="visible")


def assert_page(page: Page, label: str, width: int) -> None:
    check(page.locator(MAIN).count() == 1, f"{label} must have one main landmark")
    check(page.locator(f"{MAIN} h1:visible").count() == 1, f"{label} must have one visible h1")
    overflow = page.evaluate("document.documentElement.scrollWidth - document.documentElement.clientWidth")
    check(overflow <= 1, f"{label} overflows at {width}px: {overflow}")


SAFE_REQUEST_METHODS = {"GET", "HEAD", "OPTIONS"}


def new_diagnostics() -> dict[str, list]:
    return {
        "console_errors": [],
        "unexpected_console_errors": [],
        "page_errors": [],
        "writes": [],
        "responses": [],
        "unexpected_responses": [],
        "failed_requests": [],
        "pending_console_errors": [],
        "websocket_urls": [],
    }


def allowed_origins(base_urls: str | tuple[str, ...]) -> set[tuple[str, str]]:
    urls = (base_urls,) if isinstance(base_urls, str) else base_urls
    return {(urlparse(url).scheme, urlparse(url).netloc) for url in urls}


def known_hmr(parsed, origins: set[tuple[str, str]]) -> bool:
    return parsed.scheme == "ws" and parsed.path == "/ws" and ("http", parsed.netloc) in origins


def known_local_asset(parsed_url, origins: set[tuple[str, str]]) -> bool:
    return (
        (parsed_url.scheme, parsed_url.netloc) in origins
        and not parsed_url.query
        and not parsed_url.fragment
        and parsed_url.path in KNOWN_DEV_RESOURCE_FAILURES
    )


def known_locale(parsed_url) -> bool:
    return (
        (parsed_url.scheme, parsed_url.netloc) == ("https", "api.cortezaproject.your-domain.tld")
        and not parsed_url.query
        and not parsed_url.fragment
        and parsed_url.path in KNOWN_DEV_LOCALE_FAILURES
    )


def record_console(result: dict[str, list], message) -> None:
    if message.type != "error":
        return
    result["console_errors"].append(message.text)
    if message.text == GENERIC_CONNECTION_ERROR:
        result["pending_console_errors"].append(message.text)
    elif message.text not in KNOWN_DEV_CONSOLE_ERRORS:
        result["unexpected_console_errors"].append(message.text)


def record_response(result: dict[str, list], origins: set[tuple[str, str]], response) -> None:
    if response.status < 400:
        return
    entry = {"status": response.status, "method": response.request.method, "url": response.url}
    result["responses"].append(entry)
    parsed_url = urlparse(response.url)
    allowed = known_local_asset(parsed_url, origins) and response.status == 500
    if not allowed:
        result["unexpected_responses"].append(entry)


def record_websocket(result: dict[str, list], origins: set[tuple[str, str]], websocket) -> None:
    result["websocket_urls"].append(websocket.url)
    parsed_url = urlparse(websocket.url)
    if websocket.url != KNOWN_DEV_WEBSOCKET_URL and not known_hmr(parsed_url, origins):
        result["unexpected_console_errors"].append(f"unexpected websocket: {websocket.url}")


def record_request_failed(result: dict[str, list], origins: set[tuple[str, str]], request) -> None:
    entry = {"method": request.method, "url": request.url, "failure": request.failure}
    result["failed_requests"].append(entry)
    parsed_url = urlparse(request.url)
    is_known = request.url == KNOWN_DEV_WEBSOCKET_URL or known_hmr(parsed_url, origins) or known_local_asset(parsed_url, origins) or known_locale(parsed_url)
    if not is_known:
        result["unexpected_console_errors"].append(f"request failed: {request.url}")


def record_write(result: dict[str, list], request) -> None:
    if request.method not in SAFE_REQUEST_METHODS:
        result["writes"].append({"method": request.method, "url": request.url})


def diagnostics(page: Page, base_urls: str | tuple[str, ...]) -> dict[str, list]:
    result = new_diagnostics()
    origins = allowed_origins(base_urls)
    page.on("console", lambda message: record_console(result, message))
    page.on("pageerror", lambda error: result["page_errors"].append(str(error)))
    page.on("websocket", lambda websocket: record_websocket(result, origins, websocket))
    page.on("requestfailed", lambda request: record_request_failed(result, origins, request))
    page.on("request", lambda request: record_write(result, request))
    page.on("response", lambda response: record_response(result, origins, response))
    return result


def finalize_diagnostics(result: dict[str, list]) -> None:
    known_websocket_failures = sum(1 for url in result["websocket_urls"] if url == KNOWN_DEV_WEBSOCKET_URL or (urlparse(url).scheme == "ws" and urlparse(url).path == "/ws"))
    for message in result.pop("pending_console_errors", []):
        if known_websocket_failures:
            known_websocket_failures -= 1
        else:
            result["unexpected_console_errors"].append(message)


def fill_valid_form(page: Page) -> None:
    page.locator("#c311-service-type").select_option("GENERAL_INQUIRY")
    page.locator(SUMMARY_FIELD).fill("Need help with a city service")
    page.locator(DESCRIPTION_FIELD).fill("The resident needs help with this city service request.")
    page.locator("#c311-requester-name").fill("Fixture Resident")
    page.locator("#c311-requester-email").fill("resident@example.test")
    page.locator("#c311-requester-phone").fill("")
    page.locator(CONSENT_FIELD).check()


def check_help_and_navigation(page: Page) -> None:
    check(page.locator('[data-c311-route="/c311/account"]').count() == 0, "anonymous account link exposed")
    page.locator('[data-c311-route="/c311/services"]').click()
    page.wait_for_url("**/c311/services")
    check(page.locator('[data-c311-page="services"]').get_attribute("data-c311-content-key") == "SERVICE_CATALOGUE", "services kept stale content")
    page.locator('[data-c311-route="/c311/help"]').click()
    page.wait_for_url("**/c311/help")
    check(page.locator('[data-c311-page="help"]').get_attribute("data-c311-content-key") == "HELP", "help kept stale content")
    check("Describe the issue" in page.locator('[data-c311-page="help"]').inner_text(), "contextual help missing")
    page.locator('[data-c311-route="/c311/submit"]').click()
    page.wait_for_url(f"**{SUBMIT_PATH}")
    page.locator(SUMMARY_FIELD).fill("Draft kept")
    page.locator(DESCRIPTION_FIELD).fill("A non-sensitive draft description")
    page.locator('[data-c311-route="/c311/status"]').click()
    check(page.url.endswith(SUBMIT_PATH), "dirty cancel did not preserve form")
    page.once("dialog", lambda dialog: dialog.accept())
    page.locator('[data-c311-route="/c311/status"]').click()
    page.wait_for_url("**/c311/status")


def check_modal(page: Page, base_url: str) -> None:
    page.goto(f"{base_url}/c311/test/modal", wait_until="domcontentloaded")
    page.locator(MAIN).first.wait_for(state="visible")
    opener = page.locator('[data-c311-modal-opener]')
    opener.click()
    modal = page.locator('[data-c311-focus-modal]')
    modal.wait_for(state="visible")
    check(page.evaluate("document.activeElement.closest('[data-c311-focus-modal]') !== null"), "modal focus did not enter")
    page.keyboard.press("Tab")
    check(page.evaluate("document.activeElement.closest('[data-c311-focus-modal]') !== null"), "modal focus escaped")
    page.keyboard.press("Escape")
    modal.wait_for(state="detached")
    check(page.evaluate("document.activeElement && document.activeElement.matches('[data-c311-modal-opener]')"), "modal focus did not return")


def check_submission_and_draft(page: Page, base_url: str) -> None:
    open_page(page, base_url, SUBMIT_PATH)
    page.locator(SUBMIT_ACTION).click()
    page.locator(ERROR_SUMMARY).wait_for(state="visible")
    check(page.evaluate(f"document.activeElement && document.activeElement.matches('{ERROR_SUMMARY}')"), "validation summary did not receive focus")

    fill_valid_form(page)
    page.set_input_files("#c311-attachment-file", {"name": "fixture.txt", "mimeType": "text/plain", "buffer": b"fixture"})
    page.locator('[data-c311-attachment-list]').wait_for(state="visible")
    page.locator(SUBMIT_ACTION).click()
    page.locator(SUBMISSION_RESULT).wait_for(state="visible")
    submitted = page.locator(SUBMISSION_RESULT).inner_text()
    check("SR-2026-00002" in submitted and "SUBMITTED" in submitted, "submission result missing request number/status")
    page.locator(SUBMIT_ACTION).click()
    page.locator(SUBMISSION_RESULT).wait_for(state="visible")
    check(page.locator(SUBMISSION_RESULT).inner_text() == submitted, "equivalent replay changed the result")

    open_page(page, base_url, SUBMIT_PATH, role="constituent")
    page.locator(SUMMARY_FIELD).fill("Saved draft summary")
    page.locator(DESCRIPTION_FIELD).fill("A saved draft description that is long enough.")
    page.locator("#c311-requester-name").fill("Saved Resident")
    page.locator("#c311-requester-email").fill("saved@example.test")
    page.locator(CONSENT_FIELD).check()
    page.locator(SAVE_DRAFT_ACTION).click()
    page.locator('[data-c311-submit-success]').wait_for(state="visible")
    check("draft_id=" in page.url, "remote draft id was not added to route")
    page.reload(wait_until="domcontentloaded")
    check(page.locator(SUMMARY_FIELD).input_value() == "Saved draft summary", "remote draft was not restored")


def check_capability_controls(page: Page, base_url: str) -> None:
    open_page(page, base_url, SUBMIT_PATH, role="public_visitor")
    check(page.locator(SAVE_DRAFT_ACTION).count() == 0, "draft save action exposed to anonymous users")
    open_page(page, base_url, SUBMIT_PATH, role="service_agent")
    check(page.locator(SAVE_DRAFT_ACTION).count() == 0, "draft save action exposed without capability")
    open_page(page, base_url, SUBMIT_PATH, role="constituent")
    check(page.locator(SAVE_DRAFT_ACTION).count() == 1, "draft save action missing for constituent")


def check_restart_recovery(browser: Browser, base_url: str, viewport: tuple[int, int]) -> None:
    width, height = viewport
    context = browser.new_context(viewport={"width": width, "height": height})
    try:
        page = context.new_page()
        open_page(page, base_url, SUBMIT_PATH)
        page.locator(SUMMARY_FIELD).fill("Restart draft summary")
        page.locator(DESCRIPTION_FIELD).fill("A non-sensitive draft restored after a browser restart.")
        page.set_input_files("#c311-attachment-file", {"name": "fixture.txt", "mimeType": "text/plain", "buffer": b"fixture"})
        page.locator('[data-c311-attachment-list]').wait_for(state="visible")
        check(page.evaluate("window.localStorage.getItem('c311.portal.submit') !== null"), "portal draft was not persisted")
        stored = page.evaluate("JSON.parse(window.localStorage.getItem('c311.portal.submit'))")
        check("attachment_tokens" not in stored, "local draft persisted attachment tokens")
        check(stored.get("attachment_count") == 1, "local draft did not preserve attachment recovery metadata")
        storage_state = context.storage_state()
    finally:
        context.close()

    restored_context = browser.new_context(viewport={"width": width, "height": height}, storage_state=storage_state)
    try:
        restored_page = restored_context.new_page()
        open_page(restored_page, base_url, SUBMIT_PATH)
        check(restored_page.locator(SUMMARY_FIELD).input_value() == "Restart draft summary", "portal draft was not restored after restart")
        check(restored_page.locator('[data-c311-attachment-list]').count() == 0, "attachment token was restored from local draft")
        check(restored_page.locator(ATTACHMENT_RECOVERY).is_visible(), "attachment re-upload warning missing after restart")
    finally:
        restored_context.close()


def check_provider_errors(page: Page, base_url: str) -> None:
    for scenario in ("validation", "forbidden", "retryable", "terminal", "idempotency-conflict"):
        open_page(page, base_url, SUBMIT_PATH, scenario=scenario)
        fill_valid_form(page)
        page.locator(SUBMIT_ACTION).click()
        page.locator(ERROR_SUMMARY).wait_for(state="visible")
        check(page.locator(ERROR_SUMMARY).count() == 1, f"{scenario} error summary missing")

    open_page(page, base_url, f"{SUBMIT_PATH}?draft_id=draft-fixture-001", scenario="version-conflict", role="constituent")
    page.locator(SUMMARY_FIELD).fill("Local conflict change")
    page.locator(CONSENT_FIELD).check()
    page.locator(SUBMIT_ACTION).click()
    page.locator('[data-c311-version-conflict]').wait_for(state="visible")
    check(page.locator('[data-c311-action="reload-draft"]').count() == 1, "version reload action missing")
    check(page.locator('[data-c311-action="reapply-draft"]').count() == 1, "version reapply action missing")
    page.locator('[data-c311-action="reapply-draft"]').click()
    check(page.locator(SUMMARY_FIELD).input_value() == "Local conflict change", "reapply did not preserve local draft")
    page.locator(SUBMIT_ACTION).click()
    page.locator('[data-c311-version-conflict]').wait_for(state="visible")
    page.locator('[data-c311-action="reload-draft"]').click()
    check(page.locator(SUMMARY_FIELD).input_value() == "Saved draft", "reload did not restore server draft")


def check_access_and_admin(page: Page, base_url: str, admin_url: str) -> None:
    open_page(page, base_url, "/c311/account")
    check(page.url.endswith("/c311/401") or page.get_by_role("heading", name="Sign-in required").count() == 1, "private account did not enforce 401")
    open_page(page, base_url, "/c311/requests", role="constituent", session="expired")
    check(page.url.endswith("/c311/401") or page.get_by_role("heading", name="Sign-in required").count() == 1, "expired session did not enforce 401")
    for path, heading in (("/c311/403", "Access denied"), ("/c311/404", "Not found")):
        page.goto(f"{base_url}{path}", wait_until="domcontentloaded")
        check(page.get_by_role("heading", name=heading).count() == 1, f"{path} status page missing")

    open_page(page, admin_url, STAFF_SUBMIT_PATH, scenario="forbidden", role="service_agent")
    check(page.locator("#c311-summary").is_visible(), "staff form did not load when constituent profile was forbidden")
    check(page.locator(ERROR_SUMMARY).count() == 0 or not page.locator(ERROR_SUMMARY).is_visible(), "staff profile failure blocked the staff form")
    open_page(page, admin_url, STAFF_SUBMIT_PATH, role="service_agent")
    fill_valid_form(page)
    page.locator(STAFF_SUBMIT_ACTION).click()
    page.locator(SUBMISSION_RESULT).wait_for(state="visible")
    check("SUBMITTED" in page.locator(SUBMISSION_RESULT).inner_text(), "staff assist submission did not complete")
    open_page(page, admin_url, STAFF_SUBMIT_PATH, role="constituent")
    check(page.url.endswith("/c311/403") or page.get_by_role("heading", name="Access denied").count() == 1, "staff capability denial missing")


def run() -> dict:
    artifacts = artifact_directory()
    results = {"viewports": [list(viewport) for viewport in VIEWPORTS], "checks": [], "started_at": datetime.now(timezone.utc).isoformat()}
    with sync_playwright() as playwright:
        browser = playwright.chromium.launch(headless=True)
        try:
            for width, height in VIEWPORTS:
                for app, base_url in (("compose", COMPOSE_URL), ("admin", ADMIN_URL)):
                    context = browser.new_context(viewport={"width": width, "height": height})
                    page = context.new_page()
                    diagnostic_origins = (base_url, COMPOSE_URL) if app == "admin" else base_url
                    diagnostics_result = diagnostics(page, diagnostic_origins)
                    label = f"{app}@{width}x{height}"
                    if app == "compose":
                        open_page(page, base_url, "/c311")
                        assert_page(page, label, width)
                        check_help_and_navigation(page)
                        check_modal(page, base_url)
                        check_submission_and_draft(page, base_url)
                        check_capability_controls(page, base_url)
                        check_provider_errors(page, base_url)
                        if width == 390:
                            check_restart_recovery(browser, base_url, (width, height))
                    else:
                        open_page(page, base_url, STAFF_SUBMIT_PATH, role="service_agent")
                        assert_page(page, label, width)
                        check_access_and_admin(page, COMPOSE_URL, base_url)
                    screenshot = artifacts / f"fe03-{app}-{width}x{height}.png"
                    page.screenshot(path=str(screenshot), full_page=True)
                    finalize_diagnostics(diagnostics_result)
                    check(not diagnostics_result["page_errors"], f"{label} has uncaught page errors: {diagnostics_result['page_errors']}")
                    check(not diagnostics_result["unexpected_console_errors"], f"{label} has unexpected console errors: {diagnostics_result['unexpected_console_errors']}")
                    check(not diagnostics_result["unexpected_responses"], f"{label} has unexpected HTTP failures: {diagnostics_result['unexpected_responses']}")
                    check(not diagnostics_result["writes"], f"{label} issued real network writes in mock mode: {diagnostics_result['writes']}")
                    results["checks"].append({"label": label, "status": "passed", "screenshot": str(screenshot), "diagnostics": diagnostics_result})
                    context.close()
        finally:
            browser.close()
    results["finished_at"] = datetime.now(timezone.utc).isoformat()
    (artifacts / "fe03-matrix.json").write_text(json.dumps(results, indent=2), encoding="utf-8")
    return results


if __name__ == "__main__":
    try:
        print(json.dumps(run(), indent=2))
    except Exception as error:
        print(json.dumps({"status": "failed", "error": str(error)}), file=sys.stderr)
        raise

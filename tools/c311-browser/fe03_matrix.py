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


def diagnostics(page: Page, base_urls: str | tuple[str, ...]) -> dict[str, list]:
    result = {"console_errors": [], "unexpected_console_errors": [], "page_errors": [], "writes": [], "responses": [], "unexpected_responses": [], "failed_requests": [], "pending_console_errors": []}
    if isinstance(base_urls, str):
        base_urls = (base_urls,)
    allowed_origins = {(urlparse(url).scheme, urlparse(url).netloc) for url in base_urls}

    def record_console(message) -> None:
        if message.type == "error":
            result["console_errors"].append(message.text)
            if message.text == GENERIC_CONNECTION_ERROR:
                result["pending_console_errors"].append(message.text)
            elif message.text not in KNOWN_DEV_CONSOLE_ERRORS:
                result["unexpected_console_errors"].append(message.text)

    def record_response(response) -> None:
        if response.status < 400:
            return
        entry = {"status": response.status, "method": response.request.method, "url": response.url}
        result["responses"].append(entry)
        parsed = urlparse(response.url)
        allowed = (
            (parsed.scheme, parsed.netloc) in allowed_origins
            and not parsed.query
            and not parsed.fragment
            and parsed.path in KNOWN_DEV_RESOURCE_FAILURES
            and response.status == 500
        )
        if not allowed:
            result["unexpected_responses"].append(entry)

    page.on("console", record_console)
    page.on("pageerror", lambda error: result["page_errors"].append(str(error)))
    def record_websocket(websocket) -> None:
        result["websocket_urls"].append(websocket.url)
        parsed = urlparse(websocket.url)
        known_hmr = parsed.scheme == "ws" and parsed.path == "/ws" and ("http", parsed.netloc) in allowed_origins
        if websocket.url != KNOWN_DEV_WEBSOCKET_URL and not known_hmr:
            result["unexpected_console_errors"].append(f"unexpected websocket: {websocket.url}")
    result["websocket_urls"] = []
    page.on("websocket", record_websocket)
    def record_request_failed(request) -> None:
        entry = {"method": request.method, "url": request.url, "failure": request.failure}
        result["failed_requests"].append(entry)
        parsed = urlparse(request.url)
        known_local_asset = (
            (parsed.scheme, parsed.netloc) in allowed_origins
            and not parsed.query
            and not parsed.fragment
            and parsed.path in KNOWN_DEV_RESOURCE_FAILURES
        )
        known_locale = (
            (parsed.scheme, parsed.netloc) == ("https", "api.cortezaproject.your-domain.tld")
            and not parsed.query
            and not parsed.fragment
            and parsed.path in KNOWN_DEV_LOCALE_FAILURES
        )
        known_hmr = parsed.scheme == "ws" and parsed.path == "/ws" and ("http", parsed.netloc) in allowed_origins
        if request.url != KNOWN_DEV_WEBSOCKET_URL and not known_hmr and not known_local_asset and not known_locale:
            result["unexpected_console_errors"].append(f"request failed: {request.url}")
    page.on("requestfailed", record_request_failed)
    page.on("request", lambda request: result["writes"].append({"method": request.method, "url": request.url}) if request.method not in {"GET", "HEAD", "OPTIONS"} else None)
    page.on("response", record_response)
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
    page.locator("#c311-summary").fill("Need help with a city service")
    page.locator("#c311-description").fill("The resident needs help with this city service request.")
    page.locator("#c311-requester-name").fill("Fixture Resident")
    page.locator("#c311-requester-email").fill("resident@example.test")
    page.locator("#c311-requester-phone").fill("")
    page.locator("#c311-consent").check()


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
    page.wait_for_url("**/c311/submit")
    page.locator("#c311-summary").fill("Draft kept")
    page.locator("#c311-description").fill("A non-sensitive draft description")
    page.locator('[data-c311-route="/c311/status"]').click()
    check(page.url.endswith("/c311/submit"), "dirty cancel did not preserve form")
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
    open_page(page, base_url, "/c311/submit")
    page.locator(SUBMIT_ACTION).click()
    page.locator('[data-c311-error-summary]').wait_for(state="visible")
    check(page.evaluate("document.activeElement && document.activeElement.matches('[data-c311-error-summary]')"), "validation summary did not receive focus")

    fill_valid_form(page)
    page.set_input_files("#c311-attachment-file", {"name": "fixture.txt", "mimeType": "text/plain", "buffer": b"fixture"})
    page.locator('[data-c311-attachment-list]').wait_for(state="visible")
    page.locator(SUBMIT_ACTION).click()
    page.locator("[data-c311-submission-result]").wait_for(state="visible")
    submitted = page.locator("[data-c311-submission-result]").inner_text()
    check("SR-2026-00002" in submitted and "SUBMITTED" in submitted, "submission result missing request number/status")
    page.locator(SUBMIT_ACTION).click()
    page.locator("[data-c311-submission-result]").wait_for(state="visible")
    check(page.locator("[data-c311-submission-result]").inner_text() == submitted, "equivalent replay changed the result")

    open_page(page, base_url, "/c311/submit", role="constituent")
    page.locator("#c311-summary").fill("Saved draft summary")
    page.locator("#c311-description").fill("A saved draft description that is long enough.")
    page.locator("#c311-requester-name").fill("Saved Resident")
    page.locator("#c311-requester-email").fill("saved@example.test")
    page.locator("#c311-consent").check()
    page.locator('[data-c311-action="save-draft"]').click()
    page.locator('[data-c311-submit-success]').wait_for(state="visible")
    check("draft_id=" in page.url, "remote draft id was not added to route")
    page.reload(wait_until="domcontentloaded")
    check(page.locator("#c311-summary").input_value() == "Saved draft summary", "remote draft was not restored")


def check_capability_controls(page: Page, base_url: str) -> None:
    open_page(page, base_url, "/c311/submit", role="public_visitor")
    check(page.locator('[data-c311-action="save-draft"]').count() == 0, "draft save action exposed to anonymous users")
    open_page(page, base_url, "/c311/submit", role="service_agent")
    check(page.locator('[data-c311-action="save-draft"]').count() == 0, "draft save action exposed without capability")
    open_page(page, base_url, "/c311/submit", role="constituent")
    check(page.locator('[data-c311-action="save-draft"]').count() == 1, "draft save action missing for constituent")


def check_restart_recovery(browser: Browser, base_url: str, viewport: tuple[int, int]) -> None:
    width, height = viewport
    context = browser.new_context(viewport={"width": width, "height": height})
    try:
        page = context.new_page()
        open_page(page, base_url, "/c311/submit")
        page.locator("#c311-summary").fill("Restart draft summary")
        page.locator("#c311-description").fill("A non-sensitive draft restored after a browser restart.")
        check(page.evaluate("window.localStorage.getItem('c311.portal.submit') !== null"), "portal draft was not persisted")
        storage_state = context.storage_state()
    finally:
        context.close()

    restored_context = browser.new_context(viewport={"width": width, "height": height}, storage_state=storage_state)
    try:
        restored_page = restored_context.new_page()
        open_page(restored_page, base_url, "/c311/submit")
        check(restored_page.locator("#c311-summary").input_value() == "Restart draft summary", "portal draft was not restored after restart")
    finally:
        restored_context.close()


def check_provider_errors(page: Page, base_url: str) -> None:
    for scenario in ("validation", "forbidden", "retryable", "terminal", "idempotency-conflict"):
        open_page(page, base_url, "/c311/submit", scenario=scenario)
        fill_valid_form(page)
        page.locator(SUBMIT_ACTION).click()
        page.locator('[data-c311-error-summary]').wait_for(state="visible")
        check(page.locator('[data-c311-error-summary]').count() == 1, f"{scenario} error summary missing")

    open_page(page, base_url, "/c311/submit?draft_id=draft-fixture-001", scenario="version-conflict", role="constituent")
    page.locator("#c311-summary").fill("Local conflict change")
    page.locator("#c311-consent").check()
    page.locator(SUBMIT_ACTION).click()
    page.locator('[data-c311-version-conflict]').wait_for(state="visible")
    check(page.locator('[data-c311-action="reload-draft"]').count() == 1, "version reload action missing")
    check(page.locator('[data-c311-action="reapply-draft"]').count() == 1, "version reapply action missing")
    page.locator('[data-c311-action="reapply-draft"]').click()
    check(page.locator("#c311-summary").input_value() == "Local conflict change", "reapply did not preserve local draft")
    page.locator(SUBMIT_ACTION).click()
    page.locator('[data-c311-version-conflict]').wait_for(state="visible")
    page.locator('[data-c311-action="reload-draft"]').click()
    check(page.locator("#c311-summary").input_value() == "Saved draft", "reload did not restore server draft")


def check_access_and_admin(page: Page, base_url: str, admin_url: str) -> None:
    open_page(page, base_url, "/c311/account")
    check(page.url.endswith("/c311/401") or page.get_by_role("heading", name="Sign-in required").count() == 1, "private account did not enforce 401")
    open_page(page, base_url, "/c311/requests", role="constituent", session="expired")
    check(page.url.endswith("/c311/401") or page.get_by_role("heading", name="Sign-in required").count() == 1, "expired session did not enforce 401")
    for path, heading in (("/c311/403", "Access denied"), ("/c311/404", "Not found")):
        page.goto(f"{base_url}{path}", wait_until="domcontentloaded")
        check(page.get_by_role("heading", name=heading).count() == 1, f"{path} status page missing")

    open_page(page, admin_url, "/c311/staff/submit", role="service_agent")
    fill_valid_form(page)
    page.locator(STAFF_SUBMIT_ACTION).click()
    page.locator("[data-c311-submission-result]").wait_for(state="visible")
    check("SUBMITTED" in page.locator("[data-c311-submission-result]").inner_text(), "staff assist submission did not complete")
    open_page(page, admin_url, "/c311/staff/submit", role="constituent")
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
                        open_page(page, base_url, "/c311/staff/submit", role="service_agent")
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

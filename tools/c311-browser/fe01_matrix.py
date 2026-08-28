#!/usr/bin/env python3
"""Repeatable FE-01 browser checks for the Compose and Admin C311 surfaces."""

from __future__ import annotations

import json
import os
import sys
import tempfile
from datetime import datetime, timezone
from pathlib import Path

from playwright.sync_api import Page, sync_playwright


VIEWPORTS = (360, 767, 768, 1023, 1024, 1440, 1920)
ARTIFACT_DIR_ENV = "C311_ARTIFACT_DIR"
MAIN_LANDMARK_SELECTOR = "[data-c311-main]"
MODAL_FOCUS_EXPRESSION = "document.activeElement.closest('[data-c311-focus-modal]') !== null"
SUBMIT_PATH = "/c311/submit"
SUMMARY_SELECTOR = "#c311-summary"
STATUS_PATH = "/c311/status"
STATUS_ROUTE_SELECTOR = '[data-c311-route="/c311/status"]'
STAFF_PATH = "/c311/staff"
WORKFLOWS_PATH = "/c311/staff/workflows"
COMPOSE_URL = os.environ.get("C311_COMPOSE_URL", "http://127.0.0.1:18081").rstrip("/")
ADMIN_URL = os.environ.get("C311_ADMIN_URL", "http://127.0.0.1:18082").rstrip("/")
ACTIVE_ARTIFACT_DIR: Path | None = None


def check(condition: bool, message: str) -> None:
    if not condition:
        raise AssertionError(message)


def artifact_directory() -> Path:
    configured = os.environ.get(ARTIFACT_DIR_ENV)
    if not configured:
        return Path(tempfile.mkdtemp(prefix="c311-fe01-gate-"))

    directory = Path(configured).expanduser()
    if directory.is_symlink():
        raise ValueError(f"{ARTIFACT_DIR_ENV} must not point to a symbolic link")
    directory.mkdir(parents=True, exist_ok=True, mode=0o700)
    directory = directory.resolve()
    if directory == Path(directory.anchor):
        raise ValueError(f"{ARTIFACT_DIR_ENV} must point to a child directory")
    directory.chmod(0o700)
    return directory


def bootstrap(page: Page, role: str, scenario: str = "success", session: str = "current") -> None:
    page.add_init_script(
        "window.C311Mode = 'mock'; "
        f"window.C311MockRole = {json.dumps(role)}; "
        f"window.C311MockScenario = {json.dumps(scenario)}; "
        f"window.C311MockSession = {json.dumps(session)};"
    )


def open_c311(page: Page, base_url: str, path: str, role: str, scenario: str = "success", session: str = "current") -> None:
    bootstrap(page, role, scenario, session)
    page.goto(f"{base_url}{path}", wait_until="domcontentloaded")
    page.locator(MAIN_LANDMARK_SELECTOR).first.wait_for(state="visible")
    page.locator(f"{MAIN_LANDMARK_SELECTOR} h1, {MAIN_LANDMARK_SELECTOR}").first.wait_for(state="visible")


def assert_layout(page: Page, viewport: int, label: str) -> None:
    check(page.locator(MAIN_LANDMARK_SELECTOR).count() >= 1, f"{label} has no main landmark")
    check(page.locator(f"{MAIN_LANDMARK_SELECTOR} h1:visible").count() == 1, f"{label} must have one visible h1")
    overflow = page.evaluate("document.documentElement.scrollWidth - document.documentElement.clientWidth")
    check(overflow <= 1, f"{label} overflows horizontally at {viewport}px: {overflow}px")


def assert_route_focus(page: Page, expected_fragment: str) -> None:
    page.wait_for_timeout(150)
    focused = page.evaluate("document.activeElement && document.activeElement.outerHTML") or ""
    check("<h1" in focused or "data-c311-main" in focused, f"focus did not move to the C311 landmark: {focused[:120]}")
    check(expected_fragment in page.url, f"expected route {expected_fragment}, got {page.url}")


def check_help_and_modal(page: Page, base_url: str, label: str) -> None:
    help_trigger = page.get_by_role("button", name="Help").first
    help_trigger.wait_for(state="visible")
    help_trigger.click()
    drawer = page.locator("[data-c311-help-drawer]")
    drawer.wait_for(state="visible")
    check(page.evaluate("document.activeElement.closest('[data-c311-help-drawer]') !== null"), f"{label} help focus did not enter drawer")
    page.keyboard.press("Tab")
    check(page.evaluate("document.activeElement.closest('[data-c311-help-drawer]') !== null"), f"{label} help focus escaped on Tab")
    page.keyboard.press("Escape")
    drawer.wait_for(state="detached")
    check(page.evaluate("document.activeElement && document.activeElement.matches('[aria-controls^=\\\"c311-help-\\\"]')"), f"{label} help focus did not return")

    page.goto(f"{base_url}/c311/test/modal", wait_until="domcontentloaded")
    page.locator(MAIN_LANDMARK_SELECTOR).wait_for(state="visible")
    opener = page.get_by_role("button", name="Open dialog")
    opener.click()
    modal = page.locator("[data-c311-focus-modal]")
    modal.wait_for(state="visible")
    check(page.evaluate(MODAL_FOCUS_EXPRESSION), f"{label} modal focus did not enter")
    page.keyboard.press("Tab")
    check(page.evaluate(MODAL_FOCUS_EXPRESSION), f"{label} modal focus escaped on Tab")
    page.keyboard.press("Shift+Tab")
    check(page.evaluate(MODAL_FOCUS_EXPRESSION), f"{label} modal reverse focus escaped")
    page.keyboard.press("Escape")
    modal.wait_for(state="detached")
    check(page.evaluate("document.activeElement.matches('[data-c311-modal-opener]')"), f"{label} modal focus did not return")


def check_language_and_dirty(page: Page, base_url: str, label: str) -> None:
    open_c311(page, base_url, SUBMIT_PATH, "public_visitor")
    locale = page.locator("select[data-c311-language]")
    locale.select_option("es")
    page.wait_for_timeout(150)
    stored = page.evaluate("localStorage.getItem('c311.locale.anonymous')")
    check(stored == "es", f"{label} did not persist Spanish locale: {stored}")
    page.reload(wait_until="domcontentloaded")
    locale.wait_for(state="visible")
    check(locale.input_value() == "es", f"{label} locale was not restored after refresh")
    check("Enviar una solicitud" in page.locator(MAIN_LANDMARK_SELECTOR).inner_text(), f"{label} did not apply Spanish translation")

    summary = page.locator(SUMMARY_SELECTOR)
    description = page.locator("#c311-description")
    summary.fill("Draft kept across refresh")
    description.fill("Non-sensitive draft")
    page.reload(wait_until="domcontentloaded")
    check(page.locator(SUMMARY_SELECTOR).input_value() == "Draft kept across refresh", f"{label} draft summary was lost")
    check(page.locator("#c311-description").input_value() == "Non-sensitive draft", f"{label} draft description was lost")
    page.locator(SUMMARY_SELECTOR).fill("Draft modified before leaving")

    page.once("dialog", lambda dialog: dialog.dismiss())
    page.locator(STATUS_ROUTE_SELECTOR).click()
    page.wait_for_timeout(100)
    check(SUBMIT_PATH in page.url, f"{label} dirty cancel allowed navigation")
    page.once("dialog", lambda dialog: dialog.accept())
    page.locator(STATUS_ROUTE_SELECTOR).click()
    page.wait_for_url(f"**{STATUS_PATH}")


def check_portal_errors(page: Page, base_url: str, label: str) -> None:
    def wait_state(expected: str) -> None:
        page.wait_for_function(
            "expected => document.querySelector('[data-c311-data-state]')?.getAttribute('data-state') === expected",
            arg=expected,
            timeout=10000,
        )
        check(page.locator("[data-c311-data-state]").get_attribute("data-state") == expected, f"{label} expected {expected} state")

    open_c311(page, base_url, STATUS_PATH, "public_visitor", "retryable")
    wait_state("retryable-error")
    open_c311(page, base_url, STATUS_PATH, "public_visitor", "terminal")
    wait_state("terminal-error")
    open_c311(page, base_url, STATUS_PATH, "public_visitor", "version-conflict")
    wait_state("terminal-error")
    check(page.locator("[data-c311-server-version]").inner_text().find("2") >= 0, f"{label} did not display current server version")

    open_c311(page, base_url, SUBMIT_PATH, "public_visitor", "validation")
    page.locator(SUMMARY_SELECTOR).fill("valid summary")
    page.locator('[data-c311-action="submit-request"]').click()
    page.locator("[data-c311-error-summary]").wait_for(state="visible")
    check(page.evaluate("document.activeElement.matches('[data-c311-error-summary]')"), f"{label} validation summary did not receive focus")
    page.locator("[data-c311-error-summary] a").click()
    check(page.evaluate("document.activeElement.id") == SUMMARY_SELECTOR.removeprefix("#"), f"{label} invalid field did not receive focus")


def check_access_boundaries(page: Page, base_url: str, label: str) -> None:
    open_c311(page, base_url, STAFF_PATH, "service_agent")
    assert_route_focus(page, STAFF_PATH)
    check(page.get_by_role("link", name="Workflows").count() == 0, f"{label} exposed workflow action without capability")
    open_c311(page, base_url, WORKFLOWS_PATH, "service_agent")
    check(page.get_by_role("heading", name="Access denied").count() == 1, f"{label} capability denial was not shown")
    open_c311(page, base_url, STAFF_PATH, "service_agent", session="expired")
    check(page.get_by_role("heading", name="Sign-in required").count() == 1, f"{label} expired session did not show 401")
    open_c311(page, base_url, "/c311/does-not-exist", "public_visitor")
    check(page.get_by_role("heading", name="Not found").count() == 1, f"{label} unknown route did not show 404")


def check_admin_workflow(page: Page, base_url: str, label: str) -> None:
    open_c311(page, base_url, WORKFLOWS_PATH, "workflow_designer")
    assert_route_focus(page, WORKFLOWS_PATH)
    check(page.get_by_role("heading", name="Requests").count() == 1, f"{label} workflow fixture did not enter route")


def run() -> dict:
    global ACTIVE_ARTIFACT_DIR
    artifact_dir = artifact_directory()
    ACTIVE_ARTIFACT_DIR = artifact_dir
    results = {"viewports": list(VIEWPORTS), "checks": [], "started_at": datetime.now(timezone.utc).isoformat()}
    with sync_playwright() as playwright:
        browser_name = os.environ.get("C311_BROWSER", "chromium")
        browser_type = getattr(playwright, browser_name)
        browser = browser_type.launch(headless=True)
        try:
            for viewport in VIEWPORTS:
                for app, base_url, path, role in (
                    ("compose", COMPOSE_URL, SUBMIT_PATH, "public_visitor"),
                    ("admin", ADMIN_URL, STAFF_PATH, "service_agent"),
                ):
                    context = browser.new_context(viewport={"width": viewport, "height": 900})
                    page = context.new_page()
                    label = f"{app}@{viewport}"
                    open_c311(page, base_url, path, role)
                    assert_layout(page, viewport, label)
                    assert_route_focus(page, path)
                    if viewport == 360:
                        if app == "compose":
                            page.locator(STATUS_ROUTE_SELECTOR).click()
                            page.wait_for_url(f"**{STATUS_PATH}")
                            assert_route_focus(page, STATUS_PATH)
                        check_help_and_modal(page, base_url, label)
                        if app == "compose":
                            check_language_and_dirty(page, base_url, label)
                            check_portal_errors(page, base_url, label)
                        else:
                            check_access_boundaries(page, base_url, label)
                            check_admin_workflow(page, base_url, label)
                    screenshot = artifact_dir / f"{app}-{viewport}.png"
                    page.screenshot(path=str(screenshot), full_page=True)
                    results["checks"].append({"app": app, "viewport": viewport, "status": "passed", "screenshot": str(screenshot)})
                    context.close()
        finally:
            browser.close()
    results["finished_at"] = datetime.now(timezone.utc).isoformat()
    (artifact_dir / "fe01-matrix.json").write_text(json.dumps(results, indent=2), encoding="utf-8")
    return results


if __name__ == "__main__":
    try:
        summary = run()
        print(json.dumps(summary, indent=2))
    except Exception as error:
        failure_dir = ACTIVE_ARTIFACT_DIR or artifact_directory()
        blocker = {"status": "failed", "error": str(error), "browser": os.environ.get("C311_BROWSER", "chromium")}
        (failure_dir / "fe01-matrix-failure.json").write_text(json.dumps(blocker, indent=2), encoding="utf-8")
        print(json.dumps(blocker, indent=2), file=sys.stderr)
        raise

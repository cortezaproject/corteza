#!/usr/bin/env python3
"""Repeatable FE-02 public portal and identity checks."""

from __future__ import annotations

import json
import os
import sys
import tempfile
from datetime import datetime, timezone
from pathlib import Path
from urllib.parse import urlparse

from playwright.sync_api import Page, sync_playwright


VIEWPORTS = ((1440, 900), (768, 900), (390, 844))
ARTIFACT_DIR_ENV = "C311_ARTIFACT_DIR"
COMPOSE_URL = os.environ.get("C311_COMPOSE_URL", "http://127.0.0.1:18081").rstrip("/")
MAIN = "[data-c311-main]"
# Compose's optional development-only backend resources are proxied through
# CortezaAPI. They may be unavailable in the local mock run; every other
# client response error remains a gate failure. The origin is checked below
# so a same-named resource from another host cannot be silently ignored.
ALLOWED_HTTP_ERROR_PATHS: frozenset[str] = frozenset({"/code-snippets.js", "/custom.css"})
ALLOWED_WRITE_PATHS: frozenset[str] = frozenset({
    "/api/v1/session", "/api/v1/accounts", "/api/v1/auth/password-reset/request", "/api/v1/auth/password-reset/confirm",
    "/api/v1/auth/oidc/start", "/api/v1/auth/oidc/callback", "/api/v1/account/link/confirm", "/api/v1/account/profile",
    "/api/v1/account/password", "/api/v1/account/login-identifier", "/api/v1/preferences/language",
    "/api/v1/portal/service-requests", "/api/v1/portal/attachments",
})


def check(condition: bool, message: str) -> None:
    if not condition:
        raise AssertionError(message)


def artifact_directory() -> Path:
    configured = os.environ.get(ARTIFACT_DIR_ENV)
    directory = Path(configured).expanduser() if configured else Path(tempfile.mkdtemp(prefix="c311-fe02-gate-"))
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


def open_page(page: Page, path: str, scenario: str = "success", role: str = "public_visitor", session: str = "current") -> None:
    bootstrap(page, scenario, role, session)
    page.goto(f"{COMPOSE_URL}{path}", wait_until="domcontentloaded")
    page.locator(MAIN).first.wait_for(state="visible")


def assert_page(page: Page, label: str, width: int) -> None:
    check(page.locator(MAIN).count() == 1, f"{label} must have one main landmark")
    check(page.locator(f"{MAIN} h1:visible").count() == 1, f"{label} must have one visible h1")
    overflow = page.evaluate("document.documentElement.scrollWidth - document.documentElement.clientWidth")
    check(overflow <= 1, f"{label} overflows at {width}px: {overflow}")
    focused = page.evaluate("document.activeElement && document.activeElement.outerHTML") or ""
    check("<h1" in focused or "data-c311-main" in focused, f"{label} did not focus main content")


def attach_diagnostics(page: Page) -> dict[str, list]:
    diagnostics = {"console_errors": [], "page_errors": [], "writes": [], "http_errors": [], "unexpected_http_statuses": []}

    def on_console(message) -> None:
        if message.type == "error" and not (
            message.text.startswith("Failed to load resource:")
            or message.text.startswith("WebSocket connection to")
            or "net::ERR_CONNECTION_CLOSED" in message.text
        ):
            diagnostics["console_errors"].append(message.text)

    def on_page_error(error) -> None:
        diagnostics["page_errors"].append(str(error))

    def on_request(request) -> None:
        if request.method not in {"GET", "HEAD", "OPTIONS"}:
            parsed = urlparse(request.url)
            compose_url = urlparse(COMPOSE_URL)
            if parsed.scheme != compose_url.scheme or parsed.netloc != compose_url.netloc or parsed.path not in ALLOWED_WRITE_PATHS:
                diagnostics["writes"].append({"method": request.method, "url": request.url})

    def on_response(response) -> None:
        if response.status < 400:
            return
        error = {"status": response.status, "url": response.url}
        diagnostics["http_errors"].append(error)
        response_url = urlparse(response.url)
        compose_url = urlparse(COMPOSE_URL)
        allowed = (
            response_url.scheme == compose_url.scheme
            and response_url.netloc == compose_url.netloc
            and not response_url.query
            and not response_url.fragment
            and response_url.path in ALLOWED_HTTP_ERROR_PATHS
        )
        if not allowed:
            diagnostics["unexpected_http_statuses"].append(error)

    page.on("console", on_console)
    page.on("pageerror", on_page_error)
    page.on("request", on_request)
    page.on("response", on_response)
    return diagnostics


def check_public_navigation(page: Page) -> None:
    check(page.locator('[data-c311-route="/c311/sign-in"]').count() == 1, "anonymous sign-in link missing")
    check(page.locator('[data-c311-route="/c311/account"]').count() == 0, "anonymous account link exposed")
    page.locator('[data-c311-route="/c311/services"]').click()
    page.wait_for_url("**/c311/services")
    check(page.locator('[data-c311-page="services"]').get_attribute("data-c311-content-key") == "SERVICE_CATALOGUE", "SPA services navigation kept HOME content")
    page.locator('[data-c311-route="/c311/help"]').click()
    page.wait_for_url("**/c311/help")
    check(page.locator('[data-c311-page="help"]').get_attribute("data-c311-content-key") == "HELP", "SPA help navigation kept stale content")
    check("Describe the issue" in page.locator('[data-c311-page="help"]').inner_text(), "SPA help navigation lost contextual help")
    page.locator('[data-c311-route="/c311/sign-in"]').click()
    page.wait_for_url("**/c311/sign-in")
    check(page.locator('a[href="/c311/forgot-password"]').count() == 1, "anonymous forgot-password entry missing")


def check_focus_and_locale(page: Page) -> None:
    help_trigger = page.locator('[aria-controls^="c311-help-"]').first
    help_trigger.click()
    drawer = page.locator('[data-c311-help-drawer]')
    drawer.wait_for(state="visible")
    check(page.evaluate("document.activeElement.closest('[data-c311-help-drawer]') !== null"), "help focus did not enter drawer")
    page.keyboard.press("Tab")
    check(page.evaluate("document.activeElement.closest('[data-c311-help-drawer]') !== null"), "help focus escaped drawer")
    page.keyboard.press("Escape")
    drawer.wait_for(state="detached")
    check(page.evaluate("document.activeElement && document.activeElement.matches('[aria-controls^=\\\"c311-help-\\\"]')"), "help focus did not return")

    page.locator('[data-c311-route="/c311/submit"]').click()
    page.wait_for_url("**/c311/submit")
    page.locator("#c311-summary").fill("Draft kept")
    page.locator("#c311-description").fill("A non-sensitive draft description")
    page.locator('[data-c311-route="/c311/status"]').click()
    check("/c311/submit" in page.url, "dirty cancel did not preserve the form")
    page.once("dialog", lambda dialog: dialog.accept())
    page.locator('[data-c311-route="/c311/status"]').click()
    page.wait_for_url("**/c311/status")

    page.goto(f"{COMPOSE_URL}/c311/test/modal", wait_until="domcontentloaded")
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


def check_identity_forms(page: Page) -> None:
    open_page(page, "/c311/sign-in", "invalid-credentials")
    page.locator("#c311-login-identifier").fill("fixture")
    page.locator("#c311-login-password").fill("not-a-secret")
    page.locator('[data-c311-action="sign-in"]').click()
    page.locator('[data-c311-error-summary]').wait_for(state="visible")

    open_page(page, "/c311/register")
    page.locator("#c311-register-name").fill("Fixture visitor")
    page.locator("#c311-register-email").fill("visitor@example.test")
    page.wait_for_function("() => sessionStorage.getItem('c311.identity.register') !== null")
    check("password" not in (page.evaluate("sessionStorage.getItem('c311.identity.register')") or ""), "registration password was persisted")
    page.reload(wait_until="domcontentloaded")
    check(page.locator("#c311-register-email").input_value() == "visitor@example.test", "registration email was not restored")
    page.locator('[data-c311-action="register"]').click()
    page.locator('[data-c311-error-summary]').wait_for(state="visible")
    check(page.evaluate("document.activeElement && document.activeElement.matches('[data-c311-error-summary]')"), "registration summary did not receive focus")

    open_page(page, "/c311/forgot-password")
    page.locator("#c311-forgot-email").fill("alex@example.test")
    page.locator('[data-c311-action="forgot-password"]').click()
    known_response = page.get_by_role("status").last.inner_text()
    page.locator("#c311-forgot-email").fill("unknown@example.test")
    page.locator('[data-c311-action="forgot-password"]').click()
    check(page.get_by_role("status").last.inner_text() == known_response, "forgot-password response disclosed account existence")

    open_page(page, "/c311/sign-in", "oidc-failure")
    page.locator('[data-c311-action="oidc-sign-in"]').click()
    page.locator('[data-c311-error-summary]').wait_for(state="visible")

    open_page(page, "/c311/auth/callback?provider=oidc&error=access_denied")
    page.locator('[data-c311-error-summary]').wait_for(state="visible")

    open_page(page, "/c311/sign-in", "successful-login")
    page.locator("#c311-login-identifier").fill("fixture")
    page.locator("#c311-login-password").fill("not-a-secret")
    page.locator('[data-c311-action="sign-in"]').click()
    page.wait_for_url("**/c311")
    check(page.locator('[data-c311-route="/c311/account"]').count() == 1, "successful login did not expose account navigation")

    open_page(page, "/c311/reset-password?token=ephemeral-token")
    check("token=" not in page.url, "reset token was not removed from URL")
    page.locator("#c311-reset-password").fill("ValidPassword1!")
    page.locator('[data-c311-action="reset-password"]').click()
    check(page.get_by_role("status").count() >= 1, "reset success status missing")
    page.locator('[data-c311-action="reset-password"]').click()
    page.locator('[data-c311-error-summary]').wait_for(state="visible")
    check("INVALID_RESET_TOKEN" in page.locator('[data-c311-error-summary]').inner_text() or "invalid" in page.locator('[data-c311-error-summary]').inner_text().lower(), "used reset token was accepted twice")

    open_page(page, "/c311/reset-password?token=ephemeral-token", "expired-reset-token")
    page.locator("#c311-reset-password").fill("ValidPassword1!")
    page.locator('[data-c311-action="reset-password"]').click()
    page.locator('[data-c311-error-summary]').wait_for(state="visible")

    open_page(page, "/c311/sign-in", "link-confirmation-required")
    page.locator('[data-c311-action="oidc-sign-in"]').click()
    page.wait_for_url("**/c311/auth/link/confirm*")
    check("authorization_url" not in page.url and "token" not in page.url, "link confirmation exposed redirect data")
    check(page.locator('[data-c311-page="link-confirm"]').count() == 1, "account-link confirmation page missing")
    page.reload(wait_until="domcontentloaded")
    page.locator('[data-c311-page="link-confirm"]').wait_for(state="visible")
    page.locator('[data-c311-action="confirm-link"]').click()
    page.locator('[data-c311-page="link-confirm"] .alert-success').wait_for(state="visible")
    check("authorization_url" not in page.url and "token" not in page.url, "confirmed redirect data exposed")

    open_page(page, "/c311/sign-in", "link-confirmation-required")
    page.locator('[data-c311-action="oidc-sign-in"]').click()
    page.wait_for_url("**/c311/auth/link/confirm")
    page.locator('[data-c311-action="cancel-link-confirm"]').click()
    page.wait_for_url("**/c311/sign-in")

    open_page(page, "/c311/auth/callback?provider=oidc&error=access_denied", "identity-claims-failure")
    page.locator('[data-c311-error-summary]').wait_for(state="visible")


def check_content_states(page: Page) -> None:
    open_page(page, "/c311", "branding-failure")
    check(page.locator('[data-c311-data-state]').get_attribute("data-state") == "retryable-error", "branding failure did not expose retryable state")
    open_page(page, "/c311/services", "empty-catalogue")
    check(page.locator('[data-c311-page="services"]').count() == 1, "empty catalogue page missing")
    open_page(page, "/c311/services", "content-loading-failure")
    check(page.locator('[data-c311-data-state]').get_attribute("data-state") == "retryable-error", "content failure did not expose retryable state")
    open_page(page, "/c311/help", "help-loading-failure")
    check(page.locator('[data-c311-data-state]').get_attribute("data-state") == "retryable-error", "help failure did not expose retryable state")
    open_page(page, "/c311/help")
    check(page.locator('[data-c311-branding]').count() == 1, "branding fields were not rendered")
    check("Find answers" in page.locator('[data-c311-page="help"]').inner_text(), "public HELP content was not rendered")
    open_page(page, "/c311/services", "terminal")
    check(page.locator('[data-c311-data-state]').get_attribute("data-state") == "terminal-error", "terminal content failure did not stop with terminal state")
    open_page(page, "/c311/account", "success", "public_visitor")
    check("/c311/401" in page.url or page.get_by_role("heading", name="Sign-in required").count() == 1, "private account did not enforce 401")
    for path, heading in (("/c311/403", "Access denied"), ("/c311/404", "Not found")):
        page.goto(f"{COMPOSE_URL}{path}", wait_until="domcontentloaded")
        check(page.get_by_role("heading", name=heading).count() == 1, f"{path} status page missing")


def check_authenticated_navigation(page: Page) -> None:
    open_page(page, "/c311/requests", "empty-my-requests", "constituent")
    check(page.locator('[data-c311-data-state]').get_attribute("data-state") == "empty", "empty my-requests fixture did not render empty state")
    open_page(page, "/c311/requests", "success", "constituent")
    check(page.locator('[data-c311-page="requests"]').count() == 1, "authenticated requests page missing")
    check(page.locator('[data-c311-route="/c311/account"]').count() == 1, "authenticated account link missing")
    page.locator('[data-c311-route="/c311/account"]').click()
    page.wait_for_url("**/c311/account")
    page.locator("#c311-account-login-identifier").fill("alex.new")
    page.locator("#c311-account-current-password").fill("Current-password-1!")
    page.locator('[data-c311-action="change-login-identifier"]').click()
    page.get_by_role("status").filter(has_text="Login identifier changed").wait_for(state="visible")
    page.locator("#c311-account-new-password").fill("ValidPassword1!")
    page.locator('[data-c311-action="change-password"]').click()
    page.get_by_role("status").filter(has_text="Password changed").wait_for(state="visible")
    open_page(page, "/c311/account", "version-conflict", "constituent")
    page.locator("#c311-account-login-identifier").fill("alex.conflict")
    page.locator("#c311-account-current-password").fill("Current-password-1!")
    page.locator('[data-c311-action="change-login-identifier"]').click()
    page.locator('[data-c311-error-summary]').wait_for(state="visible")
    check(page.locator("#c311-account-login-identifier").input_value() == "alex.conflict", "login identifier conflict did not preserve input")
    page.once("dialog", lambda dialog: dialog.accept())
    page.locator('[data-c311-route="/c311/logout/callback"]').click()
    page.wait_for_url("**/c311/logout/callback*")
    check(page.get_by_role("status").count() >= 1, "logout status missing")
    check(page.locator('[data-c311-route="/c311/account"]').count() == 0, "SPA logout kept authenticated navigation")
    check(page.locator('[data-c311-route="/c311/requests"]').count() == 0, "SPA logout kept requests navigation")
    open_page(page, "/c311/logout/callback", "federated-logout-failure", "constituent")
    check(page.locator('[data-c311-data-state]').get_attribute("data-state") == "terminal-error", "federated logout failure did not stop with terminal state")


def run() -> dict:
    artifacts = artifact_directory()
    results = {"viewports": [list(viewport) for viewport in VIEWPORTS], "checks": [], "started_at": datetime.now(timezone.utc).isoformat()}
    with sync_playwright() as playwright:
        browser = playwright.chromium.launch(headless=True)
        try:
            for width, height in VIEWPORTS:
                context = browser.new_context(viewport={"width": width, "height": height})
                page = context.new_page()
                diagnostics = attach_diagnostics(page)
                label = f"compose@{width}x{height}"
                open_page(page, "/c311")
                assert_page(page, label, width)
                check_public_navigation(page)
                if width in {390, 768, 1440}:
                    page.goto(f"{COMPOSE_URL}/c311", wait_until="domcontentloaded")
                    page.locator(MAIN).first.wait_for(state="visible")
                    locale = page.locator("select[data-c311-language]")
                    locale.select_option("es")
                    check(page.evaluate("localStorage.getItem('c311.locale.anonymous')") == "es", "Spanish locale was not persisted")
                    page.reload(wait_until="domcontentloaded")
                    check(locale.input_value() == "es", "Spanish locale was not restored")
                    check_focus_and_locale(page)
                    check_identity_forms(page)
                    check_content_states(page)
                    check_authenticated_navigation(page)
                screenshot = artifacts / f"fe02-compose-{width}x{height}.png"
                page.screenshot(path=str(screenshot), full_page=True)
                check(not diagnostics["page_errors"], f"{label} has uncaught page errors: {diagnostics['page_errors']}")
                check(not diagnostics["console_errors"], f"{label} has console errors: {diagnostics['console_errors']}")
                check(not diagnostics["writes"], f"{label} issued unexpected writes: {diagnostics['writes']}")
                check(not diagnostics["unexpected_http_statuses"], f"{label} has unexpected HTTP statuses: {diagnostics['unexpected_http_statuses']}")
                results["checks"].append({"label": label, "status": "passed", "screenshot": str(screenshot), "diagnostics": diagnostics})
                context.close()
        finally:
            browser.close()
    results["finished_at"] = datetime.now(timezone.utc).isoformat()
    (artifacts / "fe02-matrix.json").write_text(json.dumps(results, indent=2), encoding="utf-8")
    return results


if __name__ == "__main__":
    try:
        print(json.dumps(run(), indent=2))
    except Exception as error:
        print(json.dumps({"status": "failed", "error": str(error)}), file=sys.stderr)
        raise

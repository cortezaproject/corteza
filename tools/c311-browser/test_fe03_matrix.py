from __future__ import annotations

import os
import stat
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

from fe03_matrix import ADMIN_URL, ARTIFACT_DIR_ENV, COMPOSE_URL, GENERIC_CONNECTION_ERROR, GENERIC_NAME_RESOLUTION_ERROR, KNOWN_DEV_CONSOLE_ERRORS, KNOWN_DEV_LOCALE_FAILURES, KNOWN_DEV_RESOURCE_FAILURES, KNOWN_DEV_WEBSOCKET_URL, STAFF_SUBMIT_ACTION, SUBMIT_ACTION, VIEWPORTS, artifact_directory, finalize_diagnostics, record_console


class Fe03MatrixTests(unittest.TestCase):
    def test_matrix_uses_contract_fields_and_stable_selectors(self) -> None:
        source = Path(__file__).with_name("fe03_matrix.py").read_text(encoding="utf-8")
        self.assertEqual(VIEWPORTS, ((1440, 900), (768, 900), (390, 844)))
        self.assertTrue(COMPOSE_URL.startswith("http"))
        self.assertTrue(ADMIN_URL.startswith("http"))
        self.assertNotEqual(COMPOSE_URL.rsplit(":", 1)[-1], "18082")
        self.assertNotEqual(ADMIN_URL.rsplit(":", 1)[-1], "18082")
        for selector in ("#c311-service-type", "#c311-requester-email", "#c311-attachment-file", "#c311-consent", '[data-c311-action="save-draft"]', SUBMIT_ACTION, STAFF_SUBMIT_ACTION, "[data-c311-submission-result]"):
            self.assertIn(selector, source)
        self.assertIn("C311Mode = 'mock'", source)
        self.assertNotIn("source_channel", source)
        self.assertIn("unexpected_responses", source)
        self.assertIn("unexpected_console_errors", source)
        self.assertIn("urlparse", source)
        self.assertEqual(KNOWN_DEV_RESOURCE_FAILURES, {"/code-snippets.js", "/custom.css"})
        self.assertIn("Failed to load resource: the server responded with a status of 500 (Internal Server Error)", KNOWN_DEV_CONSOLE_ERRORS)
        self.assertIn("WebSocket connection to 'wss://api.cortezaproject.your-domain.tld/websocket' failed: Error in connection establishment: net::ERR_NAME_NOT_RESOLVED", KNOWN_DEV_CONSOLE_ERRORS)
        self.assertEqual(GENERIC_CONNECTION_ERROR, "Failed to load resource: net::ERR_CONNECTION_CLOSED")
        self.assertEqual(GENERIC_NAME_RESOLUTION_ERROR, "Failed to load resource: net::ERR_NAME_NOT_RESOLVED")
        self.assertTrue(KNOWN_DEV_WEBSOCKET_URL.startswith("wss://"))
        self.assertIn("parsed.path == \"/ws\"", source)
        self.assertEqual(KNOWN_DEV_LOCALE_FAILURES, {
            "/system/locale/en-US/corteza-webapp-compose",
            "/system/locale/en-US/corteza-webapp-admin",
            "/system/locale/en-US+en/corteza-webapp-compose",
            "/system/locale/en-US+en/corteza-webapp-admin",
            "/system/locale/en/corteza-webapp-compose",
            "/system/locale/en/corteza-webapp-admin",
            "/system/locale/en+en-US/corteza-webapp-compose",
            "/system/locale/en+en-US/corteza-webapp-admin",
        })
        self.assertIn("set_input_files", source)
        self.assertIn("attachment_count", source)
        self.assertIn("data-c311-attachment-recovery", source)
        self.assertIn("check_capability_controls", source)
        self.assertIn("check_restart_recovery", source)
        self.assertIn('scenario="forbidden"', source)

    def test_configured_artifact_directory_is_private(self) -> None:
        with tempfile.TemporaryDirectory(prefix="c311-fe03-test-") as parent:
            configured = Path(parent) / "artifacts"
            with patch.dict(os.environ, {ARTIFACT_DIR_ENV: str(configured)}):
                directory = artifact_directory()
            self.assertEqual(stat.S_IMODE(directory.stat().st_mode), 0o700)

    def test_unmatched_generic_connection_error_is_not_whitelisted(self) -> None:
        result = {
            "console_errors": [GENERIC_CONNECTION_ERROR],
            "unexpected_console_errors": [],
            "failed_requests": [],
            "websocket_urls": [],
            "pending_console_errors": [GENERIC_CONNECTION_ERROR],
        }
        finalize_diagnostics(result)
        self.assertEqual(result["unexpected_console_errors"], [GENERIC_CONNECTION_ERROR])

    def test_known_locale_dns_console_error_is_whitelisted_by_exact_request(self) -> None:
        result = {
            "console_errors": [GENERIC_NAME_RESOLUTION_ERROR],
            "unexpected_console_errors": [],
            "failed_requests": [{"url": "https://api.cortezaproject.your-domain.tld/system/locale/en-US/corteza-webapp-compose"}],
            "websocket_urls": [],
            "pending_console_errors": [GENERIC_NAME_RESOLUTION_ERROR],
        }
        finalize_diagnostics(result)
        self.assertEqual(result["unexpected_console_errors"], [])

    def test_unmatched_generic_name_resolution_error_is_not_whitelisted(self) -> None:
        result = {
            "console_errors": [GENERIC_NAME_RESOLUTION_ERROR],
            "unexpected_console_errors": [],
            "failed_requests": [],
            "websocket_urls": [],
            "pending_console_errors": [GENERIC_NAME_RESOLUTION_ERROR],
        }
        finalize_diagnostics(result)
        self.assertEqual(result["unexpected_console_errors"], [GENERIC_NAME_RESOLUTION_ERROR])

    def test_late_console_event_after_finalize_is_safe(self) -> None:
        result = {
            "console_errors": [],
            "unexpected_console_errors": [],
            "pending_console_errors": [],
            "failed_requests": [],
            "websocket_urls": [],
        }
        finalize_diagnostics(result)

        class Message:
            type = "error"
            text = GENERIC_CONNECTION_ERROR

        record_console(result, Message())
        self.assertEqual(result["pending_console_errors"], [GENERIC_CONNECTION_ERROR])


if __name__ == "__main__":
    unittest.main()

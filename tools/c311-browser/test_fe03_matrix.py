from __future__ import annotations

import os
import stat
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

from fe03_matrix import ADMIN_URL, ARTIFACT_DIR_ENV, COMPOSE_URL, KNOWN_DEV_RESOURCE_FAILURES, STAFF_SUBMIT_ACTION, SUBMIT_ACTION, VIEWPORTS, artifact_directory


class Fe03MatrixTests(unittest.TestCase):
    def test_matrix_uses_contract_fields_and_stable_selectors(self) -> None:
        source = Path(__file__).with_name("fe03_matrix.py").read_text(encoding="utf-8")
        self.assertEqual(VIEWPORTS, ((1440, 900), (768, 900), (390, 844)))
        self.assertTrue(COMPOSE_URL.startswith("http"))
        self.assertTrue(ADMIN_URL.startswith("http"))
        for selector in ("#c311-service-type", "#c311-requester-email", "#c311-attachment-file", "#c311-consent", '[data-c311-action="save-draft"]', SUBMIT_ACTION, STAFF_SUBMIT_ACTION, "[data-c311-submission-result]"):
            self.assertIn(selector, source)
        self.assertIn("C311Mode = 'mock'", source)
        self.assertNotIn("source_channel", source)
        self.assertIn("unexpected_responses", source)
        self.assertIn("urlparse", source)
        self.assertEqual(KNOWN_DEV_RESOURCE_FAILURES, {"/code-snippets.js", "/custom.css"})
        self.assertIn("set_input_files", source)

    def test_configured_artifact_directory_is_private(self) -> None:
        with tempfile.TemporaryDirectory(prefix="c311-fe03-test-") as parent:
            configured = Path(parent) / "artifacts"
            with patch.dict(os.environ, {ARTIFACT_DIR_ENV: str(configured)}):
                directory = artifact_directory()
            self.assertEqual(stat.S_IMODE(directory.stat().st_mode), 0o700)


if __name__ == "__main__":
    unittest.main()

#!/usr/bin/env python3
"""Static and helper checks for the FE-02 browser matrix."""

from __future__ import annotations

import os
import stat
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

from fe02_matrix import ARTIFACT_DIR_ENV, VIEWPORTS, artifact_directory


class Fe02MatrixTests(unittest.TestCase):
    def test_matrix_contains_required_viewports_and_stable_selectors(self) -> None:
        source = Path(__file__).with_name("fe02_matrix.py").read_text(encoding="utf-8")
        self.assertEqual(VIEWPORTS, ((1440, 900), (768, 900), (390, 844)))
        self.assertIn('[data-c311-route="/c311/services"]', source)
        self.assertIn('[data-c311-action="sign-in"]', source)
        self.assertIn("attach_diagnostics", source)
        self.assertIn("unexpected writes", source)
        self.assertIn("unexpected HTTP statuses", source)
        self.assertIn("response.status", source)
        self.assertIn('"/code-snippets.js"', source)
        self.assertIn('"/custom.css"', source)
        self.assertIn("SERVICE_CATALOGUE", source)
        self.assertIn("SPA logout kept authenticated navigation", source)
        self.assertIn('[data-c311-action="change-login-identifier"]', source)
        self.assertIn('[data-c311-action="cancel-link-confirm"]', source)
        self.assertNotIn('/c311/security-notice', source)
        self.assertIn("C311Mode = 'mock'", source)
        self.assertNotIn("password=", source)
        sonar_properties = Path(__file__).parents[2].joinpath("sonar-project.properties").read_text(encoding="utf-8")
        self.assertIn("sonar.coverage.exclusions=tools/c311-browser/fe01_matrix.py,tools/c311-browser/fe02_matrix.py", sonar_properties)
        self.assertNotIn("sonar.coverage.exclusions=tools/c311-browser/**", sonar_properties)
        runner = Path(__file__).with_name("run-fe02.sh").read_text(encoding="utf-8")
        self.assertIn("config.example.js", runner)
        self.assertIn("config.js must not be a symbolic link", runner)

    def test_configured_artifact_directory_is_private(self) -> None:
        with tempfile.TemporaryDirectory(prefix="c311-fe02-test-") as parent:
            configured = Path(parent) / "artifacts"
            with patch.dict(os.environ, {ARTIFACT_DIR_ENV: str(configured)}):
                directory = artifact_directory()
            self.assertEqual(stat.S_IMODE(directory.stat().st_mode), 0o700)


if __name__ == "__main__":
    unittest.main()

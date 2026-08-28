#!/usr/bin/env python3
"""Unit checks for the FE-01 browser matrix helpers."""

from __future__ import annotations

import os
import stat
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

from fe01_matrix import ARTIFACT_DIR_ENV, artifact_directory


class ArtifactDirectoryTests(unittest.TestCase):
    def test_browser_dependencies_are_pinned_and_hashed(self) -> None:
        requirements = Path(__file__).with_name("requirements.txt").read_text(encoding="utf-8")
        package_lines = [
            line.strip()
            for line in requirements.splitlines()
            if "==" in line and not line.lstrip().startswith("#")
        ]
        packages = {line.split("==", 1)[0] for line in package_lines}
        self.assertEqual(packages, {"playwright", "pyee", "typing-extensions", "greenlet"})
        self.assertNotIn("playwright>=", requirements)
        self.assertGreaterEqual(requirements.count("--hash=sha256:"), len(packages))

    def test_default_directory_is_private(self) -> None:
        with patch.dict(os.environ, {ARTIFACT_DIR_ENV: ""}):
            directory = artifact_directory()
        try:
            self.assertTrue(directory.is_dir())
            self.assertEqual(stat.S_IMODE(directory.stat().st_mode), 0o700)
        finally:
            directory.rmdir()

    def test_configured_directory_is_private_and_not_a_symlink(self) -> None:
        with tempfile.TemporaryDirectory(prefix="c311-artifact-test-") as parent:
            configured = Path(parent) / "artifacts"
            with patch.dict(os.environ, {ARTIFACT_DIR_ENV: str(configured)}):
                directory = artifact_directory()
            self.assertEqual(directory, configured.resolve())
            self.assertEqual(stat.S_IMODE(directory.stat().st_mode), 0o700)

            link = Path(parent) / "link"
            link.symlink_to(directory, target_is_directory=True)
            with patch.dict(os.environ, {ARTIFACT_DIR_ENV: str(link)}):
                with self.assertRaises(ValueError):
                    artifact_directory()


if __name__ == "__main__":
    unittest.main()

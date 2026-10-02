import { defineConfig, devices } from '@playwright/test'
import { existsSync, readFileSync } from 'node:fs'

// Local settings come from .env.e2e (gitignored, see .env.e2e.example).
// The suite attaches to an already running stack (server + webapp dev
// servers); it never boots or seeds anything itself.
if (existsSync(new URL('./.env.e2e', import.meta.url))) {
  for (const line of readFileSync(new URL('./.env.e2e', import.meta.url), 'utf8').split('\n')) {
    const m = line.match(/^([A-Z0-9_]+)=(.*)$/)
    if (m && process.env[m[1]] === undefined) process.env[m[1]] = m[2]
  }
}

export default defineConfig({
  testDir: './specs',
  outputDir: './.results',
  fullyParallel: false, // the stack's data is shared between specs
  retries: 0,
  workers: 1,
  reporter: [['list']],
  use: {
    baseURL: process.env.E2E_COMPOSE_URL || 'http://127.0.0.1:8081',
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure',
  },
  projects: [
    {
      name: 'setup',
      testMatch: /auth\.setup\.ts/,
    },
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome'],
        storageState: '.auth/state.json',
      },
      dependencies: ['setup'],
    },
  ],
})

import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: true,
  webServer: {
    command: 'cd .. && make dev-playwright',
    port: 8090,
    reuseExistingServer: false, // !process.env.CI,
    timeout: 120000, // 2 minutes timeout
  },
  use: {
    baseURL: 'http://localhost:8090/ui/',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
})

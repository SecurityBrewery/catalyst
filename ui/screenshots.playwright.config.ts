import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './screenshots',
  fullyParallel: true,
  webServer: {
    command: 'make -C .. dev-10000',
    port: 8090,
    reuseExistingServer: !process.env.CI,
    timeout: 120000,
    stderr: 'ignore'
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

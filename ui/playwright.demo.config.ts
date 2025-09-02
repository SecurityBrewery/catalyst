import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './tests/demo',
  fullyParallel: true,
  webServer: {
    command: 'make -C .. dev-demo',
    port: 8090,
    reuseExistingServer: !process.env.CI,
    timeout: 120000,
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

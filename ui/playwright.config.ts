import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: true,
  webServer: {
    command: 'bash -c "rm -rf ../catalyst_data && mkdir ../catalyst_data && cp ../upgradetest/data/v0.14.1/data.db ../catalyst_data/data.db && cd .. && go run ."',
    port: 8090,
    reuseExistingServer: !process.env.CI,
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

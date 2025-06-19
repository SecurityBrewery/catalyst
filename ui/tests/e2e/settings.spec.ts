import { expect } from '@playwright/test'
import { login, test } from './util'

const goToSettings = async (page) => {
  await page.goto('settings')
  await expect(page.getByRole('heading', { name: 'Settings' })).toBeVisible()
}

test('settings page shows existing settings', async ({ page }) => {
  await login(page)
  await goToSettings(page)
  await expect(page.locator('#meta\\.appName')).toHaveValue('Catalyst')
})

const updates = [
  { field: 'app name', selector: '#meta\\.appName', value: 'Catalyst Playwright' },
  { field: 'app url', selector: '#meta\\.appUrl', value: 'https://playwright.example.com' },
  { field: 'sender name', selector: '#meta\\.senderName', value: 'Playwright' },
  { field: 'sender address', selector: '#meta\\.senderAddress', value: 'playwright@example.com' }
]

test.describe('update settings', () => {
  for (const { field, selector, value } of updates) {
    test(`can update ${field}`, async ({ page }) => {
      await login(page)
      await goToSettings(page)
      await page.locator(selector).fill(value)
      const saveBtn = page.getByRole('button', { name: 'Save' }).last()
      await expect(saveBtn).toBeEnabled()
      await saveBtn.click()
      await expect(page.locator(selector)).toHaveValue(value)
    })
  }

  test('can enable smtp', async ({ page }) => {
    await login(page)
    await goToSettings(page)
    const smtpSwitch = page.getByRole('switch').first()
    await smtpSwitch.click()
    const saveBtn = page.getByRole('button', { name: 'Save' }).last()
    await expect(saveBtn).toBeEnabled()
    await saveBtn.click()
    await expect(smtpSwitch).toHaveAttribute('data-state', 'checked')
  })
})

import { expect } from '@playwright/test'
import { login, test } from './util'

test('dashboard title is visible', async ({ page }) => {
  await login(page)
  await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible()
})

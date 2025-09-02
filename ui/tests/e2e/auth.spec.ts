import { expect } from '@playwright/test'
import { login, test } from './util'

// Verify that login stores a token and grants API access

test('login grants api access', async ({ page }) => {
  await login(page)
  await expect(page).toHaveURL(/.*\/dashboard/)

  const token = await page.evaluate(() => localStorage.getItem('token'))
  expect(token).toBeTruthy()

  const response = await page.request.get('/auth/user', {
    headers: { Authorization: `Bearer ${token}` }
  })
  expect(response.status()).toBe(200)
  const data = await response.json()
  expect(data.user.email).toBe('admin@catalyst-soar.com')
})

// Verify that logout clears the token and api requests fail without it

test('logout denies api access', async ({ page }) => {
  await login(page)
  const button = page.getByRole('button', { name: /admin/i })
  await button.click()
  await page.getByRole('menuitem', { name: 'Log out' }).click()
  await page.waitForURL('**/login')

  const token = await page.evaluate(() => localStorage.getItem('token'))
  expect(token).toBe('')

  const response = await page.request.get('/auth/user')
  const data = await response.json()
  expect(data).toBeNull()
})

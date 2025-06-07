import { test, expect } from '@playwright/test'
import { randomUUID } from 'crypto'

const login = async (page) => {
  await page.goto('login')
  await page.getByPlaceholder('Username').fill('user@catalyst-soar.com')
  await page.getByPlaceholder('Password').fill('1234567890')
  await page.getByRole('button', { name: 'Login' }).click()
  await page.waitForURL('**/dashboard')
}

const createUser = async (page, username: string) => {
  await page.goto('users')
  await page.getByRole('button', { name: 'New User' }).click()
  await page.waitForURL('**/users/new')
  await page.locator('#username').fill(username)
  await page.locator('#email').fill(`${username}@example.com`)
  await page.locator('#name').fill(username)
  await page.getByRole('button', { name: 'Save' }).last().click()
  await page.waitForURL('**/users/*')
}

test('users list shows existing users', async ({ page }) => {
  await login(page)
  await page.goto('users')
  await expect(page.getByRole('heading', { name: 'Users' })).toBeVisible()
  await expect(page.getByText('user@catalyst-soar.com')).toBeVisible()
})

test('can create a user', async ({ page }) => {
  await login(page)
  const username = `playwright-${randomUUID()}`
  await createUser(page, username)
  await expect(page.locator('#username')).toHaveValue(username)
})

test.describe('update a user', () => {
  const updates = [
    {
      field: 'email',
      selector: '#email',
      value: 'updated@example.com'
    },
    {
      field: 'name',
      selector: '#name',
      value: 'Updated Name'
    },
  ]

  for (const { field, selector, value } of updates) {
    test(`can update ${field}`, async ({ page }) => {
      await login(page)
      const username = `playwright-${randomUUID()}`
      await createUser(page, username)
      await page.waitForSelector(selector)
      await page.locator(selector).fill(value)
      const saveBtn = page.getByRole('button', { name: 'Save' }).last()
      await expect(saveBtn).toBeEnabled()
      await saveBtn.click()
      await expect(page.locator(selector)).toHaveValue(value)
    })
  }
})

test('can update username', async ({ page }) => {
  await login(page)
  const username = `playwright-${randomUUID()}`
  await createUser(page, username)
  await page.waitForSelector('#username')
  const newUsername = `playwright-${randomUUID()}`
  await page.locator('#username').fill(newUsername)
  const saveBtn = page.getByRole('button', { name: 'Save' }).last()
  await expect(saveBtn).toBeEnabled()
  await saveBtn.click()
  await expect(page.locator('#username')).toHaveValue(newUsername)
})

test('can delete a user', async ({ page }) => {
  await login(page)
  const username = `playwright-${randomUUID()}`
  await createUser(page, username)
  await page.waitForURL('**/users/*')
  await page.getByRole('button', { name: 'Delete User' }).click()
  await page.getByRole('dialog').getByRole('button', { name: 'Delete' }).click()
  await page.waitForURL('**/users')
  await expect(page.locator(`text=${username}`)).toHaveCount(0)
})

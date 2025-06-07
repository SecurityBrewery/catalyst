import { test, expect } from '@playwright/test'
import { randomUUID } from 'crypto'

const login = async (page) => {
  await page.goto('login')
  await page.getByPlaceholder('Username').fill('user@catalyst-soar.com')
  await page.getByPlaceholder('Password').fill('1234567890')
  await page.getByRole('button', { name: 'Login' }).click()
  await page.waitForURL('**/dashboard')
}

const createReaction = async (page, name: string) => {
  await page.goto('reactions')
  await page.getByRole('button', { name: 'New Reaction' }).click()
  await page.waitForURL('**/reactions/new')
  await page.locator('#name').fill(name)
  await page.locator('#trigger').selectOption('Schedule')
  await page.locator('#expression').fill('* * * * *')
  await page.locator('#action').selectOption('Python')
  await page.locator('#script').fill('print("Hello, world!")')
  await page.getByRole('button', { name: 'Save' }).last().click()
  await page.waitForURL('**/reactions/*')
}

test('reactions list shows existing reactions', async ({ page }) => {
  await login(page)
  await page.goto('reactions')
  await expect(page.getByRole('heading', { name: 'Reactions' })).toBeVisible()
})

test('can create a reaction', async ({ page }) => {
  await login(page)
  const name = `playwright-${randomUUID()}`
  await createReaction(page, name)
  await expect(page.locator('#name')).toHaveValue(name)
  await expect(page.locator('#trigger')).toHaveValue('schedule')
  await expect(page.locator('#expression')).toHaveValue('* * * * *')
  await expect(page.locator('#action')).toHaveValue('python')
  await expect(page.locator('#script')).toHaveValue('print("Hello, world!")')
})

test.describe('update a reaction', () => {
  const updates = [
    {
      field: 'name',
      update: async (page) => {
        await page.waitForSelector('#name')
        await page.locator('#name').fill("Updated Reaction")
      },
      selector: '#name',
      value: 'Updated Reaction'
    },
    {
      field: 'trigger',
      update: async (page) => {
        await page.waitForSelector('#trigger')
        await page.locator('button').filter({ hasText: 'Schedule' }).click()
        await page.getByRole('option', { name: 'HTTP / Webhook' }).click()
        await page.locator('#path').fill('webhook')
      },
      selector: '#trigger',
      value: 'webhook'
    },
    {
      field: 'action',
      update: async (page) => {
        await page.waitForSelector('#action')
        await page.locator('button').filter({ hasText: 'Python' }).click()
        await page.getByRole('option', { name: 'HTTP / Webhook' }).click()
        await page.locator('#url').fill('https://example.com')
      },
      selector: '#action',
      value: 'webhook'
    },
    {
      field: 'script',
      update: async (page) => {
        await page.waitForSelector('#script')
        await page.locator('#script').fill('pass')
      },
      selector: '#script',
      value: 'pass'
    },
  ]

  for (const { field, update, selector, value } of updates) {
    test(`can update ${field}`, async ({ page }) => {
      await login(page)
      const name = `playwright-${randomUUID()}`
      await createReaction(page, name)
      await page.waitForURL('**/reactions/*')
      await update(page)
      const saveBtn = page.getByRole('button', { name: 'Save' }).last()
      await expect(saveBtn).toBeEnabled()
      await saveBtn.click()
      await expect(page.locator(selector)).toHaveValue(value)
    })
  }
})

test('can delete a reaction', async ({ page }) => {
  await login(page)
  const name = `playwright-${randomUUID()}`
  await createReaction(page, name)
  await page.waitForURL('**/reactions/*')
  await page.getByRole('button', { name: 'Delete Reaction' }).click()
  await page.getByRole('dialog').getByRole('button', { name: 'Delete' }).click()
  await page.waitForURL('**/reactions')
  await expect(page.locator(`text=${name}`)).toHaveCount(0)
})

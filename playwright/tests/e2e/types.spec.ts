import { test, expect } from '@playwright/test'
import { randomUUID } from 'crypto'

const login = async (page) => {
  await page.goto('login')
  await page.getByPlaceholder('Username').fill('user@catalyst-soar.com')
  await page.getByPlaceholder('Password').fill('1234567890')
  await page.getByRole('button', { name: 'Login' }).click()
  await page.waitForURL('**/dashboard')
}

const createType = async (page, name: string) => {
  await page.goto('types')
  await page.getByRole('button', { name: 'New Type' }).click()
  await page.waitForURL('**/types/new')
  await page.locator('#singular').fill(name)
  await page.locator('#plural').fill(`${name}s`)
  await page.locator('#icon input').fill('Bug')
  await page.locator('#schema').fill('{}')
  await page.getByRole('button', { name: 'Save' }).last().click()
  await page.waitForURL('**/types/*')
}

test('types list shows incident', async ({ page }) => {
  await login(page)
  await page.goto('types')
  await expect(page.getByRole('heading', { name: 'Types' })).toBeVisible()
  await expect(page.getByText('Incident')).toBeVisible()
})

test('can create a type', async ({ page }) => {
  await login(page)
  const name = `Playwright-${randomUUID()}`
  await createType(page, name)
  await expect(page.locator('#singular')).toHaveValue(name)
  await expect(page.locator('#plural')).toHaveValue(`${name}s`)
  await expect(page.locator('#icon input')).toHaveValue('Bug')
  await expect(page.locator('#schema')).toHaveValue('{}')
})

test.describe('update a type', () => {
  const updates = [
    { field: 'singular', selector: '#singular', value: 'UpdatedSingular' },
    { field: 'plural', selector: '#plural', value: 'UpdatedPlural' },
    { field: 'icon', selector: '#icon input', value: 'Activity' },
    { field: 'schema', selector: '#schema', value: '{"foo":"bar"}' },
  ]

  for (const { field, selector, value } of updates) {
    test(`can update ${field}`, async ({ page }) => {
      await login(page)
      const name = `Playwright-${randomUUID()}`
      await createType(page, name)
      await page.waitForSelector(selector)
      await page.locator(selector).fill(value)
      const saveBtn = page.getByRole('button', { name: 'Save' }).last()
      await expect(saveBtn).toBeEnabled()
      await saveBtn.click()
      await expect(page.locator(selector)).toHaveValue(value)
    })
  }
})

test('can delete a type', async ({ page }) => {
  await login(page)
  const name = `Playwright-${randomUUID()}`
  await createType(page, name)
  await page.waitForURL('**/types/*')
  await page.getByRole('button', { name: 'Delete Type' }).click()
  await page.getByRole('dialog').getByRole('button', { name: 'Delete' }).click()
  await page.waitForURL('**/types')
  await expect(page.locator(`text=${name}`)).toHaveCount(0)
})
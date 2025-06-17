import { expect } from '@playwright/test'
import { randomUUID } from 'crypto'
import { login, test } from './util'

const createGroup = async (page, name: string) => {
  await page.goto('groups')
  await page.getByRole('button', { name: 'New Group' }).click()
  await page.waitForURL('**/groups/new')
  await page.locator('#name').fill(name)
  await page.getByRole('combobox', { name: 'Permissions' }).click()
  await page.getByRole('option', { name: 'ticket:read' }).click()
  const saveBtn = page.getByRole('button', { name: 'Save' }).last()
  await expect(saveBtn).toBeEnabled()
  await saveBtn.click()
  await page.waitForURL('**/groups/g*')
}

test('groups list shows existing groups', async ({ page }) => {
  await login(page)
  await page.goto('groups')
  await expect(page.getByRole('heading', { name: 'Groups' })).toBeVisible()
})

test('can create a group', async ({ page }) => {
  await login(page)
  const name = `playwright-${randomUUID()}`
  await createGroup(page, name)
  await expect(page.locator('#name')).toHaveValue(name)
})

test.describe('update a group', () => {
  const updates = [
    { field: 'name', selector: '#name', value: 'Updated Group' },
    // { field: 'permissions', selector: '#permissions', value: 'reaction:write' },
  ]

  for (const { field, selector, value } of updates) {
    test(`can update ${field}`, async ({ page }) => {
      await login(page)
      const name = `playwright-${randomUUID()}`
      await createGroup(page, name)
      await page.waitForSelector(selector)
      await page.locator(selector).fill(value)
      const saveBtn = page.getByRole('button', { name: 'Save' }).last()
      await expect(saveBtn).toBeEnabled()
      await saveBtn.click()
      await expect(page.locator(selector)).toHaveValue(value)
    })
  }
})

test('can add a permission', async ({ page }) => {
  await login(page)
  const name = `playwright-${randomUUID()}`
  await createGroup(page, name)
  await page.getByRole('combobox', { name: 'Permissions' }).click()
  await page.getByRole('option', { name: 'reaction:write' }).click()
  const saveBtn = page.getByRole('button', { name: 'Save' }).last()
  await expect(saveBtn).toBeEnabled()
  await saveBtn.click()
  await page.waitForURL('**/groups/*')
  await expect(page.locator('#permissions')).toHaveText('Permissionsticket:readreaction:write')
})

test('can delete a group', async ({ page }) => {
  await login(page)
  const name = `playwright-${randomUUID()}`
  await createGroup(page, name)
  await page.getByRole('button', { name: 'Delete Group' }).click()
  await page.getByRole('dialog').getByRole('button', { name: 'Delete' }).click()
  await page.waitForURL('**/groups')
  await expect(page.locator(`text=${name}`)).toHaveCount(0)
})

import { expect } from '@playwright/test'
import { randomUUID } from 'crypto'
import { login, test } from './util'

const createTicket = async (page, name: string) => {
  await page.goto('tickets/incident')
  await page.getByRole('button', { name: 'New Ticket' }).click()
  await page.locator('#name').fill(name)
  await page.locator('#description').fill('Test description')
  await page.locator('#severity').selectOption('Low')
  await page.getByRole('button', { name: 'Save' }).click()
  await page.waitForURL('**/tickets/incident/incident*')
}

test('can create a ticket', async ({ page }) => {
  await login(page)
  const name = `playwright-${randomUUID()}`
  await createTicket(page, name)
  await expect(page.locator('#app #name').getByText(name)).toBeVisible()
})

test.describe('update a ticket', () => {
  const updates = [
    {
      field: 'description',
      update: async (page) => {
        await page.getByRole('button', { name: 'Edit' }).click()
        await page.getByRole('application').getByRole('textbox').fill('Updated description')
        await page.getByRole('button', { name: 'Save' }).last().click()
      },
      assert: async (page) => {
        await expect(page.getByText('Updated description')).toBeVisible()
      }
    },
    {
      field: 'severity',
      update: async (page) => {
        await page.locator('#app #severity').selectOption('High')
      },
      assert: async (page) => {
        await expect(page.locator('button').filter({ hasText: 'High' })).toBeVisible()
      }
    }
  ]

  for (const { field, update, assert } of updates) {
    test(`can update ${field}`, async ({ page }) => {
      await login(page)
      const name = `playwright-${randomUUID()}`
      await createTicket(page, name)
      await update(page)
      await assert(page)
    })
  }
})

test('can delete a ticket', async ({ page }) => {
  await login(page)
  const name = `playwright-${randomUUID()}`
  await createTicket(page, name)
  await page.getByRole('button', { name: 'Delete Ticket' }).click()
  await page.getByRole('dialog').getByRole('button', { name: 'Delete' }).click()
  await page.waitForURL('**/tickets/incident')
  await expect(page.locator(`text=${name}`)).toHaveCount(0)
})
